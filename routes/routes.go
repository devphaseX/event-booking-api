package routes

import (
	"github.com/devphaseX/event-booking-api/middlewares"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(s *gin.Engine) {
	//user

	s.POST("/sign-up", signUp)
	s.POST("/sign-in", SignIn)

	authenicated := s.Group("/")
	authenicated.Use(middlewares.Authenicate)
	authenicated.GET("/events", getEvents)
	authenicated.GET("/events/:id", getEvent)
	authenicated.POST("/events", createEvents)
	authenicated.PUT("/events/:id", updateEvent)
	authenicated.DELETE("/events/:id", deleteEvent)

	//ticket
	authenicated.POST("/events/:id/register", eventRegisterUser)
	authenicated.DELETE("/events/:id/unregister", removeRegisterEventUser)

}
