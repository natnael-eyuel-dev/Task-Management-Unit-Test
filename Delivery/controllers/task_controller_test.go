package controllers

// imports
import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
	"github.com/gin-gonic/gin"
	"github.com/natnael-eyuel-dev/Task-Management-Unit-Test/Domain"
	"github.com/natnael-eyuel-dev/Task-Management-Unit-Test/Usecases/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

// test suite of TaskController
type TaskControllerTestSuite struct {
	suite.Suite
	router     *gin.Engine               			// gin router instance 
	mockUC     *mock_usecases.MockTaskUseCase    			// mock task usecase instance
	controller *TaskController // task controller instance being tested
}

// intialize the test suite before each test
func (suite *TaskControllerTestSuite) SetupTest() {
	
	gin.SetMode(gin.TestMode)                              // set gin to test mode
	suite.mockUC = new(mock_usecases.MockTaskUseCase)              // create new mock usecase
	suite.controller = NewTaskController(suite.mockUC)       // create controller with mock usecase

	// setup test router with all task routes
	router := gin.Default()      // create new gin router
	router.POST("/tasks", suite.controller.CreateTask)          // create task route
	router.GET("/tasks", suite.controller.GetAllTasks)          // get all tasks route
	router.GET("/tasks/:id", suite.controller.GetTaskByID)      // get task by ID route
	router.PUT("/tasks/:id", suite.controller.UpdateTask)       // update task route
	router.DELETE("/tasks/:id", suite.controller.DeleteTask)    // delete task route

	suite.router = router
}

// tests successful task creation
func (suite *TaskControllerTestSuite) TestCreateTask_Success() {
	
	// create mock task with fixed due date
	mockTask := &domain.Task{
		Title:       "Test Task",
		Description: "A test task",
		DueDate:     time.Now().Add(24 * time.Hour),
		Status:      "pending",
	}

	// mock CreateTask method to return the mock task
	suite.mockUC.On("CreateTask", mock.MatchedBy(func(t *domain.Task) bool {
		return t.Title == mockTask.Title &&
			t.Description == mockTask.Description &&
			t.Status == mockTask.Status &&
			t.DueDate.Equal(mockTask.DueDate)
	})).Return(mockTask, nil)

	// create test request with JSON body
	body, _ := json.Marshal(mockTask)
	req, _ := http.NewRequest(http.MethodPost, "/tasks", bytes.NewBuffer(body))      // create test request
	req.Header.Set("Content-Type", "application/json")       // set content type header
	w := httptest.NewRecorder()

	// serve the request using the router
	suite.router.ServeHTTP(w, req)

	// verify response
	suite.Equal(http.StatusCreated, w.Code)          				  // status should be 201
	suite.mockUC.AssertExpectations(suite.T())                        // verify mock was called as expected
}

// tests task creation with invalid input
func (suite *TaskControllerTestSuite) TestCreateTask_InvalidInput() {
	
	// create invalid task data - invalid title type
	body := []byte(`{"title":123}`)

	req, _ := http.NewRequest(http.MethodPost, "/tasks", bytes.NewBuffer(body))      // create test request
	req.Header.Set("Content-Type", "application/json")       // set content type header
	w := httptest.NewRecorder()

	// serve the request using the router
	suite.router.ServeHTTP(w, req)

	// verify response
	suite.Equal(http.StatusBadRequest, w.Code)    	       // status should be 400
}

// tests task creation with missing required fields
func (suite *TaskControllerTestSuite) TestCreateTask_MissingFields() {

    // missing title and description
    body := []byte(`{"
		due_date":"2025-07-30T00:00:00Z",
		"status":"pending
	"}`)

    req, _ := http.NewRequest(http.MethodPost, "/tasks", bytes.NewBuffer(body))
    req.Header.Set("Content-Type", "application/json")
    w := httptest.NewRecorder()

    suite.router.ServeHTTP(w, req)
    suite.Equal(http.StatusBadRequest, w.Code)        // status should be 400
    suite.Contains(w.Body.String(), "error")          // should contain error message
}

// tests getting all tasks when empty
func (suite *TaskControllerTestSuite) TestGetAllTasks_Empty() {
	
	// mock GetAllTasks to return empty slice
	suite.mockUC.
		On("GetAllTasks").
		Return([]domain.Task{}, nil)

	// create test request
	req, _ := http.NewRequest(http.MethodGet, "/tasks", nil)      // create test request
	w := httptest.NewRecorder()

	// serve the request using the router
	suite.router.ServeHTTP(w, req)

	// verify response
	suite.Equal(http.StatusOK, w.Code)            // status should be 200
	suite.Contains(w.Body.String(), "[]")         // reponse body should be empty array
}

