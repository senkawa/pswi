package umi

import (
	"errors"
	"fmt"
	"io"
	"strconv"

	"github.com/lukasjarosch/go-docx"
)

type ReportProcessor struct {
	templateData []byte
	reportTypes  []string

	parser *ReportParser
}

func NewReportProcessor(templateData []byte, reportTypes []string, file io.Reader) (*ReportProcessor, error) {
	if len(reportTypes) <= 0 {
		return nil, errors.New("one or more report types must be specified")
	}

	return &ReportProcessor{
		templateData: templateData,
		reportTypes:  reportTypes,
		parser:       NewReportParser(file),
	}, nil
}

func (rp *ReportProcessor) ReplacePlaceholders(report *Report) docx.PlaceholderMap {
	return docx.PlaceholderMap{
		"background": report.Background,
		"task":       report.Task,
		"result":     report.Result,
		"feedback":   report.Feedback,
		"week":       report.WeekNumber,
		"date":       report.Date,
	}
}

func (rp *ReportProcessor) GenerateFileName(report *Report, reportType string) (string, error) {
	reportNumber, err := strconv.Atoi(report.ReportNumber)
	if err != nil {
		return "", fmt.Errorf("got non-numeric report number: %s", report.ReportNumber)
	}

	return fmt.Sprintf("output/%s %02d.docx", reportType, reportNumber), nil
}

func (rp *ReportProcessor) Process() ([]Report, error) {
	reports, err := rp.parser.Parse()
	if err != nil {
		return nil, err
	}

	for _, reportType := range rp.reportTypes {
		for _, report := range reports {
			if err := rp.replace(report, reportType); err != nil {
				return nil, err
			}
		}
	}

	return reports, nil
}

func (rp *ReportProcessor) replace(report Report, reportType string) error {
	doc, err := docx.OpenBytes(rp.templateData)
	if err != nil {
		return err
	}

	placeholders := rp.ReplacePlaceholders(&report)
	if err = doc.ReplaceAll(placeholders); err != nil {
		return err
	}

	filename, err := rp.GenerateFileName(&report, reportType)
	if err != nil {
		return err
	}

	err = doc.WriteToFile(filename)
	doc.Close()
	if err != nil {
		return err
	}
	return nil
}
