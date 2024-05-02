package routes

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/devphaseX/event-booking-api/models"
	"github.com/gin-gonic/gin"
)

func eventRegisterUser(ctx *gin.Context) {
	var eventRegUserTicket models.EventRegisterUserInput
	err := ctx.ShouldBindJSON(&eventRegUserTicket)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "invalid data payload"})
		return
	}

	_, err = models.FindUserById(eventRegUserTicket.RegisteredUserID)

	if err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"message": "event registered user not found"})
		return

	}

	eventId, err := strconv.ParseInt(ctx.Param("id"), 10, 64)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "invalid event id"})
		return
	}

	userId := ctx.GetInt64("userId")
	_, err = models.GetEventById(eventId, userId)

	if err != nil {
		fmt.Println(err)

		if err.Error() == "sql: no rows in result set" {
			ctx.JSON(http.StatusNotFound, gin.H{"message": "event not exist"})
			return
		}
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "invalid event id"})
		return
	}

	ticket := models.Ticket{
		UserId:  eventRegUserTicket.RegisteredUserID,
		EventId: eventId,
	}

	err = ticket.Save()

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "an error occurred while registering user to event"})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "user registered to event"})
}

func removeRegisterEventUser(ctx *gin.Context) {
	eventId, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "invalid event id"})
		return
	}

	var formInput struct {
		RegUserId int64 `json:"reg_user_id" binding:"required"`
	}

	err = ctx.ShouldBindJSON(&formInput)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "reg_user_id is invalid or missing"})
		return
	}

	userId := ctx.GetInt64("userId")
	_, err = models.GetEventById(eventId, userId)

	if err != nil {

		if err.Error() == "sql: no rows in result set" {
			ctx.JSON(http.StatusNotFound, gin.H{"message": "event not exist"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "could not unregister user from event"})
		return
	}

	ticket, err := models.GetRegUserTicket(eventId, formInput.RegUserId, userId)

	if err != nil {
		fmt.Println("first", err)

		if err.Error() == "sql: no rows in result set" {
			ctx.JSON(http.StatusNotFound, gin.H{"message": "user not registered to event"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "could not unregister user from event"})
		return
	}

	err = models.DeleteTicketById(ticket.ID)

	if err != nil {
		fmt.Println(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "could not unregister user from event"})
		return
	}

	ctx.JSON(http.StatusNoContent, gin.H{"message": "unregister user from event succesfully"})
}
