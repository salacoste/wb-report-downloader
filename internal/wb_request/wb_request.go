package wb_request

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"wb-report-downloader/internal/ziptool"
)

const (
	kGetDetailedReportUrl = "https://seller.wildberries.ru/ns/realization-reports/suppliers-portal-analytics/api/v1/reports/%d/details/archived-excel"
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

