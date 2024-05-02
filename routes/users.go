package routes

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/devphaseX/event-booking-api/models"
	"github.com/devphaseX/event-booking-api/utils"
	"github.com/gin-gonic/gin"
)

func signUp(ctx *gin.Context) {

	var userFormInput models.RegisterUserInput

	err := ctx.ShouldBindJSON(&userFormInput)

	if err != nil {
		fmt.Println(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "invalid form data payload"})
		return
	}

	newUser := models.User{
		Email:    userFormInput.Email,
		Password: userFormInput.Password,
	}

	err = newUser.Save()

	if err != nil {
		fmt.Println(err)
		if strings.Contains(err.Error(), "constraint") && strings.Contains(err.Error(), "email") {
			ctx.JSON(http.StatusConflict, gin.H{"message": "User already registered with this email"})
			return

		}

		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "We encountered an error while registering your account"})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"data": newUser})
}

func SignIn(ctx *gin.Context) {
	var cred models.SignInUserInput

	err := ctx.ShouldBindJSON(&cred)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid credentials"})
		return
	}

	user, err := models.FindUserByEmail(cred.Email)

	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Account credential not a match"})
			return
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "We encountered an error while signing your account"})
		return
	}

	err = user.ComparePassword(cred)

	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Account credential not a match"})
		return
	}

	token, err := utils.GenerateJwtToken(user.ID, user.Email)

	if err != nil {
		fmt.Println(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "We encountered an error while signing your account"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "User signed in", "token": token})
}
