package rest

import (
	"net/http"
	"strconv"

	"github.com/CnTeng/rx-todo/internal/model"
	"github.com/gin-gonic/gin"
)

func (h *handler) registerTaskRoutes() {
	group := h.Engine.Group("/tasks")

	group.POST("", h.createTask)
	group.GET(":id", h.getTask)
	group.GET("", h.getTasks)
	group.PUT(":id", h.updateTask)
	group.DELETE(":id", h.deleteTask)
}

func (h *handler) createTask(c *gin.Context) {
	userID := c.GetInt64("user_id")

	r := &model.CreateTaskRequest{}
	task := &model.Task{}
	if err := c.BindJSON(r); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	inboxID, err := h.GetUserInboxID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	r.Patch(task, userID, inboxID)

	task, err = h.CreateTask(userID, task)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusCreated, task)
}

func (h *handler) getTask(c *gin.Context) {
	userID := c.GetInt64("user_id")
	idStr := c.Param("id")

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	task, err := h.GetTaskByID(userID, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, task)
}

func (h *handler) getTasks(c *gin.Context) {
	userID := c.GetInt64("user_id")

	tasks, err := h.GetTasks(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, tasks)
}

func (h *handler) updateTask(c *gin.Context) {
	userID := c.GetInt64("user_id")
	idStr := c.Param("id")

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	task, err := h.GetTaskByID(userID, id)
	if err != nil {
		c.Status(http.StatusNotFound)
		return
	}

	r := &model.UpdateTaskRequest{}
	if err := c.BindJSON(r); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	r.Patch(task)

	task, err = h.UpdateTask(task)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, task)
}

func (h *handler) deleteTask(c *gin.Context) {
	userID := c.GetInt64("user_id")
	idStr := c.Param("id")

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	err = h.DeleteTask(userID, id)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusNoContent)
}
