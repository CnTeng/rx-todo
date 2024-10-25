package api

import (
	"strconv"

	"github.com/CnTeng/rx-todo/internal/model"
	"github.com/gofiber/fiber/v2"
)

func (h *handler) registerReminderRoutes() {
	group := h.Group("/reminders")

	group.Post("", h.createReminder)
	group.Get(":id", h.getReminder)
	group.Get("", h.getReminders)
	group.Put(":id", h.updateReminder)
	group.Delete(":id", h.deleteReminder)
}

func (h *handler) createReminder(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int64)

	req := &model.ReminderCreationRequest{}
	if err := h.parse(c, req); err != nil {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"error": err.Error()})
	}

	reminder := &model.Reminder{UserID: userID}
	req.Patch(reminder)

	reminder, err := h.CreateReminder(reminder)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).
			JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(reminder)
}

func (h *handler) getReminder(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int64)

	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"error": err.Error()})
	}

	reminder, err := h.GetReminderByID(userID, id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).
			JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(reminder)
}

func (h *handler) getReminders(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int64)

	reminders, err := h.GetReminders(userID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).
			JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(reminders)
}

func (h *handler) updateReminder(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int64)

	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"error": err.Error()})
	}

	req := &model.ReminderUpdateRequest{}
	if err := h.parse(c, req); err != nil {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"error": err.Error()})
	}

	reminder, err := h.GetReminderByID(userID, id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).
			JSON(fiber.Map{"error": err.Error()})
	} else {
		req.Patch(reminder)
	}

	reminder, err = h.UpdateReminder(reminder)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).
			JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(reminder)
}

func (h *handler) deleteReminder(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int64)

	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"error": err.Error()})
	}

	reminder, err := h.GetReminderByID(userID, id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).
			JSON(fiber.Map{"error": err.Error()})
	}

	if err := h.DeleteReminder(reminder); err != nil {
		return c.Status(fiber.StatusInternalServerError).
			JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(reminder)
}
