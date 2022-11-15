package wb_request

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"wb-report-downloader/internal/ziptool"
	"wb-report-downloader/internal/report"
)

const (
	kGetDetailedReportUrl = "https://seller.wildberries.ru/ns/realization-reports/suppliers-portal-analytics/api/v1/reports/%d/details/archived-excel"
	kGetReportsUrl = "https://seller.wildberries.ru/ns/realization-reports/suppliers-portal-analytics/api/v1/reports"
)

type ArchivedExcelResponse struct {
	Data struct {
		File string `json:"file"`
	} `json:"data"`
}

func GetDetailedReport(reportID uint64, rawCookies string) []byte {
	url := fmt.Sprintf(kGetDetailedReportUrl, reportID)
	log.Printf("url: %v\n", url)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatalf("http.NewRequest: %v\n", err)
	}

	req.Header.Add("Cookie", rawCookies)

	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()
	log.Printf("archived-excel Status code: %v\n", response.StatusCode)
	
	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	// log.Printf("archived-excel Body: %v\n", string(body))

	var archivedExcelResponse ArchivedExcelResponse
	err = json.Unmarshal(body, &archivedExcelResponse)
	if err != nil {
		log.Fatal(err)
	}

	// log.Printf("archivedExcelResponse: %v\n", archivedExcelResponse.Data.File)
	return ziptool.Unbase64(archivedExcelResponse.Data.File)
}

func GetReports(rawCookies string) report.ReportsResponse {
	req, err := http.NewRequest("GET", kGetReportsUrl, nil)
	if err != nil {
		log.Fatalf("http.NewRequest: %v\n", err)
	}
	req.Header.Add("Cookie", rawCookies)

	q := req.URL.Query()
	q.Add("limit", "108")
	q.Add("skip", "0")
	q.Add("type", "2")
	req.URL.RawQuery = q.Encode()
	
	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()
	log.Printf("http request 'reports' status code: %v\n", response.StatusCode)

	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	// log.Printf("response body: %v\n", string(body))

	var reports report.ReportsResponse
	err = json.Unmarshal(body, &reports)
	if err != nil {
		log.Fatal(err)
	}
	return reports
}

