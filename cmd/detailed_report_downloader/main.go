package main

import (
	"context"
	"log"
	"os"
	"time"
	"wb-report-downloader/internal/config"
	"wb-report-downloader/internal/task"

	"wb-report-downloader/internal/cookies/db"
	"wb-report-downloader/internal/detreport"
	"wb-report-downloader/internal/detreport/db"
	"wb-report-downloader/internal/task/db"

	"wb-report-downloader/internal/wb_request"
	"wb-report-downloader/internal/ziptool"
	"wb-report-downloader/pkg/client/postgresql"
)

func main() {
	log.Printf("Detailed report downloader\n")

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

	for {
		if(!workIteration(postgreSQLClient)) {
			log.Printf("New tasks not found. Sleep for %v secs", c.SleepOnTaskNotFoundSec)
			time.Sleep(time.Duration(c.SleepOnTaskNotFoundSec) * time.Second)
		}
	}
}

func workIteration(db_client postgresql.Client) bool {
	const (
		kTasksLimitPerIteration = 2
	)

	log.Printf("Searching for report download task...\n")
	taskRepository := taskdb.NewRepository(db_client)
	tasks, err := taskRepository.GetDownloadTasks(context.TODO(), kTasksLimitPerIteration)
	if err != nil {
		log.Fatalf("GetDownloadTask: %v\n", err)
	}
	if len(tasks) <= 0 {
		return false
	}

	for _, task := range tasks {
		log.Printf("Handle task. report_id: %v seller_id: %v \n", task.ReportID, task.SellerID)
		handleTask(task, taskRepository, db_client)
	}
	return true
}

func handleTask(taskData task.Task, taskrep task.Repository, db_client postgresql.Client) {
	log.Printf("Getting cookies for seller: %s\n", taskData.SellerName)
	cookiesRepository := cookiesdb.NewRepository(db_client)
	cookies, err := cookiesRepository.GetCookies(context.TODO(), taskData.SellerID)
	if err != nil {
		log.Fatalf("GetCookies: %v\n", err)
	}
	log.Printf("Getting cookies... OK\n")
	log.Printf("Sending http request to wb...\n")
	zippedExcelReport := wb_request.GetDetailedReport(uint64(taskData.ReportID), cookies.RawCookies)
	log.Printf("Response received. File size: %v\n", len(zippedExcelReport))

	log.Printf("Unzipping file...\n")
	excelReportBytes := ziptool.DecompressData(zippedExcelReport)
	log.Printf("Decompressed data size: %v\n", len(excelReportBytes))

	os.WriteFile("excelReport_downloaded.xlsx", excelReportBytes, 0644)

	log.Printf("Parsing xlsx...\n")
	report, err := detreport.ParseReportDetailesXlsx(excelReportBytes)
	if err != nil {
		log.Fatalf("report parsing error: %v\n", err)
	}
	log.Printf("Parsing xlsx... OK\n")

	if report.IsEmpty {
		log.Printf("Empty report. Set 'empty' status and skip data saving\n")
		taskData.Status = task.Empty
		taskrep.UpdateTaskStatus(context.TODO(), taskData)
		return
	}

	log.Printf("Saving detailed report...\n")
	reportRepository := db.NewRepository(db_client)
	for i := 0; i < report.Data.Len(); i++ {
		reportRow := report.Data.Index(i)
		reportRow.FieldByName("ReportID").SetUint(uint64(taskData.ReportID))
		reportRow.FieldByName("SellerID").SetUint(taskData.SellerID)
	}
	err = reportRepository.Create(context.TODO(), report)
	if err != nil {
		log.Fatalf("Insert detailed report error: %s\n", err)
	}
	log.Printf("Saving detailed report... OK\n")

	taskData.Status = task.Downloaded
	log.Printf("Updating task status...\n")
	err = taskrep.UpdateTaskStatus(context.TODO(), taskData)
	if err != nil {
		log.Fatalf("Could not update wb_reports.status: %v\n", taskData)
	}
	log.Printf("Updating task status... OK\n")
}