package vault

import (
	"context"
	"fmt"
	"runtime/debug" // For panic stack trace
	"time"

	"github.com/opengovern/og-util/pkg/vault"
	"github.com/opengovern/opensecurity/services/core/config"
	"go.uber.org/zap"
	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

const (
	// componentName is used for structured logging context.
	componentName = "VaultSealHandler"
	// secretName is the name of the Kubernetes Secret storing unseal keys and root token.
	secretName = "vault-unseal-keys"
	// unsealCheckInterval defines how often to attempt unsealing Vault.
	// TODO: Make this configurable via cfg.
	unsealCheckInterval = 1 * time.Minute
)

// SealHandler manages the initialization and unsealing process of a Vault instance,
// coordinating with Kubernetes secrets for key storage.
type SealHandler struct {
	logger           *zap.Logger // Logger includes namespace and component context
	cfg              config.Config
	vaultSealHandler *vault.HashiCorpVaultSealHandler // Use interface type for easier testing later
	kubeClientset    kubernetes.Interface             // Use interface type for easier testing
}

// NewSealHandler creates and initializes a new SealHandler.
func NewSealHandler(ctx context.Context, logger *zap.Logger, cfg config.Config) (*SealHandler, error) {
	// Create a logger with persistent context for this component and namespace.
	componentLogger := logger.With(
		zap.String("component", componentName),
		zap.String("namespace", cfg.OpengovernanceNamespace),
	)
	componentLogger.Debug("Initializing VaultSealHandler...")

	// Validate essential configuration
	if cfg.OpengovernanceNamespace == "" {
		componentLogger.Error("Configuration error: OpengovernanceNamespace is empty")
		return nil, fmt.Errorf("opengovernance namespace is required in configuration")
	}
	// Add other config validations as needed

	// Initialize Vault Seal Handler utility
	hashiCorpVaultSealHandler, err := vault.NewHashiCorpVaultSealHandler(ctx, componentLogger, cfg.Vault.HashiCorp)
	if err != nil {
		componentLogger.Error("Failed to create underlying HashiCorp Vault seal handler", zap.Error(err))
		return nil, fmt.Errorf("new hashicorp vault seal handler: %w", err)
	}
	componentLogger.Debug("Underlying HashiCorp Vault seal handler created")

	// Initialize Kubernetes client
	kuberConfig, err := rest.InClusterConfig()
	if err != nil {
		componentLogger.Error("Failed to get in-cluster Kubernetes config", zap.Error(err))
		return nil, fmt.Errorf("failed to get kubernetes config: %w", err)
	}
	clientset, err := kubernetes.NewForConfig(kuberConfig)
	if err != nil {
		componentLogger.Error("Failed to create Kubernetes clientset", zap.Error(err))
		return nil, fmt.Errorf("failed to create kubernetes clientset: %w", err)
	}
	componentLogger.Debug("Kubernetes clientset created")

	handler := &SealHandler{
		logger:           componentLogger,
		cfg:              cfg,
		vaultSealHandler: hashiCorpVaultSealHandler,
		kubeClientset:    clientset,
	}

	componentLogger.Info("VaultSealHandler initialized successfully")
	return handler, nil
}

// initVault checks Vault initialization status, initializes if necessary,
// and ensures the keys secret state in Kubernetes is consistent.
// Returns true if a new Vault initialization was performed, false otherwise.
// It will Fatal if Vault is initialized but the keys secret is missing,
// if storing keys fails after initialization, or if status check fails.
func (s *SealHandler) initVault(ctx context.Context) bool {
	s.logger.Debug("Entering initVault")
	defer s.logger.Debug("Exiting initVault")

	// 1. Check Vault initialization status
	s.logger.Debug("Checking Vault initialization status...")
	initRes, err := s.vaultSealHandler.TryInit(ctx)
	if err != nil {
		// Failure to even check the status is critical
		s.logger.Fatal("Failed checking vault init status", zap.Error(err))
	}

	// 2. Handle based on initialization status
	if initRes != nil {
		// Vault was NOT initialized, TryInit performed the initialization and returned keys/token.
		s.logger.Info("Vault was not initialized. Initialization performed, proceeding to store keys.")

		// Prepare the secret structure
		keysSecret := corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      secretName,
				Namespace: s.cfg.OpengovernanceNamespace,
				// Consider adding labels/annotations
			},
			// Ensure Type is Opaque or a suitable custom type if needed
			Type:       corev1.SecretTypeOpaque,
			StringData: make(map[string]string),
		}
		for i, key := range initRes.Keys {
			// CRITICAL: Avoid logging keys/token in production code, even at Debug/Fatal.
			keysSecret.StringData[fmt.Sprintf("key-%d", i)] = key
		}
		keysSecret.StringData["root-token"] = initRes.RootToken
		s.logger.Debug("Prepared Kubernetes secret structure", zap.String("secretName", secretName))

		// Attempt to Create or Update the secret
		secretsClient := s.kubeClientset.CoreV1().Secrets(s.cfg.OpengovernanceNamespace)
		s.logger.Info("Attempting to create Kubernetes secret for unseal keys", zap.String("secretName", secretName))
		_, err = secretsClient.Create(ctx, &keysSecret, metav1.CreateOptions{})
		if err != nil {
			if k8serrors.IsAlreadyExists(err) {
				// If it already exists (unexpected after fresh init, but possible race/leftover), try updating.
				s.logger.Warn("Vault unseal keys secret already existed unexpectedly after init, attempting update.", zap.String("secretName", secretName))
				_, updateErr := secretsClient.Update(ctx, &keysSecret, metav1.UpdateOptions{})
				if updateErr != nil {
					// If update fails, it's fatal - we can't guarantee keys are stored.
					s.logger.Fatal("Failed to update existing vault unseal keys secret after initialization",
						zap.String("secretName", secretName),
						zap.Error(updateErr))
				}
				s.logger.Info("Successfully updated existing vault unseal keys secret.", zap.String("secretName", secretName))
			} else {
				// Any other error on Create is fatal.
				s.logger.Fatal("Failed to create vault unseal keys secret after initialization",
					zap.String("secretName", secretName),
					zap.Error(err))
			}
		} else {
			s.logger.Info("Successfully created vault unseal keys secret.", zap.String("secretName", secretName))
		}

		s.logger.Info("Vault initialization complete and keys stored in Kubernetes secret.")
		return true // Signify that a new initialization occurred.

	} else {
		// Vault WAS already initialized.
		s.logger.Info("Vault already initialized. Verifying keys secret exists in Kubernetes.")

		// 3. Verify the keys secret exists in Kubernetes since Vault is initialized.
		secretsClient := s.kubeClientset.CoreV1().Secrets(s.cfg.OpengovernanceNamespace)
		_, err := secretsClient.Get(ctx, secretName, metav1.GetOptions{})
		if err != nil {
			if k8serrors.IsNotFound(err) {
				// Inconsistent State: Vault initialized, but keys secret is missing! This is fatal.
				s.logger.Fatal("Inconsistent State: Vault is initialized, but the required keys secret is missing in Kubernetes. Manual intervention required (restore secret or reset Vault).",
					zap.String("secretName", secretName),
					zap.Error(err)) // Error provides context "secrets ... not found"
			} else {
				// Other error trying to Get the secret (permissions, network, etc.) - also fatal.
				s.logger.Fatal("Failed to get vault unseal keys secret for verification",
					zap.String("secretName", secretName),
					zap.Error(err))
			}
		}

		// If Get succeeded without error, the secret exists. State is consistent.
		s.logger.Info("Vault is initialized and keys secret exists in Kubernetes. State is consistent.", zap.String("secretName", secretName))
		return false // Signify no new initialization occurred.
	}
}

