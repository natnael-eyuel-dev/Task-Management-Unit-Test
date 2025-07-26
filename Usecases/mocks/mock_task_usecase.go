package mock_usecases

// imports
import (
	"github.com/natnael-eyuel-dev/Task-Management-Unit-Test/Domain"
	"github.com/stretchr/testify/mock"
)

// mocks the TaskUseCase interface for testing
type MockTaskUseCase struct {
	mock.Mock
}

// mocks CreateTask method of TaskUseCase interface
func (mctuc *MockTaskUseCase) CreateTask(task *domain.Task) (*domain.Task, error) {
	
	// call the mocked method and return the result
	args := mctuc.Called(task)
	var result *domain.Task
	if args.Get(0) != nil {
		result = args.Get(0).(*domain.Task)
	}

	return result, args.Error(1)
}

// mocks DeleteTask method of TaskUseCase interface
func (mctuc *MockTaskUseCase) DeleteTask(taskID string) error {
	
	// call the mocked method and return the result
	args := mctuc.Called(taskID)

	return args.Error(0)
}

// mocks GetAllTasks method of TaskUseCase interface
func (mctuc *MockTaskUseCase) GetAllTasks() ([]domain.Task, error) {
	
	// call the mocked method and return the result
	args := mctuc.Called()
	var result []domain.Task
	if args.Get(0) != nil {
		result = args.Get(0).([]domain.Task)
	}

	return result, args.Error(1)
}

// mocks GetTaskByID method of TaskUseCase interface
func (mctuc *MockTaskUseCase) GetTaskByID(taskID string) (*domain.Task, error) {
	
	// call the mocked method and return the result
	args := mctuc.Called(taskID)
	var result *domain.Task
	if args.Get(0) != nil {
		result = args.Get(0).(*domain.Task)
	}

	return result, args.Error(1)
}

// mocks UpdateTask method of TaskUseCase interface
func (mctuc *MockTaskUseCase) UpdateTask(taskID string, task *domain.Task) (*domain.Task, error) {
	
	// call the mocked method and return the result
	args := mctuc.Called(taskID, task)
	var result *domain.Task
	if args.Get(0) != nil {
		result = args.Get(0).(*domain.Task)
	}

	return result, args.Error(1)
}
