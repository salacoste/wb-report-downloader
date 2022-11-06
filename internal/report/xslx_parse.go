package report

import (
	"bytes"
	"fmt"
	"log"
	"reflect"
	"strconv"
	"strings"

	"github.com/xuri/excelize/v2"
)

const (
	kXlsxPageName string = "Sheet1"
)

func ParseReportDetailesXlsx(excelReportBytes []byte) ([]ReportDetailes, error) {
	excelReportReqder := bytes.NewReader(excelReportBytes)
	excelReport, err := excelize.OpenReader(excelReportReqder)
	if err != nil {
		return nil, fmt.Errorf("could not open xlsx report")
	}

	rows, err := excelReport.GetRows(kXlsxPageName)
	if err != nil {
		return nil, fmt.Errorf("could not get rows from %v", kXlsxPageName)
	}
	if len(rows) <= 1 {
		return nil, fmt.Errorf("empty report")
	}

	headers := rows[0]
	data := rows[1:]

	err = validateColumnsCount(headers)
	if err != nil {
		return nil, err
	}


	ReportDetailesList := make([]ReportDetailes, len(data))
	for rowNumber, row := range data {
		row = row[:len(headers)]
		for colIndex := 0; colIndex < len(headers); colIndex++ {
			columnName := headers[colIndex]
			fieldName, fieldType := getFieldNameByTag(columnName)
			structField := reflect.ValueOf(&ReportDetailesList[rowNumber]).Elem().FieldByName(fieldName)
			if reflect.TypeOf(row[colIndex]).ConvertibleTo(fieldType) {
				structField.Set(reflect.ValueOf(row[colIndex]).Convert(fieldType))
			} else if fieldType.Kind() == reflect.Uint64 {
				value, _ := strconv.ParseUint(row[colIndex], 10, 64)
				structField.SetUint(value)
			} else if fieldType.Kind() == reflect.Float64 {
				value, _ := strconv.ParseFloat(row[colIndex], 64)
				structField.SetFloat(value)
			} else {
				return nil, fmt.Errorf("unexpected type '%v' field '%v'", fieldType.Kind(), fieldName)
			}
		}
	}

	return ReportDetailesList, nil
}

func validateColumnsCount(headers []string) error {
	expectedColCount := 0
	rt := reflect.TypeOf(ReportDetailes{})
	var structFields []string
	for i := 0; i < rt.NumField(); i++ {
		name, containXlsxTag := rt.Field(i).Tag.Lookup("xlsx")
		if containXlsxTag {
			structFields = append(structFields, name)
			expectedColCount++
		}
	}
	if unexpectedHeadersFound := difference(headers, structFields); len(unexpectedHeadersFound) > 0 {
		return fmt.Errorf("theese unexpected xlsx headers are found: %v", strings.Join(unexpectedHeadersFound, ", "))
	}
	if headersNotFound := difference(structFields, headers); len(headersNotFound) > 0 {
		return fmt.Errorf("theese expected xlsx headers not found: %v", strings.Join(headersNotFound, ", "))
	}
	return nil
}

func getFieldNameByTag(name string) (string, reflect.Type) {
	rt := reflect.TypeOf(ReportDetailes{})
	for i := 0; i < rt.NumField(); i++ {
		tagName := rt.Field(i).Tag.Get("xlsx")
		if tagName == name {
			return rt.Field(i).Name, rt.Field(i).Type
		}
	}
	log.Fatalf("Tag '%v' not found ", name)
	return "", nil
}

func difference(a, b []string) []string {
    mb := make(map[string]struct{}, len(b))
    for _, x := range b {
        mb[x] = struct{}{}
    }
    var diff []string
    for _, x := range a {
        if _, found := mb[x]; !found {
            diff = append(diff, x)
        }
    }
    return diff
}