// unsealChecker periodically attempts to unseal Vault using keys fetched
// from the Kubernetes Secret. It also handles initial Kubernetes auth setup.
func (s *SealHandler) unsealChecker(ctx context.Context, isNewInit bool, unsealed chan<- struct{}) {
	s.logger.Debug("Starting unsealChecker goroutine", zap.Bool("isNewInit", isNewInit))

	// Gracefully handle panics within the goroutine
	defer func() {
		if r := recover(); r != nil {
			// Log the panic with stack trace for detailed debugging.
			s.logger.Error("Panic recovered in unsealChecker. Restarting checker.",
				zap.Any("panicValue", r),
				zap.String("stacktrace", string(debug.Stack())), // Get stack trace
			)
			// Restart the checker. Consider adding backoff or failure counter in production.
			go s.unsealChecker(ctx, isNewInit, unsealed) // Pass original channel (might be nil if already signaled)
		} else {
			s.logger.Debug("unsealChecker goroutine finished cleanly.")
		}
	}()

	// 1. Get Unseal Keys from Kubernetes Secret
	s.logger.Debug("Attempting to get unseal keys from Kubernetes secret", zap.String("secretName", secretName))
	secretsClient := s.kubeClientset.CoreV1().Secrets(s.cfg.OpengovernanceNamespace)
	keysSecret, err := secretsClient.Get(ctx, secretName, metav1.GetOptions{})
	if err != nil {
		// If secret cannot be fetched at startup, log error and stop this attempt.
		// initVault should have ensured it exists if Vault was initialized, but it could be deleted later or permissions change.
		s.logger.Error("Failed to get vault unseal keys secret at start of unsealChecker. Stopping unseal attempts.",
			zap.String("secretName", secretName),
			zap.Error(err),
		)
		// Cannot proceed without keys, exit the goroutine.
		// Consider signaling an error state back to the main application if needed.
		return
	}
	s.logger.Debug("Successfully retrieved unseal keys secret", zap.String("secretName", secretName))

	// Extract keys, ignoring the root token for unsealing purposes.
	keys := make([]string, 0, len(keysSecret.Data))
	var rootTokenBytes []byte
	var rootTokenFound bool
	for k, v := range keysSecret.Data {
		if k == "root-token" {
			rootTokenBytes = v
			rootTokenFound = true
			continue
		}
		// Assume other keys are unseal keys (key-0, key-1, etc.)
		keys = append(keys, string(v))
	}
	s.logger.Debug("Extracted unseal keys from secret", zap.Int("keyCount", len(keys)), zap.Bool("rootTokenFound", rootTokenFound))
	if len(keys) == 0 {
		s.logger.Error("No unseal keys found within the Kubernetes secret data. Cannot proceed.", zap.String("secretName", secretName))
		return
	}

	// 2. Initial Unseal Attempt & Kube Auth Setup (if needed)
	if unsealed != nil { // Only attempt initial setup if the channel is provided and open
		s.logger.Debug("Performing initial unseal attempt...")
		err = s.vaultSealHandler.TryUnseal(ctx, keys)
		if err != nil {
			s.logger.Error("Initial unseal attempt failed", zap.Error(err))
			// Proceed to ticker loop, will retry periodically.
		} else {
			s.logger.Info("Initial unseal successful")

			// Setup Kubernetes auth using the root token *only once* after initial unseal.
			if rootTokenFound && len(rootTokenBytes) > 0 {
				s.logger.Info("Attempting to setup Kubernetes auth in Vault...")
				kubeAuthErr := s.vaultSealHandler.SetupKuberAuth(ctx, string(rootTokenBytes))
				if kubeAuthErr != nil {
					s.logger.Error("Failed to setup Kubernetes auth", zap.Error(kubeAuthErr))
					// Log error but continue; unseal is more critical.
				} else {
					s.logger.Info("Kubernetes auth setup successful or already configured.")
				}
			} else {
				s.logger.Warn("Root token not found or empty in secret, skipping Kubernetes auth setup.", zap.String("secretName", secretName))
			}

			// Signal that initial unseal (and optional Kube auth setup) is complete.
			s.logger.Debug("Signaling initial unseal complete.")
			unsealed <- struct{}{}
			close(unsealed)
			unsealed = nil // Prevent further sends/closes
		}
	} else {
		s.logger.Debug("Initial unseal channel was nil, skipping initial unseal specific logic.")
	}

	// 3. Periodic Unseal Check Loop
	ticker := time.NewTicker(unsealCheckInterval)
	defer ticker.Stop()
	s.logger.Info("Starting periodic unseal check loop", zap.Duration("interval", unsealCheckInterval))

	for {
		select {
		case <-ctx.Done():
			// Context canceled (e.g., application shutting down)
			s.logger.Info("Context cancelled, stopping unseal checker loop.", zap.Error(ctx.Err()))
			return // Exit the goroutine

		case tickTime := <-ticker.C:
			// Periodically try to unseal. TryUnseal should be idempotent if already unsealed.
			s.logger.Debug("Ticker received, attempting periodic unseal check", zap.Time("tickTime", tickTime))
			err = s.vaultSealHandler.TryUnseal(ctx, keys)
			if err != nil {
				// Log periodic failures but continue checking.
				s.logger.Error("Periodic unseal attempt failed", zap.Error(err))
				// Consider adding metrics/alerting for persistent failures here.
				continue // Wait for next tick
			}
			s.logger.Debug("Periodic unseal check successful (or Vault already unsealed)")
		}
	}
}

