// Testing utils for DeatailedReport
package detreporttest

import (
	"bytes"
	"fmt"
	"os"

	"github.com/xuri/excelize/v2"
)

const (
	kXlsxPageName string = "Sheet1"
)

func readToExcelize(testFileName string) (excelReport *excelize.File, err error) {
	b, err := os.ReadFile(testFileName)
	if err != nil {
		return nil, fmt.Errorf("Reading testdata by path: '%s' caused error: '%s'", testFileName, err)
	}
	r := bytes.NewReader(b)
	excelReport, err = excelize.OpenReader(r)
	return excelReport, err
}

func Headers(testFileName string) ([]string, error) {
	e, err := readToExcelize(testFileName)
	if err != nil {
		return nil, err
	}
	rows, err := e.Rows(kXlsxPageName)
	if err != nil {
		return nil, fmt.Errorf("Parsing xlsx: '%s'", err)
	}
	rows.Next()
	h, err := rows.Columns()
	if err != nil {
		return nil, fmt.Errorf("Reading headers: '%s'", err)
	}
	return h, nil
}

func Data(testFileName string) ([][]string, error) {
	e, err := readToExcelize(testFileName)
	if err != nil {
		return nil, err
	}
	rows, err := e.GetRows(kXlsxPageName)
	if err != nil {
		return nil, err
	}
	data := rows[1:]
	return data, err
}

