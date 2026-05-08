package middlewares

import (
	"boilerplate-api/lib/constants"
	"boilerplate-api/services"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockAuthService is a mock implementation of the go_firebase_service.IAuthService for testing.
type MockAuthService struct {
	mock.Mock
	services.IFirebaseMiddlewareService
}

func (m *MockAuthService) VerifyToken(idToken string) (*services.FirebaseToken, *services.AuthErrorResponse) {
	args := m.Called(idToken)
	if args.Get(0) == nil {
		return nil, args.Get(1).(*services.AuthErrorResponse)
	}
	return args.Get(0).(*services.FirebaseToken), nil
}

func TestFirebaseAuthMiddleware_HandleAuth(t *testing.T) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Test case 1: Successful authentication with a single SetClaims function
	t.Run("Successful auth with single SetClaims", func(t *testing.T) {
		// Setup
		mockService := new(MockAuthService)
		middleware := NewFirebaseAuthMiddleware(mockService)

		// Mock VerifyToken to return a valid token
		mockService.On("VerifyToken", "test_token").Return(&services.FirebaseToken{UID: "test_uid", Claims: map[string]interface{}{
			constants.Roles.Key:                constants.Roles.User,
			constants.Claims.UserId.ToString(): 123.0,
		}}, nil)

		// Create a new Gin router
		r := gin.New()
		r.Use(middleware.HandleAuth())
		r.GET("/test", func(c *gin.Context) {
			c.Status(http.StatusOK)
		})

		// Create a new HTTP request
		req, _ := http.NewRequest(http.MethodGet, "/test", nil)
		req.Header.Set(constants.Headers.Authorization.ToString(), "Bearer test_token")

		// Create a new HTTP recorder
		w := httptest.NewRecorder()

		// Serve the HTTP request
		r.ServeHTTP(w, req)

		// Assertions
		assert.Equal(t, http.StatusOK, w.Code)
	})

	// Test case 2: Successful authentication with multiple SetClaims functions (one passes)
	t.Run("Successful auth with multiple SetClaims (one passes)", func(t *testing.T) {
		// Setup
		mockService := new(MockAuthService)
		middleware := NewFirebaseAuthMiddleware(mockService)

		// Mock VerifyToken to return a valid token for a super admin
		mockService.On("VerifyToken", "test_token").Return(&services.FirebaseToken{UID: "test_uid", Claims: map[string]interface{}{
			constants.Roles.Key:                 constants.Roles.SuperAdmin.ToString(),
			constants.Claims.AdminId.ToString(): 123.0,
		}}, nil)

		// Create a new Gin router
		r := gin.New()
		// This should pass because AdminAuth will succeed
		r.Use(middleware.HandleAuth())
		r.GET("/test", func(c *gin.Context) {
			c.Status(http.StatusOK)
		})

		// Create a new HTTP request
		req, _ := http.NewRequest(http.MethodGet, "/test", nil)
		req.Header.Set(constants.Headers.Authorization.ToString(), "Bearer test_token")

		// Create a new HTTP recorder
		w := httptest.NewRecorder()

		// Serve the HTTP request
		r.ServeHTTP(w, req)

		// Assertions
		assert.Equal(t, http.StatusOK, w.Code)
	})

	// Test case 3: Failed authentication with multiple SetClaims functions (all fail)
	t.Run("Failed auth with multiple SetClaims (all fail)", func(t *testing.T) {
		// Setup
		mockService := new(MockAuthService)
		middleware := NewFirebaseAuthMiddleware(mockService)

		// Mock VerifyToken to return a valid token for a user
		mockService.On("VerifyToken", "test_token").Return(&services.FirebaseToken{UID: "test_uid", Claims: map[string]interface{}{
			constants.Roles.Key: constants.Roles.User.ToString(),
		}}, nil)

		// Create a new Gin router
		r := gin.New()
		// This should fail because the user role is not a company PIC or a super admin
		r.Use(middleware.HandleAuth())
		r.GET("/test", func(c *gin.Context) {
			c.Status(http.StatusOK)
		})

		// Create a new HTTP request
		req, _ := http.NewRequest(http.MethodGet, "/test", nil)
		req.Header.Set(constants.Headers.Authorization.ToString(), "Bearer test_token")

		// Create a new HTTP recorder
		w := httptest.NewRecorder()

		// Serve the HTTP request
		r.ServeHTTP(w, req)

		// Assertions
		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	// Test case 4: Failed authentication due to invalid token
	t.Run("Failed auth with invalid token", func(t *testing.T) {
		// Setup
		mockService := new(MockAuthService)
		middleware := NewFirebaseAuthMiddleware(mockService)

		// Mock VerifyToken to return an error
		mockService.On("VerifyToken", "invalid_token").Return(nil, &services.AuthErrorResponse{Message: "invalid token"})

		// Create a new Gin router
		r := gin.New()
		r.Use(middleware.HandleAuth())
		r.GET("/test", func(c *gin.Context) {
			c.Status(http.StatusOK)
		})

		// Create a new HTTP request
		req, _ := http.NewRequest(http.MethodGet, "/test", nil)
		req.Header.Set(constants.Headers.Authorization.ToString(), "Bearer invalid_token")

		// Create a new HTTP recorder
		w := httptest.NewRecorder()

		// Serve the HTTP request
		r.ServeHTTP(w, req)

		// Assertions
		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	// Test case 5: Failed authentication with multiple SetClaims functions with roles, but the token is of the user
	t.Run("Failed auth with multiple SetClaims functions with roles, but the token is of the user", func(t *testing.T) {
		// Setup
		mockService := new(MockAuthService)
		middleware := NewFirebaseAuthMiddleware(mockService)

		// Mock VerifyToken to return a valid token for a super admin
		mockService.On("VerifyToken", "test_token").Return(&services.FirebaseToken{UID: "test_uid", Claims: map[string]interface{}{
			constants.Roles.Key:                 constants.Roles.User.ToString(),
			constants.Claims.AdminId.ToString(): 123.0,
		}}, nil)

		// Create a new Gin router
		r := gin.New()
		// This should pass because AdminAuth will succeed
		r.Use(middleware.HandleAuth())
		r.GET("/test", func(c *gin.Context) {
			c.Status(http.StatusOK)
		})

		// Create a new HTTP request
		req, _ := http.NewRequest(http.MethodGet, "/test", nil)
		req.Header.Set(constants.Headers.Authorization.ToString(), "Bearer test_token")

		// Create a new HTTP recorder
		w := httptest.NewRecorder()

		// Serve the HTTP request
		r.ServeHTTP(w, req)

		// Assertions
		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}
