package usecases

// imports
import (
	"errors"
	"testing"

	"github.com/natnael-eyuel-dev/Task-Management-Unit-Test/Domain"
	"github.com/natnael-eyuel-dev/Task-Management-Unit-Test/Infrastructure/mocks"
	"github.com/natnael-eyuel-dev/Task-Management-Unit-Test/Repositories/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// test suite for UserUseCase
type UserUseCaseTestSuite struct {
	suite.Suite
	userRepo     *mock_repositories.MockUserRepository         // mock user repository instance
	jwtService   *mock_infrastructure.MockJWTService           // mock JWT service instance
	pwdService   *mock_infrastructure.MockPasswordService      // mock password service instance
	usecase      domain.UserUseCase                          // user usecase instance being tested
}

// initializes the test environment before each test
func (suite *UserUseCaseTestSuite) SetupTest() {
	suite.userRepo = new(mock_repositories.MockUserRepository)            // create new mock user repository
	suite.jwtService = new(mock_infrastructure.MockJWTService)            // create new mock JWT service
	suite.pwdService = new(mock_infrastructure.MockPasswordService)       // create new mock password service
	suite.usecase = NewUserUseCase(                              // create new usecase with mocks
		suite.userRepo, suite.jwtService, suite.pwdService,
	)       
}

// tests successful user registration where first user becomes admin
func (suite *UserUseCaseTestSuite) TestRegister_SuccessFirstUserBecomesAdmin() {
	
	// create test user
	user := &domain.User{
		Username: "testuser", 
		Password: "password123",
	}

	// mock GetByUsername of the repository to return error 
	suite.userRepo.
		On("GetByUsername", user.Username).
		Return(nil, domain.ErrUserNotFound)
	// mock HashPassword of the password service to return hashed password
	suite.pwdService.
		On("HashPassword", user.Password).
		Return("hashedpass", nil)
	// mock GetUserCount of the repository to return 0 - first user
	suite.userRepo.
		On("GetUserCount").
		Return(int64(0), nil)
	// mock CreateUser of the repository to return nil - successful creation
	suite.userRepo.
		On("CreateUser", mock.AnythingOfType("*domain.User")).
		Return(nil)

	// call the Register method on usecase
	err := suite.usecase.Register(user)

	// verify results
	assert.NoError(suite.T(), err)                             // no error expected
	assert.Equal(suite.T(), "admin", user.Role)                // first user should be admin
	suite.userRepo.AssertExpectations(suite.T())               // verify all mock expectations were met
	suite.pwdService.AssertExpectations(suite.T())             // verify password service was called
}

// tests registration with existing username
func (suite *UserUseCaseTestSuite) TestRegister_AlreadyExists() {

	// mock GetByUsername of the repository to return error
	suite.userRepo.
		On("GetByUsername", "testuser").
		Return(nil, domain.ErrUserExists)

	// call the Register method on usecase
	err := suite.usecase.Register(&domain.User{
		Username: "testuser", 
		Password: "somepass",
	})

	// verify error response
	assert.ErrorIs(suite.T(), err, domain.ErrUserExists)       // error should be user exists
}

// tests registration with invalid password
func (suite *UserUseCaseTestSuite) TestRegister_InvalidPassword() {
	
	// call the Register method with invalid password on usecase
	err := suite.usecase.Register(&domain.User{
		Username: "user", 
		Password: "123",
	})

	// verify error response
	assert.ErrorContains(suite.T(), err, "at least 8 characters")      // error should match expected message
}

// tests registration with empty username
func (suite *UserUseCaseTestSuite) TestRegister_EmptyUsername() {
    
	// create test user with empty username
	user := &domain.User{
        Username: "",
        Password: "password123",
    }

	// call the Register method on usecase
    err := suite.usecase.Register(user)
    assert.EqualError(suite.T(), err, "username cannot be empty")      // error should match expected message
}

// tests registration with empty password
func (suite *UserUseCaseTestSuite) TestRegister_EmptyPassword() {

	// create test user with empty password
    user := &domain.User{
        Username: "user",
        Password: "",
    }

	// call the Register method on usecase
    err := suite.usecase.Register(user)
    assert.EqualError(suite.T(), err, "password cannot be empty")       // error should match expected message
}

// tests registration with password less than 8 characters
func (suite *UserUseCaseTestSuite) TestRegister_ShortPassword() {

	// create test user with short password
    user := &domain.User{
        Username: "user",
        Password: "short",
    }

	// call the Register method on usecase
    err := suite.usecase.Register(user)
    assert.EqualError(suite.T(), err, "password must be at least 8 characters")      // error should match expected message
}

// tests Register when repository returns unexpected error on GetByUsername
func (suite *UserUseCaseTestSuite) TestRegister_RepoErrorOnGetByUsername() {
    
	// create test user
	user := &domain.User{
        Username: "user",
        Password: "password123",
    }
	
	// mock GetByUsername of the repository to return nil and error
    suite.userRepo.
        On("GetByUsername", user.Username).
        Return(nil, errors.New("db error"))

	// call the Register method on usecase
    err := suite.usecase.Register(user)
    assert.EqualError(suite.T(), err, "db error")        // error should match expected message
}

// tests Register when password hashing fails
func (suite *UserUseCaseTestSuite) TestRegister_HashPasswordError() {
    
	// create test user
	user := &domain.User{
        Username: "user",
        Password: "password123",
    }

	// mock GetByUsername of the repository to return and error
    suite.userRepo.
        On("GetByUsername", user.Username).
        Return(nil, domain.ErrUserNotFound)
	// mock HashPassword of the repository to return empty string and error
    suite.pwdService.
        On("HashPassword", user.Password).
        Return("", errors.New("hash error"))
	// mock GetUserCount of the repository to return number and nil
    suite.userRepo.
        On("GetUserCount").
        Return(int64(1), nil)

	// call the Register method on usecase
    err := suite.usecase.Register(user)
    assert.EqualError(suite.T(), err, "hash error")       // error should match expected message
}

// tests Register when GetUserCount fails
func (suite *UserUseCaseTestSuite) TestRegister_GetUserCountError() {
    
	// create test user
	user := &domain.User{
        Username: "user",
        Password: "password123",
    }

	// mock GetByUsername of the repository to return nil and error
    suite.userRepo.
        On("GetByUsername", user.Username).
        Return(nil, domain.ErrUserNotFound)
	// mock HashPassword of the repository to return error
    suite.pwdService.
        On("HashPassword", user.Password).
        Return("hashedpass", nil)
	// mock GetUserCount of the repository to return error
    suite.userRepo.
        On("GetUserCount").
        Return(int64(0), errors.New("count error"))

	// call the Register method on usecase
    err := suite.usecase.Register(user)
    assert.EqualError(suite.T(), err, "count error")       // error should match expected message
}

// tests successful user login
func (suite *UserUseCaseTestSuite) TestLogin_Success() {

	// create test user 
	user := &domain.User{
		ID: primitive.NewObjectID(), 
		Username: "testuser", 
		Password: "hashedpass", 
		Role: "user",
	}

	// create test credentials
	credentials := &domain.Credentials{
		Username: "testuser", 
		Password: "password123",
	}

	// mock GetByUsername of the repository to return the test user
	suite.userRepo.
		On("GetByUsername", credentials.Username).
		Return(user, nil)
	// mock GetByUsername of the respository to return true
	suite.pwdService.
		On("CheckPassword", user.Password, credentials.Password).
		Return(true)
	// mock GenerateToken of the JWT service to return a token
	suite.jwtService.
		On("GenerateToken", user.ID.Hex(), user.Username, user.Role).
		Return("token123", nil)

	// call the Login method on usecase
	token, returnUser, err := suite.usecase.Login(credentials)

	// verify results
	assert.NoError(suite.T(), err)                                 // no error expected
	assert.Equal(suite.T(), "token123", token)                 	   // token should match mock response
	assert.Equal(suite.T(), user.ID, returnUser.ID)            	   // returned user should match
	assert.Equal(suite.T(), "testuser", returnUser.Username)       // username should match
}

// tests login with invalid password
func (suite *UserUseCaseTestSuite) TestLogin_InvalidPassword() {
	
	// create test user with hashed password
	user := &domain.User{
		Username: "user", 
		Password: "hashed",
	}
	// create test credentials with wrong password
	creds := &domain.Credentials{
		Username: "user", 
		Password: "wrong",
	}

	// mock GetByUsername of the repository to return the test user
	suite.userRepo.
		On("GetByUsername", creds.Username).Return(user, nil)
	// mock CheckPassword of the password service to return false
	suite.pwdService.
		On("CheckPassword", user.Password, creds.Password).
		Return(false)

	// call the Login method on usecase
	_, _, err := suite.usecase.Login(creds)

	// verify error response
	assert.ErrorIs(suite.T(), err, domain.ErrInvalidCredentials)      // error should be invalid credentials
}

// tests login with non-existent user
func (suite *UserUseCaseTestSuite) TestLogin_UserNotFound() {
	
	// create credentials for non-existent user
	creds := &domain.Credentials{
		Username: "nouser", 
		Password: "pass",
	}
	
	// mock GetByUsername of the repository to return error
	suite.userRepo.
		On("GetByUsername", creds.Username).
		Return(nil, domain.ErrUserNotFound)

	// call the Login method on usecase
	_, _, err := suite.usecase.Login(creds)

	// verify error response
	assert.ErrorIs(suite.T(), err, domain.ErrInvalidCredentials)      // error should be invalid credentials
}

// tests Login with empty username or password
func (suite *UserUseCaseTestSuite) TestLogin_EmptyCredentials() {
    
	// create test empty login credentials
	creds := &domain.Credentials{
        Username: "",
        Password: "",
    }

	// call the Login method on usecase
    token, user, err := suite.usecase.Login(creds)
    assert.Empty(suite.T(), token)                            // token should be empty 
    assert.Nil(suite.T(), user)                               // user should be nil
    assert.EqualError(suite.T(), err, "username and password are required")      // error should match expected message
}

// tests Login when repository returns error other than ErrUserNotFound
func (suite *UserUseCaseTestSuite) TestLogin_RepoErrorOnGetByUsername() {
    
	// create test login credentials
	creds := &domain.Credentials{
        Username: "user",
        Password: "password123",
    }

	// mock GetByUsername of the repository to return nil and error
    suite.userRepo.
        On("GetByUsername", creds.Username).
        Return(nil, errors.New("db error"))

	// call the Login method on usecase
    token, user, err := suite.usecase.Login(creds)
    assert.Empty(suite.T(), token)                       // token should be empty
    assert.Nil(suite.T(), user)                          // user should be nil
    assert.EqualError(suite.T(), err, "db error")        // error should match expected message
}

// tests Login when JWT generation fails
func (suite *UserUseCaseTestSuite) TestLogin_JWTGenerationError() {
    
	// create test user
	user := &domain.User{
        ID:       primitive.NewObjectID(),
        Username: "user",
        Password: "hashedpass",
        Role:     "user",
    }
	// create test login credentials
    creds := &domain.Credentials{
        Username: "user",
        Password: "password123",
    }

	// mock GetByUsername of the repository to return user and error
    suite.userRepo.
        On("GetByUsername", creds.Username).
        Return(user, nil)
	// mock CheckPassword of the repository to return true
    suite.pwdService.
        On("CheckPassword", user.Password, creds.Password).
        Return(true)
	// mock GenerateToken of the repository to return empty string and error
    suite.jwtService.
        On("GenerateToken", user.ID.Hex(), user.Username, user.Role).
        Return("", errors.New("jwt error"))

	// call the Login method on usecase
    token, returnUser, err := suite.usecase.Login(creds)
    assert.Empty(suite.T(), token)                        // token should be empty 
    assert.Nil(suite.T(), returnUser)                     // return user should be nil
    assert.EqualError(suite.T(), err, "jwt error")        // error should match expected message
}

// tests successful user promotion to admin
func (suite *UserUseCaseTestSuite) TestPromoteToAdmin_Success() {
	
	// create test user ID
	id := primitive.NewObjectID()
	
	// mock GetUserById of the repository to return a user
	suite.userRepo.
		On("GetUserById", id).
		Return(&domain.User{ID: id}, nil)
	// mock UpdateRole of the repository to return nil - successful promotion
	suite.userRepo.
		On("UpdateRole", id, "admin").
		Return(nil)

	// call the PromoteToAdmin method on usecase
	err := suite.usecase.PromoteToAdmin(id.Hex())

	// verify results
	assert.NoError(suite.T(), err)      // no error expected
}

// tests PromoteToAdmin with empty user ID
func (suite *UserUseCaseTestSuite) TestPromoteToAdmin_EmptyID() {
    
	// call the PromoteToAdmin method on usecase
	err := suite.usecase.PromoteToAdmin("")
    assert.EqualError(suite.T(), err, "user ID cannot be empty")        // error should match expected message
}

// tests promotion with non-existent user
func (suite *UserUseCaseTestSuite) TestPromoteToAdmin_UserNotFound() {
	
	// create test user ID
	id := primitive.NewObjectID()

	// mock GetUserById of the repository to return error
	suite.userRepo.
		On("GetUserById", id).
		Return(nil, domain.ErrUserNotFound)

	// call the PromoteToAdmin method on usecase
	err := suite.usecase.PromoteToAdmin(id.Hex())

	// verify error response
	assert.ErrorIs(suite.T(), err, domain.ErrUserNotFound)       // error should be user not found
}

// tests promotion with invalid user ID format
func (suite *UserUseCaseTestSuite) TestPromoteToAdmin_InvalidID() {

	// mock GetUserById of the repository to return error
	suite.userRepo.
		On("GetUserById", mock.Anything).
		Return(nil, domain.ErrInvalidUserID)

	// call the PromoteToAdmin method with invalid ID format
	err := suite.usecase.PromoteToAdmin("invalid")

	// verify error response
	assert.ErrorIs(suite.T(), err, domain.ErrInvalidUserID)      // error should be invalid user ID
}

// tests PromoteToAdmin when UpdateRole fails
func (suite *UserUseCaseTestSuite) TestPromoteToAdmin_UpdateRoleError() {
    
	// mock user id
	id := primitive.NewObjectID()

	// mock GetUserById of the repository to return user and nil
    suite.userRepo.
        On("GetUserById", id).
        Return(&domain.User{ID: id}, nil)
	// mock UpdateRole of the repository to return error
    suite.userRepo.
        On("UpdateRole", id, "admin").
        Return(errors.New("update error"))

	// call the PromoteToAdmin method on usecase
    err := suite.usecase.PromoteToAdmin(id.Hex())
    assert.EqualError(suite.T(), err, "update error")       // error should match expected message
}

// runs the test suite for UserUseCase
func TestUserUseCaseTestSuite(t *testing.T) {
	suite.Run(t, new(UserUseCaseTestSuite))       // run the test suite
}