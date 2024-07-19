package rest

import (
	"net/http"
	"strconv"

	"github.com/CnTeng/rx-todo/internal/model"
	"github.com/gin-gonic/gin"
)

func (h *handler) registerReminderRoutes() {
	group := h.Engine.Group("/reminders")

	group.POST("", h.createReminder)
	group.GET(":id", h.getReminder)
	group.GET("", h.getReminders)
	group.PUT(":id", h.updateReminder)
	group.DELETE(":id", h.deleteReminder)
}

func (h *handler) createReminder(c *gin.Context) {
	userID := c.GetInt64("user_id")

	r := &model.CreateReminderRequest{}
	reminder := &model.Reminder{}
	if err := c.BindJSON(r); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	r.Patch(userID, reminder)

	reminder, err := h.CreateReminder(reminder)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusCreated, reminder)
}

func (h *handler) getReminders(c *gin.Context) {
	userID := c.GetInt64("user_id")

	reminders, err := h.GetReminders(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, reminders)
}

func (h *handler) getReminder(c *gin.Context) {
	userID := c.GetInt64("user_id")
	idStr := c.Param("id")

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	reminder, err := h.GetReminderByID(userID, id)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, reminder)
}

func (h *handler) updateReminder(c *gin.Context) {
	userID := c.GetInt64("user_id")
	idStr := c.Param("id")

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	reminder, err := h.GetReminderByID(userID, id)
	if err != nil {
		c.JSON(http.StatusNotFound, err.Error())
		return
	}

	r := &model.UpdateReminderRequest{}
	if err := c.BindJSON(r); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	r.Patch(userID, reminder)

	reminder, err = h.UpdateReminder(reminder)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, reminder)
}

func (h *handler) deleteReminder(c *gin.Context) {
	userID := c.GetInt64("user_id")
	idStr := c.Param("id")

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	err = h.DeleteReminder(userID, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.Status(http.StatusNoContent)
}
