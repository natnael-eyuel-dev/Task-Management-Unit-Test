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

type taskRepository struct {
	collection domain.MongoCollection
}

// creates a new user repository instance
func NewTaskRepository() domain.TaskRepository {
	// setup mongodb
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)       // set timeout
	defer cancel()

	// connect
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}

	db := client.Database("taskmanager")
	taskCol := db.Collection("tasks")         // initialize task collection
	return &taskRepository{&adapters.MongoCollectionAdapter{Collection: taskCol}}
}

// this is used for testing purposes to inject a mock collection
func NewTaskRepositoryWithCollection(coll domain.MongoCollection) domain.TaskRepository {
	return &taskRepository{coll}
}

func (taskRepo *taskRepository) CreateTask(task *domain.Task) (*domain.Task, error) {
	
	contx, cancel := context.WithTimeout(context.Background(), 5*time.Second)     // set timeout
	defer cancel()

	task.ID = primitive.NewObjectID()                         // create a unique id for the new task
	_, err := taskRepo.collection.InsertOne(contx, task)      // create the new task with error handling
	if err != nil {
        return nil, err
    }

	return task, nil       // return the new created task and nil
}

func (taskRepo *taskRepository) DeleteTask(taskID string) error {
	
	contx, cancel := context.WithTimeout(context.Background(), 5*time.Second)        // set timeout
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(taskID)       // convert string id to mongodb's id format with error handling 
	if err != nil {
		return domain.ErrInvalidTaskID
	}

	result, err := taskRepo.collection.DeleteOne(contx, bson.M{"_id": objID})       // delete the task with error handling
	if err != nil {
		return err
	}

	if result == nil {
    	return errors.New("delete error")
	}

	// verify task deleted
	if result.DeletedCount == 0 {
		return domain.ErrTaskNotFound
	}

	return nil
}

func (taskRepo *taskRepository) GetAllTasks() ([]domain.Task, error) {
	
	var allTasks []domain.Task
	contx, cancel := context.WithTimeout(context.Background(), 5*time.Second)        // set timeout
	defer cancel()

	cursor, err := taskRepo.collection.Find(contx, bson.M{})      // find all documents in the collection
	if err != nil {
		return nil, err
	}

	if cursor == nil {
		return nil, errors.New("find error")
	}

	defer cursor.Close(contx)      // close cursor when done

	err = cursor.All(contx, &allTasks)      // read all result into our slice
	if err != nil {  
		return nil, err
	}

	if allTasks == nil {
		return []domain.Task{}, nil
	}

	return allTasks, nil
}

func (taskRepo *taskRepository) GetTaskByID(taskID string) (*domain.Task, error) {
	
	var task domain.Task
	contx, cancel := context.WithTimeout(context.Background(), 5*time.Second)        // set timeout
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(taskID)      // convert string id to mongodb's format with error handling 
	if err != nil {
		return nil, domain.ErrInvalidTaskID
	}

	err = taskRepo.collection.FindOne(contx, bson.M{"_id": objID}).Decode(&task)       // check if task exists
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, domain.ErrTaskNotFound
		}
		return nil, err
	}

	return &task, nil
}

func (taskRepo *taskRepository) UpdateTask(taskID string, taskUpdate *domain.Task) (*domain.Task, error) {
	
	var updatedTask domain.Task
	contx, cancel := context.WithTimeout(context.Background(), 5*time.Second)        // set timeout
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(taskID)      // convert string id to mongodb's format with error handling 
	if err != nil {
		return nil, domain.ErrInvalidTaskID
	}

	update := bson.M{"$set": bson.M{}}
	setFields := update["$set"].(bson.M)        // prepare what we want to change

	// only update fields that were actually provided
	if taskUpdate.Title != "" {
		setFields["title"] = taskUpdate.Title
	}
	if taskUpdate.Description != "" {
		setFields["description"] = taskUpdate.Description
	}
	if !taskUpdate.DueDate.IsZero() {
		setFields["due_date"] = taskUpdate.DueDate
	}
	if taskUpdate.Status != "" {
		setFields["status"] = taskUpdate.Status
	}

	// stop if nothing valid to update
	if len(setFields) == 0 {
		return nil, errors.New("no valid fields provided for update")
	}
 
	opts := options.FindOneAndUpdate().         // to get updated document back
		SetReturnDocument(options.After)

	// perform update and get the updated task
	err = taskRepo.collection.FindOneAndUpdate(
		contx,
		bson.M{"_id": objID},
		update,
		opts,
	).Decode(&updatedTask)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, domain.ErrTaskNotFound
		}
		return nil, err
	}

	return &updatedTask, nil       // return the updated task and nil
}

