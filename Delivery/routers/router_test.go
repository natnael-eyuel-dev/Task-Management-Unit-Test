package routers

// imports
import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/natnael-eyuel-dev/Task-Management-Unit-Test/Domain"
	"github.com/natnael-eyuel-dev/Task-Management-Unit-Test/Infrastructure/mocks"
	"github.com/natnael-eyuel-dev/Task-Management-Unit-Test/Usecases/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// test suite for the Router
type RouterTestSuite struct {
	suite.Suite                                               // embed the suite.Suite type
	router         *gin.Engine                                // gin router instance
	mockTaskUC     *mock_usecases.MockTaskUseCase             // mock task usecase
	mockUserUC     *mock_usecases.MockUserUseCase             // mock user usecase
	mockJWT        *mock_infrastructure.MockJWTService        // mock JWT service
}

// initializes the test suite
func (suite *RouterTestSuite) SetupTest() {
	gin.SetMode(gin.TestMode)                         			   // set gin to test mode
	suite.mockTaskUC = new(mock_usecases.MockTaskUseCase)          // create new mock task usecase
	suite.mockUserUC = new(mock_usecases.MockUserUseCase)          // create new mock user usecase
	suite.mockJWT = new(mock_infrastructure.MockJWTService)        // create new mock JWT service
	suite.router = SetupRouter(									   // setup router with mocks
		suite.mockTaskUC, suite.mockUserUC, suite.mockJWT,
	) 
}

// tests authenticated GetTaskByID 
func (suite *RouterTestSuite) TestGetTaskByID_Authenticated() {
	
	// generate valid task ID
	validTaskID := primitive.NewObjectID().Hex()  
	// test token    
	validToken := "valid.token.here"       

	// mock ValidateToken 
	suite.mockJWT.
		On("ValidateToken", validToken).
		Return(&jwt.Token{Valid: true}, nil)

	// mock task retrieval
	suite.mockTaskUC.
		On("GetTaskByID", validTaskID).
		Return(&domain.Task{}, nil)

	// create test request 
	req, _ := http.NewRequest("GET", "/tasks/"+validTaskID, nil)      // create test request
	req.Header.Set("Authorization", validToken)      // set auth header
	w := httptest.NewRecorder()

	// serve the request using the router
	suite.router.ServeHTTP(w, req)                  

	assert.Equal(suite.T(), http.StatusOK, w.Code)         // status should be 200
	suite.mockJWT.AssertExpectations(suite.T())            // verify mock was called
	suite.mockTaskUC.AssertExpectations(suite.T())         // verify mock was called
}

// tests unauthorized GetTaskAllTasks
func (suite *RouterTestSuite) TestGetAllTasks_Unauthorized() {

	// create test request without token
	req, _ := http.NewRequest("GET", "/tasks", nil)  	// create test request 
	w := httptest.NewRecorder()  

	// serve the request using the router
	suite.router.ServeHTTP(w, req)                  

	assert.Equal(suite.T(), http.StatusUnauthorized, w.Code) 	   // status should be 404
}

// tests admin route: POST /tasks - create task
func (suite *RouterTestSuite) TestCreateTask_AdminSuccess() {

	// test admin token
    adminToken := "admin.token.here"
	
	// mock admin claims
    claims := jwt.MapClaims{"role": "admin"}

    // mock ValidateToken to return admin claims
    suite.mockJWT.
        On("ValidateToken", adminToken).
        Return(&jwt.Token{Valid: true, Claims: claims}, nil)

    // mock CreateTask to return a new task and no error
    suite.mockTaskUC.
        On("CreateTask", mock.AnythingOfType("*domain.Task")).
        Return(&domain.Task{}, nil)

	// create test task
	task := &domain.Task{
		Title:       "New Task",
		Description: "Task description",
		DueDate:     time.Now().Add(24 * time.Hour),
		Status:      "pending",
	}

	taskJSON, err := json.Marshal(task)
	if err != nil {
		suite.T().Fatal("Failed to marshal task:", err)
	}

	// create request with JSON body
    req, err := http.NewRequest("POST", "/tasks", bytes.NewReader(taskJSON))       // create test request
	if err != nil {
		suite.T().Fatal("Failed to create request:", err)
	}
    req.Header.Set("Authorization", adminToken)                 // set auth header
    req.Header.Set("Content-Type", "application/json")          // set content type header
    w := httptest.NewRecorder()

    // serve the request using the router
    suite.router.ServeHTTP(w, req)

    assert.Equal(suite.T(), http.StatusCreated, w.Code)       // status should be 201
    suite.mockJWT.AssertExpectations(suite.T())               // verify mock was called
    suite.mockTaskUC.AssertExpectations(suite.T())            // verify mock was called
}

// tests admin route: PUT /tasks/:id - update task
func (suite *RouterTestSuite) TestUpdateTask_AdminSuccess() {

	// test admin token
    adminToken := "admin.token.here"
	// test task id
    taskID := primitive.NewObjectID().Hex()

	// mock admin claims
    claims := jwt.MapClaims{"role": "admin"}

    // mock ValidateToken to return admin claims
    suite.mockJWT.
        On("ValidateToken", adminToken).
        Return(&jwt.Token{Valid: true, Claims: claims}, nil)

    // mock UpdateTask to return updated task and no error
    suite.mockTaskUC.
        On("UpdateTask", taskID, mock.AnythingOfType("*domain.Task")).
        Return(&domain.Task{}, nil)

	// create test request with request body
    reqBody := `{
        "title":"Updated Task",
        "description":"Updated description",
        "due_date":"2025-07-26T00:00:00Z",
        "status":"completed"
    }`
    req, _ := http.NewRequest("PUT", "/tasks/"+taskID, strings.NewReader(reqBody))       // create test request
    req.Header.Set("Authorization", adminToken)                 // set auth header
    req.Header.Set("Content-Type", "application/json")          // set content type header
    w := httptest.NewRecorder()

    // serve the request using the router
    suite.router.ServeHTTP(w, req)

    assert.Equal(suite.T(), http.StatusOK, w.Code)        // status should be 200
    suite.mockJWT.AssertExpectations(suite.T())           // verify mock was called
    suite.mockTaskUC.AssertExpectations(suite.T())        // verify mock was called
}

