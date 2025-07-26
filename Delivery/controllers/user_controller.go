package controllers

// imports
import (
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/natnael-eyuel-dev/Task-Management-Unit-Test/Domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// user controller
type UserController struct {
	userUseCase domain.UserUseCase        // user usecase for user operations 
}

// new user controller
func NewUserController(uc domain.UserUseCase) *UserController {
	return &UserController{userUseCase: uc}        // return new user controller instance
}

func (uc *UserController) Register(c *gin.Context) {
	
	var user domain.User
	err := c.ShouldBindJSON(&user)       // parse request body into user struct
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if user.Username == "" || user.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username and password must be set"})
		return
	}

	// create user through usecase layer
	if err := uc.userUseCase.Register(&user); err != nil {
		if err == domain.ErrUserExists {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "user created successfully"})       // success response
}

func (uc *UserController) Login(c *gin.Context) {
	
	var creds domain.Credentials
	err := c.ShouldBindJSON(&creds)        // parse request body into user struct
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// authenticate user through usecase layer
	token, user, err := uc.userUseCase.Login(&creds)
	if err != nil {
		if err == domain.ErrInvalidCredentials {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// return token, user info (excluding sensitive data)
	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"role":     user.Role,
		},
	})
}

func (uc *UserController) PromoteToAdmin(c *gin.Context) {
	
	userID := c.Param("id")       // get user id from request parameter
	 
	_, err := primitive.ObjectIDFromHex(userID)       // validate it is a valid ObjectID
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
		return
	}

	// promote user through usecase layer
	err = uc.userUseCase.PromoteToAdmin(userID) 
	if err != nil {
		if err == domain.ErrUserNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "user promoted to admin successfully"})       // success response
}