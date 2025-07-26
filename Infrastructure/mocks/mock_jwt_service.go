package mock_infrastructure

// imports
import (
	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/mock"
)

// mocks JWTservice for testing
type MockJWTService struct {
	mock.Mock
}

// mocks GenerateToken method of JWTService
func (mcjwts *MockJWTService) GenerateToken(userID, username, role string) (string, error) {
	
	// call the mocked method and return the results
	args := mcjwts.Called(userID, username, role)

	return args.String(0), args.Error(1)
}

// mocks ValidateToken method of JWTService
func (mcjwts *MockJWTService) ValidateToken(token string) (*jwt.Token, error) {
	
	// call the mocked method and return the results
	args := mcjwts.Called(token)
	jwtToken, _ := args.Get(0).(*jwt.Token)
	
	return jwtToken, args.Error(1)
}

// mocks GetSecret method of JWTService
func (m *MockJWTService) GetSecret() string {

	// call the mocked method and return the results
	args := m.Called()
	
	return args.String(0)
}
