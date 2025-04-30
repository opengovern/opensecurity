// Package authcache provides an embedded in-process cache for user authorization info,
// backed by in-memory go-cache and instrumented with Prometheus metrics.
package authcache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	// Use correct import path for go-cache
	gocache "github.com/patrickmn/go-cache"

	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/zap"
)

// --- Configurable constants ---

// DefaultCacheTTL is the duration that entries live in the cache by default.
const DefaultCacheTTL = 5 * time.Minute

// DefaultCleanupInterval is how often expired items are purged from the cache.
// Making it longer than TTL is standard practice for go-cache.
const DefaultCleanupInterval = 10 * time.Minute

// UserEmailKeyPrefix is the prefix for all cache keys storing user info by email.
const UserEmailKeyPrefix = "user:email:"

// externalIDKeyPrefix is the prefix for cache keys mapping external ID to email.
const externalIDKeyPrefix = "extid:email:" // <-- New prefix for the ID map cache

// --- Errors ---

// ErrUserInfoNotFound indicates that a user's info was not present in cache,
// or the data found was invalid/corrupted.
var ErrUserInfoNotFound = errors.New("authcache: user info not found")

// ErrEmailNotFound indicates that an email mapping was not found for an external ID.
var ErrEmailNotFound = errors.New("authcache: email not found for external ID") // <-- New error

// --- Prometheus metrics ---

var (
	metricHits = prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: "authcache",
		Name:      "userinfo_hits_total", // Added suffix for clarity
		Help:      "Total number of user info cache hits (by email).",
	})
	metricMisses = prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: "authcache",
		Name:      "userinfo_misses_total", // Added suffix for clarity
		Help:      "Total number of user info cache misses (by email).",
	})
	metricIdMapHits = prometheus.NewCounter(prometheus.CounterOpts{ // <-- New metric
		Namespace: "authcache",
		Name:      "id_map_hits_total",
		Help:      "Total number of ExternalID-to-Email cache hits.",
	})
	metricIdMapMisses = prometheus.NewCounter(prometheus.CounterOpts{ // <-- New metric
		Namespace: "authcache",
		Name:      "id_map_misses_total",
		Help:      "Total number of ExternalID-to-Email cache misses.",
	})
	metricErrors = prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: "authcache",
		Name:      "errors_total",
		Help:      "Total number of cache operation errors (e.g., marshal/unmarshal).",
	})
)

func init() {
	// Register metrics with Prometheus default registry.
	prometheus.MustRegister(metricHits, metricMisses, metricIdMapHits, metricIdMapMisses, metricErrors) // Added new metrics
}

// --- Data structures ---

// CachedUserInfo holds the minimal authorization details for a user.
type CachedUserInfo struct {
	ID         uint      `json:"id"`          // Internal DB ID
	Role       string    `json:"role"`        // User role (e.g., "admin", "viewer")
	ExternalID string    `json:"external_id"` // External subject/UUID
	IsActive   bool      `json:"is_active"`   // Account enabled flag
	LastLogin  time.Time `json:"last_login"`  // Needed for direct last login check
	// Add FullName/CreatedAt if needed by /me cache logic
	// FullName   string    `json:"full_name"`
	// CreatedAt  time.Time `json:"created_at"`
}

// --- CacheClient abstraction ---

// CacheClient defines the methods required to interact with an underlying cache.
// Context is included for potential future use or different implementations.
type CacheClient interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key, value string, ttl time.Duration) error
	Del(ctx context.Context, key string) error
	Close() error // For potential cleanup in other implementations
}

// --- In-memory adapter (go-cache) ---

// inMemoryClient wraps a go-cache.Cache to satisfy CacheClient.
type inMemoryClient struct {
	inner *gocache.Cache
}

// NewInMemoryClient returns a CacheClient backed by go-cache.
// cleanupInterval controls how often expired items are purged.
func NewInMemoryClient(defaultTTL, cleanupInterval time.Duration) CacheClient {
	if defaultTTL <= 0 {
		defaultTTL = DefaultCacheTTL
	}
	if cleanupInterval <= 0 {
		cleanupInterval = DefaultCleanupInterval
	}
	return &inMemoryClient{
		inner: gocache.New(defaultTTL, cleanupInterval),
	}
}

