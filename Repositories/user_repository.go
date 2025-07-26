package repositories

// imports
import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/natnael-eyuel-dev/Task-Management-Unit-Test/Domain"
	"github.com/natnael-eyuel-dev/Task-Management-Unit-Test/Repositories/adapters"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type userRepository struct {
	collection domain.MongoCollection
}

// creates a new user repository instance
func NewUserRepository() domain.UserRepository {
	// setup mongodb
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)       // set timeout
	defer cancel()

	// connect
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}

	db := client.Database("taskmanager")
	userCol := db.Collection("users")         // initialize user collection
	return &userRepository{&adapters.MongoCollectionAdapter{Collection: userCol}}
}

// this is used for testing purposes to inject a mock collection
func NewUserRepositoryWithCollection(coll domain.MongoCollection) domain.UserRepository {
	return &userRepository{coll}
}

//  register user in to database
func (userRepo *userRepository) CreateUser(user *domain.User) error {
	
	contx, cancel := context.WithTimeout(context.Background(), 5*time.Second)        // set timeout
	defer cancel()

	// generate new ObjectID if not set
	if user.ID.IsZero() {
		user.ID = primitive.NewObjectID()
	}

	// save user to database
	_, err := userRepo.collection.InsertOne(contx, user)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return domain.ErrUserExists
		}
		return err
	}

	return nil        // success
}

// find user from database by username
func (userRepo *userRepository) GetByUsername(username string) (*domain.User, error) {

	// check username
	if username == "" {
		return nil, errors.New("username cannot be empty")
	}
	
	var user domain.User
	contx, cancel := context.WithTimeout(context.Background(), 5*time.Second)        // set timeout
	defer cancel()
	
	// find user by username
	err := userRepo.collection.FindOne(contx, bson.M{"username": username}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}

	return &user, nil        // success
}

// find user from database by id
func (userRepo *userRepository) GetUserById(userID primitive.ObjectID) (*domain.User, error) {
	
	var user domain.User
	contx, cancel := context.WithTimeout(context.Background(), 5*time.Second)        // set timeout
	defer cancel()
	
	// find user by id
	err := userRepo.collection.FindOne(contx, bson.M{"_id": userID}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}

	return &user, nil         // success
}

// count users in the database currently
func (userRepo *userRepository) GetUserCount() (int64, error) {
	
	contx, cancel := context.WithTimeout(context.Background(), 5*time.Second)        // set timeout
	defer cancel()

	// count users in user collection currently
	count, err := userRepo.collection.CountDocuments(contx, bson.M{})
	if err != nil {
		return 0, err
	}

	return count, nil        // success
}

// update user role to admin in database (only admins can perform this operation)
func (userRepo *userRepository) UpdateRole(id primitive.ObjectID, role string) error {
	
	if role == "" {
		return errors.New("role cannot be empty")
	}

	contx, cancel := context.WithTimeout(context.Background(), 5*time.Second)        // set timeout
	defer cancel()

	// update user's role to admin
	result := userRepo.collection.FindOneAndUpdate(
		contx,
		bson.M{"_id": id},
		bson.M{"$set": bson.M{"role": role}},
	)

	var updated domain.User

	if err := result.Decode(&updated); err != nil {
		if err == mongo.ErrNoDocuments {
			return domain.ErrUserNotFound
		}
		return err
	}

	return nil        // success
}