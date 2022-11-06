package report

import "context"

type Repository interface {
	Create(ctx context.Context, report *ReportDetailes) error
}
