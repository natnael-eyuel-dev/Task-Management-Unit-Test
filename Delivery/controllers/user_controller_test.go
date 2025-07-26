package controllers

// imports 
import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"github.com/gin-gonic/gin"
	"github.com/natnael-eyuel-dev/Task-Management-Unit-Test/Domain"
	"github.com/natnael-eyuel-dev/Task-Management-Unit-Test/Usecases/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// test suite for UserController
type UserControllerTestSuite struct {
	suite.Suite
	router       *gin.Engine               			// gin router instance 
	mockUseCase  *mock_usecases.MockUserUseCase    			// mock user usecase instance
	controller   *UserController 		// user controller instance being tested
}

// intilize the test suite before each test
func (suite *UserControllerTestSuite) SetupTest() {
	gin.SetMode(gin.TestMode)                              // set gin to test mode
	suite.router = gin.Default()                           // create new gin router
	suite.mockUseCase = new(mock_usecases.MockUserUseCase)         // create new mock usecase
	suite.controller = NewUserController(suite.mockUseCase)     // create controller with mock usecase

	// setup test router with all user routes
	suite.router.POST("/register", suite.controller.Register)             // user registration route
	suite.router.POST("/login", suite.controller.Login)                   // user login route
	suite.router.PUT("/promote/:id", suite.controller.PromoteToAdmin)     // promote user to admin route
}

// tests successful user registration
func (suite *UserControllerTestSuite) TestRegister_Success() {
	
	// create test user
	user := domain.User{
		Username: "john", 
		Password: "password123",
		Role: "user",
	}

	// mock Register method to return no error
	suite.mockUseCase.
		On("Register", &user).
		Return(nil)

	// create test request with JSON body
	body, _ := json.Marshal(user)
	req, _ := http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(body))      // create test request
	req.Header.Set("Content-Type", "application/json")      // set content type header
	resp := httptest.NewRecorder()

	// serve the request using the router
	suite.router.ServeHTTP(resp, req)

	// verify response
	assert.Equal(suite.T(), http.StatusCreated, resp.Code)             // status should be 201
	suite.mockUseCase.AssertCalled(suite.T(), "Register", &user)       // verify mock was called
}

// tests registration with existing username
func (suite *UserControllerTestSuite) TestRegister_Conflict() {
	
	// create test user
	user := domain.User{
		Username: "john", 
		Password: "password123",
		Role:     "user",
	}

	// mock Register method to return error
	suite.mockUseCase.
		On("Register", &user).
		Return(domain.ErrUserExists)

	// create test request with JSON body
	body, _ := json.Marshal(user)
	req, _ := http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(body))        // create test request
	req.Header.Set("Content-Type", "application/json")        // set content type header
	resp := httptest.NewRecorder()

	// serve the request using the router
	suite.router.ServeHTTP(resp, req)

	// verify response
	assert.Equal(suite.T(), http.StatusConflict, resp.Code) 	  // status should be 409
}

// tests registration with missing username field
func (suite *UserControllerTestSuite) TestRegister_MissingUsername() {
    
	// create test user with missing username
    user := domain.User{
        Password: "password123",
        Role:     "user",
    }

    // create test request with JSON body
    body, _ := json.Marshal(user)
    req, _ := http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(body))       // create test request
    req.Header.Set("Content-Type", "application/json")       // set content type header
    resp := httptest.NewRecorder()

    // serve the request using the router
    suite.router.ServeHTTP(resp, req)

    // verify response
    assert.Equal(suite.T(), http.StatusBadRequest, resp.Code)      // status should be 400
    assert.Contains(suite.T(), resp.Body.String(), "error")
}

// tests registration with missing password field
func (suite *UserControllerTestSuite) TestRegister_MissingPassword() {
    
	// create test user with missing password
    user := domain.User{
        Username: "john",
        Role:     "user",
    }

    body, _ := json.Marshal(user)
    req, _ := http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(body))       // create test request 
    req.Header.Set("Content-Type", "application/json")         // set content type header
    resp := httptest.NewRecorder()

	// serve the request using the router
    suite.router.ServeHTTP(resp, req)
    assert.Equal(suite.T(), http.StatusBadRequest, resp.Code)       // status should be 400
}

