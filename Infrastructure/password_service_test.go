package infrastructure

// imports
import (
	"testing"
	"github.com/natnael-eyuel-dev/Task-Management-Unit-Test/Domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"golang.org/x/crypto/bcrypt"
)

// test suite for PasswordService
type PasswordServiceTestSuite struct {
	suite.Suite
	service domain.PasswordService      // password service instance
}

// initializes the PasswordService before each test
func (suite *PasswordServiceTestSuite) SetupTest() {
	suite.service = NewPasswordService()      // create a new PasswordService instance
}

// tests the HashPassword method of PasswordService
func (suite *PasswordServiceTestSuite) TestHashPassword() {
	
	// test cases for hashing passwords
	tests := []struct {
		name      string
		password  string
		wantError bool
	}{
		{
			name:      "success with normal password",
			password:  "securePassword123!",
			wantError: false,
		},
		{
			name:      "success with empty password",
			password:  "",
			wantError: false,
		},
		{
			name:      "fail with very long password",
			password:  string(make([]byte, 100)),      
			wantError: true,
		},
	}

	// iterate over each test case
	for _, tt := range tests {
		// run each test case
		suite.Run(tt.name, func() {
			// call the HashPassword method
			hashed, err := suite.service.HashPassword(tt.password)

			// check if the error matches the expected outcome
			if tt.wantError {
				suite.Error(err)
				suite.Empty(hashed)
				suite.Contains(err.Error(), "password length exceeds 72 bytes")
			} else {
				suite.NoError(err)
				suite.NotEmpty(hashed)
				
				
				_, err = bcrypt.Cost([]byte(hashed))      // verify it's a valid bcrypt hash
				suite.NoError(err)                        // no error should occur
			}
		})
	}
}

// tests the CheckPassword method of PasswordService
func (suite *PasswordServiceTestSuite) TestCheckPassword() {
	
	// setup test passwords
	plain := "correctPassword"
	wrong := "wrongPassword"
	empty := ""
	
	// generate hashes for testing
	hashedPlain, err := suite.service.HashPassword(plain)
	require.NoError(suite.T(), err)
	hashedEmpty, err := suite.service.HashPassword(empty)
	require.NoError(suite.T(), err)

	// test cases for checking passwords
	tests := []struct {
		name     string
		hashed   string
		plain    string
		expected bool
	}{
		{
			name:     "correct password",
			hashed:   hashedPlain,
			plain:    plain,
			expected: true,
		},
		{
			name:     "wrong password",
			hashed:   hashedPlain,
			plain:    wrong,
			expected: false,
		},
		{
			name:     "empty password correct",
			hashed:   hashedEmpty,
			plain:    empty,
			expected: true,
		},
		{
			name:     "empty password wrong",
			hashed:   hashedEmpty,
			plain:    " ",
			expected: false,
		},
		{
			name:     "malformed hash",
			hashed:   "not-a-real-hash",
			plain:    plain,
			expected: false,
		},
	}

	// iterate over each test case
	for _, tt := range tests {
		// run each test case
		suite.Run(tt.name, func() {
			// call the CheckPassword method
			result := suite.service.CheckPassword(tt.hashed, tt.plain)
			assert.Equal(suite.T(), tt.expected, result)
		})
	}
}

// tests the consistency of password hashing
func (suite *PasswordServiceTestSuite) TestPasswordHashingConsistency() {
	
	// test password for consistency
	password := "consistentHashingTest"
	
	// generate two hashes of the same password
	hash1, err := suite.service.HashPassword(password)
	require.NoError(suite.T(), err)
	
	hash2, err := suite.service.HashPassword(password)
	require.NoError(suite.T(), err)
	
	// hashes should be different - due to random salt
	assert.NotEqual(suite.T(), hash1, hash2)
	
	// both should verify correctly
	assert.True(suite.T(), suite.service.CheckPassword(hash1, password))     // check first hash 
	assert.True(suite.T(), suite.service.CheckPassword(hash2, password))     // check second hash
}

// tests the password length limits
func (suite *PasswordServiceTestSuite) TestPasswordLengthLimits() {
	
	// test right at the bcrypt limit - 72 bytes
	maxLengthPassword := string(make([]byte, 72))
	hashed, err := suite.service.HashPassword(maxLengthPassword)

	suite.NoError(err)          // no error should occur
	suite.True(suite.service.CheckPassword(hashed, maxLengthPassword))       // check the hash
	
	// test 1 byte over the limit - 73 bytes
	tooLongPassword := string(make([]byte, 73))
	_, err = suite.service.HashPassword(tooLongPassword)
	suite.Error(err)            // error should occur
	suite.Contains(err.Error(), "password length exceeds 72 bytes")          // check error message
}

// runs the test suite for PasswordService
func TestPasswordServiceSuite(t *testing.T) {
	suite.Run(t, new(PasswordServiceTestSuite))     // run the test suite
}