package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (app *application) routes() http.Handler {
	g := gin.Default()
	v1 := g.Group("/api/v1")

	{
		v1.POST("/events", app.createEvent)
		v1.GET("/events", app.getAllEvents)
		v1.GET("/events/:id", app.getEventByID)
		v1.PUT("/events/:id", app.updateEvent)
		v1.DELETE("/events/:id", app.deleteEvent)

		//attendee routes
		v1.POST("/events/:id/attendees/:userId", app.addAttendeeToEvent)
		v1.GET("/events/:id/attendees", app.getAttendeesForEvent)

		//user routes
		v1.POST("/auth/register", app.createUser)
		// v1.POST("/auth/login", app.loginUser)
		// v1.GET("/auth/users/:id", app.getUserByID)
		// v1.PUT("/auth/users/:id", app.updateUser)
		// v1.DELETE("/auth/users/:id", app.deleteUser)

	}

	return g

}