// tests getting all tasks with usecase error
func (suite *TaskControllerTestSuite) TestGetAllTasks_Error() {
    
	// mock GetAllTasks to return nil and error
	suite.mockUC.
        On("GetAllTasks").
        Return(nil, errors.New("db error"))

    req, _ := http.NewRequest(http.MethodGet, "/tasks", nil)
    w := httptest.NewRecorder()

    suite.router.ServeHTTP(w, req)
    suite.Equal(http.StatusInternalServerError, w.Code)       // status should be 500
    suite.Contains(w.Body.String(), "db error")               // should contain error message
}

// tests getting a task with invalid ID format
func (suite *TaskControllerTestSuite) TestGetTaskByID_InvalidID() {

    req, _ := http.NewRequest(http.MethodGet, "/tasks/invalid-id", nil)
    w := httptest.NewRecorder()

    suite.router.ServeHTTP(w, req)
    suite.Equal(http.StatusBadRequest, w.Code)                       // status should be 400
    suite.Contains(w.Body.String(), "Invalid task ID format")        // should contain error message
}

// tests getting a task with usecase error
func (suite *TaskControllerTestSuite) TestGetTaskByID_Error() {

    id := "60d5ec49f9a3c7001c5b2b0d"
    suite.mockUC.
        On("GetTaskByID", id).
        Return(nil, errors.New("db error"))

    req, _ := http.NewRequest(http.MethodGet, "/tasks/"+id, nil)
    w := httptest.NewRecorder()

    suite.router.ServeHTTP(w, req)
    suite.Equal(http.StatusInternalServerError, w.Code)       // status should be 500
    suite.Contains(w.Body.String(), "db error")               // should contain error message
}

// tests getting a non-existent task
func (suite *TaskControllerTestSuite) TestGetTaskByID_NotFound() {
	
	// test task ID
	id := "60d5ec49f9a3c7001c5b2b0d"
	
	// mock GetTaskByID to return not found error
	suite.mockUC.
		On("GetTaskByID", id).
		Return(nil, domain.ErrTaskNotFound)

	// create test request
	req, _ := http.NewRequest(http.MethodGet, "/tasks/"+id, nil)      // create test request
	w := httptest.NewRecorder()

	// serve the request using the router
	suite.router.ServeHTTP(w, req)

	// verify response
	suite.Equal(http.StatusNotFound, w.Code)                  // status should be 404
	suite.Contains(w.Body.String(), "task not found") 		  // should contain error message
}

// tests successful task update
func (suite *TaskControllerTestSuite) TestUpdateTask_Success() {
	
	// create updated task data
	task := domain.Task{
		Status:      "completed",
		Title:       "Updated Task",
		Description: "Updated description",
		DueDate:     time.Now().Add(24 * time.Hour),
	}

	// mock task ID to update
	id := "60d5ec49f9a3c7001c5b2b0d" 

	// mock UpdateTask method to return the updated task
	suite.mockUC.On("UpdateTask", id, mock.MatchedBy(func(t *domain.Task) bool {
        return t.Title == task.Title &&
            t.Description == task.Description &&
            t.Status == task.Status &&
            t.DueDate.Round(time.Second).Equal(task.DueDate.Round(time.Second))
    })).Return(&task, nil)

	// create test request with JSON body
	body, _ := json.Marshal(task)
	req, _ := http.NewRequest(http.MethodPut, "/tasks/"+id, bytes.NewBuffer(body))      // create test request
	req.Header.Set("Content-Type", "application/json") 	     // set content type header
	w := httptest.NewRecorder()

	// serve the request using the router
	suite.router.ServeHTTP(w, req)

	// verify response
	suite.Equal(http.StatusOK, w.Code)                                // status should be 200
	suite.Contains(w.Body.String(), "task updated successfully")      // message should be in response body
}

// tests updating a task with invalid ID format
func (suite *TaskControllerTestSuite) TestUpdateTask_InvalidID() {

    body := []byte(`{"title":"Updated"}`)
    req, _ := http.NewRequest(http.MethodPut, "/tasks/invalid-id", bytes.NewBuffer(body))
    req.Header.Set("Content-Type", "application/json")
    w := httptest.NewRecorder()

    suite.router.ServeHTTP(w, req)
    suite.Equal(http.StatusBadRequest, w.Code)                     // status should be 400
    suite.Contains(w.Body.String(), "Invalid task ID format")      // should contain error message
}

