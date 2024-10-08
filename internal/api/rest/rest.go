package rest

import (
	"github.com/CnTeng/rx-todo/internal/database"
	"github.com/CnTeng/rx-todo/internal/model"
	"github.com/gofiber/fiber/v2"
)

type handler struct {
	*database.DB
	*fiber.App
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

func Serve(db *database.DB, app *fiber.App) {
	h := handler{DB: db, App: app}

	h.registerProjectRoutes()
	h.registerLabelRoutes()
	h.registerTaskRoutes()
	h.registerUserRoutes()
	h.registerTokenRoutes()
	h.registerReminderRoutes()
}
