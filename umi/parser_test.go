package umi

import (
	"strings"
	"testing"
)

func TestReportParser_Parse(t *testing.T) {
	testInput := `
# 2023-01-01 | 1 | 2
Background: some background info
Task: some task info
Results: some result info
Feedback: some feedback info

# 2023-01-02 | 2 | 3
Background: some other background info
Task: some other task info
Results: some other result info
Feedback: some other feedback info`

	expectedOutput := []Report{
		{
			Date:         "2023-01-01",
			ReportNumber: "1",
			WeekNumber:   "2",
			Background:   "some background info",
			Task:         "some task info",
			Result:       "some result info",
			Feedback:     "some feedback info\n",
		},
		{
			Date:         "2023-01-02",
			ReportNumber: "2",
			WeekNumber:   "3",
			Background:   "some other background info",
			Task:         "some other task info",
			Result:       "some other result info",
			Feedback:     "some other feedback info",
		},
	}

	parser := NewReportParser(strings.NewReader(testInput))
	actualOutput, err := parser.Parse()
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(actualOutput) != len(expectedOutput) {
		t.Fatalf("Report count mismatch. Expected: %v, got: %v", len(expectedOutput), len(actualOutput))
	}

	for i, report := range actualOutput {
		expectedReport := expectedOutput[i]

		if report != expectedReport {
			t.Errorf("Report mismatch at index %v.\nExpected: %+v\nGot: %+v", i, expectedReport, report)
		}
	}
}
