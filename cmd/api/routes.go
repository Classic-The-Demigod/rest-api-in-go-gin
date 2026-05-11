package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func (app *application) routes() http.Handler {
	g := gin.Default()
	v1 := g.Group("/api/v1")

	{

		v1.GET("/events", app.getAllEvents)
		v1.GET("/events/:id", app.getEventByID)

		//attendee routes

		v1.GET("/events/:id/attendees", app.getAttendeesForEvent)

		v1.GET("/attendees/:id/events", app.getEventsByAttendee)

		//user routes
		v1.POST("/auth/register", app.createUser)
		v1.POST("/auth/login", app.loginUser)
		// v1.GET("/auth/users/:id", app.getUserByID)
		// v1.PUT("/auth/users/:id", app.updateUser)
		// v1.DELETE("/auth/users/:id", app.deleteUser)

	}

	authGroup := v1.Group("/")
	authGroup.Use(app.AuthMiddleWare())
	{
		authGroup.POST("/events", app.createEvent)
		authGroup.PUT("/events/:id", app.updateEvent)
		authGroup.DELETE("/events/:id", app.deleteEvent)
		authGroup.POST("/events/:id/attendees/:userId", app.addAttendeeToEvent)
		authGroup.DELETE("/events/:id/attendees/:userId", app.deleteAttendeeFromEvent)
	}

	g.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return g

}
