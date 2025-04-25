// Package client provides a Go client library for interacting with the
// OpenSecurity Authentication Service's HTTP API. It simplifies making requests
// to common endpoints like listing users and connectors.
package client

import (
	"fmt"
	"net/http" // Import standard http package

	"github.com/labstack/echo/v4" // Used for creating HTTP errors
	// Import the shared httpclient utility for making requests
	"github.com/opengovern/og-util/pkg/httpclient"
	// Import the API definitions from the auth service to use response types
	"github.com/opengovern/opensecurity/services/auth/api"
)

// AuthServiceClient defines the interface for interacting with the Auth Service API.
// It abstracts the underlying HTTP calls.
type AuthServiceClient interface {
	// ListUsers retrieves a list of all users from the Auth Service.
	// It requires an httpclient.Context which contains the necessary request context
	// (like context.Context and potentially authorization headers).
	ListUsers(ctx *httpclient.Context) ([]api.GetUsersResponse, error)

	// GetConnectors retrieves a list of configured identity provider connectors
	// from the Auth Service. It requires an httpclient.Context.
	GetConnectors(ctx *httpclient.Context) ([]api.GetConnectorsResponse, error)
}

// authClient is the concrete implementation of the AuthServiceClient interface.
// It holds the base URL for the target Auth Service API.
type authClient struct {
	baseURL string // Base URL of the auth service (e.g., "http://auth-service.namespace.svc.cluster.local:5555")
}

// NewAuthClient creates a new instance of the AuthServiceClient.
// It takes the baseURL of the target Auth Service API as an argument.
func NewAuthClient(baseURL string) AuthServiceClient {
	// Return a pointer to the concrete implementation, satisfying the interface.
	return &authClient{baseURL: baseURL}
}

// ListUsers implements the AuthServiceClient interface method to fetch users.
func (s *authClient) ListUsers(ctx *httpclient.Context) ([]api.GetUsersResponse, error) {
	// Construct the full URL for the /users endpoint.
	url := fmt.Sprintf("%s/api/v1/users", s.baseURL)

	// Prepare a slice to hold the response data.
	var users []api.GetUsersResponse

	// Use the shared httpclient utility to perform the GET request.
	// - ctx.Ctx: The context.Context for cancellation/deadlines.
	// - http.MethodGet: The HTTP method.
	// - url: The target URL.
	// - ctx.ToHeaders(): Extracts relevant headers (like Authorization) from the httpclient.Context.
	// - nil: No request body for GET.
	// - &users: Pointer to the slice where the JSON response body should be unmarshaled.
	if statusCode, err := httpclient.DoRequest(ctx.Ctx, http.MethodGet, url, ctx.ToHeaders(), nil, &users); err != nil {
		// Check if the error is likely a client-side error (4xx).
		if 400 <= statusCode && statusCode < 500 {
			// Wrap the original error message in an echo.HTTPError with the received status code.
			// This allows upstream handlers to easily return the correct HTTP status.
			return nil, echo.NewHTTPError(statusCode, fmt.Sprintf("auth service client error: %v", err))
		}
		// For other errors (5xx, network errors, etc.), return the error directly.
		return nil, fmt.Errorf("failed to list users from auth service: %w", err)
	}

	// If the request was successful, return the populated slice of users and nil error.
	return users, nil
}

// GetConnectors implements the AuthServiceClient interface method to fetch connectors.
func (s *authClient) GetConnectors(ctx *httpclient.Context) ([]api.GetConnectorsResponse, error) {
	// Construct the full URL for the /connectors endpoint.
	url := fmt.Sprintf("%s/api/v1/connectors", s.baseURL)

	// Prepare a slice to hold the response data.
	var connectors []api.GetConnectorsResponse

	// Use the shared httpclient utility to perform the GET request.
	if statusCode, err := httpclient.DoRequest(ctx.Ctx, http.MethodGet, url, ctx.ToHeaders(), nil, &connectors); err != nil {
		// Handle client-side errors (4xx).
		if 400 <= statusCode && statusCode < 500 {
			return nil, echo.NewHTTPError(statusCode, fmt.Sprintf("auth service client error: %v", err))
		}
		// Handle other errors.
		return nil, fmt.Errorf("failed to get connectors from auth service: %w", err)
	}

	// Return the populated slice of connectors on success.
	return connectors, nil
}
