package routes

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/devphaseX/event-booking-api/models"
	"github.com/gin-gonic/gin"
)

func getEvents(ctx *gin.Context) {
	userId := ctx.GetInt64("userId")
	events, err := models.GetAllEvents(userId)

	if err != nil {

		fmt.Println(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "could not fetch events"})
		return
	}

	ctx.JSON(http.StatusOK, events)
}

func getEvent(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "invalid event id"})
		return
	}

	userId := ctx.GetInt64("userId")
	event, err := models.GetEventById(id, userId)

	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			ctx.JSON(http.StatusNotFound, gin.H{"message": "event with id not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "error occurred while getting event"})
		return
	}

	ctx.JSON(http.StatusOK, event)
}

func createEvents(ctx *gin.Context) {
	var event models.Event
	err := ctx.ShouldBindJSON(&event)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "could not parse request data", "error": err})
		return
	}

	userId := ctx.GetInt64("userId")
	event.UserID = userId

	err = event.Save()

	if err != nil {

		ctx.JSON(http.StatusBadRequest, gin.H{"message": "could not create event", "error": err})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "Event created", "event": event})
}

func updateEvent(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "invalid event id"})
		return
	}

	userId := ctx.GetInt64("userId")
	_, err = models.GetEventById(id, userId)

	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			ctx.JSON(http.StatusNotFound, gin.H{"message": "event with id not found"})
			return
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "could not fetch the event"})
		return
	}

	var updatedEvent models.Event

	err = ctx.ShouldBindJSON(&updatedEvent)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "invalid update data payload"})
		return
	}

	updatedEvent.ID = id
	err = updatedEvent.Update(userId)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "could not update event"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Event updated successfully"})
}

func deleteEvent(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "invalid event id"})
		return
	}

	userId := ctx.GetInt64("userId")
	_, err = models.GetEventById(id, userId)

	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			ctx.JSON(http.StatusNotFound, gin.H{"message": "event with id not found"})
			return
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "could not fetch the event"})
		return
	}

	err = models.DeleteEventById(id)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "could not delete event"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Event deleted successfully"})

}
