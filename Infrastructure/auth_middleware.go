package infrastructure

// imports
import (
	"net/http"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/natnael-eyuel-dev/Task-Management-Unit-Test/Domain"
)

type AuthMiddleWare struct {
	jwtService domain.JWTService
}

func NewAuthMiddleware(jwtServ domain.JWTService) *AuthMiddleWare {
	return &AuthMiddleWare{jwtService: jwtServ}
}

// auth handler
func (authmidlw *AuthMiddleWare) Handler() gin.HandlerFunc {
	
	return func(c *gin.Context) {

		tokenStr := c.GetHeader("Authorization")        // get token from authorization header
		// reject if empty
		if tokenStr == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "authorization header required"})
			c.Abort()
			return
		}
		
		// validate token structure/signature with error handling 
		token, err := authmidlw.jwtService.ValidateToken(tokenStr)     
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			c.Abort()
			return
		}

		// if token is valid, extract claims and store in request context
		claims, ok := token.Claims.(jwt.MapClaims)      
		if ok {
			c.Set("userID", claims["sub"])             // user id
			c.Set("username", claims["username"])      // username 
			c.Set("role", claims["role"])              // user role (admin/user)
		}

		c.Next()       // proceed to next handler
	}
}

func AdminOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		
		role, exists := c.Get("role")          // get role from context 

		// block if either role doesn't exist in context or role isn't "admin"
		if !exists || role != "admin" {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "admin access required",
			})
			
			c.Abort()
			return
		}

		c.Next()       // allow admin to proceed
	}
}
