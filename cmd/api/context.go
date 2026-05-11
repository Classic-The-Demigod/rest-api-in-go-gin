package main

import (
	"rest-api-in-go/internal/database"

	"github.com/gin-gonic/gin"
)

func (app *application) GetUserContext(c *gin.Context) *database.User {
	contextUser, exists := c.Get("user")
	if !exists {
		return &database.User{}
	}

	user, ok := contextUser.(*database.User)
	if !ok {
		return &database.User{}
	}

	return user
}
