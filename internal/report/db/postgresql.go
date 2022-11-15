package db

import (
	"context"
	"fmt"
	"strings"
	"wb-report-downloader/internal/report"
	"wb-report-downloader/pkg/client/postgresql"
)

type repository struct {
	client postgresql.Client
}

func (r* repository) Save(ctx context.Context, sellerID uint64, report *report.Report) error {
	q := `
		INSERT INTO wb_reports (id, date_from, date_to, "createDate", "totalSale", for_pay, "bankPaymentSum", "deliveryRub",
							"paidStorageSum", "additionalPayment", "paidWithholdingSum", "paidAcceptanceSum", penalty,
							seller_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
	`

	_, err := r.client.Exec(ctx, q, 
		report.Id,
		report.DateFrom,
		report.DateTo,
		report.CreateDate,
		report.TotalSale,
		report.ForPay,
		report.BankPaymentSum,
		report.DeliveryRub,
		report.PaidStorageSum,
		report.AdditionalPayment,
		report.PaidWithholdingSum,
		report.PaidAcceptanceSum,
		report.Penalty,
		sellerID)

	return err
}

func (r* repository) FindAll(ctx context.Context, reports_ids []uint64) (found []uint64, err error) {
	q := `SELECT id
		FROM wb_reports
		WHERE id IN (%s)
	`
	idsStr := strings.Trim(strings.Replace(fmt.Sprint(reports_ids), " ", ",", -1), "[]")
	q = fmt.Sprintf(q, idsStr)

	rows, err := r.client.Query(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var reportID uint64
		rows.Scan(&reportID)
		found = append(found, reportID)
	}
	return found, nil
}

func NewRepository(client postgresql.Client) report.Repository {
	return &repository{
		client: client,
	}
}
