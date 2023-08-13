package umi

import "time"

// DefaultBaselineDate represents the first date of the first report: 2023-01-13
// This is the reference date to calculate the bi-weekly ranges (used in report number etc.)
var DefaultBaselineDate = time.Date(2023, 1, 13, 0, 0, 0, 0, time.UTC)

// every two weeks
const daysInABiWeekly = 14

// DiscoverReportNumber returns the report number for a given date.
func DiscoverReportNumber(forDate time.Time) int {
	difference := forDate.Sub(DefaultBaselineDate)

	// time.Duration doesn't have a Weeks()/Days() method. The largest unit is Hours().
	biWeeks := int(difference.Hours() / 24 / daysInABiWeekly)

	// Since we start counting reports from 1 (not 0), add 1 to the result
	biWeeks++

	return biWeeks
}

// WeekOfTheYear returns the week of the year for a given date.
func WeekOfTheYear(date time.Time) int {
	_, week := date.ISOWeek()
	return week
}
