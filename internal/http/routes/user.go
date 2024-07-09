package routes

import (
	"net/http"

	"github.com/CnTeng/rx-todo/internal/database"
	"github.com/CnTeng/rx-todo/internal/model"
	"github.com/gin-gonic/gin"
)

type UserRoutes struct {
	db database.DB
}

func NewUserRoutes(db database.DB) *UserRoutes {
	return &UserRoutes{db}
}

func (ur *UserRoutes) RegisterRoutes(group *gin.RouterGroup) {
	group.POST("add", ur.addUser)
}

func (ur *UserRoutes) addUser(c *gin.Context) {
	var request model.User

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := ur.db.AddUser(&request); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}
