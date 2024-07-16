package rest

import (
	"net/http"
	"strconv"

	"github.com/CnTeng/rx-todo/internal/model"
	"github.com/gin-gonic/gin"
)

func (h *handler) registerTokenRoutes() {
	group := h.Engine.Group("/tokens")

	h.Engine.POST("/token", h.createToken)
	group.GET("", h.getTokens)
	group.PUT(":id", h.updateToken)
	group.DELETE(":id", h.deleteToken)
}

func (h *handler) createToken(c *gin.Context) {
	r := &model.CreateTokenRequest{}
	token := &model.Token{}
	if err := c.BindJSON(r); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	r.Patch(token)

	if err := h.VerifyUser(token.UserID, r.Password); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	token, err := h.CreateToken(token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusCreated, token)
}

func (h *handler) getTokens(c *gin.Context) {
	userID := c.GetInt64("user_id")

	tokens, err := h.GetTokens(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, tokens)
}

func (h *handler) updateToken(c *gin.Context) {
	userID := c.GetInt64("user_id")
	idStr := c.Param("id")

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	token, err := h.GetTokenByID(userID, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	r := &model.UpdateTokenRequest{}
	if err := c.BindJSON(r); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	r.Patch(token)

	if err := h.VerifyUser(userID, r.Password); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	token, err = h.UpdateToken(token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, token)
}

func (h *handler) deleteToken(c *gin.Context) {
	userID := c.GetInt64("user_id")
	idStr := c.Param("id")

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	if err := h.DeleteToken(userID, id); err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.Status(http.StatusNoContent)
}