// tests successful user login
func (suite *UserControllerTestSuite) TestLogin_Success() {
	
	// create test credentials
	creds := domain.Credentials{
		Username: "john", 
		Password: "password123",
	}

	// create mock user response
	user := &domain.User{
		ID: primitive.NewObjectID(), 
		Username: "john", 
		Role: "user",
	}

	// create mock token
	token := "mocked.jwt.token"     // mock token

	// mock Login method to return token, user and no error
	suite.mockUseCase.
		On("Login", &creds).
		Return(token, user, nil)

	// create test request with JSON body
	body, _ := json.Marshal(creds)
	req, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(body))       // create test request
	req.Header.Set("Content-Type", "application/json")        // set content type header 
	resp := httptest.NewRecorder()

	// serve the request using the router
	suite.router.ServeHTTP(resp, req)

	// verify response
	assert.Equal(suite.T(), http.StatusOK, resp.Code)       // status should be 200
}

// tests login with invalid credentials
func (suite *UserControllerTestSuite) TestLogin_InvalidCredentials() {
	
	// create test credentials with wrong password
	creds := domain.Credentials{
		Username: "john", 
		Password: "wrongpass",
	}
	
	// mock Login method to return empty, nil and  error 
	suite.mockUseCase.
		On("Login", &creds).
		Return("", nil, domain.ErrInvalidCredentials)

	// create test request with JSON body
	body, _ := json.Marshal(creds)
	req, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(body))       // create test request
	req.Header.Set("Content-Type", "application/json")        // set content type header 
	resp := httptest.NewRecorder()

	// serve the request using the router
	suite.router.ServeHTTP(resp, req)

	// verify response
	assert.Equal(suite.T(), http.StatusUnauthorized, resp.Code)       // status should be 401
}

// tests login with empty credentials
func (suite *UserControllerTestSuite) TestLogin_EmptyCredentials() {
    
	// create test empty credentials
	emptyCreds := domain.Credentials{
        Username: "",
        Password: "",
    }

    body, _ := json.Marshal(emptyCreds)
    req, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(body))        // create test request
    req.Header.Set("Content-Type", "application/json")         // set content type header
    resp := httptest.NewRecorder()

	// serve the request using the router
    suite.router.ServeHTTP(resp, req)
    assert.Equal(suite.T(), http.StatusBadRequest, resp.Code)         // status should be 400
}

// tests successful user promotion to admin
func (suite *UserControllerTestSuite) TestPromoteToAdmin_Success() {

	// mock user ID
	id := primitive.NewObjectID().Hex()

	// mock PromoteToAdmin to return no error
	suite.mockUseCase.
		On("PromoteToAdmin", id).
		Return(nil)

	// create test request
	req, _ := http.NewRequest(http.MethodPut, "/promote/"+id, nil)       // create test request
	resp := httptest.NewRecorder()

	// serve the request using the router
	suite.router.ServeHTTP(resp, req)

	// verify response
	assert.Equal(suite.T(), http.StatusOK, resp.Code)       // status should be 200
}

// tests promotion with invalid user ID format
func (suite *UserControllerTestSuite) TestPromoteToAdmin_InvalidID() {

	// mock PromoteToAdmin method to return error 
	suite.mockUseCase.
		On("PromoteToAdmin", "invalid-id").
		Return(domain.ErrInvalidUserID)

	// create test request with invalid ID
	req, _ := http.NewRequest(http.MethodPut, "/promote/invalid-id", nil)      // create test request
	resp := httptest.NewRecorder()

	// serve the request using the router 
	suite.router.ServeHTTP(resp, req)
	// verify response
	assert.Equal(suite.T(), http.StatusBadRequest, resp.Code) 	     // status should be 400
}

// tests promotion when user is not found
func (suite *UserControllerTestSuite) TestPromoteToAdmin_UserNotFound() {
    
	// mock valid user id
	validID := primitive.NewObjectID().Hex()

    // mock PromoteToAdmin to return user not found
    suite.mockUseCase.
        On("PromoteToAdmin", validID).
        Return(domain.ErrUserNotFound)

	// create test request with valid ID
    req, _ := http.NewRequest(http.MethodPut, "/promote/"+validID, nil)
    resp := httptest.NewRecorder()

	// serve the request using the router
    suite.router.ServeHTTP(resp, req)
	// verify response
    assert.Equal(suite.T(), http.StatusNotFound, resp.Code)         // status should be 404
}

// tests promotion with empty ID parameter
func (suite *UserControllerTestSuite) TestPromoteToAdmin_EmptyID() {

	// create test request with empty id
    req, _ := http.NewRequest(http.MethodPut, "/promote/", nil)
    resp := httptest.NewRecorder()

	// serve the request using the router
    suite.router.ServeHTTP(resp, req)
	// verify response
    assert.Equal(suite.T(), http.StatusNotFound, resp.Code)        // status should be 404
}

// runs the test suite for UserController
func TestUserController(t *testing.T) {
	suite.Run(t, new(UserControllerTestSuite))       // run the test suite
}