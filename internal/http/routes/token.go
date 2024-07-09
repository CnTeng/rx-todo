package routes

import (
	"net/http"

	"github.com/CnTeng/rx-todo/internal/database"
	"github.com/CnTeng/rx-todo/internal/model"
	"github.com/gin-gonic/gin"
)

type TokenRoutes struct {
	db *database.DB
}

func NewTokenRoutes(db *database.DB) *TokenRoutes {
	return &TokenRoutes{db}
}

func (tr *TokenRoutes) RegisterRoutes(group *gin.RouterGroup) {
	group.POST("add", tr.addToken)
}

func (tr *TokenRoutes) addToken(c *gin.Context) {
	var r model.TokenAddRequest

	if err := c.ShouldBindJSON(&r); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := tr.db.AddToken(&r); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}
