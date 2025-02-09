package api

import (
	"github.com/CnTeng/rx-todo/internal/database"
	"github.com/CnTeng/rx-todo/internal/model"
	"github.com/gofiber/fiber/v2"
)

type handler struct {
	*database.DB
	fiber.Router
}

func (*handler) parse(c *fiber.Ctx, req any) error {
	if err := c.BodyParser(req); err != nil {
		return err
	}

	if err := model.Validate(req); err != nil {
		return err
	}

	return nil
}

func nilIfZero(value int64) *int64 {
	if value == 0 {
		return nil
	}
	return &value
}

func RegisterV1Routes(db *database.DB, router fiber.Router) {
	h := handler{DB: db, Router: router}

	h.registerProjectRoutes()
	h.registerLabelRoutes()
	h.registerTaskRoutes()
	h.registerUserRoutes()
	h.registerTokenRoutes()
	h.registerReminderRoutes()
	h.registerResourceRoutes()
}
