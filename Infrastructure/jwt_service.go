package infrastructure

// imports
import (
	"errors"
	"log"			
	"path/filepath"
	"runtime"
	"time"							
	"github.com/dgrijalva/jwt-go"
	"github.com/spf13/viper"
)

type JWTService struct {
	secret []byte
}

func NewJWTService() (*JWTService, error) {
	
	// intialize viper
	viper.AutomaticEnv() 
	viper.BindEnv("JWT_SECRET") 
	
	_, filename, _, _ := runtime.Caller(0)
	rootDir := filepath.Dir(filepath.Dir(filename))
	
	// configure viper
	viper.SetConfigName(".env")               // set config name
	viper.SetConfigType("env")                // set config type
	viper.AddConfigPath(".")                  // current directory
	viper.AddConfigPath(rootDir)              // project root
	
	err := viper.ReadInConfig(); 
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			log.Printf("error reading config: %v", err)
		}
	}
    
	// get from JWT_SECRET variable in .env
	secret := viper.GetString("JWT_SECRET")
	if secret == "" {
		return nil, errors.New("JWT_SECRET must be set in .env or environment variables")
	}

	return &JWTService{secret: []byte(secret)}, nil        // success 
}

func (jwtServ *JWTService) GenerateToken(userID, username, role string) (string, error) {
	
	// input validation
	if userID == "" {
		return "", errors.New("userID cannot be empty")
	}
	if username == "" {
		return "", errors.New("username cannot be empty")
	}
	if role == "" {
		return "", errors.New("role cannot be empty")
	}

	// create token with claims 
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId": userID,            // user id          
		"username": username,        // username
		"role": role,                // user role (admin/user)
		"exp": time.Now().Add(time.Hour * 24).Unix(),      // expires in 24h
	})

	// sign with secret key
	return token.SignedString(jwtServ.secret)         // success 
}

func (jwtServ *JWTService) ValidateToken(tokenStr string) (*jwt.Token, error) {
	
	// input validation
	if tokenStr == "" {
		return nil, errors.New("token cannot be empty")
	}

	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {	
		_, ok := token.Method.(*jwt.SigningMethodHMAC)    // check if token uses HMAC signing  
		if !ok {
			return nil, jwt.ErrSignatureInvalid      // block invalid signing 
		}
		return jwtServ.secret, nil     // return secret to verify signature
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	// check if token expired
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok {
		exp, ok := claims["exp"].(float64); 
		if ok {
			if time.Now().Unix() > int64(exp) {
				return nil, errors.New("Token is expired")
			}
		} else {
			return nil, errors.New("invalid expiration claim")
		}
	}

	return token, nil       // success 
} 

func (jwtServ *JWTService) GetSecret() string {
	return string(jwtServ.secret)
}
