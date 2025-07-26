package infrastructure

// imports
import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/natnael-eyuel-dev/Task-Management-Unit-Test/Infrastructure/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

// test suite for AuthMiddleware
type AuthMiddlewareTestSuite struct {
	suite.Suite
	mockJWTService  *mock_infrastructure.MockJWTService        // mock JWT service instance
	router          *gin.Engine          	   				   // gin router for testing
}

// initializes the test environment before each test
func (suite *AuthMiddlewareTestSuite) SetupTest() {
	gin.SetMode(gin.TestMode)                     					     // set gin to test mode
	suite.mockJWTService = new(mock_infrastructure.MockJWTService)       // create new mock JWT service
	suite.router = gin.New()                      					     // create new gin router
}

// tests the AuthHandler with a valid token
func (suite *AuthMiddlewareTestSuite) TestAuthHandler_ValidToken() {
	
	// setup test claims
	claims := jwt.MapClaims{
		"sub":      "user123",
		"username": "testuser",
		"role":     "admin",
	}
	
	// create a valid token
	token := &jwt.Token{
		Valid:  true,
		Claims: claims,
	}
	
	// mock the ValidateToken method
	suite.mockJWTService.On("ValidateToken", "valid.token").Return(token, nil)

	// setup router with auth middleware
	auth := NewAuthMiddleware(suite.mockJWTService)
	suite.router.Use(auth.Handler())
	suite.router.GET("/protected", func(c *gin.Context) {
		// extract claims from context
		userID, _ := c.Get("userID")
		username, _ := c.Get("username")
		role, _ := c.Get("role")

		// return user info in response
		c.JSON(http.StatusOK, gin.H{
			"userID":   userID,
			"username": username,
			"role":     role,
		})
	})

	// create test request with valid token
	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "valid.token")      // set the authorization header with valid token
	w := httptest.NewRecorder()

	// serve the request using the router
	suite.router.ServeHTTP(w, req)

	// verify response
	require.Equal(suite.T(), http.StatusOK, w.Code)       // status should be 200
	suite.Contains(w.Body.String(), "user123")            // check userID in response
	suite.Contains(w.Body.String(), "testuser")           // check username in response
	suite.Contains(w.Body.String(), "admin")              // check role in response
}

// tests the AuthHandler with missing token
func (suite *AuthMiddlewareTestSuite) TestAuthHandler_MissingToken() {
	
	// setup router with auth middleware
	auth := NewAuthMiddleware(suite.mockJWTService)
	// use the auth middleware
	suite.router.Use(auth.Handler())
	// define a protected route
	suite.router.GET("/protected", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	// create test request without token
	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	w := httptest.NewRecorder()

	// serve the request using the router
	suite.router.ServeHTTP(w, req)          

	// verify unauthorized response
	assert.Equal(suite.T(), http.StatusUnauthorized, w.Code)                          // status should be 404
	assert.Contains(suite.T(), w.Body.String(), "authorization header required")      // check response body
}

// tests the AuthHandler with invalid token
func (suite *AuthMiddlewareTestSuite) TestAuthHandler_InvalidToken() {
	
	// mock the ValidateToken method to return error
	suite.mockJWTService.
		On("ValidateToken", "invalid.token").
		Return(nil, errors.New("invalid token"))

	// setup router with auth middleware
	auth := NewAuthMiddleware(suite.mockJWTService)
	// use the auth middleware
	suite.router.Use(auth.Handler())
	// define a protected route
	suite.router.GET("/protected", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	// create test request with invalid token
	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "invalid.token")      // set the auth header with invalid token
	w := httptest.NewRecorder()

	// serve the request using the router
	suite.router.ServeHTTP(w, req)          

	// verify unauthorized response
	assert.Equal(suite.T(), http.StatusUnauthorized, w.Code)            // status should be 404
	assert.Contains(suite.T(), w.Body.String(), "invalid token")        // check response body
}

// tests the AdminOnly middleware with admin role
func (suite *AuthMiddlewareTestSuite) TestAdminOnly_AllowAdmin() {
	
	// setup router with admin role in context
	suite.router.Use(func(c *gin.Context) {
		c.Set("role", "admin")
	})
	// use the AdminOnly middleware
	suite.router.Use(AdminOnly())
	// define an admin route
	suite.router.GET("/admin", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "welcome admin"})
	})

	// create test request
	req := httptest.NewRequest(http.MethodGet, "/admin", nil)
	w := httptest.NewRecorder()

	// serve the request using the router
	suite.router.ServeHTTP(w, req)           
 
	// verify successful response
	assert.Equal(suite.T(), http.StatusOK, w.Code)                      // status should be 200
	assert.Contains(suite.T(), w.Body.String(), "welcome admin")       	// check response body
}

// tests the AdminOnly middleware with non-admin role
func (suite *AuthMiddlewareTestSuite) TestAdminOnly_RejectNonAdmin() {
	
	// setup router with user role in context
	suite.router.Use(func(c *gin.Context) {
		c.Set("role", "user")
	})
	// use the AdminOnly middleware
	suite.router.Use(AdminOnly())
	// define an admin route
	suite.router.GET("/admin", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "unauthorized"})
	})

	// create test request
	req := httptest.NewRequest(http.MethodGet, "/admin", nil)
	w := httptest.NewRecorder()

	// serve the request using the router
	suite.router.ServeHTTP(w, req)       

	// verify forbidden response
	assert.Equal(suite.T(), http.StatusForbidden, w.Code)                      // status should be 403
	assert.Contains(suite.T(), w.Body.String(), "admin access required")       // check response body
}

// tests the AdminOnly middleware with no role in context
func (suite *AuthMiddlewareTestSuite) TestAdminOnly_NoRoleInContext() {
	
	// setup router without setting role in context
	suite.router.Use(func(c *gin.Context) {
		// no role set
	})
	// use the AdminOnly middleware
	suite.router.Use(AdminOnly())
	// define an admin route
	suite.router.GET("/admin", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "unauthorized"})
	})

	// create test request
	req := httptest.NewRequest(http.MethodGet, "/admin", nil)
	w := httptest.NewRecorder()

	// serve the request using the router
	suite.router.ServeHTTP(w, req)          

	// verify forbidden response
	assert.Equal(suite.T(), http.StatusForbidden, w.Code)                     // status should be 404
	assert.Contains(suite.T(), w.Body.String(), "admin access required")      // check response body
}

// runs the test suite for AuthMiddleware
func TestAuthMiddlewareTestSuite(t *testing.T) {
	suite.Run(t, new(AuthMiddlewareTestSuite))     // run the test suite
}