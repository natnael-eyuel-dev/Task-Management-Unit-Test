package mock_infrastructure

// imports
import (
	"github.com/stretchr/testify/mock"
)

// mocks PasswordService for testing
type MockPasswordService struct {
	mock.Mock
}

// mocks HashPassword method of PasswordService
func (m *MockPasswordService) HashPassword(password string) (string, error) {
	
	// call the mocked method and return the results
	args := m.Called(password)
	
	return args.String(0), args.Error(1)
}

// mocks CheckPassword method of PasswordService
func (m *MockPasswordService) CheckPassword(hashedPassword, plainPassword string) bool {
	
	// call the mocked method and return the results
	args := m.Called(hashedPassword, plainPassword)
	
	return args.Bool(0)
}
