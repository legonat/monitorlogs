package middleware

import (
	"github.com/gin-gonic/gin"
	"monitorlogs/internal/auth"
	"monitorlogs/pkg/erx"
	"monitorlogs/pkg/tools"
	"net/http"
)

func AuthJWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		const BEARER_SCHEMA = "Bearer "
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || len(authHeader) < len(BEARER_SCHEMA) {
			tools.LogErr(erx.NewError(601, "User is unauthorized"))
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		tokenString := authHeader[len(BEARER_SCHEMA):]
		login, err := auth.ParseToken(tokenString)
		if login != "" {
			tools.LogInfo("Token is valid")
		} else {
			tools.LogErr(err)
			c.AbortWithStatus(http.StatusUnauthorized)
		}
	}
}