// Get retrieves a string value from the cache.
// Returns ErrUserInfoNotFound if the key doesn't exist or the value is not a string.
func (c *inMemoryClient) Get(_ context.Context, key string) (string, error) {
	value, found := c.inner.Get(key)
	if !found {
		// Use the specific error defined for this implementation/package
		return "", ErrUserInfoNotFound
	}
	stringValue, ok := value.(string)
	if !ok {
		c.inner.Delete(key)
		return "", fmt.Errorf("invalid type found in cache for key %s: expected string, got %T", key, value)
	}
	return stringValue, nil
}

// Set stores a string value in the cache with a specific TTL.
func (c *inMemoryClient) Set(_ context.Context, key, value string, ttl time.Duration) error {
	expiration := ttl
	if ttl == 0 {
		expiration = gocache.DefaultExpiration
	}
	c.inner.Set(key, value, expiration)
	return nil
}

// Del removes a key from the cache.
func (c *inMemoryClient) Del(_ context.Context, key string) error {
	c.inner.Delete(key)
	return nil
}

// Close is a no-op for the basic in-memory go-cache.
func (c *inMemoryClient) Close() error {
	// go-cache Stop() method exists for janitor goroutine, but typically not needed for simple use.
	// c.inner.Flush() // Optionally flush if needed before shutdown
	return nil
}

// --- AuthCacheService ---

// AuthCacheService manages caching of user authorization info.
type AuthCacheService struct {
	userInfoCache  CacheClient   // Cache for user:email:<email> -> CachedUserInfo JSON string
	idToEmailCache CacheClient   // Cache for extid:email:<extid> -> email string
	logger         *zap.Logger   // structured logger
	ttl            time.Duration // default entry TTL used by AddUserToCache
}

// NewAuthCacheService creates and initializes the cache service.
// It uses two in-memory go-cache instances under the hood.
func NewAuthCacheService(logger *zap.Logger) (*AuthCacheService, error) {
	if logger == nil {
		return nil, errors.New("logger cannot be nil")
	}

	// Initialize the main user info cache
	userInfoClient := NewInMemoryClient(DefaultCacheTTL, DefaultCleanupInterval)
	// Initialize the ID-to-Email mapping cache (could potentially have a different TTL/cleanup)
	idToEmailClient := NewInMemoryClient(DefaultCacheTTL*2, DefaultCleanupInterval*2) // Example: longer TTL for mapping

	namedLogger := logger.Named("authcache")
	namedLogger.Info("AuthCacheService initialized (using go-cache)",
		zap.Duration("user_info_ttl", DefaultCacheTTL),
		zap.Duration("id_map_ttl", DefaultCacheTTL*2), // Log TTL for the new cache
	)

	return &AuthCacheService{
		userInfoCache:  userInfoClient,
		idToEmailCache: idToEmailClient, // Assign the second cache client
		logger:         namedLogger,
		ttl:            DefaultCacheTTL, // Store the default TTL for AddUserToCache
	}, nil
}

// formatUserEmailKey builds the cache key for user info based on email.
func formatUserEmailKey(email string) string {
	return UserEmailKeyPrefix + strings.ToLower(strings.TrimSpace(email))
}

// formatExternalIDKey builds the cache key for the ID-to-Email map.
func formatExternalIDKey(externalID string) string {
	return externalIDKeyPrefix + externalID // Assuming externalID is already consistent case/format
}

// GetUser fetches a user's authorization info from cache using their email.
// Returns ErrUserInfoNotFound if no valid entry exists.
func (s *AuthCacheService) GetUser(ctx context.Context, email string) (*CachedUserInfo, error) {
	if email == "" {
		return nil, errors.New("email cannot be empty")
	}
	key := formatUserEmailKey(email)
	s.logger.Debug("Attempting UserInfo cache GET", zap.String("key", key))

	rawValue, err := s.userInfoCache.Get(ctx, key) // Use userInfoCache
	if err != nil {
		if errors.Is(err, ErrUserInfoNotFound) { // Check against our specific error
			metricMisses.Inc()
			s.logger.Debug("UserInfo cache miss", zap.String("key", key))
			return nil, ErrUserInfoNotFound
		}
		metricErrors.Inc()
		s.logger.Error("UserInfo cache GET error", zap.String("key", key), zap.Error(err))
		return nil, fmt.Errorf("user info cache GET failed for %q: %w", key, err)
	}

	var info CachedUserInfo
	if err := json.Unmarshal([]byte(rawValue), &info); err != nil {
		metricErrors.Inc()
		s.logger.Error("UserInfo cache data unmarshal error",
			zap.String("key", key),
			zap.String("rawData", rawValue),
			zap.Error(err),
		)
		_ = s.userInfoCache.Del(ctx, key) // Delete corrupted entry from userInfoCache
		return nil, ErrUserInfoNotFound   // Treat unmarshal error as not found
	}

	metricHits.Inc()
	s.logger.Debug("UserInfo cache hit", zap.String("key", key))
	return &info, nil
}

