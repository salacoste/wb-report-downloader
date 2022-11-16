package detreport

import (
	"bytes"
	"fmt"
	"log"
	"reflect"
	"strconv"
	"strings"
	"wb-report-downloader/pkg/slice"

	"github.com/xuri/excelize/v2"
)

const (
	kXlsxPageName string = "Sheet1"
)

func ParseReportDetailesXlsx(excelReportBytes []byte) (*DetailedReport, error) {
	excelReportReqder := bytes.NewReader(excelReportBytes)
	excelReport, err := excelize.OpenReader(excelReportReqder)
	if err != nil {
		return nil, fmt.Errorf("could not open xlsx report")
	}

	rows, err := excelReport.GetRows(kXlsxPageName)
	if err != nil {
		return nil, fmt.Errorf("could not get rows from %v", kXlsxPageName)
	}
	if len(rows) == 0 {
		return nil, fmt.Errorf("empty report")
	}

	headers := rows[0]
	// определяем версию отчета по заголовкам в xlsx
	interfaceReportType, err := recognizeReportVersion(headers)
	log.Printf("recognized report version: %T", interfaceReportType)
	if err != nil {
		return nil, err
	}
	repType := reflect.TypeOf(interfaceReportType)

	data := rows[1:]
	if len(data) == 0 {
		return &DetailedReport{Data: reflect.Value{}, IsEmpty: true}, nil
	}
	return parseReportData(repType, headers, &data)
}

func recognizeReportVersion(headers []string) (any, error) {
	var notFoundCols, newCols []string
	var fit bool
	for _, kind := range repKinds {
		fit, notFoundCols, newCols = compareHeadersWithReportType(headers, kind)
		if fit {
			return kind, nil
		}
	}
	return nil, fmt.Errorf("incompatible report version.\nDeprecated: %s\nNew: %s",
		strings.Join(notFoundCols, "; "), strings.Join(newCols, "; "))
}

func compareHeadersWithReportType(headers []string, structType any) (matched bool, notFoundCols []string, newCols []string) {
	f := structFieldsWithTagXlsx(structType)
	notFoundCols = slice.Difference(headers, f)
	newCols = slice.Difference(f, headers)
	if len(notFoundCols) == 0 && len(newCols) == 0 {
		return true, nil, nil
	}
	return false, notFoundCols, newCols
}

func structFieldsWithTagXlsx(structType any) []string {
	rt := reflect.TypeOf(structType)
	var structFields []string
	for i := 0; i < rt.NumField(); i++ {
		name, containXlsxTag := rt.Field(i).Tag.Lookup("xlsx")
		if containXlsxTag {
			structFields = append(structFields, name)
		}
	}
	return structFields
}

func parseReportData(repType reflect.Type, headers []string, data *[][]string) (dr *DetailedReport, err error) {
	reportDataset := reflect.MakeSlice(reflect.SliceOf(repType), len(*data), len(*data))
	for rowNumber, row := range *data {
		row = row[:len(headers)]
		reportElem := reportDataset.Index(rowNumber)
		for colIndex := 0; colIndex < len(headers); colIndex++ {
			columnName := headers[colIndex]
			fieldName, fieldType := getFieldNameByTag(columnName, reportElem.Type())
			if fieldName == "" {
				return nil, fmt.Errorf("could not found xlsx Tag: %s in struct: %s",
					columnName, reflect.TypeOf(reportElem).Name())
			}
			structField := reportElem.FieldByName(fieldName)
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
	dr = new(DetailedReport)
	dr.Data = reportDataset
	return dr, nil
}

func getFieldNameByTag(name string, rt reflect.Type) (string, reflect.Type) {
	for i := 0; i < rt.NumField(); i++ {
		tagName := rt.Field(i).Tag.Get("xlsx")
		if tagName == name {
			return rt.Field(i).Name, rt.Field(i).Type
		}
	}
	return "", nil
}
