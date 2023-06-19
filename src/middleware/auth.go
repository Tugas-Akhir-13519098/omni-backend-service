package middleware

import (
	"context"
	"net/http"
	"omni-backend-service/config"
	"omni-backend-service/src/model"
	"omni-backend-service/src/util"

	"firebase.google.com/go/auth"
	"github.com/gin-gonic/gin"
)

func Authentication(authClient *auth.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		idToken, err := util.ExtractBearerToken(c.GetHeader("Authorization"))
		if err != nil {
			c.AbortWithError(http.StatusUnauthorized, model.UnAuthorizedError)
			return
		}

		if idToken == config.Get().AdminToken {
			c.Set("userID", "")
			c.Set("jwt", idToken)
			c.Next()
			return
		}

		token, err := authClient.VerifyIDToken(context.Background(), idToken)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		c.Set("userID", token.UID)
		c.Set("jwt", idToken)
		c.Next()
	}
}
