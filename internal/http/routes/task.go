package routes

import (
	"net/http"

	"github.com/CnTeng/rx-todo/internal/database"
	"github.com/CnTeng/rx-todo/internal/model"
	"github.com/gin-gonic/gin"
)

type TaskRoutes struct {
	db database.DB
}

func NewTaskRoutes(db database.DB) *TaskRoutes {
	return &TaskRoutes{db}
}

func (tr *TaskRoutes) ResgisterRoutes(group *gin.RouterGroup) {
	group.POST("add", tr.addTask)
}

func (tr *TaskRoutes) addTask(c *gin.Context) {
	var r model.TaskAddRequest
	user := c.GetInt64("user")

	if err := c.ShouldBindJSON(&r); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := tr.db.AddTask(user, &r); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}
