package usecases

// imports
import (
	"errors"
	"github.com/natnael-eyuel-dev/Task-Management-Unit-Test/Domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)


type userUseCase struct {
	userRepo     domain.UserRepository
	jwtService  domain.JWTService
	pwdService   domain.PasswordService
}

// creates new UserUseCase instance
func NewUserUseCase(userRepo domain.UserRepository, jwtServ domain.JWTService, pwdServ domain.PasswordService,) domain.UserUseCase {
	return &userUseCase{ userRepo:userRepo, jwtService:jwtServ, pwdService:pwdServ}
}

// register user
func (userUsc *userUseCase) Register(user *domain.User) error {
	
	// validate input
	if user.Username == "" {
		return errors.New("username cannot be empty")
	}
	if user.Password == "" {
		return errors.New("password cannot be empty")
	}
	if len(user.Password) < 8 {
		return errors.New("password must be at least 8 characters")
	}
	// check if user already exists
	existing, err := userUsc.userRepo.GetByUsername(user.Username)
	if err != nil && err != domain.ErrUserNotFound {
		return err
	}
	if existing != nil {
		return domain.ErrUserExists
	}

	// hash password securely 
	hashed, err := userUsc.pwdService.HashPassword(user.Password)
	if err != nil {
		return err
	}
	user.Password = hashed       // set user password to hashed password

	// set default role
	user.Role = "user"

	// first user becomes admin
	count, err := userUsc.userRepo.GetUserCount()
	if err != nil {
		return err
	}
	if count == 0 {
		user.Role = "admin"
	}

	return userUsc.userRepo.CreateUser(user)
}

// authenticate user
func (userUsc *userUseCase) Login(credentials *domain.Credentials) (string, *domain.User, error) {
	
	// validate input
	if credentials.Username == "" || credentials.Password == "" {
		return "", nil, errors.New("username and password are required")
	}

	// get user from repository
	user, err := userUsc.userRepo.GetByUsername(credentials.Username)
	if err != nil {
		if err == domain.ErrUserNotFound {
			return "", nil, domain.ErrInvalidCredentials
		}
		return "", nil, err
	}

	// verify password
	if !userUsc.pwdService.CheckPassword(user.Password, credentials.Password) {
		return "", nil, domain.ErrInvalidCredentials
	}

	// generate jwt token
	token, err := userUsc.jwtService.GenerateToken(user.ID.Hex(), user.Username, user.Role)
	if err != nil {
		return "", nil, err
	}

	// return token and user (without sensitive data)
	returnUser := &domain.User{
		ID:       user.ID,
		Username: user.Username,
		Role:     user.Role,
	}

	return token, returnUser, nil
}

// promote a user to admin role (only admin can do this)
func (userUsc *userUseCase) PromoteToAdmin(userID string) error {
	
	// validate input
	if userID == "" {
		return errors.New("user ID cannot be empty")
	}

	objID, err := primitive.ObjectIDFromHex(userID)        // convert string id to ObjectID
	if err != nil {
		return domain.ErrInvalidUserID
	}

	// check if user exists
	_, err = userUsc.userRepo.GetUserById(objID)
	if err != nil {
		if err == domain.ErrUserNotFound {
			return domain.ErrUserNotFound
		}
		return err
	}

	// update role
	return userUsc.userRepo.UpdateRole(objID, "admin")
}