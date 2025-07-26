package mock_usecases

// imports
import (
	"github.com/natnael-eyuel-dev/Task-Management-Unit-Test/Domain"
	"github.com/stretchr/testify/mock"
)

// mock the UserUseCase interface for testing
type MockUserUseCase struct {
	mock.Mock
}

// mocks Register method of UserUseCase interface
func (mcuuc *MockUserUseCase) Register(user *domain.User) error {
	
	// call the mocked method and return the error if any
	args := mcuuc.Called(user)
	
	return args.Error(0)
}

// mocks Login method of UserUseCase interface
func (mcuuc *MockUserUseCase) Login(credentials *domain.Credentials) (string, *domain.User, error) {
	
	// call the mocked method and return the results
	args := mcuuc.Called(credentials)

	var user *domain.User
	if u := args.Get(1); u != nil {
		user = u.(*domain.User)
	}

	return args.String(0), user, args.Error(2)
}

// mocks PromoteToAdmin method of UserUseCase interface
func (mcuuc *MockUserUseCase) PromoteToAdmin(userID string) error {
	
	// call the mocked method and return the error if any
	args := mcuuc.Called(userID)

	return args.Error(0)
}