// AddUserToCache stores a user's authorization info AND their ExternalID->Email mapping.
func (s *AuthCacheService) AddUserToCache(ctx context.Context, email string, info *CachedUserInfo) error {
	if email == "" {
		return errors.New("email cannot be empty for cache add")
	}
	if info == nil {
		return errors.New("cannot cache nil user info")
	}
	if info.ExternalID == "" || info.Role == "" {
		return errors.New("invalid user info: missing ExternalID or Role")
	}

	// 1. Add/Update main user info cache
	userInfoKey := formatUserEmailKey(email)
	s.logger.Debug("Attempting UserInfo cache SET", zap.String("key", userInfoKey), zap.Duration("ttl", s.ttl))
	data, err := json.Marshal(info)
	if err != nil {
		metricErrors.Inc()
		s.logger.Error("UserInfo cache marshal error before SET", zap.String("key", userInfoKey), zap.Error(err))
		return fmt.Errorf("cache marshal failed for %q: %w", userInfoKey, err)
	}
	if err := s.userInfoCache.Set(ctx, userInfoKey, string(data), s.ttl); err != nil {
		metricErrors.Inc()
		s.logger.Error("UserInfo cache SET error", zap.String("key", userInfoKey), zap.Error(err))
		// Continue to set the ID map even if this fails? Or return error? Returning error for now.
		return fmt.Errorf("user info cache SET failed for %q: %w", userInfoKey, err)
	}
	s.logger.Debug("UserInfo cache SET successful", zap.String("key", userInfoKey))

	// 2. Add/Update ExternalID -> Email mapping cache
	err = s.SetExternalIDToEmail(ctx, info.ExternalID, email)
	if err != nil {
		// Log the error but don't necessarily fail the whole operation,
		// as the main user info was cached. The ID map can be repopulated later.
		s.logger.Error("Failed to set ExternalID-to-Email mapping in cache",
			zap.String("externalID", info.ExternalID),
			zap.String("email", email),
			zap.Error(err))
		// Decide whether to return this error or just log it. Logging only for now.
	}

	return nil
}

// RemoveUserFromCache deletes a user's entries from BOTH caches.
// It first tries to fetch the user info by email to get the ExternalID.
func (s *AuthCacheService) RemoveUserFromCache(ctx context.Context, email string) error {
	if email == "" {
		return errors.New("email cannot be empty for cache remove")
	}

	userInfoKey := formatUserEmailKey(email)
	externalID := "" // Initialize externalID

	// 1. Attempt to get user info to find the ExternalID
	s.logger.Debug("Fetching user info before cache DEL to get ExternalID", zap.String("emailKey", userInfoKey))
	cachedInfo, err := s.GetUser(ctx, email) // Use GetUser which handles errors/misses
	if err != nil {
		if errors.Is(err, ErrUserInfoNotFound) {
			// User info not found by email, likely ID map entry is also gone or stale.
			// Log this and proceed to attempt deletion by email key only.
			s.logger.Debug("UserInfo not found during pre-delete check, proceeding with email key deletion only", zap.String("emailKey", userInfoKey))
		} else {
			// Unexpected error fetching user info, log it but still attempt deletion by email key
			s.logger.Error("Error fetching user info before cache DEL", zap.String("emailKey", userInfoKey), zap.Error(err))
			// Do not return here, still try to delete the email key entry
		}
	} else if cachedInfo != nil {
		// Found user info, store the external ID
		externalID = cachedInfo.ExternalID
		s.logger.Debug("Found ExternalID during pre-delete check", zap.String("externalID", externalID))
	}

	// 2. Remove UserInfo cache entry (by email)
	s.logger.Debug("Attempting UserInfo cache DEL", zap.String("key", userInfoKey))
	if errDel := s.userInfoCache.Del(ctx, userInfoKey); errDel != nil {
		// Log unexpected errors (go-cache Del usually doesn't error)
		metricErrors.Inc()
		s.logger.Error("UserInfo cache DEL error (unexpected)", zap.String("key", userInfoKey), zap.Error(errDel))
		// Do not return error, proceed to delete ID map entry if possible
	} else {
		s.logger.Debug("UserInfo cache DEL complete", zap.String("key", userInfoKey))
	}

	// 3. Remove ExternalID -> Email mapping if we found the ExternalID
	if externalID != "" {
		externalIDKey := formatExternalIDKey(externalID)
		s.logger.Debug("Attempting ID->Email cache DEL", zap.String("key", externalIDKey))
		if errDel := s.idToEmailCache.Del(ctx, externalIDKey); errDel != nil {
			metricErrors.Inc()
			s.logger.Error("ID->Email cache DEL error (unexpected)", zap.String("key", externalIDKey), zap.Error(errDel))
			// Do not return error
		} else {
			s.logger.Debug("ID->Email cache DEL complete", zap.String("key", externalIDKey))
		}
	} else {
		s.logger.Debug("Skipping ID->Email cache DEL because ExternalID was not retrieved", zap.String("email", email))
	}

	return nil // Return nil as deletion is best-effort
}

