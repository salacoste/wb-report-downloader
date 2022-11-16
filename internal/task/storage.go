package task

import "context"

type Repository interface {
	GetDownloadTask(ctx context.Context) (Task, error)
	GetDownloadTasks(ctx context.Context, limit uint32) ([]Task, error)
	UpdateTaskStatus(ctx context.Context, task Task) (error)
}
