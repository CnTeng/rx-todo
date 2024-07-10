package rest

import (
	"net/http"
	"strconv"

	"github.com/CnTeng/rx-todo/internal/model"
	"github.com/gin-gonic/gin"
)

func (h *handler) createLabel(c *gin.Context) {
	user := c.GetInt64("user")

	label := &model.Label{}
	if err := c.BindJSON(label); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	label, err := h.CreateLabel(user, label)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusCreated, label)
}

func (h *handler) getLabel(c *gin.Context) {
	user := c.GetInt64("user")
	idStr := c.Param("id")

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	label, err := h.GetLabelByID(user, id)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, label)
}

func (h *handler) getLabels(c *gin.Context) {
	user := c.GetInt64("user")

	labels, err := h.GetLabels(user)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, labels)
}

func (h *handler) updateLabel(c *gin.Context) {
	user := c.GetInt64("user")
	idStr := c.Param("id")

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	label, err := h.GetLabelByID(user, id)
	if err != nil {
		c.Status(http.StatusNotFound)
		return
	}

	r := new(model.UpdateLabelRequest)
	if err := c.BindJSON(r); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	r.Patch(label)

	label, err = h.UpdateLabel(label)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, label)
}

func (h *handler) deleteLabel(c *gin.Context) {
	user := c.GetInt64("user")
	idStr := c.Param("id")

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	if err := h.DeleteLabel(user, id); err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
}
