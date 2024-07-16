package http

import (
	"github.com/CnTeng/rx-todo/internal/api/rest"
	"github.com/CnTeng/rx-todo/internal/database"
	"github.com/CnTeng/rx-todo/internal/http/middleware"
	"github.com/gin-gonic/gin"
)

type Server struct {
	*gin.Engine
	*database.DB
}

func NewServer(db *database.DB) *Server {
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())
	r.Use(middleware.AuthMiddleware(db))

	rest.Serve(db, r)

	return &Server{r, db}
}

func (s *Server) Start() error {
	return s.Run(":8080")
}
