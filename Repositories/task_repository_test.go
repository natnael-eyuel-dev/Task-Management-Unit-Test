package repositories

// imports
import (
	"context"
	"errors"
	"testing"
	"time"

	domain "github.com/natnael-eyuel-dev/Task-Management-Unit-Test/Domain"
	mock_repositories "github.com/natnael-eyuel-dev/Task-Management-Unit-Test/Repositories/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// test suite for the TaskRepository
type TaskRepositoryTestSuite struct {
	suite.Suite                                      // embed the suite.Suite type
	mockCollection *mock_repositories.MockCollection // mock collection for testing
	repo           domain.TaskRepository             // task repository to be tested
}

// initializes the test suite
func (suite *TaskRepositoryTestSuite) SetupTest() {
	suite.mockCollection = new(mock_repositories.MockCollection)       // create a new mock collection
	suite.repo = NewTaskRepositoryWithCollection(suite.mockCollection) // create a new task repository with mock collection
}

// tests CreateTask method of the TaskRepository
func (suite *TaskRepositoryTestSuite) TestCreateTask_Success() {

	// create a new task
	task := &domain.Task{
		Title:       "Test Task",
		Description: "A task to test",
		DueDate:     time.Now().Add(24 * time.Hour),
		Status:      "Pending",
	}

	// mock the InsertOne method of the collection
	suite.mockCollection.
		On("InsertOne", mock.Anything, mock.MatchedBy(func(t interface{}) bool {
			_, ok := t.(*domain.Task)
			return ok
		})).Return(&mongo.InsertOneResult{}, nil)

	result, err := suite.repo.CreateTask(task) // call CreateTask method
	assert.NoError(suite.T(), err)             // assert no error
	assert.NotNil(suite.T(), result)           // assert result is not nil
	assert.NotEmpty(suite.T(), result.ID)      // assert ID is not empty
}

// tests CreateTask method of the TaskRepository for error case
func (suite *TaskRepositoryTestSuite) TestCreateTask_Error() {

	// create a new task
	task := &domain.Task{
		Title:       "Test Task",
		Description: "A task to test",
	}

	// mock the InsertOne method of the collection
	suite.mockCollection.
		On("InsertOne", mock.Anything, mock.Anything).
		Return(nil, errors.New("insert error"))

	result, err := suite.repo.CreateTask(task)        // call CreateTask method
	assert.Nil(suite.T(), result)                     // assert result is nil
	assert.EqualError(suite.T(), err, "insert error") // assert error message
}

// tests CreateTask method of the TaskRepository for context timeout
func (suite *TaskRepositoryTestSuite) TestCreateTask_ContextTimeout() {

	// create a context with a very short timeout
    _, cancel := context.WithTimeout(context.Background(), 1*time.Nanosecond)
    defer cancel()

	// create a user
    task := &domain.Task{Title: "timeout"}

	// mock the InsertOne method of the collection
    suite.mockCollection.
        On("InsertOne", mock.Anything, task).
        Return(nil, context.DeadlineExceeded)

    result, err := suite.repo.CreateTask(task)       // call CreateTask method
	assert.ErrorIs(suite.T(), err, context.DeadlineExceeded)      // assert error is context deadline exceeded
    assert.Nil(suite.T(), result)                    // assert result is nil
}

// tests GetTaskByID method of the TaskRepository for non-existing task
func (suite *TaskRepositoryTestSuite) TestGetTaskByID_NotFound() {

	// create a new object ID
	objID := primitive.NewObjectID()
	// create a mock result
	mockResult := &mock_repositories.MockSingleResult{
		Err: mongo.ErrNoDocuments,
	}

	// mock the FindOne method of the collection
	suite.mockCollection.
		On("FindOne", mock.Anything, bson.M{"_id": objID}).
		Return(mockResult)

	task, err := suite.repo.GetTaskByID(objID.Hex())       // call GetTaskByID method
	assert.Nil(suite.T(), task)                            // assert task is nil
	assert.ErrorIs(suite.T(), err, domain.ErrTaskNotFound) // assert error is ErrTaskNotFound
}

// tests GetTaskByID method of the TaskRepository for invalid ID
func (suite *TaskRepositoryTestSuite) TestGetTaskByID_InvalidID() {

	task, err := suite.repo.GetTaskByID("invalid-id")       // call GetTaskByID with invalid ID
	assert.Nil(suite.T(), task)                             // assert task is nil
	assert.ErrorIs(suite.T(), err, domain.ErrInvalidTaskID) // assert error is ErrInvalidTaskID
}

// tests GetTaskByID method of the TaskRepository for error case
func (suite *TaskRepositoryTestSuite) TestGetTaskByID_Error() {

	// create a new object ID
	objID := primitive.NewObjectID()
	// create a mock result
	mockResult := &mock_repositories.MockSingleResult{
		Err: errors.New("find error"),
	}

	// mock the FindOne method of the collection to return error
	suite.mockCollection.
		On("FindOne", mock.Anything, bson.M{"_id": objID}).
		Return(mockResult)

	task, err := suite.repo.GetTaskByID(objID.Hex()) // call GetTaskByID method
	assert.Nil(suite.T(), task)                      // assert task is nil
	assert.EqualError(suite.T(), err, "find error")  // assert error message
}

// tests DeleteTask method of the TaskRepository with invalid ID
func (suite *TaskRepositoryTestSuite) TestDeleteTask_InvalidID() {

	// mock the DeleteOne method of the collection
	suite.mockCollection.
		On("DeleteOne", mock.Anything, bson.M{"_id": mock.AnythingOfType("primitive.ObjectID")}).
		Return(&mongo.DeleteResult{}, nil)

	err := suite.repo.DeleteTask("invalid-id")              // call DeleteTask with an invalid ID
	assert.Error(suite.T(), err)                            // assert error is returned
	assert.ErrorIs(suite.T(), err, domain.ErrInvalidTaskID) // assert error is ErrInvalidTaskID
}

// tests DeleteTask method of the TaskRepository for non-existing task
func (suite *TaskRepositoryTestSuite) TestDeleteTask_NotFound() {

	// create a new object ID
	objID := primitive.NewObjectID()

	// mock the DeleteOne method of the collection
	suite.mockCollection.
		On("DeleteOne", mock.Anything, bson.M{"_id": objID}).
		Return(&mongo.DeleteResult{DeletedCount: 0}, nil)

	err := suite.repo.DeleteTask(objID.Hex())              // call DeleteTask method
	assert.Error(suite.T(), err)                           // assert error is returned
	assert.ErrorIs(suite.T(), err, domain.ErrTaskNotFound) // assert error is ErrTaskNotFound
}

// tests DeleteTask method of the TaskRepository for success case
func (suite *TaskRepositoryTestSuite) TestDeleteTask_Success() {

	// create a new object ID
	objID := primitive.NewObjectID()

	// mock the DeleteOne method of collection
	suite.mockCollection.
		On("DeleteOne", mock.Anything, bson.M{"_id": objID}).
		Return(&mongo.DeleteResult{DeletedCount: 1}, nil)

	err := suite.repo.DeleteTask(objID.Hex()) // call DeleteTask method
	assert.NoError(suite.T(), err)            // assert no error
}

// tests UpdateTask method of the TaskRepository with no fields provided
func (suite *TaskRepositoryTestSuite) TestUpdateTask_NoFieldsProvided() {

	// create a new object ID
	objID := primitive.NewObjectID()
	// create a mock result
	task := &domain.Task{}

	// mock the UpdateOne method of the collection
	suite.mockCollection.
		On("UpdateOne", mock.Anything, bson.M{"_id": objID}, mock.Anything).
		Return(&mongo.UpdateResult{}, nil)

	updated, err := suite.repo.UpdateTask(objID.Hex(), task)                    // call UpdateTask method with no fields provided
	assert.Nil(suite.T(), updated)                                              // assert updated task is nil
	assert.Error(suite.T(), err)                                                // assert error is returned
	assert.Equal(suite.T(), "no valid fields provided for update", err.Error()) // assert error message
}

// tests UpdateTask method of the TaskRepository for invalid ID
func (suite *TaskRepositoryTestSuite) TestUpdateTask_InvalidID() {

	// create test task
	task := &domain.Task{Title: "Update"}
	updated, err := suite.repo.UpdateTask("invalid-id", task) // call UpdateTask with invalid ID
	assert.Nil(suite.T(), updated)                            // assert updated task is nil
	assert.ErrorIs(suite.T(), err, domain.ErrInvalidTaskID)   // assert error is ErrInvalidTaskID
}

// tests UpdateTask method of the TaskRepository for not found
func (suite *TaskRepositoryTestSuite) TestUpdateTask_NotFound() {

	// create a new object ID
	objID := primitive.NewObjectID()
	// create test task
	task := &domain.Task{Title: "Update"}
	// create a mock result
	mockResult := &mock_repositories.MockSingleResult{
		Err: mongo.ErrNoDocuments,
	}

	// mock the FindOneAndUpdate method of the collection
	suite.mockCollection.
		On("FindOneAndUpdate", mock.Anything, bson.M{"_id": objID}, mock.Anything).
		Return(mockResult)

	updated, err := suite.repo.UpdateTask(objID.Hex(), task)
	assert.Nil(suite.T(), updated)                         // assert updated task is nil
	assert.ErrorIs(suite.T(), err, domain.ErrTaskNotFound) // assert error is ErrTaskNotFound
}

// tests UpdateTask method of the TaskRepository for error case
func (suite *TaskRepositoryTestSuite) TestUpdateTask_Error() {

	// create a new object ID
	objID := primitive.NewObjectID()
	// create test task
	task := &domain.Task{Title: "Update"}
	// create a mock result
	mockResult := &mock_repositories.MockSingleResult{
		Err: errors.New("update error"),
	}

	// mock the FindOneAndUpdate method of the collection
	suite.mockCollection.
		On("FindOneAndUpdate", mock.Anything, bson.M{"_id": objID}, mock.Anything).
		Return(mockResult)

	updated, err := suite.repo.UpdateTask(objID.Hex(), task) // call UpdateTask method
	assert.Nil(suite.T(), updated)                           // assert updated task is nil
	assert.EqualError(suite.T(), err, "update error")        // assert error message
}

// suite entry point for running the tests
func TestTaskRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(TaskRepositoryTestSuite)) // run the test suite
}