// Start initiates the Vault initialization and unsealing process in the background.
// It blocks until the initial unseal signal is received.
func (s *SealHandler) Start(ctx context.Context) {
	s.logger.Debug("Entering Start")

	// Perform initialization checks and potentially initialize Vault + create secret.
	// This will Fatal if state is inconsistent or critical errors occur.
	isNewInit := s.initVault(ctx)
	s.logger.Debug("initVault completed", zap.Bool("isNewInit", isNewInit))

	// Channel to signal completion of the *initial* unseal.
	unsealChan := make(chan struct{})

	s.logger.Info("Starting Vault unseal checker background goroutine...")
	go s.unsealChecker(ctx, isNewInit, unsealChan)

	s.logger.Debug("Waiting for initial unseal signal...")
	// Block until the unsealChecker signals the initial unseal is done.
	// Add a timeout here if needed for application startup guarantees.
	// select {
	// case <-unsealChan:
	// 	s.logger.Info("Initial unseal signal received. Vault should be ready.")
	// case <-time.After(5 * time.Minute): // Example timeout
	// 	s.logger.Fatal("Timeout waiting for initial Vault unseal signal.")
	// case <-ctx.Done(): // Handle application context cancellation during wait
	//  s.logger.Warn("Context cancelled while waiting for initial unseal signal.", zap.Error(ctx.Err()))
	//  return // Or handle appropriately
	// }
	<-unsealChan // Using the simpler blocking wait for now
	s.logger.Info("Initial unseal signal received. Vault should be ready.")
	s.logger.Debug("Exiting Start")
}
