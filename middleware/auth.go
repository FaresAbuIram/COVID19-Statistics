package middleware

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(context *gin.Context) {
		// Get the JWT token from the request header
		tokenString := context.Request.Header.Get("Authorization")
		if tokenString == "" {
			context.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization token not provided"})
			context.Abort()
			return
		}

		// Parse and validate the token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Return the secret key used to sign the token
			return []byte("secreatetoken"), nil
		})
		if err != nil {
			context.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization token"})
			context.Abort()
			return
		}

		// Check if the token is valid and not expired
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			// Get the user ID from the token
			userID, err := strconv.Atoi(fmt.Sprintf("%.0f", claims["user_id"]))
			if err != nil {
				context.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID in authorization token"})
				context.Abort()
				return
			}

			// // Store the token in a cookie
			// cookie := &http.Cookie{
			// 	Name:     "access_token",
			// 	Value:    tokenString,
			// 	Path:     "/",
			// 	HttpOnly: true,
			// 	MaxAge:   3600, // Set the cookie to expire in 1 hour
			// 	Secure:   true, // Set the cookie to be secure (HTTPS only)
			// }
			// http.SetCookie(context.Writer, cookie)

			// Set the user ID in the request context
			context.Set("user_id", userID)
			
			// Call the next middleware/handler in the chain
			context.Next()
		} else {
			context.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization token"})
			context.Abort()
			return
		}

	}
}

func GetUserID(context *gin.Context) int {
	if userID, ok := context.Get("user_id"); ok {
		return userID.(int)
	}
	return 0
}
