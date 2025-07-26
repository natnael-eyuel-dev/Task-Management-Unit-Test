package routers

// imports
import (
	"github.com/gin-gonic/gin"
	"github.com/natnael-eyuel-dev/Task-Management-Unit-Test/Delivery/controllers"
	"github.com/natnael-eyuel-dev/Task-Management-Unit-Test/Domain"
	"github.com/natnael-eyuel-dev/Task-Management-Unit-Test/Infrastructure"
)

// setup router
func SetupRouter( taskUsc domain.TaskUseCase, userUsc domain.UserUseCase, jwtServ domain.JWTService) *gin.Engine {

	router := gin.Default()     // create default gin router

	taskContrl := controllers.NewTaskController(taskUsc)        // initialize task controller with task usecase
	userContrl := controllers.NewUserController(userUsc)        // initialize user controller with user usecase

	// public routes
	router.POST("/register", userContrl.Register)         // register new user
	router.POST("/login", userContrl.Login)               // authenticate a user

	// authenticated routes
	authMiddleware := infrastructure.NewAuthMiddleware(jwtServ)

	authGroup := router.Group("")
	authGroup.Use(authMiddleware.Handler())
	{
		authGroup.GET("/tasks", taskContrl.GetAllTasks)             // get all tasks
		authGroup.GET("/tasks/:id", taskContrl.GetTaskByID)         // get specific task by id
	}

	// admin routes
	adminMiddleware := infrastructure.AdminOnly()

	adminGroup := router.Group("")
	adminGroup.Use(authMiddleware.Handler(), adminMiddleware)
	{
		adminGroup.POST("/tasks", taskContrl.CreateTask)                 // create new task
		adminGroup.PUT("/tasks/:id", taskContrl.UpdateTask)              // update existing task by id
		adminGroup.DELETE("/tasks/:id", taskContrl.DeleteTask)           // delete existing task by id
		adminGroup.PUT("/promote/:id", userContrl.PromoteToAdmin)        // promote user to admin by id
	}

	return router        // return configured router
}