// --- New methods for ID -> Email mapping ---

// GetEmailByExternalID fetches the email associated with an external ID from the cache.
// Returns ErrEmailNotFound if the mapping doesn't exist.
func (s *AuthCacheService) GetEmailByExternalID(ctx context.Context, externalID string) (string, error) {
	if externalID == "" {
		return "", errors.New("external ID cannot be empty")
	}
	key := formatExternalIDKey(externalID)
	s.logger.Debug("Attempting ID->Email cache GET", zap.String("key", key))

	email, err := s.idToEmailCache.Get(ctx, key) // Use idToEmailCache
	if err != nil {
		// Check if the error is our specific not found error from the adapter
		if errors.Is(err, ErrUserInfoNotFound) { // Using ErrUserInfoNotFound from adapter here
			metricIdMapMisses.Inc()
			s.logger.Debug("ID->Email cache miss", zap.String("key", key))
			return "", ErrEmailNotFound // Return our specific email not found error
		}
		// Log other unexpected errors
		metricErrors.Inc()
		s.logger.Error("ID->Email cache GET error", zap.String("key", key), zap.Error(err))
		return "", fmt.Errorf("id map cache GET failed for %q: %w", key, err)
	}

	metricIdMapHits.Inc()
	s.logger.Debug("ID->Email cache hit", zap.String("key", key))
	return email, nil
}

// SetExternalIDToEmail stores the mapping from external ID to email in the cache.
// Uses a potentially longer TTL suitable for mappings.
func (s *AuthCacheService) SetExternalIDToEmail(ctx context.Context, externalID, email string) error {
	if externalID == "" || email == "" {
		return errors.New("external ID and email cannot be empty for mapping cache")
	}
	key := formatExternalIDKey(externalID)
	// Use a potentially longer TTL for this mapping cache
	mappingTTL := s.ttl * 2 // Example: double the user info TTL
	if mappingTTL <= 0 {
		mappingTTL = DefaultCacheTTL * 2 // Ensure positive
	}

	s.logger.Debug("Attempting ID->Email cache SET", zap.String("key", key), zap.Duration("ttl", mappingTTL))

	// Store the email string directly
	if err := s.idToEmailCache.Set(ctx, key, email, mappingTTL); err != nil { // Use idToEmailCache
		metricErrors.Inc()
		s.logger.Error("ID->Email cache SET error", zap.String("key", key), zap.Error(err))
		return fmt.Errorf("id map cache SET failed for %q: %w", key, err)
	}

	s.logger.Debug("ID->Email cache SET successful", zap.String("key", key))
	return nil
}

// Close shuts down the cache service and underlying client.
func (s *AuthCacheService) Close() error {
	s.logger.Info("Shutting down AuthCacheService...")
	// Close both underlying clients if they have Close methods
	var errs []string
	if err := s.userInfoCache.Close(); err != nil {
		s.logger.Error("UserInfo cache client close error", zap.Error(err))
		errs = append(errs, fmt.Sprintf("user info cache: %v", err))
	}
	if err := s.idToEmailCache.Close(); err != nil {
		s.logger.Error("ID->Email cache client close error", zap.Error(err))
		errs = append(errs, fmt.Sprintf("id map cache: %v", err))
	}

	if len(errs) > 0 {
		return fmt.Errorf("cache client close failed: %s", strings.Join(errs, "; "))
	}

	s.logger.Info("AuthCacheService closed successfully.")
	return nil
}
