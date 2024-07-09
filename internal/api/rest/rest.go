package rest

import (
	"github.com/CnTeng/rx-todo/internal/database"
	"github.com/gin-gonic/gin"
)

type handler struct {
	*database.DB
}

func Serve(db *database.DB, r *gin.Engine) {
	h := handler{db}

	labelGroup := r.Group("/labels")

	labelGroup.POST("", h.createLabel)
	labelGroup.GET(":id", h.getLabel)
	labelGroup.GET("", h.getLabels)
	labelGroup.PUT(":id", h.updateLabel)
	labelGroup.DELETE(":id", h.deleteLabel)
}
