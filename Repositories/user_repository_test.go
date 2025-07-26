package repositories

// imports
import (
	"context"
	"errors"
	"testing"
	"time"
	"github.com/natnael-eyuel-dev/Task-Management-Unit-Test/Domain"
	"github.com/natnael-eyuel-dev/Task-Management-Unit-Test/Repositories/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// test suite for the UserRepository
type UserRepositoryTestSuite struct {
    suite.Suite                                             // embed the suite.Suite type
    mockCollection *mock_repositories.MockCollection        // mock collection for testing
    repo           domain.UserRepository                    // user repository to be tested
}

// initializes the test suite
func (suite *UserRepositoryTestSuite) SetupTest() {
    suite.mockCollection = new(mock_repositories.MockCollection)                           // create a new mock collection
    suite.repo = NewUserRepositoryWithCollection(suite.mockCollection)        // create a new user repository with mock collection
}

// tests CreateUser method of the UserRepository
func (suite *UserRepositoryTestSuite) TestCreateUser_Success() {
    
	// create a new user
	user := &domain.User{
		Username: "testuser",
		Password: "securepass123",
		Role:     "user",
    }

	// mock the InsertOne method of the collection
    suite.mockCollection.
        On("InsertOne", mock.Anything, mock.Anything).
        Return(&mongo.InsertOneResult{}, nil)

    err := suite.repo.CreateUser(user)       // call CreateUser method
    assert.NoError(suite.T(), err)           // assert no error
    assert.NotZero(suite.T(), user.ID)       // assert ID is not zero
}

// tests CreateUser method of the UserRepository for duplicate user
func (suite *UserRepositoryTestSuite) TestCreateUser_Duplicate() {
    
	// create a new user
	user := &domain.User{
        Username: "existing",
		Password: "securepass123",
		Role:     "user",
    }

	// mock the InsertOne method of the collection
    suite.mockCollection.
        On("InsertOne", mock.Anything, mock.Anything).
        Return(nil, mongo.WriteException{
            WriteErrors: []mongo.WriteError{{Code: 11000}},
        })

    err := suite.repo.CreateUser(user)                          // call CreateUser method
    assert.ErrorIs(suite.T(), err, domain.ErrUserExists)        // assert error is ErrUserExists
}

// tests CreateUser method of the UserRepository for error case
func (suite *UserRepositoryTestSuite) TestCreateUser_Error() {
    
    // create a new user
    user := &domain.User{
        Username: "erroruser",
        Password: "securepass123",
        Role:     "user",
    }

    // mock the InsertOne method of the collection
    suite.mockCollection.
        On("InsertOne", mock.Anything, mock.Anything).
        Return(nil, errors.New("insert error"))

    err := suite.repo.CreateUser(user)                     // call CreateUser method
    assert.EqualError(suite.T(), err, "insert error")      // assert error message matches
}

// tests CreateUser method of the UserRepository for context timeout
func (suite *UserRepositoryTestSuite) TestCreateUser_ContextTimeout() {
    
    // create a context with a very short timeout
    _, cancel := context.WithTimeout(context.Background(), 1*time.Nanosecond)
    defer cancel()

    // create a user
    user := &domain.User{Username: "timeout"}

    // mock the InsertOne method of the collection 
    suite.mockCollection.
        On("InsertOne", mock.Anything, user).
        Return(nil, context.DeadlineExceeded)

    err := suite.repo.CreateUser(user)                            // call CreateUser method 
    assert.ErrorIs(suite.T(), err, context.DeadlineExceeded)      // assert error is context deadline exceeded
}

// tests GetByUsername method of the UserRepository for existing user
func (suite *UserRepositoryTestSuite) TestGetByUsername_Success() {
    
	// create a new username 
	username := "john"
	// create a mock user
    expected := domain.User{
		ID:       primitive.NewObjectID(),
		Username: username,
		Role:     "user",
	}

	// mock the FindOne method of the collection
    suite.mockCollection.
        On("FindOne", mock.Anything, bson.M{"username": username}).
        Return(&mock_repositories.MockSingleResult{Err: nil, Result: &expected})

    user, err := suite.repo.GetByUsername(username)        // call GetByUsername method
    assert.NoError(suite.T(), err)                         // assert no error
    assert.Equal(suite.T(), username, user.Username)       // assert username matches
}

// tests GetByUsername method of the UserRepository for non-existing user
func (suite *UserRepositoryTestSuite) TestGetByUsername_NotFound() {
    
	// create a new username
	username := "not_found"

	// mock the FindOne method of the collection
    suite.mockCollection.
        On("FindOne", mock.Anything, bson.M{"username": username}).
        Return(&mock_repositories.MockSingleResult{Err: mongo.ErrNoDocuments})

    user, err := suite.repo.GetByUsername(username)              // call GetByUsername method
    assert.Nil(suite.T(), user)                                  // assert user is nil
    assert.ErrorIs(suite.T(), err, domain.ErrUserNotFound)       // assert error is ErrUserNotFound
}

// tests GetByUsername method of the UserRepository for empty username
func (suite *UserRepositoryTestSuite) TestGetByUsername_EmptyUsername() {
    
    user, err := suite.repo.GetByUsername("")                               // call GetByUsername method
    assert.Nil(suite.T(), user)                                             // assert user is nil
    assert.ErrorContains(suite.T(), err, "username cannot be empty")        // assert error contains message
}

// tests GetUserById method of the UserRepository for existing user
func (suite *UserRepositoryTestSuite) TestGetUserById_Success() {
    
    // create a new object ID
    id := primitive.NewObjectID()
    expected := domain.User{ID: id}

    // mock the FindOne method of the collection
    suite.mockCollection.
        On("FindOne", mock.Anything, bson.M{"_id": id}).
        Return(&mock_repositories.MockSingleResult{Err: nil, Result: &expected})

    user, err := suite.repo.GetUserById(id)              // call GetUserById method
    assert.NoError(suite.T(), err)                       // assert no error
    assert.Equal(suite.T(), id, user.ID)                 // assert ID matches
}

// tests GetUserById method of the UserRepository for non-existing user
func (suite *UserRepositoryTestSuite) TestGetUserById_NotFound() {

    // create a new object ID
    id := primitive.NewObjectID()

    // mock the FindOne method of the collection
    suite.mockCollection.
        On("FindOne", mock.Anything, bson.M{"_id": id}).
        Return(&mock_repositories.MockSingleResult{Err: mongo.ErrNoDocuments})

    user, err := suite.repo.GetUserById(id)                      // call GetUserById method
    assert.Nil(suite.T(), user)                                  // assert user is nil
    assert.ErrorIs(suite.T(), err, domain.ErrUserNotFound)       // assert error is ErrUserNotFound
}

// tests GetUserCount method of the UserRepository
func (suite *UserRepositoryTestSuite) TestGetUserCount_Success() {

	// mock the CountDocuments method of the collection
	suite.mockCollection.
        On("CountDocuments", mock.Anything, bson.M{}).
        Return(int64(42), nil)

    count, err := suite.repo.GetUserCount()         // call GetUserCount method
    assert.NoError(suite.T(), err)                  // assert no error
    assert.Equal(suite.T(), int64(42), count)       // assert count matches
}

// tests GetUserCount method of the UserRepository for error case
func (suite *UserRepositoryTestSuite) TestGetUserCount_Error() {
    
    // mock the CountDocuments method of the collection
    suite.mockCollection.
        On("CountDocuments", mock.Anything, bson.M{}).
        Return(int64(0), errors.New("count error"))

    count, err := suite.repo.GetUserCount()               // call GetUserCount method
    assert.Equal(suite.T(), int64(0), count)              // assert count is zero
    assert.EqualError(suite.T(), err, "count error")      // assert error message matches
}

// tests GetUserCount method of the UserRepository for zero users
func (suite *UserRepositoryTestSuite) TestGetUserCount_ZeroUsers() {
    
    // mock the CountDocuments method of the collection
    suite.mockCollection.
        On("CountDocuments", mock.Anything, bson.M{}).
        Return(int64(0), nil)

    count, err := suite.repo.GetUserCount()      // call GetUserCount method
    assert.NoError(suite.T(), err)               // assert no error
    assert.Zero(suite.T(), count)                // assert count zero
}

// tests GetByUsername method of the UserRepository for error case
func (suite *UserRepositoryTestSuite) TestGetByUsername_Error() {

    // test username
    username := "erroruser"

    // mock the FindOne method of the collection
    suite.mockCollection.
        On("FindOne", mock.Anything, bson.M{"username": username}).
        Return(&mock_repositories.MockSingleResult{Err: errors.New("find error")})

    user, err := suite.repo.GetByUsername(username)       // call GetByUsername method
    assert.Nil(suite.T(), user)                           // assert user is nil
    assert.EqualError(suite.T(), err, "find error")       // assert error message matches
}

// tests GetUserById method of the UserRepository for error case
func (suite *UserRepositoryTestSuite) TestGetUserById_Error() {

    // create a new object ID
    id := primitive.NewObjectID()

    // mock the FindOne method of the collection
    suite.mockCollection.
        On("FindOne", mock.Anything, bson.M{"_id": id}).
        Return(&mock_repositories.MockSingleResult{Err: errors.New("find error")})

    user, err := suite.repo.GetUserById(id)               // call GetUserById method
    assert.Nil(suite.T(), user)                           // assert user is nil
    assert.EqualError(suite.T(), err, "find error")       // assert error message matches
}

// tests UpdateRole method of the UserRepository for existing user
func (suite *UserRepositoryTestSuite) TestUpdateRole_Success() {
    
	// create a new object ID 
	id := primitive.NewObjectID()
	// create a new role
    role := "admin"

	// mock the FindOneAndUpdate method of the collection
    suite.mockCollection.
        On("FindOneAndUpdate", mock.Anything, bson.M{"_id": id}, bson.M{"$set": bson.M{"role": role}}).
        Return(&mock_repositories.MockSingleResult{Err: nil, Result: &domain.User{ID: id, Role: role}})

    err := suite.repo.UpdateRole(id, role)        // call UpdateRole method
	assert.NoError(suite.T(), err)                // assert no error
    assert.NoError(suite.T(), err)                // assert no error
}

// tests UpdateRole method of the UserRepository for non-existing user
func (suite *UserRepositoryTestSuite) TestUpdateRole_NotFound() {
    
	// create a new object ID
	id := primitive.NewObjectID()
	// create a new role
    role := "admin"

	// mock the FindOneAndUpdate method of the collection
    suite.mockCollection.
        On("FindOneAndUpdate", mock.Anything, bson.M{"_id": id}, bson.M{"$set": bson.M{"role": role}}).
        Return(&mock_repositories.MockSingleResult{Err: mongo.ErrNoDocuments})

    err := suite.repo.UpdateRole(id, role)                       // call UpdateRole method
	assert.Error(suite.T(), err)                                 // assert error is returned
    assert.ErrorIs(suite.T(), err, domain.ErrUserNotFound)       // assert error is ErrUserNotFound
}

// tests UpdateRole method of the UserRepository for error case
func (suite *UserRepositoryTestSuite) TestUpdateRole_Error() {
    
	// create a new object ID
	id := primitive.NewObjectID()
	// create a new role
    role := "admin"

	// mock the FindOneAndUpdate method of the collection
    suite.mockCollection.
        On("FindOneAndUpdate", mock.Anything, bson.M{"_id": id}, bson.M{"$set": bson.M{"role": role}}).
        Return(&mock_repositories.MockSingleResult{Err: errors.New("db error")})

    err := suite.repo.UpdateRole(id, role)                       // call UpdateRole method
    assert.Error(suite.T(), err)                                 // assert error is returned
    assert.Equal(suite.T(), err.Error(), "db error")             // assert error message
}

// tests UpdateRole method of the UserRepository for empty role
func (suite *UserRepositoryTestSuite) TestUpdateRole_EmptyRole() {

    err := suite.repo.UpdateRole(primitive.NewObjectID(), "")           // call UpdateRole method 
    assert.ErrorContains(suite.T(), err, "role cannot be empty")        // assert error contains message
}

// tests UpdateRole method of the UserRepository for invalid role value
func (suite *UserRepositoryTestSuite) TestUpdateRole_InvalidRole() {

    // create a new object ID
    id := primitive.NewObjectID()
    // create invalid role not empty 
    role := "invalid_role"     

    // mock the FindOneAndUpdate method of the collection
    suite.mockCollection.
        On("FindOneAndUpdate", mock.Anything, bson.M{"_id": id}, bson.M{"$set": bson.M{"role": role}}).
        Return(&mock_repositories.MockSingleResult{Err: errors.New("invalid role")})

    err := suite.repo.UpdateRole(id, role)                     // call UpdateRole method
    assert.Error(suite.T(), err)                               // assert error is returned
    assert.Equal(suite.T(), err.Error(), "invalid role")       // assert error message
}

// suite entry point for running the tests
func TestUserRepositoryTestSuite(t *testing.T) {
    suite.Run(t, new(UserRepositoryTestSuite))        // run the test suite
}
