package rest

import (
	"strconv"

	"github.com/CnTeng/rx-todo/internal/model"
	"github.com/gofiber/fiber/v2"
)

func (h *handler) registerLabelRoutes() {
	group := h.Group("/labels")

	group.Post("", h.createLabel)
	group.Get(":id", h.getLabel)
	group.Get("", h.getLabels)
	group.Put(":id", h.updateLabel)
	group.Delete(":id", h.deleteLabel)
}

func (h *handler) createLabel(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int64)

	req := &model.LabelCreationRequest{}
	if err := h.parse(c, req); err != nil {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"parser": err.Error()})
	}

	label := &model.Label{UserID: userID}
	req.Patch(label)

	label, err := h.CreateLabel(label)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).
			JSON(fiber.Map{"db": err.Error()})
	}

	return c.JSON(label)
}

func (h *handler) getLabel(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int64)

	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"parser": err.Error()})
	}

	label, err := h.GetLabelByID(userID, id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).
			JSON(fiber.Map{"db": err.Error()})
	}

	return c.JSON(label)
}

func (h *handler) getLabels(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int64)

	labels, err := h.GetLabels(userID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).
			JSON(fiber.Map{"db": err.Error()})
	}

	return c.JSON(labels)
}

func (h *handler) updateLabel(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int64)

	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"parser": err.Error()})
	}

	req := &model.LabelUpdateRequest{}
	if err := h.parse(c, req); err != nil {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"parser": err.Error()})
	}

	label, err := h.GetLabelByID(userID, id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).
			JSON(fiber.Map{"db": err.Error()})
	} else {
		req.Patch(label)
	}

	label, err = h.UpdateLabel(label)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).
			JSON(fiber.Map{"db": err.Error()})
	}

	return c.JSON(label)
}

func (h *handler) deleteLabel(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int64)

	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"parser": err.Error()})
	}

	label, err := h.GetLabelByID(userID, id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).
			JSON(fiber.Map{"db": err.Error()})
	}

	if err := h.DeleteLabel(label); err != nil {
		return c.Status(fiber.StatusInternalServerError).
			JSON(fiber.Map{"db": err.Error()})
	}

	return c.SendStatus(fiber.StatusNoContent)
}
