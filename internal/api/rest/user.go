package rest

import (
	"net/http"

	"github.com/CnTeng/rx-todo/internal/model"
	"github.com/gin-gonic/gin"
)

func (h *handler) registerUserRoutes() {
	group := h.Engine.Group("/users")

	h.Engine.POST("/registry", h.createUser)
	group.GET(":id", h.getUser)
	group.PUT(":id", h.updateUser)
	group.DELETE(":id", h.deleteUser)
}

func (h *handler) createUser(c *gin.Context) {
	r := &model.CreateUserRequest{}
	user := &model.User{}
	if err := c.BindJSON(r); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	r.Patch(user)

	var err error
	user.Password, err = model.HashPassword(*r.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	user, err = h.CreateUser(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusCreated, user)
}

func (h *handler) getUser(c *gin.Context) {
	id := c.GetInt64("user_id")

	user, err := h.GetUserByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *handler) updateUser(c *gin.Context) {
	id := c.GetInt64("user_id")

	user, err := h.GetUserByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, err.Error())
		return
	}

	r := &model.UpdateUserRequest{}
	if err := c.BindJSON(r); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	newPassword, err := model.HashPassword(*r.NewPassword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	if err := h.VerifyUser(id, *r.OldPassword); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	r.Patch(user)
	user.Password = newPassword

	user, err = h.UpdateUser(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *handler) deleteUser(c *gin.Context) {
	id := c.GetInt64("user_id")

	if err := h.DeleteUser(id); err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.Status(http.StatusOK)
}
