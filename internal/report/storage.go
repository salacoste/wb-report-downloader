package report

import "context"

type Repository interface {
	Save(ctx context.Context, sellerID uint64, report *Report) error
	FindAll(ctx context.Context, reports_ids []uint64) (found []uint64, err error)
}
