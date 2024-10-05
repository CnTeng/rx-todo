package rest

import (
	"github.com/CnTeng/rx-todo/model"
	"github.com/gofiber/fiber/v2"
)

func (h *handler) registerUserRoutes() {
	group := h.Group("/users")

	h.Post("/registry", h.createUser)
	group.Get(":id", h.getUser)
	group.Put(":id", h.updateUser)
	group.Delete(":id", h.deleteUser)
}

func (h *handler) createUser(c *fiber.Ctx) error {
	req := &model.CreateUserRequest{}
	if err := h.parse(c, req); err != nil {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"error": err.Error()})
	}

	hashedPassword, err := model.HashPassword(*req.Password)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).
			JSON(fiber.Map{"error": err.Error()})
	}

	user := &model.User{Password: hashedPassword}
	req.Patch(user)

	user, err = h.CreateUser(user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).
			JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(user)
}

func (h *handler) getUser(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int64)

	user, err := h.GetUserByID(userID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).
			JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(user)
}

func (h *handler) updateUser(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int64)

	req := &model.UpdateUserRequest{}
	if err := h.parse(c, req); err != nil {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"error": err.Error()})
	}

	if err := h.VerifyUser(userID, *req.OldPassword); err != nil {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"error": err.Error()})
	}

	user, err := h.GetUserByID(userID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).
			JSON(fiber.Map{"error": err.Error()})
	}

	*req.NewPassword, err = model.HashPassword(*req.NewPassword)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).
			JSON(fiber.Map{"error": err.Error()})
	} else {
		req.Patch(user)
	}

	user, err = h.UpdateUser(user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).
			JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(user)
}

func (h *handler) deleteUser(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int64)

	if _, err := h.GetUserByID(userID); err != nil {
		return c.Status(fiber.StatusNotFound).
			JSON(fiber.Map{"error": err.Error()})
	}

	if err := h.DeleteUser(userID); err != nil {
		return c.Status(fiber.StatusInternalServerError).
			JSON(fiber.Map{"error": err.Error()})
	}

	return c.SendStatus(fiber.StatusNoContent)
}
