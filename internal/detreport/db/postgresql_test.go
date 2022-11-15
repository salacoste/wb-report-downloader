package db

import (
	"reflect"
	"testing"
	"wb-report-downloader/internal/detreport"

	"github.com/go-faker/faker/v4"
)

func TestMakeInsertQuery(t *testing.T) {
	v3 := [5]detreport.ReportRowV3{}
	// v3 := make([]detreport.ReportRowV3, 5)
	err := faker.FakeData(&v3)
	if err != nil {
		t.Fatalf("could not make fake data: %s", err)
	}

	r := detreport.DetailedReport{
		Data: reflect.ValueOf(v3[:]),
	}

	insertQuery, err := makeDetailedReportInsertQuery(&r)
	if err != nil {
		t.Fatalf("could not make insert query: %s", err)
	}
	t.Log(insertQuery)
}