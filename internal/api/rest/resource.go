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

	var err error

	var labels []*model.Label
	if req.LabelSyncedAt == nil {
		labels, err = h.GetLabels(userID)
	} else {
		labels, err = h.GetLabelsByUpdateTime(userID, req.LabelSyncedAt)
	}
	if err != nil {
		return c.Status(fiber.StatusNotFound).
			JSON(fiber.Map{"db": err.Error()})
	}

	var projects []*model.Project
	if req.ProjectSyncedAt == nil {
		projects, err = h.GetProjects(userID)
	} else {
		projects, err = h.GetProjectsByUpdateTime(userID, req.ProjectSyncedAt)
	}
	if err != nil {
		return c.Status(fiber.StatusNotFound).
			JSON(fiber.Map{"error": err.Error()})
	}

	var reminders []*model.Reminder
	if req.ReminderSyncedAt == nil {
		reminders, err = h.GetReminders(userID)
	} else {
		reminders, err = h.GetRemindersByUpdateTime(userID, req.ReminderSyncedAt)
	}
	if err != nil {
		return c.Status(fiber.StatusNotFound).
			JSON(fiber.Map{"error": err.Error()})
	}

	var tasks []*model.Task
	if req.TaskSyncedAt == nil {
		tasks, err = h.GetTasks(userID)
	} else {
		tasks, err = h.GetTasksByUpdateTime(userID, req.TaskSyncedAt)
	}
	if err != nil {
		return c.Status(fiber.StatusNotFound).
			JSON(fiber.Map{"error": err.Error()})
	}

	var user *model.User
	if req.UserSyncedAt == nil {
		user, err = h.GetUserByID(userID)
	} else {
		user, err = h.GetUserByUpdateTime(userID, req.UserSyncedAt)
	}
	if err != nil {
		return c.Status(fiber.StatusNotFound).
			JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(model.NewResponse(labels, projects, reminders, tasks, user))
}
