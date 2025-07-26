package domain

// imports
import (
	"context"
	"errors"
	"time"
	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"				
	"go.mongodb.org/mongo-driver/mongo/options"
)

// task item
type Task struct {
	ID              primitive.ObjectID         // unique identifier of task 
	Title           string                     // title of task
	Description     string                     // description of task
	DueDate         time.Time                  // due date of task 
	Status          string                     // status of task
}

// user item
type User struct {
	ID              primitive.ObjectID         // unique identifier for users 
	Username     	string                     // username 
	Password     	string                     // password - hashed before storage
	Role         	string                     // user role - role/user 
}

// credential item
type Credentials struct {
	Username 	 string        `binding:"required"`      // login username - required
    Password 	 string 	   `binding:"required"`      // login password - required
}

// claim item
type Claims struct {
	ID           primitive.ObjectID         // id for claim
	Username     string                     // username for claim
	Role         string      			    // role for claim
}

// task repository interface 
type TaskRepository interface {
	CreateTask(task *Task) (*Task, error)                     // create new task with validation
	DeleteTask(taskID string) error                 		  // delete existing task or return error if not found
	GetAllTasks() ([]Task, error)         					  // get all tasks in the system
	GetTaskByID(taskID string) (*Task, error) 				  // get specific task by id or return error if not found
	UpdateTask(taskID string, task *Task) (*Task, error)      // update existing task or return error if not found
}

// user repository interface
type UserRepository interface {    
	CreateUser(user *User) error                              // create new user with validation
	GetByUsername(username string) (*User, error)             // get specific user by username or return error if not found
	GetUserById(id primitive.ObjectID) (*User, error)         // get specific user by id or return error if not found
	GetUserCount() (int64, error)                             // get total user count or return error 
	UpdateRole(id primitive.ObjectID, role string) error      // update user's role to admin or return error if not found                            
}

// task usecase interface
type TaskUseCase interface {
	CreateTask(task *Task) (*Task, error)                     // create new task with validation
	DeleteTask(taskID string) error                 		  // delete existing task or return error if not found
	GetAllTasks() ([]Task, error)         					  // get all tasks in the system
	GetTaskByID(taskID string) (*Task, error) 				  // get specific task by id or return error if not found
	UpdateTask(taskID string, task *Task) (*Task, error)      // update existing task or return error if not found
}

// user usecase interface
type UserUseCase interface {
	Register(user *User) error                                 // register new user with validation
	Login(credentials *Credentials) (string, *User, error)     // authenticate user and return token, user or error
	PromoteToAdmin(userID string) error                        // promote user to admin role or return error if not found
}

// jwt service interface
type JWTService interface {
	GenerateToken(userID, username, role string) (string, error)       	// generate token or return error
	ValidateToken(tokenStr string) (*jwt.Token, error)                 	// validate token or return error
}

// password service interface
type PasswordService interface {
	HashPassword(password string) (string, error)       	   // hash password or return error
	CheckPassword(hashed, plain string) bool            	   // check password and return bool (true/false)
}

// single result interface 
type SingleResult interface {
	Decode(v interface{}) error           // decode single result into provided interface
}

// mongo collection interface
type MongoCollection interface {
	InsertOne(context.Context, interface{}, ...*options.InsertOneOptions) (*mongo.InsertOneResult, error)       		// insert one document into collection         
	Find(context.Context, interface{}, ...*options.FindOptions) (*mongo.Cursor, error)                          		// find documents in collection
	FindOne(context.Context, interface{}, ...*options.FindOneOptions) SingleResult                              		// find one document in collection
	FindOneAndUpdate(context.Context, interface{}, interface{}, ...*options.FindOneAndUpdateOptions) SingleResult       // find one document and update it
	DeleteOne(context.Context, interface{}, ...*options.DeleteOptions) (*mongo.DeleteResult, error)                     // delete one document from collection
	CountDocuments(context.Context, interface{}, ...*options.CountOptions) (int64, error)                               // count documents in collection
}

// custom errors
var (
	ErrTaskNotFound     	 = errors.New("task not found")              		 // custom task not found error
	ErrInvalidTaskID     	 = errors.New("invalid task ID")             		 // custom invalid task id error
	ErrUserExists            = errors.New("user already exists")         		 // custom user exists error
	ErrUserNotFound          = errors.New("user not found")              		 // custom user not found error
	ErrInvalidUserID         = errors.New("invalid user ID")             		 // custom invalid user id error
	ErrInvalidCredentials    = errors.New("invalid credentials")        	     // custom invalid credentials error
	ErrUnauthorized          = errors.New("unauthorized access")         		 // custom unauthorized access error
	ErrInvalidDueDate        = errors.New("due date must be in the future")      // custom invalid due date error
)