// tests updating a task with invalid input
func (suite *TaskControllerTestSuite) TestUpdateTask_InvalidInput() {

    id := "60d5ec49f9a3c7001c5b2b0d"
    body := []byte(`{"title":123}`)
    req, _ := http.NewRequest(http.MethodPut, "/tasks/"+id, bytes.NewBuffer(body))
    req.Header.Set("Content-Type", "application/json")
    w := httptest.NewRecorder()

    suite.router.ServeHTTP(w, req)
    suite.Equal(http.StatusBadRequest, w.Code)       // status should be 400
}

// tests updating a non-existent task
func (suite *TaskControllerTestSuite) TestUpdateTask_NotFound() {

    id := "60d5ec49f9a3c7001c5b2b0d"
    task := &domain.Task{Title: "Updated"}

    suite.mockUC.
        On("UpdateTask", id, mock.AnythingOfType("*domain.Task")).
        Return(nil, domain.ErrTaskNotFound)

    body, _ := json.Marshal(task)
    req, _ := http.NewRequest(http.MethodPut, "/tasks/"+id, bytes.NewBuffer(body))
    req.Header.Set("Content-Type", "application/json")
    w := httptest.NewRecorder()

    suite.router.ServeHTTP(w, req)
    suite.Equal(http.StatusNotFound, w.Code)                // status should be 404
    suite.Contains(w.Body.String(), "task not found")
}

// tests updating a task with usecase error
func (suite *TaskControllerTestSuite) TestUpdateTask_Error() {

    id := "60d5ec49f9a3c7001c5b2b0d"
    task := &domain.Task{Title: "Updated"}

    suite.mockUC.
        On("UpdateTask", id, mock.AnythingOfType("*domain.Task")).
        Return(nil, errors.New("update error"))

    body, _ := json.Marshal(task)
    req, _ := http.NewRequest(http.MethodPut, "/tasks/"+id, bytes.NewBuffer(body))
    req.Header.Set("Content-Type", "application/json")
    w := httptest.NewRecorder()

    suite.router.ServeHTTP(w, req)
    suite.Equal(http.StatusBadRequest, w.Code)           // status should be 400
    suite.Contains(w.Body.String(), "update error")
}

// tests task deletion failure
func (suite *TaskControllerTestSuite) TestDeleteTask_Error() {
	
	// mock task ID to delete
	id := "60d5ec49f9a3c7001c5b2b0d" 
	
	// mock DeleteTask method to return an error
	suite.mockUC.
		On("DeleteTask", id).
		Return(errors.New("failed to delete"))

	// create test request
	req, _ := http.NewRequest(http.MethodDelete, "/tasks/"+id, nil)      // create test request
	w := httptest.NewRecorder()

	// serve the request using the router
	suite.router.ServeHTTP(w, req)

	// verify response
	suite.Equal(http.StatusInternalServerError, w.Code)        // status should be 500
	suite.Contains(w.Body.String(), "failed to delete")        // should contain error message
}

// tests deleting a task with invalid ID format
func (suite *TaskControllerTestSuite) TestDeleteTask_InvalidID() {
    
	req, _ := http.NewRequest(http.MethodDelete, "/tasks/invalid-id", nil)
    w := httptest.NewRecorder()

    suite.router.ServeHTTP(w, req)
    suite.Equal(http.StatusBadRequest, w.Code)                      // status should be 400
    suite.Contains(w.Body.String(), "Invalid task ID format")       // should contain error message
}

// tests deleting a non-existent task
func (suite *TaskControllerTestSuite) TestDeleteTask_NotFound() {
    
	id := "60d5ec49f9a3c7001c5b2b0d"
    
	suite.mockUC.
        On("DeleteTask", id).
        Return(domain.ErrTaskNotFound)

    req, _ := http.NewRequest(http.MethodDelete, "/tasks/"+id, nil)
    w := httptest.NewRecorder()

    suite.router.ServeHTTP(w, req)
    suite.Equal(http.StatusNotFound, w.Code)                // status should be 404
    suite.Contains(w.Body.String(), "task not found")       // should contain error message
}

// runs the test suite for TaskController
func TestTaskControllerTestSuite(t *testing.T) {
	suite.Run(t, new(TaskControllerTestSuite))        // run the test suite
}