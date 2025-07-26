package infrastructure

// imports
import (
	"github.com/natnael-eyuel-dev/Task-Management-Unit-Test/Domain"
	"golang.org/x/crypto/bcrypt"
)

// implements the domain.PasswordService interface
type passwordService struct{}

// creates a new instance of passwordService
func NewPasswordService() domain.PasswordService {
	return &passwordService{}
}

// hashes a password using bcrypt
func (pswserv *passwordService) HashPassword(password string) (string, error) {
	
	// generate a bcrypt hash from the password
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	
	return string(bytes), err
}

// checks the plain text password against the hashed password
func (pswserv *passwordService) CheckPassword(hashed, plain string) bool {
	
	// compare the hashed password with the plain password
	err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(plain))
	
	return err == nil
}