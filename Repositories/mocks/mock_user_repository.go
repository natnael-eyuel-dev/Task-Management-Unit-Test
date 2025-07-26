package mock_repositories

// imports
import (
	"github.com/natnael-eyuel-dev/Task-Management-Unit-Test/Domain"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// mocks the UserRepository interface for testing
type MockUserRepository struct {
	mock.Mock
}

// mocks CreateUser method
func (mctr *MockUserRepository) CreateUser(user *domain.User) error {
	
	// call the mocked method and return the result
	args := mctr.Called(user)

	return args.Error(0)
}

// mocks GetByUsername method
func (mctr *MockUserRepository) GetByUsername(username string) (*domain.User, error) {
	
	// call the mocked method and return the result
	args := mctr.Called(username)
	if args.Get(0) != nil {
		return args.Get(0).(*domain.User), args.Error(1)
	}

	return nil, args.Error(1)
}

// mocks GetUserCount method
func (mctr *MockUserRepository) GetUserCount() (int64, error) {
	
	// call the mocked method and return the result
	args := mctr.Called()

	return args.Get(0).(int64), args.Error(1)
}

// mocks GetUserById method
func (mctr *MockUserRepository) GetUserById(id primitive.ObjectID) (*domain.User, error) {
	
	// call the mocked method and return the result
	args := mctr.Called(id)
	if args.Get(0) != nil {
		return args.Get(0).(*domain.User), args.Error(1)
	}

	return nil, args.Error(1)
}

// mocks UpdateRole method
func (mctr *MockUserRepository) UpdateRole(id primitive.ObjectID, role string) error {
	
	// call the mocked method and return the result
	args := mctr.Called(id, role)
	
	return args.Error(0)
}
