package infrastructure

// imports
import (
	"testing"
	"time"
	"github.com/dgrijalva/jwt-go"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

// test suite for JWTService
type JWTServiceTestSuite struct {
	suite.Suite
	service *JWTService      // JWT service instance
}

// initializes the JWTService before running tests
func (suite *JWTServiceTestSuite) SetupSuite() {

	service, err := NewJWTService()      // create a new JWT service instance
	require.NoError(suite.T(), err)                     // check for errors
	suite.service = service                             // assign to the test suite
}

// resets the viper configuration after tests
func (suite *JWTServiceTestSuite) TearDownSuite() {
	viper.Reset()
}

// tests the creation of a new JWTService instance
func (suite *JWTServiceTestSuite) TestNewJWTService() {
	
	// test cases for creating a new JWTService
	tests := []struct {
		name      string
		secret    string
		wantError bool
	}{
		{
			name:      "success with valid secret",
			secret:    "valid-secret-123",
			wantError: false,
		},
		{
			name:      "fail with empty secret",
			secret:    "",
			wantError: true,
		},
	}

	// iterate over each test case
	for _, tt := range tests {
		// run each test case
		suite.Run(tt.name, func() {
			// set the environment variable for each test case
			viper.Reset()
			if tt.secret != "" {
				_ = viper.BindEnv("JWT_SECRET")
				viper.Set("JWT_SECRET", tt.secret)
			} else {
				_ = viper.BindEnv("JWT_SECRET")
				viper.Set("JWT_SECRET", "")
			}

			// create a new JWTService instance
			service, err := NewJWTService()

			// check if the error matches the expected outcome
			if tt.wantError {
				require.Error(suite.T(), err)
				require.Nil(suite.T(), service)
			} else {
				require.NoError(suite.T(), err)
				require.NotNil(suite.T(), service)
				assert.NotEmpty(suite.T(), service.GetSecret())         // check if secret is set
			}
		})
	}
}

// tests the GenerateToken method of JWTService
func (suite *JWTServiceTestSuite) TestGenerateToken() {
	
	// test cases for generating tokens
	tests := []struct {
		name      string
		userID    string
		username  string
		role      string
		wantError bool
		errMsg    string
	}{
		{
			name:      "success with valid claims",
			userID:    "user123",
			username:  "testuser",
			role:      "user",
			wantError: false,
		},
		{
			name:      "fail with empty userID",
			userID:    "",
			username:  "testuser",
			role:      "user",
			wantError: true,
			errMsg:    "userID cannot be empty",
		},
		{
			name:      "fail with empty username",
			userID:    "user123",
			username:  "",
			role:      "user",
			wantError: true,
			errMsg:    "username cannot be empty",
		},
		{
			name:      "fail with empty role",
			userID:    "user123",
			username:  "testuser",
			role:      "",
			wantError: true,
			errMsg:    "role cannot be empty",
		},
	}

	// iterate over each test case
	for _, tt := range tests {
		// run each test case
		suite.Run(tt.name, func() {
			// call the GenerateToken method
			token, err := suite.service.GenerateToken(tt.userID, tt.username, tt.role)

			// check if the error matches the expected outcome
			if tt.wantError {
				require.Error(suite.T(), err)
				require.Empty(suite.T(), token)
				if tt.errMsg != "" {
					assert.Contains(suite.T(), err.Error(), tt.errMsg)
				}
			} else {
				require.NoError(suite.T(), err)
				require.NotEmpty(suite.T(), token)

				// verify token can be parsed
				parsed, err := suite.service.ValidateToken(token)
				require.NoError(suite.T(), err)

				// verify claims
				claims, ok := parsed.Claims.(jwt.MapClaims)
				require.True(suite.T(), ok)
				assert.Equal(suite.T(), tt.userID, claims["userId"])             // check userId
				assert.Equal(suite.T(), tt.username, claims["username"])	     // check username
				assert.Equal(suite.T(), tt.role, claims["role"])                 // check role
			}
		})
	}
}

// tests the ValidateToken method of JWTService
func (suite *JWTServiceTestSuite) TestValidateToken() {
	
	// generate a valid token 
	validToken, err := suite.service.GenerateToken("user123", "testuser", "user")
	require.NoError(suite.T(), err)

	// generate an expired token
	expiredToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId":   "user123",
		"username": "testuser",
		"role":     "user",
		"exp":      time.Now().Add(-time.Hour).Unix(),
	}).SignedString([]byte(suite.service.GetSecret()))
	require.NoError(suite.T(), err)

	// generate a token with invalid signature
	invalidSigToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId":   "user123",
		"username": "testuser",
		"role":     "user",
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	}).SignedString([]byte("wrong-secret"))
	require.NoError(suite.T(), err)

	// test cases for validating tokens
	tests := []struct {
		name      string
		token     string
		wantError bool
		errMsg    string
	}{
		{
			name:      "valid token",
			token:     validToken,
			wantError: false,
		},
		{
			name:      "empty token",
			token:     "",
			wantError: true,
			errMsg:    "token cannot be empty",
		},
		{
			name:      "expired token",
			token:     expiredToken,
			wantError: true,
			errMsg:    "Token is expired",
		},
		{
			name:      "invalid signature",
			token:     invalidSigToken,
			wantError: true,
			errMsg:    "signature is invalid",
		},
		{
			name:      "malformed token",
			token:     "invalid.token.here",
			wantError: true,
		},
	}

	// iterate over each test case
	for _, tt := range tests {
		// run each test case
		suite.Run(tt.name, func() {
			// call the ValidateToken method
			token, err := suite.service.ValidateToken(tt.token)

			// check if the error matches the expected outcome
			if tt.wantError {
				require.Error(suite.T(), err)
				require.Nil(suite.T(), token)
				if tt.errMsg != "" {
					assert.Contains(suite.T(), err.Error(), tt.errMsg)
				}
			} else {
				require.NoError(suite.T(), err)
				require.NotNil(suite.T(), token)
				assert.True(suite.T(), token.Valid)
			}
		})
	}
}

// tests the token expiration functionality of JWTService
func (suite *JWTServiceTestSuite) TestTokenExpiration() {

	// generate token with short expiration
	shortExpToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId":   "user123",
		"username": "testuser",
		"role":     "user",
		"exp":      time.Now().Add(time.Second).Unix(),
	}).SignedString([]byte(suite.service.GetSecret()))
	require.NoError(suite.T(), err)

	// immediately validate - should work
	_, err = suite.service.ValidateToken(shortExpToken)
	require.NoError(suite.T(), err)

	// after waiting for expiration - it should fail
	time.Sleep(2 * time.Second) 
	_, err = suite.service.ValidateToken(shortExpToken)               // validate the token
	require.Error(suite.T(), err)                                     // check for error
	assert.Contains(suite.T(), err.Error(), "Token is expired")       // check for expiration error
}

// runs the test suite for JWTService
func TestJWTServiceSuite(t *testing.T) {
	suite.Run(t, new(JWTServiceTestSuite))     // run the test suite
}
