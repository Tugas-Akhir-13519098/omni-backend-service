package middleware

import (
	"net/http"
	"omni-backend-service/src/model"

	"github.com/gin-gonic/gin"
)

func ErrorHandler(c *gin.Context) {
	c.Next()

	for _, err := range c.Errors {
		switch err.Err {
		case model.DuplicateUserError:
			c.JSON(-1, gin.H{"error": err.Err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
	}

}
