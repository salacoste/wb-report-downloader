package detreport

import (
	"os"
	"reflect"
	"testing"
	"wb-report-downloader/internal/detreporttest"
)

const (
	kTestDataDir = "../detreporttest/testdata/"
)

func TestParseReportDetailesXlsx(t *testing.T) {
	var testFiles = []string{
		kTestDataDir + "v1.xlsx", 
		kTestDataDir + "v2.xlsx", 
		kTestDataDir + "v3.xlsx",
	}

	for _, file := range testFiles {
		b, err := os.ReadFile(file)
		if err != nil {
			t.Errorf("Reading testdata by path: '%s' caused error: '%s'\n", file, err)
		}
	
		r, err := ParseReportDetailesXlsx(b)
		if err != nil {
			t.Errorf("For file '%s' ParseReportDetailesXlsx: %s\n", file, err)
		}
		if r == nil {
			t.Fatalf("For file '%s' parsed report is nil\n", file)
		}

		e, err := detreporttest.Data(file)
		if err != nil {
			t.Fatalf("Reading test data error: %s\n", err)
		}
		if len(e) != r.Data.Len() {
			t.Errorf("For file: '%s' expected len: %d, actual: %d\n", file, len(e), r.Data.Len())
		}
	}
}

type versionRecognizerTest struct {
	testFile     string
	expectedType any
}

var vRecognizerTests = []versionRecognizerTest{
	{kTestDataDir + "v1.xlsx", ReportRowV1{}},
	{kTestDataDir + "v2.xlsx", ReportRowV2{}},
	{kTestDataDir + "v3.xlsx", ReportRowV3{}},
}

func TestVersionRecognizer(t *testing.T) {
	for _, test := range vRecognizerTests {
		h, err := detreporttest.Headers(test.testFile)
		if err != nil {
			t.Error(err)
		}
		v, err := recognizeReportVersion(h)
		if err != nil {
			t.Errorf("recognizeReportVersion: %s\n", err)
		}

		if reflect.TypeOf(test.expectedType) != reflect.TypeOf(v) {
			t.Errorf("For file '%s' recogized report type expected: %T, actual: %T\n", test.testFile, test.expectedType, v)
		}
	}
}

func TestGetFieldByName(t *testing.T) {
	rt := reflect.TypeOf(ReportRowV1{})
	fName, fType := getFieldNameByTag("Номер поставки", rt)
	if fName != "SupplyNumber" {
		t.Errorf("Field name not found")
	}
	if fType == nil {
		t.Errorf("Nil fType")
	} else if fType.Kind() != reflect.Uint64 {
		t.Errorf("Incorrect field type")
	}
}

