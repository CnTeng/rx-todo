package rest

import (
	"net/http"
	"strconv"

	"github.com/CnTeng/rx-todo/internal/model"
	"github.com/gin-gonic/gin"
)

func (h *handler) registerLabelRoutes() {
	group := h.Engine.Group("/labels")

	group.POST("", h.createLabel)
	group.GET(":id", h.getLabel)
	group.GET("", h.getLabels)
	group.PUT(":id", h.updateLabel)
	group.DELETE(":id", h.deleteLabel)
}

func (h *handler) createLabel(c *gin.Context) {
	userID := c.GetInt64("user_id")

	req := &model.CreateLabelRequest{}
	label := &model.Label{}
	if err := c.BindJSON(req); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	req.Patch(userID, label)

	label, err := h.CreateLabel(label)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusCreated, label)
}

func (h *handler) getLabel(c *gin.Context) {
	userID := c.GetInt64("user_id")
	idStr := c.Param("id")

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	label, err := h.GetLabelByID(userID, id)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, label)
}

func (h *handler) getLabels(c *gin.Context) {
	userID := c.GetInt64("user_id")

	labels, err := h.GetLabels(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, labels)
}

func (h *handler) updateLabel(c *gin.Context) {
	userID := c.GetInt64("user_id")
	idStr := c.Param("id")

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	label, err := h.GetLabelByID(userID, id)
	if err != nil {
		c.Status(http.StatusNotFound)
		return
	}

	req := new(model.UpdateLabelRequest)
	if err := c.BindJSON(req); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	req.Patch(label)

	label, err = h.UpdateLabel(label)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, label)
}

func (h *handler) deleteLabel(c *gin.Context) {
	userID := c.GetInt64("user_id")
	idStr := c.Param("id")

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	if err := h.DeleteLabel(userID, id); err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
}
