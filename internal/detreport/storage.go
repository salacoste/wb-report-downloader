package detreport

import "context"

type Repository interface {
	Create(ctx context.Context, report *DetailedReport) error
}
