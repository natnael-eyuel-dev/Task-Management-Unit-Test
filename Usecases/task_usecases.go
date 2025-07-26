package usecases

// imports
import (
	"errors"
	"time"
	"github.com/natnael-eyuel-dev/Task-Management-Unit-Test/Domain"
)

type taskUseCase struct {
	taskRepo domain.TaskRepository
}

// creates new TaskUseCase instance
func NewTaskUseCase(repo domain.TaskRepository) domain.TaskUseCase {
	return &taskUseCase{taskRepo: repo}
}

// create a task
func (taskUsc *taskUseCase) CreateTask(task *domain.Task) (*domain.Task, error) {
	
	// validate task fields before creation
	if task.Title == "" {
		return nil, errors.New("task title cannot be empty")
	}
	if task.Description == "" {
		return nil, errors.New("task description cannot be empty")
	}
	if task.DueDate.IsZero() {
		return nil, errors.New("due date cannot be empty")
	}
	if task.Status == "" {
		task.Status = "pending"      // default status
	}
	// validate due date is in the future
	if time.Until(task.DueDate) < 0 {
		return nil, errors.New("due date must be in the future")
	}
	// validate status is one of allowed values
	validStatuses := map[string]bool{
		"pending":      true,
		"in_progress":  true,
		"completed":    true,
	}
	if !validStatuses[task.Status] {
		return nil, errors.New("invalid task status")
	}

	return taskUsc.taskRepo.CreateTask(task)
}

// remove task by its id
func (taskUsc *taskUseCase) DeleteTask(id string) error {
	
	// validate id field 
	if id == "" {
		return errors.New("task ID cannot be empty")
	}
	// verify task exists first
	_, err := taskUsc.taskRepo.GetTaskByID(id)
	if err != nil {
		if err == domain.ErrTaskNotFound {
			return domain.ErrTaskNotFound
		}
		return err
	}

	return taskUsc.taskRepo.DeleteTask(id)
}

// get all tasks 
func (taskUsc *taskUseCase) GetAllTasks() ([]domain.Task, error) {
	
	tasks, err := taskUsc.taskRepo.GetAllTasks()
	if err != nil {
		return nil, err
	}
	// return empty slice 
	if tasks == nil {
		return []domain.Task{}, nil
	}

	return tasks, nil
}

// find task by its id
func (taskUsc *taskUseCase) GetTaskByID(id string) (*domain.Task, error) {
	
	// validate id field 
	if id == "" {
		return nil, errors.New("task ID cannot be empty")
	}

	task, err := taskUsc.taskRepo.GetTaskByID(id)
	if err != nil {
		return nil, err
	}
	if task == nil {
		return nil, domain.ErrTaskNotFound
	}

	return task, nil
}

// update task by its id
func (taskUsc *taskUseCase) UpdateTask(id string, task *domain.Task) (*domain.Task, error) {
	
	// validate id field 
	if id == "" {
		return nil, errors.New("task ID cannot be empty")
	}
	// stop if nothing valid to update
	if task.Title == "" && task.Description == "" && 
	   task.DueDate.IsZero() && task.Status == "" {
		return nil, errors.New("no valid fields provided for update")
	}
	// validate status if provided
	if task.Status != "" {
		validStatuses := map[string]bool{
			"pending":      true,
			"in_progress":  true,
			"completed":    true,
		}
		if !validStatuses[task.Status] {
			return nil, errors.New("invalid task status")
		}
	}
	// validate due date if provided
	if !task.DueDate.IsZero() && time.Until(task.DueDate) < 0 {
		return nil, errors.New("due date must be in the future")
	}

	return taskUsc.taskRepo.UpdateTask(id, task)
}