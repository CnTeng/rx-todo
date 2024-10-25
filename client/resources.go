package client

import "github.com/CnTeng/rx-todo/internal/model"

type resources struct {
	*model.Resources
}

func newResources() *resources {
	return &resources{model.NewResources()}
}

func (r *resources) GetLabel(id int64) *model.Label {
	return r.Labels[id]
}

func (r *resources) GetLabels() LabelSlice {
	labels := make(LabelSlice, 0, len(r.Labels))

	for _, label := range r.Labels {
		labels = append(labels, label)
	}

	return labels
}

func (r *resources) GetProject(id int64) *model.Project {
	return r.Projects[id]
}

func (r *resources) GetProjects() ProjectSlice {
	projects := make(ProjectSlice, 0, len(r.Projects))

	for _, project := range r.Projects {
		projects = append(projects, project)
	}

	return projects
}

func (r *resources) GetProjectProgress() {
	progressMap := make(map[int64]model.Progress, len(r.Projects))

	for _, t := range r.Tasks {
		if t.ProjectID == nil {
			continue
		}

		progress := progressMap[*t.ProjectID]
		progress.Total++

		if t.Done {
			progress.Done++
		}

		progressMap[*t.ProjectID] = progress
	}

	for id, progress := range progressMap {
		r.Projects[id].Progress = &progress
	}
}

// TODO:
func (r *resources) GetReminders() []*model.Reminder {
	reminders := make([]*model.Reminder, 0, len(r.Reminders))

	for _, reminder := range r.Reminders {
		reminders = append(reminders, reminder)
	}

	return reminders
}

func (r *resources) GetTask(id int64) *model.Task {
	return r.Tasks[id]
}

func (r *resources) GetTasks() TaskSlice {
	tasks := make(TaskSlice, 0, len(r.Tasks))

	for _, task := range r.Tasks {
		tasks = append(tasks, task)
	}

	return tasks
}

func (r *resources) GetTaskProgress() {
	progressMap := make(map[int64]model.Progress)

	for _, t := range r.Tasks {
		if t.ParentID == nil {
			continue
		}

		progress := progressMap[*t.ParentID]
		progress.Total++

		if t.Done {
			progress.Done++
		}

		progressMap[*t.ParentID] = progress
	}

	for id, progress := range progressMap {
		r.Tasks[id].Progress = &progress
	}
}
