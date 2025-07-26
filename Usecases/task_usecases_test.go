package usecases

// imports
import (
	"testing"
	"time"
	"github.com/natnael-eyuel-dev/Task-Management-Unit-Test/Domain"
	"github.com/natnael-eyuel-dev/Task-Management-Unit-Test/Repositories/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// test suite for TaskUseCase
type TaskUseCaseTestSuite struct {
	suite.Suite
	mockRepo     *mock_repositories.MockTaskRepository      // mock task repository instance
	taskUsecase  domain.TaskUseCase                         // task usecase instance being tested
}

// intialize the test suite before each test
func (suite *TaskUseCaseTestSuite) SetupTest() {
	suite.mockRepo = new(mock_repositories.MockTaskRepository)      // create new mock repository
	suite.taskUsecase = NewTaskUseCase(suite.mockRepo)     // create new usecase with mock repo
}

// tests successful creation of a task
func (suite *TaskUseCaseTestSuite) TestCreateTask_Success() {
	
	// create test task
	task := &domain.Task{
		Title:       "Test",
		Description: "Test description",
		DueDate:     time.Now().Add(48 * time.Hour),
		Status:      "pending",
	}
	expected := &domain.Task{ID: task.ID}

	// mock CreateTask of the repository to return expected task
	suite.mockRepo.
		On("CreateTask", task).        
		Return(expected, nil)          

	// call the CreateTask method on usecase
	result, err := suite.taskUsecase.CreateTask(task)

	// verify the results
	assert.NoError(suite.T(), err)                                  // no error expected
	assert.Equal(suite.T(), expected, result)                       // result should match expected task
	suite.mockRepo.AssertCalled(suite.T(), "CreateTask", task)      // verify CreateTask was called with correct task
}

// tests task creation with invalid due date - past date
func (suite *TaskUseCaseTestSuite) TestCreateTask_InvalidDueDate() {
	
	// create test task with past due date
	task := &domain.Task{
		Title:       "Bad Task",
		Description: "Past due",
		DueDate:     time.Now().Add(-1 * time.Hour),
		Status:      "pending",
	}

	// mock CreateTask of the repository to return an error
	suite.mockRepo.
		On("CreateTask", task).
		Return(nil, domain.ErrInvalidDueDate)

	// call the CreateTask method on usecase
	result, err := suite.taskUsecase.CreateTask(task)

	// verify error response
	assert.Nil(suite.T(), result)                                             // result should be nil
	assert.EqualError(suite.T(), err, "due date must be in the future")       // error message should match expected
}

// tests task creation with empty title
func (suite *TaskUseCaseTestSuite) TestCreateTask_EmptyTitle() {

	// create tets task with empty title 
    task := &domain.Task{
        Title:       "",
        Description: "desc",
        DueDate:     time.Now().Add(24 * time.Hour),
        Status:      "pending",
    }
	
	// call the CreateTask method on usecase
    result, err := suite.taskUsecase.CreateTask(task)
    assert.Nil(suite.T(), result)                                             // result should be nil
    assert.EqualError(suite.T(), err, "task title cannot be empty")           // error message should match expected 
}

// tests task creation with empty description
func (suite *TaskUseCaseTestSuite) TestCreateTask_EmptyDescription() {

	// create tets task with empty description
    task := &domain.Task{
        Title:       "title",
        Description: "",
        DueDate:     time.Now().Add(24 * time.Hour),
        Status:      "pending",
    }

	// call the CreateTask method on usecase
    result, err := suite.taskUsecase.CreateTask(task)
    assert.Nil(suite.T(), result)                                                // result should be nil
    assert.EqualError(suite.T(), err, "task description cannot be empty")        // error message should match expected 
}

// tests task creation with empty due date
func (suite *TaskUseCaseTestSuite) TestCreateTask_EmptyDueDate() {

	// create tets task with empty due date
    task := &domain.Task{
        Title:       "title",
        Description: "desc",
        DueDate:     time.Time{},
        Status:      "pending",
    }

	// call the CreateTask method on usecase
    result, err := suite.taskUsecase.CreateTask(task)
    assert.Nil(suite.T(), result)                                         // result should be nil
    assert.EqualError(suite.T(), err, "due date cannot be empty")         // error message should match expected 
}  

// tests task creation with empty status (should default to pending)
func (suite *TaskUseCaseTestSuite) TestCreateTask_EmptyStatusDefaultsPending() {

	// create tets task with empty status
    task := &domain.Task{
        Title:       "title",
        Description: "desc",
        DueDate:     time.Now().Add(24 * time.Hour),
        Status:      "",
    }

    expected := &domain.Task{ID: task.ID}

	// mock CreateTask of the repository to return an expected task and nil
    suite.mockRepo.
        On("CreateTask", task).
        Return(expected, nil)

	// call the CreateTask method on usecase
    result, err := suite.taskUsecase.CreateTask(task)
    assert.NoError(suite.T(), err)                           // should be no error
    assert.Equal(suite.T(), expected, result)                // result should match expected
    assert.Equal(suite.T(), "pending", task.Status)          // task status should match pending 
}

// tests deletion of a non-existent task
func (suite *TaskUseCaseTestSuite) TestDeleteTask_NotFound() {
	
	// create a task ID that does not exist
	id := "nonexistent-id"  

	// mock GetTaskByID of the repository to return an error
	suite.mockRepo.
		On("GetTaskByID", id).
		Return(nil, domain.ErrTaskNotFound)

	// call the DeleteTask method on usecase
	err := suite.taskUsecase.DeleteTask(id)

	// verify error response
	assert.ErrorIs(suite.T(), err, domain.ErrTaskNotFound)      // should return task not found error
}

// tests task update with invalid status
func (suite *TaskUseCaseTestSuite) TestUpdateTask_InvalidStatus() {
	
	// valid id and invalid status 
	id := "some-task-id"       
	task := &domain.Task{Status: "invalid_status"}      // invalid status

	// call the UpdateTask method on usecase
	result, err := suite.taskUsecase.UpdateTask(id, task)

	// verify error response
	assert.Nil(suite.T(), result)                                  // result should be nil
	assert.EqualError(suite.T(), err, "invalid task status")       // error message should match expected
}

// tests DeleteTask with empty id
func (suite *TaskUseCaseTestSuite) TestDeleteTask_EmptyID() {

	// call the DeleteTask method on usecase
    err := suite.taskUsecase.DeleteTask("")
    assert.EqualError(suite.T(), err, "task ID cannot be empty")        // error message should match expected
}

// tests GetTaskByID with empty id
func (suite *TaskUseCaseTestSuite) TestGetTaskByID_EmptyID() {

	// call the GetTaskByID method on usecase
    result, err := suite.taskUsecase.GetTaskByID("")
    assert.Nil(suite.T(), result)                                        // result should be nil
    assert.EqualError(suite.T(), err, "task ID cannot be empty")         // error message should match expected
}

// tests GetTaskByID for not found
func (suite *TaskUseCaseTestSuite) TestGetTaskByID_NotFound() {
    
	// non-existent id
	id := "notfound-id"

	// mock GetTaskByID of the repository to return an nil and nil
    suite.mockRepo.
        On("GetTaskByID", id).
        Return(nil, nil)

	// call the GetTaskByID method on usecase
    result, err := suite.taskUsecase.GetTaskByID(id)
    assert.Nil(suite.T(), result)                                    // result should be nil
    assert.ErrorIs(suite.T(), err, domain.ErrTaskNotFound)           // error message should match expected
}

// tests GetAllTasks returns empty slice if repo returns nil
func (suite *TaskUseCaseTestSuite) TestGetAllTasks_RepoReturnsNil() {
    
	// mock GetTaskByID of the repository to return an nil and nil
	suite.mockRepo.
        On("GetAllTasks").
        Return(nil, nil)

	// call the GetTaskByID method on usecase
    result, err := suite.taskUsecase.GetAllTasks()
    assert.NoError(suite.T(), err)                 // no error should exist
    assert.NotNil(suite.T(), result)               // result should not be nil
    assert.Len(suite.T(), result, 0)               // length of result should be 0
}

// tests UpdateTask with empty id
func (suite *TaskUseCaseTestSuite) TestUpdateTask_EmptyID() {

	// test task
    task := &domain.Task{Title: "title"}

	// call the UpdateTask method on usecase
    result, err := suite.taskUsecase.UpdateTask("", task)
    assert.Nil(suite.T(), result)                                        // result should be nil
    assert.EqualError(suite.T(), err, "task ID cannot be empty")         // error message should match expected
}

// tests UpdateTask with no valid fields provided
func (suite *TaskUseCaseTestSuite) TestUpdateTask_NoValidFields() {

	// test task id
    id := "some-id"
	// test task
    task := &domain.Task{}

	// call the UpdateTask method on usecase
    result, err := suite.taskUsecase.UpdateTask(id, task)
    assert.Nil(suite.T(), result)                                                    // result should be nil
    assert.EqualError(suite.T(), err, "no valid fields provided for update")         // error message should match expected
}

// tests UpdateTask with invalid due date
func (suite *TaskUseCaseTestSuite) TestUpdateTask_InvalidDueDate() {
    
	// test task id
	id := "some-id"
	// test task
    task := &domain.Task{DueDate: time.Now().Add(-1 * time.Hour)}

	// call the UpdateTask method on usecase
    result, err := suite.taskUsecase.UpdateTask(id, task)
    assert.Nil(suite.T(), result)                                              // result should be nil
    assert.EqualError(suite.T(), err, "due date must be in the future")        // error message should match expected
}

// runs the test suite for TaskUseCase
func TestTaskUseCaseTestSuite(t *testing.T) {
	suite.Run(t, new(TaskUseCaseTestSuite))        // run the test suite
}