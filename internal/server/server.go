package server

import (
	"github.com/CnTeng/rx-todo/internal/api/v1"
	"github.com/CnTeng/rx-todo/internal/database"
	"github.com/CnTeng/rx-todo/internal/server/middleware"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

type Server struct {
	*fiber.App
	*database.DB
}

func NewServer(db *database.DB) *Server {
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.Status(fiber.StatusBadRequest).
				JSON(fiber.Map{"error": err.Error()})
		},
	})

	app.Use(logger.New(), recover.New())
	app.Use(middleware.AuthMiddleware(db))

	v1 := app.Group("/v1")
	api.RegisterV1Routes(db, v1)

	return &Server{app, db}
}

func (s *Server) Start() error {
	return s.Listen(":8080")
}
