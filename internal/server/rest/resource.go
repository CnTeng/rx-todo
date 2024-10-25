package rest

import (
	"github.com/CnTeng/rx-todo/internal/model"
	"github.com/gofiber/fiber/v2"
)

func (h *handler) registerResourceRoutes() {
	group := h.Group("/resources")

	group.Post("/sync", h.syncResources)
}

func (h *handler) syncResources(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int64)

	req := &model.ResourceSyncRequest{}
	if err := h.parse(c, req); err != nil {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"parser": err.Error()})
	}

	labels, err := h.GetLabelsByUpdateTime(req.LabelSyncedAt, userID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).
			JSON(fiber.Map{"db": err.Error()})
	}

	projects, err := h.GetProjectsByUpdateTime(userID, req.ProjectSyncedAt)
	if err != nil {
		return c.Status(fiber.StatusNotFound).
			JSON(fiber.Map{"error": err.Error()})
	}

	reminders, err := h.GetRemindersByUpdateTime(userID, req.ReminderSyncedAt)
	if err != nil {
		return c.Status(fiber.StatusNotFound).
			JSON(fiber.Map{"error": err.Error()})
	}

	tasks, err := h.GetTasksByUpdateTime(userID, req.TaskSyncedAt)
	if err != nil {
		return c.Status(fiber.StatusNotFound).
			JSON(fiber.Map{"error": err.Error()})
	}

	user, err := h.GetUserByUpdateTime(userID, req.UserSyncedAt)
	if err != nil {
		return c.Status(fiber.StatusNotFound).
			JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(model.NewResources(labels, projects, reminders, tasks, user))
}
