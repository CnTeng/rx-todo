package rest

import (
	"strconv"

	"github.com/CnTeng/rx-todo/model"
	"github.com/gofiber/fiber/v2"
)

func (h *handler) registerTaskRoutes() {
	group := h.Group("/tasks")

	group.Post("", h.createTask)
	group.Get(":id", h.getTask)
	group.Get("", h.getTasks)
	group.Put(":id", h.updateTask)
	group.Put(":id/open", h.openTask)
	group.Put(":id/close", h.closeTask)
	group.Put(":id/archive", h.archiveTask)
	group.Put(":id/unarchive", h.unarchiveTask)
	group.Delete(":id", h.deleteTask)
}

func (h *handler) createTask(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int64)

	req := &model.TaskCreationRequest{}
	if err := h.parse(c, req); err != nil {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"error": err.Error()})
	}

	inboxID, err := h.GetUserInboxID(userID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).
			JSON(fiber.Map{"error": err.Error()})
	}

	task := &model.Task{UserID: userID, ProjectID: &inboxID}
	req.Patch(task)

	task, err = h.CreateTask(task)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).
			JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(task)
}

func (h *handler) getTask(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int64)

	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"error": err.Error()})
	}

	task, err := h.GetTaskByID(userID, id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).
			JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(task)
}

func (h *handler) getTasks(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int64)

	tasks, err := h.GetTasks(userID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).
			JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(tasks)
}

func (h *handler) updateTask(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int64)

	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"error": err.Error()})
	}

	req := &model.TaskUpdateRequest{}
	if err := h.parse(c, req); err != nil {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"error": err.Error()})
	}

	task, err := h.GetTaskByID(userID, id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).
			JSON(fiber.Map{"error": err.Error()})
	} else {
		req.Patch(task)
	}

	task, err = h.UpdateTask(task)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).
			JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(task)
}

func (h *handler) openTask(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int64)

	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"error": err.Error()})
	}

	task, err := h.GetTaskByID(userID, id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).
			JSON(fiber.Map{"error": err.Error()})
	}

	if !task.Done {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"error": "task already open"})
	} else {
		task.Done = false
	}

	task, err = h.UpdateTaskDoneStatus(task)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).
			JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(task)
}

func (h *handler) closeTask(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int64)

	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"error": err.Error()})
	}

	task, err := h.GetTaskByID(userID, id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).
			JSON(fiber.Map{"error": err.Error()})
	}

	if task.Done {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"error": "task already done"})
	} else {
		task.Done = true
	}

	task, err = h.UpdateTaskDoneStatus(task)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).
			JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(task)
}

func (h *handler) archiveTask(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int64)

	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"error": err.Error()})
	}

	task, err := h.GetTaskByID(userID, id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).
			JSON(fiber.Map{"error": err.Error()})
	}

	if task.Archived {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"error": "task already archived"})
	} else {
		task.Archived = true
	}

	task, err = h.UpdateTaskArchivedStatus(task)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).
			JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(task)
}

func (h *handler) unarchiveTask(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int64)

	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"error": err.Error()})
	}

	task, err := h.GetTaskByID(userID, id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).
			JSON(fiber.Map{"error": err.Error()})
	}

	if !task.Archived {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"error": "task already unarchived"})
	} else {
		task.Archived = false
	}

	task, err = h.UpdateTaskArchivedStatus(task)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).
			JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(task)
}

func (h *handler) deleteTask(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int64)

	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"error": err.Error()})
	}

	task, err := h.GetTaskByID(userID, id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).
			JSON(fiber.Map{"error": err.Error()})
	}

	if err := h.DeleteTask(task); err != nil {
		return c.Status(fiber.StatusInternalServerError).
			JSON(fiber.Map{"error": err.Error()})
	}

	return c.SendStatus(fiber.StatusNoContent)
}
