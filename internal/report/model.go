package report

import "time"

type ReportsResponse struct {
	Data struct {
		Reports []Report `json:"reports"`
	} `json:"data"`
}

type Report struct {
	DetailsCount uint32 `json:"detailsCount"`
	Type uint32 `json:"type"`
	Id uint64 `json:"id"`
	DateFrom time.Time `json:"dateFrom"`
	DateTo time.Time `json:"dateTo"`
	CreateDate time.Time `json:"createDate"`
	TotalSale float64 `json:"totalSale"`
	AvgSalePercent float64 `json:"avgSalePercent"`
	ForPay float64 `json:"forPay"`
	DeliveryRub float64 `json:"deliveryRub"`
	PaidStorageSum float64 `json:"paidStorageSum"`
	PaidAcceptanceSum float64 `json:"paidAcceptanceSum"`
	PaidWithholdingSum float64 `json:"paidWithholdingSum"`
	BankPaymentSum float64 `json:"bankPaymentSum"`
	CurrentStatusDocument struct {
		Id uint64 `json:"id"`
		Name string `json:"name"`
	} `json:"currentStatusDocument"`
	BankPaymentStatusId int32 `json:"bankPaymentStatusId"`
	BankPaymentStatusName string `json:"bankPaymentStatusName"`
	BankPaymentStatusDescription string `json:"bankPaymentStatusDescription"`
	Penalty float64 `json:"penalty"`
	AdditionalPayment float64 `json:"additionalPayment"`
	Currency string `json:"currency"`
	BanReason string `json:"banReason"`	
}

func GetReportIds (reports ReportsResponse) (ids []uint64) {
	for _, report := range reports.Data.Reports {
		ids = append(ids, report.Id)
	}
	return ids
}

func GetReportsByIds (r ReportsResponse, ids []uint64) (reports []Report) {
	for _, v := range r.Data.Reports {
		for _, reportId := range ids {
			if v.Id == reportId {
				reports = append(reports, v)
			}				
		}
	}
	return reports
}

