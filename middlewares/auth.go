package middlewares

import (
	"net/http"
	"strings"

	"github.com/devphaseX/event-booking-api/utils"
	"github.com/gin-gonic/gin"
)

func Authenicate(ctx *gin.Context) {
	token := strings.TrimSpace(ctx.Request.Header.Get("Authorization"))

	if token == "" {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "user not authenicated"})
		return
	}

	userId, err := utils.VerifyToken(token)

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "user not authenicated"})
		return
	}

	ctx.Set("userId", userId)
	ctx.Next()
}
