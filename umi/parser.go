package umi

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

type Report struct {
	Date         string
	ReportNumber string
	WeekNumber   string
	Background   string
	Task         string
	Result       string
	Feedback     string
}

const (
	headerPrefix     = "# "
	backgroundPrefix = "Background: "
	taskPrefix       = "Task: "
	resultPrefix     = "Results: "
	feedbackPrefix   = "Feedback: "
)

type Line int

const (
	NoneLine Line = iota
	HeaderLine
	BackgroundLine
	TaskLine
	ResultLine
	FeedbackLine
)

var prefixes = map[string]Line{
	headerPrefix:     HeaderLine,
	backgroundPrefix: BackgroundLine,
	taskPrefix:       TaskLine,
	resultPrefix:     ResultLine,
	feedbackPrefix:   FeedbackLine,
}

func NewLine(line string) Line {
	for prefix, lineType := range prefixes {
		if strings.HasPrefix(line, prefix) {
			return lineType
		}
	}

	return NoneLine
}

type ReportParser struct {
	scanner *bufio.Scanner
}

func NewReportParser(file io.Reader) *ReportParser {
	return &ReportParser{scanner: bufio.NewScanner(file)}
}

func (r *ReportParser) Parse() ([]Report, error) {
	var reports []Report

	var current Report
	var currentField *string

	for r.scanner.Scan() {
		line := r.scanner.Text()
		lineType := NewLine(line)

		switch lineType {
		case HeaderLine:
			// A header is matched, but the date is not empty.
			// This means that we have a new report to parse, and the current one is finished.
			if current.Date != "" {
				// complete the current report
				reports = append(reports, current)
				// start a new one
				current = Report{}
			}

			// expected format: # 2023-01-01 | 1 | 2
			split := strings.Split(line[2:], " | ")
			if len(split) != 3 {
				return nil, fmt.Errorf("invalid report header: %s", line)
			}

			current.Date = split[0]
			current.ReportNumber = split[1]
			current.WeekNumber = split[2]

			// reset the current field
			currentField = nil

		case BackgroundLine:
			currentField = &current.Background
			*currentField = strings.TrimPrefix(line, backgroundPrefix)
		case TaskLine:
			currentField = &current.Task
			*currentField = strings.TrimPrefix(line, taskPrefix)
		case ResultLine:
			currentField = &current.Result
			*currentField = strings.TrimPrefix(line, resultPrefix)
		case FeedbackLine:
			currentField = &current.Feedback
			*currentField = strings.TrimPrefix(line, feedbackPrefix)

		// no headers found, so we're in a field
		case NoneLine:
			// simple guard to protect against empty lines at the start of the file
			if currentField != nil {

				// this helps make
				if *currentField != "" {
					// only append a newline if the field is not empty
					*currentField += "\n"
				}

				// append the current line to the current field
				*currentField += line
			}
		}

	}

	// append the final report
	if current.Date != "" {
		reports = append(reports, current)
	}

	return reports, r.scanner.Err()
}
