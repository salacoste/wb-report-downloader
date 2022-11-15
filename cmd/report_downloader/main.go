package main

import (
	"context"
	"log"
	"wb-report-downloader/internal/config"
	cookiesdb "wb-report-downloader/internal/cookies/db"
	"wb-report-downloader/internal/report"
	reportdb "wb-report-downloader/internal/report/db"
	"wb-report-downloader/internal/wb_request"
	"wb-report-downloader/pkg/client/postgresql"
	"wb-report-downloader/pkg/slice"
)

func main()  {
	const (
		kSellerID = 1
	)

	log.Printf("Report downloader\n")

	log.Printf("Reading config...")
	c, err := config.LoadConfig(".")
	if err != nil {
		log.Fatalf("LoadConfig: %v\n", err)
	}
	log.Printf("Reading config... OK")

	log.Printf("Connecting database...\n")
	postgreSQLClient, err := postgresql.NewClient(context.TODO(),
		c.Database.User, c.Database.Password, c.Database.Host, c.Database.Port, c.Database.Name)
	if err != nil {
		log.Fatalf("%v\n", err)
	}
	log.Printf("Connecting database... OK\n")

	log.Printf("Getting cookies for seller: %v\n", kSellerID)
	cookiesRepository := cookiesdb.NewRepository(postgreSQLClient)
	cookies, err := cookiesRepository.GetCookies(context.TODO(), kSellerID)
	if err != nil {
		log.Fatalf("GetCookies: %v\n", err)
	}
	log.Printf("Getting cookies... OK\n")

	log.Printf("Sending http request 'report'... \n")
	r := wb_request.GetReports(cookies.RawCookies)
	log.Printf("Http request ... OK \n")

	allReportsIds := report.GetReportIds(r)
	reportRepository := reportdb.NewRepository(postgreSQLClient)
	foundIds, err := reportRepository.FindAll(context.TODO(), allReportsIds)
	if err != nil {
		log.Fatalf("report FindAll: %v\n", err)
	}

	newReportsIds := slice.Difference(allReportsIds, foundIds)

	log.Printf("newReports: %v\n", newReportsIds)

	reportsForSave := report.GetReportsByIds(r, newReportsIds)

	for _, newReport := range reportsForSave {
		log.Printf("Saving report: %v ...\n", newReport.Id)
		err := reportRepository.Save(context.TODO(), kSellerID, &newReport)
		if err != nil {
			log.Fatalf("Error save report: %v", err)
		}
		log.Printf("Saving report: %v ... OK\n", newReport.Id)
	}
	
}

