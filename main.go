package main

import (
	"github.com/devphaseX/event-booking-api/db"
	"github.com/devphaseX/event-booking-api/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	db.InitDB()
	server := gin.Default()
	routes.RegisterRoutes(server)
	server.Run(":8080")
}