// tests admin route: DELETE /tasks/:id - delete task
func (suite *RouterTestSuite) TestDeleteTask_AdminSuccess() {

	// test admin token
    adminToken := "admin.token.here"
	// test task id
    taskID := primitive.NewObjectID().Hex()

	// mock admin claims
    claims := jwt.MapClaims{"role": "admin"}

    // mock ValidateToken to return admin claims
    suite.mockJWT.
        On("ValidateToken", adminToken).
        Return(&jwt.Token{Valid: true, Claims: claims}, nil)

    // mock DeleteTask to return no error
    suite.mockTaskUC.
        On("DeleteTask", taskID).
        Return(nil)

	// create test request
    req, _ := http.NewRequest("DELETE", "/tasks/"+taskID, nil)      // create test request
    req.Header.Set("Authorization", adminToken)       // set auth header
    w := httptest.NewRecorder()

    // serve the request using the router
    suite.router.ServeHTTP(w, req)

    assert.Equal(suite.T(), http.StatusOK, w.Code)          // status should be 200
    suite.mockJWT.AssertExpectations(suite.T())             // verify mock was called
    suite.mockTaskUC.AssertExpectations(suite.T())          // verify mock was called
}

// tests admin routes with admin user - success
func (suite *RouterTestSuite) TestPromoteToAdmin_Success() {
	
	// generate valid user ID
	validUserID := primitive.NewObjectID().Hex()    
	// test admin token 
	adminToken := "admin.token.here"                

	// mock admin claims
	claims := jwt.MapClaims{"role": "admin"}

	// mock ValidateToken to return token and nil
	suite.mockJWT.
		On("ValidateToken", adminToken).
		Return(&jwt.Token{Valid: true, Claims: claims}, nil)

	// mock PromoteToAdmin to return nil - successful promotion
	suite.mockUserUC.
		On("PromoteToAdmin", validUserID).
		Return(nil)

	// create test request
	req, _ := http.NewRequest("PUT", "/promote/"+validUserID, nil)     // create test request
	req.Header.Set("Authorization", adminToken)            // set auth header
	req.Header.Set("Content-Type", "application/json")     // set content type
	w := httptest.NewRecorder()

	// serve the request using the router
	suite.router.ServeHTTP(w, req)                  

	assert.Equal(suite.T(), http.StatusOK, w.Code)         // status should be 200
	suite.mockJWT.AssertExpectations(suite.T())            // verify mock was called
	suite.mockUserUC.AssertExpectations(suite.T())         // verify mock was called
}

// tests admin routes with non-admin user
func (suite *RouterTestSuite) TestAdminRoutes_NonAdmin() {

	// test user token
	userToken := "user.token.here"                  

	// mock user claims
	claims := jwt.MapClaims{"role": "user"}

	// mock ValidateToken to return token and nil
	suite.mockJWT.
		On("ValidateToken", userToken).
		Return(&jwt.Token{Valid: true, Claims: claims}, nil)

	// create test request 
	req, _ := http.NewRequest("PUT", "/promote/123", nil)      // create test request
	req.Header.Set("Authorization", userToken)     // set auth header
	w := httptest.NewRecorder()

	// serve the request using the router
	suite.router.ServeHTTP(w, req)                  

	assert.Equal(suite.T(), http.StatusForbidden, w.Code)      // status should be 403
}

// tests successful register - public route
func (suite *RouterTestSuite) TestRegister_Success() {
	
	// mock Register to return no error
	suite.mockUserUC.
		On("Register", mock.AnythingOfType("*domain.User")).
		Return(nil)

	// create test request with request body
	reqBody := `{
		"username":"test",
		"password":"pass123"
	}`
	req, _ := http.NewRequest("POST", "/register", strings.NewReader(reqBody))     // create test request
	req.Header.Set("Content-Type", "application/json")      // set content type header
	w := httptest.NewRecorder()

	// serve the request using the router
	suite.router.ServeHTTP(w, req)                  

	assert.Equal(suite.T(), http.StatusCreated, w.Code)          // status should be 201
	suite.mockUserUC.AssertExpectations(suite.T())               // verify mock was called
}

// tests successful login - public route
func (suite *RouterTestSuite) TestLogin_Success() {

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

    // mock Login to return token and no error
    suite.mockUserUC.
        On("Login", &creds).
        Return("mock.jwt.token", user, nil)

	// create test request with JSON body
	body, _ := json.Marshal(creds)
    req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(body))        // create test request
    req.Header.Set("Content-Type", "application/json")        // set content type header
    w := httptest.NewRecorder()

    // serve the request using the router
    suite.router.ServeHTTP(w, req)

    assert.Equal(suite.T(), http.StatusOK, w.Code)        // status should be 200
    suite.mockUserUC.AssertExpectations(suite.T())        // verify mock was called
}

// suite entry point for running the tests
func TestRouterTestSuite(t *testing.T) {
	suite.Run(t, new(RouterTestSuite))         // run the test suite
}