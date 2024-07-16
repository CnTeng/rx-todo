package rest

import (
	"net/http"
	"strconv"

	"github.com/CnTeng/rx-todo/internal/model"
	"github.com/gin-gonic/gin"
)

func (h *handler) registerProjectRoutes() {
	group := h.Engine.Group("/projects")

	group.POST("", h.createProject)
	group.GET(":id", h.getProject)
	group.GET("", h.getProjects)
	group.PUT(":id", h.updateProject)
	group.PUT("reorder", h.reorderProject)
	group.PUT(":id/archive", h.archiveProject)
	group.PUT(":id/unarchive", h.unarchiveProject)
	group.DELETE(":id", h.deleteProject)
}

func (h *handler) createProject(c *gin.Context) {
	userID := c.GetInt64("user_id")

	r := &model.CreateProjectRequest{}
	project := &model.Project{}
	if err := c.BindJSON(r); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	r.Patch(project)
	project.UserID = userID

	project, err := h.CreateProject(project)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusCreated, project)
}

func (h *handler) getProjects(c *gin.Context) {
	userID := c.GetInt64("user_id")

	projects, err := h.GetProjects(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, projects)
}

func (h *handler) getProject(c *gin.Context) {
	userID := c.GetInt64("user_id")
	idStr := c.Param("id")

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	project, err := h.GetProjectByID(userID, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, project)
}

func (h *handler) updateProject(c *gin.Context) {
	userID := c.GetInt64("user_id")
	idStr := c.Param("id")

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	project, err := h.GetProjectByID(userID, id)
	if err != nil {
		c.JSON(http.StatusNotFound, err.Error())
		return
	}

	r := &model.UpdateProjectRequest{}
	if err := c.BindJSON(r); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	r.Patch(project)

	project, err = h.UpdateProject(project)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, project)
}

func (h *handler) reorderProject(c *gin.Context) {
	userID := c.GetInt64("user_id")
	projects := []*model.Project{}

	r := &model.ReorderProjectRequest{}
	if err := c.BindJSON(r); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	for _, child := range r.Children {
		project, err := h.GetProjectByID(userID, child.ID)
		if err != nil || project.UserID != userID {
			c.JSON(http.StatusNotFound, err.Error())
			return
		}
		child.Patch(project)

		projects = append(projects, project)
	}

	if err := h.UpdateProjects(projects); err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.Status(http.StatusOK)
}

func (h *handler) archiveProject(c *gin.Context) {
	userID := c.GetInt64("user_id")
	idStr := c.Param("id")

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	_, err = h.GetProjectByID(userID, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	err = h.ArchiveProject(userID, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.Status(http.StatusOK)
}

func (h *handler) unarchiveProject(c *gin.Context) {
	userID := c.GetInt64("user_id")
	idStr := c.Param("id")

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	_, err = h.GetProjectByID(userID, id)
	if err != nil {
		c.JSON(http.StatusNotFound, err.Error())
		return
	}

	err = h.UnarchiveProject(userID, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.Status(http.StatusOK)
}

func (h *handler) deleteProject(c *gin.Context) {
	userID := c.GetInt64("user_id")
	idStr := c.Param("id")

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	_, err = h.GetProjectByID(userID, id)
	if err != nil {
		c.JSON(http.StatusNotFound, err.Error())
		return
	}

	err = h.DeleteProject(userID, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.Status(http.StatusOK)
}
