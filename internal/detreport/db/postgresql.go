package db

import (
	"context"
	"fmt"
	"log"
	"reflect"
	"wb-report-downloader/internal/detreport"
	"wb-report-downloader/pkg/client/postgresql"
)

type repository struct {
	client postgresql.Client
}

func NewRepository(client postgresql.Client) detreport.Repository {
	return &repository{
		client: client,
	}
}

func (r *repository) Create(ctx context.Context, report *detreport.DetailedReport) error {
	if(report == nil) {
		return fmt.Errorf("report is nil")
	}
	q, err := makeDetailedReportInsertQuery(report)
	if err != nil {
		return err
	}
	_, err = r.client.Exec(ctx, q)
	return err
}

func makeDetailedReportInsertQuery(report *detreport.DetailedReport) (string, error) {
	const (
		kTableName = "wb_reports_details_v2"
	)

	insert := "INSERT INTO " + kTableName
	values := " VALUES "

	switch k := report.Data.Kind(); k {
	case reflect.Slice:
		insert += "(" + makeColumns(report.Data.Type().Elem()) + ")"
		
		for i := 0; i < report.Data.Len(); i++ {
			row := report.Data.Index(i)
			values += "(" + makeDataRow(row) + ")"
			if i != report.Data.Len()-1 {
				values += ", "
			} 
		}
	default:
		return "", fmt.Errorf("could not convert to reflect.Slice: %s", k)
	}
	log.Printf("Insert query (witout values): %s", insert)

	return insert + values, nil
}

func makeColumns(reportStruct reflect.Type) string {
	var cols string
	for i := 0; i < reportStruct.NumField(); i++ {
		col := reportStruct.Field(i).Tag.Get("db")
		cols += "\"" + col + "\""
		if i != reportStruct.NumField()-1 {
			cols += ", "
		}
	}
	return cols
}

func makeDataRow(rv reflect.Value) string {
	var row string
	for i := 0; i < rv.NumField(); i++ {
		value := fmt.Sprintf("%v", rv.Field(i).Interface())
		if rv.Field(i).Type().Kind() == reflect.String {
			row += "'" + value + "'"
		} else {
			row += value
		}
		if i != rv.NumField()-1 {
			row += ", "
		}
	}
	return row
}