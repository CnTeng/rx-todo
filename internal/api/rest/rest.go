package rest

import (
	"github.com/CnTeng/rx-todo/internal/database"
	"github.com/gin-gonic/gin"
)

type handler struct {
	*database.DB
	Engine *gin.Engine
}

func Serve(db *database.DB, r *gin.Engine) {
	h := handler{DB: db, Engine: r}

	h.registerProjectRoutes()
	h.registerLabelRoutes()
	h.registerTaskRoutes()
	h.registerUserRoutes()
	h.registerTokenRoutes()
	h.registerReminderRoutes()
}
