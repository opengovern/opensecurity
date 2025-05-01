// health.go
package core

import (
	"context"
	"net/http"
	"time"

	vaultapi "github.com/hashicorp/vault/api" // Import for HealthResponse type
	"github.com/labstack/echo/v4"
	"github.com/opengovern/og-util/pkg/vault" // Assuming HashiCorpVault constant is here
	"github.com/opengovern/opensecurity/services/core/db"
	"go.uber.org/zap"
	// <<< Ensure this import exists
)

// --- Health Status Structs & Constants ---

// ComponentHealthStatus represents the status of a single dependency.
type ComponentHealthStatus struct {
	Status string `json:"status"`          // e.g., "ok", "error", "unsealed", "disabled", "uninitialized"
	Error  string `json:"error,omitempty"` // Error message if status is "error" or provides context
}

// OverallHealthStatus represents the overall health report.
type OverallHealthStatus struct {
	OverallStatus string                           `json:"overall_status"` // "healthy" or "unhealthy"
	Timestamp     time.Time                        `json:"timestamp"`
	Components    map[string]ComponentHealthStatus `json:"components"` // Status of individual components
}

// Constants for status values
const (
	StatusOk            = "ok"
	StatusError         = "error"
	StatusHealthy       = "healthy"
	StatusUnhealthy     = "unhealthy"
	StatusDisabled      = "disabled"
	StatusUnsealed      = "unsealed"
	StatusUninitialized = "uninitialized"
)

// healthCheck handles the /health endpoint, checking critical dependencies like DB and Vault.
// It returns HTTP 200 if all critical checks pass, otherwise HTTP 503.
// Assumes h.db, h.cfg.Vault, h.logger, and h.vaultSealHandler (type *vault.HashiCorpVaultSealHandler)
// are initialized and available on the HttpHandler struct.
func (h *HttpHandler) healthCheck(c echo.Context) error {
	// Overall timeout for all checks. Adjust duration as needed.
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	overallHealthy := true // Assume healthy until a check fails
	componentStatus := make(map[string]ComponentHealthStatus)

	// --- Check Database Connectivity ---
	dbStatus := StatusOk
	dbErrorMsg := ""
	// Check if DB handler itself is initialized (basic check against zero value)
	if h.db == (db.Database{}) { // Compare against the zero value of the Database struct
		dbStatus = StatusError
		dbErrorMsg = "Database handler not initialized"
		h.logger.Error("Health check: " + dbErrorMsg)
		overallHealthy = false
	} else {
		// Use the Ping method added to db.Database
		pingCtx, pingCancel := context.WithTimeout(ctx, 3*time.Second) // Timeout for the ping itself
		err := h.db.Ping(pingCtx)
		pingCancel()

		if err != nil {
			dbStatus = StatusError
			dbErrorMsg = "Database ping failed" // Keep error generic for external response
			h.logger.Error("Health check: DB ping failed", zap.Error(err))
			overallHealthy = false
		}
		// Optional: Log success at debug level if needed
		// else { h.logger.Debug("Health check: DB ping successful.") }
	}
	componentStatus["database"] = ComponentHealthStatus{Status: dbStatus, Error: dbErrorMsg}

	// --- Check Vault Status (if HashiCorp Vault is configured) ---
	if h.cfg.Vault.Provider == vault.HashiCorpVault {
		vaultStatus := StatusOk
		vaultErrorMsg := ""
		// h.logger.Debug("Health check: Checking Vault status...") // Optional debug log

		// Check if the seal handler itself is initialized/passed correctly
		if h.vaultSealHandler == nil {
			vaultStatus = StatusError
			vaultErrorMsg = "Vault Seal Handler not available in HttpHandler"
			h.logger.Error("Health check: " + vaultErrorMsg)
			// If Vault is configured but handler isn't present, consider it unhealthy
			overallHealthy = false
		} else {
			// Call Health() method on the vaultSealHandler
			healthCtx, healthCancel := context.WithTimeout(ctx, 4*time.Second) // Shorter timeout for Vault call
			var healthResp *vaultapi.HealthResponse
			healthResp, err := h.vaultSealHandler.Health(healthCtx) // Use the correct handler

			healthCancel()

			// Process response or error
			if err != nil {
				vaultStatus = StatusError
				vaultErrorMsg = "Vault health check API call failed" // Mask specific error details
				h.logger.Error("Health check: Vault health check API call failed", zap.Error(err))
				overallHealthy = false
			} else if healthResp == nil {
				// Defensive check in case API returns nil response without error
				vaultStatus = StatusError
				vaultErrorMsg = "Vault health check returned nil response without error"
				h.logger.Error("Health check: " + vaultErrorMsg)
				overallHealthy = false
			} else if !healthResp.Initialized {
				vaultStatus = StatusUninitialized
				vaultErrorMsg = "Vault is not initialized"
				h.logger.Warn("Health check: " + vaultErrorMsg)
				overallHealthy = false // Treat as unhealthy for readiness
			} else if healthResp.Sealed {
				vaultStatus = StatusUnsealed
				vaultErrorMsg = "Vault is sealed"
				h.logger.Warn("Health check: " + vaultErrorMsg)
				overallHealthy = false // Treat as unhealthy for readiness
			} else if healthResp.Standby {
				// Vault is initialized, unsealed, but in standby mode. Usually considered healthy.
				h.logger.Info("Health check: Vault is ready (Standby).")
				// vaultStatus remains StatusOk
			} else {
				// Initialized, Unsealed, and Active - Healthy state.
				h.logger.Info("Health check: Vault is ready (Active).")
				// vaultStatus remains StatusOk
			}
		}
		componentStatus["vault"] = ComponentHealthStatus{Status: vaultStatus, Error: vaultErrorMsg}
	} else {
		// Vault is not configured or using a different provider
		componentStatus["vault"] = ComponentHealthStatus{Status: StatusDisabled}
		// h.logger.Debug("Health check: Vault check skipped (provider not HashiCorp).") // Optional log
	}

	// --- Add other critical dependency checks here ---
	// ...

	// --- Determine Overall Status and HTTP Code ---
	finalStatus := StatusHealthy
	httpCode := http.StatusOK
	if !overallHealthy {
		finalStatus = StatusUnhealthy
		httpCode = http.StatusServiceUnavailable // Use 503 for unhealthy
	}

	// --- Prepare and Return Response ---
	response := OverallHealthStatus{
		OverallStatus: finalStatus,
		Timestamp:     time.Now().UTC().Truncate(time.Second), // Use UTC, truncate ms
		Components:    componentStatus,
	}

	// Log the final outcome
	if !overallHealthy {
		// Log as Warn or Error when unhealthy
		h.logger.Warn("Health check determined service state is unhealthy", zap.String("overall_status", finalStatus), zap.Any("components", componentStatus))
	} else {
		// Log as Info when healthy
		h.logger.Info("Health check determined service state is healthy")
	}

	return c.JSON(httpCode, response)
}
