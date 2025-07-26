package main

// imports
import (
	"log"
	"github.com/natnael-eyuel-dev/Task-Management-Unit-Test/Delivery/routers"
	"github.com/natnael-eyuel-dev/Task-Management-Unit-Test/Infrastructure"
	"github.com/natnael-eyuel-dev/Task-Management-Unit-Test/Repositories"
	"github.com/natnael-eyuel-dev/Task-Management-Unit-Test/Usecases"
)

// entry point of the Task Management application
func main() {

	jwtservice, _ := infrastructure.NewJWTService()              // setup jwt service infrastructure
	passwordService := infrastructure.NewPasswordService()       // setup password service infrastructure

	taskRepo := repositories.NewTaskRepository()       // setup task repositorie
	userRepo := repositories.NewUserRepository()       // setup user repositorie

	taskUC := usecases.NewTaskUseCase(taskRepo)                                    // setup task use case
	userUC := usecases.NewUserUseCase(userRepo, jwtservice, passwordService)       // setup user use case

	router := routers.SetupRouter(taskUC, userUC, jwtservice)       // initialize the router with all configured routes

	// start the server on port 8080
	router.Run(":8080")                        
	log.Println("Starting server on :8080")
}
