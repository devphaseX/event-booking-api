package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/devphaseX/event-booking-api/db"
	"github.com/devphaseX/event-booking-api/models"
	"github.com/gin-gonic/gin"
)

func main() {
	db.InitDB()
	server := gin.Default()

	server.GET("/events", getEvents)
	server.GET("/events/:id", getEvent)
	server.POST("/events", createEvents)
	server.Run(":8080")
}

func getEvents(ctx *gin.Context) {
	events, err := models.GetAllEvents()

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

	event, err := models.GetEventById(id)

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

	event.ID = 1
	event.UserID = 1

	err = event.Save()

	if err != nil {

		ctx.JSON(http.StatusBadRequest, gin.H{"message": "could not create event", "error": err})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "Event created", "event": event})
}
