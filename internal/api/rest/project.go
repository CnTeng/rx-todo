package rest

import (
	"strconv"

	"github.com/CnTeng/rx-todo/internal/model"
	"github.com/gofiber/fiber/v2"
)

func (h *handler) registerProjectRoutes() {
	group := h.Group("/projects")

	group.Post("", h.createProject)
	group.Get(":id", h.getProject)
	group.Get("", h.getProjects)
	group.Put("reorder", h.reorderProject)
	group.Put(":id", h.updateProject)
	group.Put(":id/archive", h.archiveProject)
	group.Put(":id/unarchive", h.unarchiveProject)
	group.Delete(":id", h.deleteProject)
}

func (h *handler) createProject(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int64)

	req := &model.ProjectCreationRequest{}
	if err := h.parse(c, req); err != nil {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"parser": err.Error()})
	}

	project := &model.Project{UserID: userID}
	req.Patch(project)

	project, err := h.CreateProject(project)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).
			JSON(fiber.Map{"db": err.Error()})
	}

	return c.JSON(project)
}

func (h *handler) getProject(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int64)

	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"error": err.Error()})
	}

	project, err := h.GetProjectByID(userID, id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).
			JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(project)
}

func (h *handler) getProjects(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int64)

	projects, err := h.GetProjects(userID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).
			JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(projects)
}

func (h *handler) updateProject(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int64)

	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"error": err.Error()})
	}

	req := &model.ProjectUpdateRequest{}
	if err := h.parse(c, req); err != nil {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"error": err.Error()})
	}

	project, err := h.GetProjectByID(userID, id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).
			JSON(fiber.Map{"error": err.Error()})
	} else {
		req.Patch(project)
	}

	project, err = h.UpdateProject(project)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).
			JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(project)
}

func (h *handler) reorderProject(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int64)

	req := &model.ProjectReorderRequest{}
	if err := h.parse(c, req); err != nil {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"error": "error"})
	}

	projects := []*model.Project{}
	for _, child := range req.Children {
		project, err := h.GetProjectByID(userID, child.ID)
		if err != nil {
			return c.Status(fiber.StatusNotFound).
				JSON(fiber.Map{"error": err.Error()})
		} else {
			child.Patch(project)
		}

		projects = append(projects, project)
	}

	if err := h.UpdateProjects(projects); err != nil {
		return c.Status(fiber.StatusInternalServerError).
			JSON(fiber.Map{"error": err.Error()})
	}

	return c.SendStatus(fiber.StatusOK)
}

func (h *handler) archiveProject(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int64)

	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"error": err.Error()})
	}

	project, err := h.GetProjectByID(userID, id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).
			JSON(fiber.Map{"error": err.Error()})
	}

	if err := h.ArchiveProject(project); err != nil {
		return c.Status(fiber.StatusInternalServerError).
			JSON(fiber.Map{"error": err.Error()})
	}

	return c.SendStatus(fiber.StatusOK)
}

func (h *handler) unarchiveProject(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int64)

	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"error": err.Error()})
	}

	project, err := h.GetProjectByID(userID, id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).
			JSON(fiber.Map{"error": err.Error()})
	}

	if err := h.UnarchiveProject(project); err != nil {
		return c.Status(fiber.StatusInternalServerError).
			JSON(fiber.Map{"error": err.Error()})
	}

	return c.SendStatus(fiber.StatusOK)
}

func (h *handler) deleteProject(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int64)

	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"error": err.Error()})
	}

	project, err := h.GetProjectByID(userID, id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).
			JSON(fiber.Map{"error": err.Error()})
	}

	if err := h.DeleteProject(project); err != nil {
		return c.Status(fiber.StatusInternalServerError).
			JSON(fiber.Map{"error": err.Error()})
	}

	return c.SendStatus(fiber.StatusOK)
}
