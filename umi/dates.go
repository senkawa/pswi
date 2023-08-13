package umi

import (
	"bytes"
	"fmt"
	"time"
)

type Dates []time.Time

func (d Dates) String() string {
	var buffer bytes.Buffer

	for _, date := range d {
		_, err := fmt.Fprintf(
			&buffer,
			"%s | ISO: %d | Report Number: %d\n",
			date.Format("2006-01-02 Monday"),
			WeekOfTheYear(date),
			DiscoverReportNumber(date),
		)
		if err != nil {
			panic(err) // this should never happen
		}
	}

	return buffer.String()
}
