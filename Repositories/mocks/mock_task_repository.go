package mock_repositories

// imports
import (
	"github.com/natnael-eyuel-dev/Task-Management-Unit-Test/Domain"
	"github.com/stretchr/testify/mock"
)

// mocks the TaskRepository interface for testing
type MockTaskRepository struct {
	mock.Mock
}

// mocks CreateTask method
func (mctr *MockTaskRepository) CreateTask(task *domain.Task) (*domain.Task, error) {
	
	// call the mocked method and return the result
	args := mctr.Called(task)
	if args.Get(0) != nil {
		return args.Get(0).(*domain.Task), args.Error(1)
	}

	return nil, args.Error(1)
}

func (mctr *MockTaskRepository) DeleteTask(id string) error {
	
	// call the mocked method and return the result
	args := mctr.Called(id)

	return args.Error(0)
}

func (mctr *MockTaskRepository) GetAllTasks() ([]domain.Task, error) {

	// call the mocked method and return the result
	args := mctr.Called()
	if args.Get(0) != nil {
		return args.Get(0).([]domain.Task), args.Error(1)
	}

	return nil, args.Error(1)
}

func (mctr *MockTaskRepository) GetTaskByID(id string) (*domain.Task, error) {
	
	// call the mocked method and return the result
	args := mctr.Called(id)
	if args.Get(0) != nil {
		return args.Get(0).(*domain.Task), args.Error(1)
	}

	return nil, args.Error(1)
}

func (mctr *MockTaskRepository) UpdateTask(id string, task *domain.Task) (*domain.Task, error) {
	
	// call the mocked method and return the result
	args := mctr.Called(id, task)
	if args.Get(0) != nil {
		return args.Get(0).(*domain.Task), args.Error(1)
	}

	return nil, args.Error(1)
}
