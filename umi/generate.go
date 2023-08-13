package umi

import (
	"bytes"
	"fmt"
	"os"
	"time"
)

const templateFormat = `# %s | %d | %d
Background:

Task:

Results:

Feedback:
`

type DateGenerator struct {
	FilenameToWrite string
	StartDate       time.Time
}

func NewDateGenerator(filename string, startDate string) (*DateGenerator, error) {
	startingDate, err := startingDate(startDate)
	if err != nil {
		return nil, err
	}

	return &DateGenerator{
		FilenameToWrite: filename,
		StartDate:       startingDate,
	}, nil
}

func (d DateGenerator) Persist(dates Dates) error {
	var buffer bytes.Buffer
	for _, date := range dates {
		_, err := fmt.Fprintf(
			&buffer,
			templateFormat,
			date.Format("2006-01-02"),
			DiscoverReportNumber(date),
			WeekOfTheYear(date),
		)
		if err != nil {
			return fmt.Errorf("error while writing to buffer: %w", err)
		}
	}

	// write to file
	err := os.WriteFile(d.FilenameToWrite, buffer.Bytes(), 0o644)
	if err != nil {
		return fmt.Errorf("error while writing to file: %w", err)
	}

	return nil
}

// RangeOfDates generates a list of dates for bi-weekly reports.
func (d DateGenerator) RangeOfDates() (Dates, error) {
	var dates Dates
	for d.StartDate.Before(time.Now()) {
		dates = append(dates, d.StartDate)
		d.StartDate = d.StartDate.AddDate(0, 0, 14)
	}

	return dates, nil
}

// parseStartDate parses a date string in the format of "2006-01-02" and returns a time.Time.
// If the date is not a Friday, the closest Friday will be returned.
func startingDate(date string) (time.Time, error) {
	t, err := time.Parse("2006-01-02", date)
	if err != nil {
		return time.Time{}, fmt.Errorf("error while parsing date: %w", err)
	}

	if t.Weekday() != time.Friday {
		return time.Time{}, fmt.Errorf("starting date must be a Friday")
	}

	return t, nil
}
