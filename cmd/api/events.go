package main

import (
	"fmt"
	"net/http"
	"rest-api-in-go/internal/database"
	"strconv"

	"github.com/gin-gonic/gin"
)

// @Summary Create an event
// @Description Create a new event
// @Tags Events
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body database.Event true "Event Request"
// @Success 201 {object} database.Event
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/events [post]
func (app *application) createEvent(c *gin.Context) {
	var event database.Event

	if err := c.ShouldBindJSON(&event); err != nil {
		fmt.Println("Error binding JSON:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := app.GetUserContext(c)
	event.OwnerId = user.Id

	err := app.models.Events.Create(&event)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create event"})
		return
	}

	c.JSON(http.StatusCreated, event)
}

// @Summary Get event by ID
// @Description Returns a single event by ID
// @Tags Events
// @Produce json
// @Param id path int true "Event ID"
// @Success 200 {object} database.Event
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /api/v1/events/{id} [get]
func (app *application) getEventByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID"})
		return
	}

	event, err := app.models.Events.GetByID(id)

	if event == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Event not found"})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve event"})
		return
	}
	c.JSON(http.StatusOK, event)
}

// @Summary Get all events
// @Description Returns all events
// @Tags Events
// @Produce json
// @Success 200 {array} database.Event
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/events [get]
func (app *application) getAllEvents(c *gin.Context) {
	events, err := app.models.Events.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve events"})
		return
	}

	c.JSON(http.StatusOK, events)
}

// @Summary Update an event
// @Description Update an existing event by ID
// @Tags Events
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Event ID"
// @Param request body database.Event true "Event Request"
// @Success 200 {object} database.Event
// @Failure 400 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /api/v1/events/{id} [put]
func (app *application) updateEvent(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID"})
		return
	}

	user := app.GetUserContext(c)

	existingEvent, err := app.models.Events.GetByID(id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve event"})
		return
	}

	if existingEvent == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Event not found"})
		return
	}

	if existingEvent.OwnerId != user.Id {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not authorized to update this event"})
		return
	}

	var updatedEvent database.Event

	if err := c.ShouldBindJSON(&updatedEvent); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	updatedEvent.Id = id
	err = app.models.Events.Update(&updatedEvent)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update event"})
		return
	}
	c.JSON(http.StatusOK, updatedEvent)
}

// @Summary Delete an event
// @Description Delete an event by ID
// @Tags Events
// @Produce json
// @Security BearerAuth
// @Param id path int true "Event ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /api/v1/events/{id} [delete]
func (app *application) deleteEvent(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID"})
		return
	}

	user := app.GetUserContext(c)
	existingEvent, err := app.models.Events.GetByID(id)

	if err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve event"})
		return
	}

	if existingEvent == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Event not found"})
		return
	}
	if existingEvent.OwnerId != user.Id {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not authorized to delete this event"})
		return
	}

	err = app.models.Events.Delete(id)
	if err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete event"})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}

// @Summary Add attendee to event
// @Description Add a user as an attendee to an event
// @Tags Attendees
// @Produce json
// @Security BearerAuth
// @Param id path int true "Event ID"
// @Param userId path int true "User ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 409 {object} map[string]interface{}
// @Router /api/v1/events/{id}/attendees/{userId} [post]
func (app *application) addAttendeeToEvent(c *gin.Context) {
	eventId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID"})
		return
	}
	userId, err := strconv.Atoi(c.Param("userId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}
	// Implementation for adding attendee to event

	event, err := app.models.Events.GetByID(eventId)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve event"})
		return
	}

	if event == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Event not found"})
		return
	}

	userToAdd, err := app.models.Users.GetUserByID(userId)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user"})
		return
	}

	if userToAdd == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	user := app.GetUserContext(c)

	if event.OwnerId != user.Id {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not authorized to add an attendee to this event"})
		return
	}

	existingAttendees, err := app.models.Attendees.GetByEventAndAttendee(eventId, userId)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve attendee"})
		return
	}

	if existingAttendees != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Attendee is already registered for this event"})
		return
	}

	attendee := database.Attendee{
		EventId: eventId,
		UserId:  userToAdd.Id,
	}

	_, err = app.models.Attendees.AddAttendee(&attendee)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add attendee to event"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Attendee added to event successfully"})
}

// @Summary Get attendees for an event
// @Description Returns all attendees for a specific event
// @Tags Attendees
// @Produce json
// @Param id path int true "Event ID"
// @Success 200 {array} database.Attendee
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/events/{id}/attendees [get]
func (app *application) getAttendeesForEvent(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID"})
		return
	}

	attendees, err := app.models.Attendees.GetAttendeesByEvent(id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve attendees for event"})
		return
	}

	c.JSON(http.StatusOK, attendees)
}

// @Summary Remove attendee from event
// @Description Remove a user from an event's attendee list
// @Tags Attendees
// @Produce json
// @Security BearerAuth
// @Param id path int true "Event ID"
// @Param userId path int true "User ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /api/v1/events/{id}/attendees/{userId} [delete]
func (app *application) deleteAttendeeFromEvent(c *gin.Context) {
	eventId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID"})
		return
	}
	userId, err := strconv.Atoi(c.Param("userId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	event, err := app.models.Events.GetByID(eventId)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong"})
		return
	}

	if event == nil {
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Event not found"})
		}
	}

	user := app.GetUserContext(c)

	if event.OwnerId != user.Id {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "You are not authorized to delete an attendee from this event"})
	}

	// Implementation for deleting attendee from event
	err = app.models.Attendees.DeleteAttendeeFromEvent(eventId, userId)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete attendee from event"})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

// @Summary Get events for an attendee
// @Description Returns all events a user is attending
// @Tags Attendees
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {array} database.Event
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/attendees/{id}/events [get]
func (app *application) getEventsByAttendee(c *gin.Context) {

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid attendee ID"})
		return
	}

	events, err := app.models.Attendees.GetEventsByAttendee(id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve events for attendee"})
		return
	}

	c.JSON(http.StatusOK, events)

}
