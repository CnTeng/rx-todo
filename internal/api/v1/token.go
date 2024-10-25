package api

import (
	"strconv"

	"github.com/CnTeng/rx-todo/internal/model"
	"github.com/gofiber/fiber/v2"
)

func (h *handler) registerTokenRoutes() {
	group := h.Group("/tokens")

	h.Post("/token", h.createToken)
	group.Get("", h.getTokens)
	group.Put(":id", h.updateToken)
	group.Delete(":id", h.deleteToken)
}

func (h *handler) createToken(c *fiber.Ctx) error {
	req := &model.CreateTokenRequest{}
	if err := h.parse(c, req); err != nil {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"error": err.Error()})
	}

	if err := h.VerifyUser(req.UserID, req.Password); err != nil {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"error": err.Error()})
	}

	token := &model.Token{}
	req.Patch(token)

	token, err := h.CreateToken(token)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).
			JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(token)
}

func (h *handler) getTokens(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int64)

	tokens, err := h.GetTokens(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).
			JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(tokens)
}

func (h *handler) updateToken(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int64)

	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"error": err.Error()})
	}

	req := &model.UpdateTokenRequest{}
	if err := h.parse(c, req); err != nil {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"error": err.Error()})
	}

	if req.UserID != id {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"error": "user_id in request body must match id in URL"})
	}

	if err := h.VerifyUser(req.UserID, req.Password); err != nil {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"error": err.Error()})
	}

	token, err := h.GetTokenByID(userID, id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).
			JSON(fiber.Map{"error": err.Error()})
	} else {
		req.Patch(token)
	}

	token, err = h.UpdateToken(token)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).
			JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(token)
}

func (h *handler) deleteToken(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int64)

	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"error": err.Error()})
	}

	token, err := h.GetTokenByID(userID, id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).
			JSON(fiber.Map{"error": err.Error()})
	}

	if err := h.DeleteToken(token); err != nil {
		return c.Status(fiber.StatusInternalServerError).
			JSON(fiber.Map{"error": err.Error()})
	}

	return c.SendStatus(fiber.StatusOK)
}
