package taskdb

import (
	"context"
	"fmt"
	"wb-report-downloader/internal/task"
	"wb-report-downloader/pkg/client/postgresql"
)

type repository struct {
	client postgresql.Client
}

func (r *repository) GetDownloadTask(ctx context.Context) (task.Task, error) {
	q := `
		SELECT r.id, r.seller_id, s.name
		FROM wb_reports r
		JOIN sellers s ON s.id = r.seller_id
		LEFT JOIN wb_reports_details_v2 wrd ON r.id = wrd.report_id
		WHERE wrd.report_id IS NULL
		LIMIT 1
	`
	// r.logger.Trace(fmt.Sprintf("SQL Query: %s", formatQuery(q)))

	var t task.Task
	err := r.client.QueryRow(ctx, q).Scan(&t.ReportID, &t.SellerID, &t.SellerName)
	if err != nil {
		return task.Task{}, err
	}

	return t, nil
}

func (r *repository) GetDownloadTasks(ctx context.Context, limit uint32) (tasks []task.Task, err error) {
	q := `
		SELECT r.id, r.seller_id, s.name
		FROM wb_reports r
		JOIN sellers s ON s.id = r.seller_id
		LEFT JOIN wb_reports_details_v2 wrd ON r.id = wrd.report_id
		WHERE wrd.report_id IS NULL
		LIMIT %v
	`

	q = fmt.Sprintf(q, limit)
	// r.logger.Trace(fmt.Sprintf("SQL Query: %s", formatQuery(q)))

	rows, err := r.client.Query(ctx, q)
	if err != nil {
		return []task.Task{}, err
	}
	defer rows.Close()
	for rows.Next() {
		var t task.Task
		rows.Scan(&t.ReportID, &t.SellerID, &t.SellerName)
		tasks = append(tasks, t)
	}
	return tasks, nil
}

func NewRepository(client postgresql.Client) task.Repository {
	return &repository{
		client: client,
	}
}
