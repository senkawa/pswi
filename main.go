package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/terminal"
	"github.com/fatih/color"
	"github.com/lukasjarosch/go-docx"
)

type answers struct {
	Background string
	Task       string
	Result     string
	Feedback   string
	Week       string
	Date       string
}

type questions struct {
	Name      string
	Message   string
	Default   string
	Multiline bool
}

var toAsk = []questions{
	{Name: "Background", Message: "Background? {background}", Default: "N/A", Multiline: true},
	{Name: "Task", Message: "Task? {task}", Default: "N/A", Multiline: true},
	{Name: "Result", Message: "Result? {result}", Default: "N/A", Multiline: true},
	{Name: "Feedback", Message: "Feedback? {feedback}", Default: "N/A", Multiline: true},
	{Name: "Week", Message: "Week? {week} (Defaults to current ISO week)", Default: getWeekOfTheYear()},
	{Name: "Date", Message: "Date? {date} (Defaults to current date)", Default: time.Now().Format("2006-01-02")},
}

func getWeekOfTheYear() string {
	_, week := time.Now().ISOWeek()
	return fmt.Sprintf("%d", week)
}

func parseQuestions() []*survey.Question {
	var qs []*survey.Question
	for _, q := range toAsk {
		if q.Multiline {
			qs = append(qs, &survey.Question{
				Name:   q.Name,
				Prompt: &survey.Multiline{Message: q.Message, Default: q.Default},
			})
			continue
		}

		qs = append(qs, &survey.Question{
			Name:   q.Name,
			Prompt: &survey.Input{Message: q.Message, Default: q.Default},
		})
	}
	return qs
}

var reportTypes = []string{"IWSP", "Capstone"}

const help = `
Prepare a "template.docx" in the current directory with the following placeholders: {background}, {task}, {result}, {feedback}, {week}, {date}
If you're lazy, you can use the same contents for both report types.

- Leave any of the placeholders empty for N/A. For background/task/result/feedback, multiple lines are supported.
- Leave {date} empty for the current date.
- Leave {week} empty for the current ISO week.
- Outputs will follow the naming convention: {type}-{week}.docx

`

func main() {
	fmt.Print(help)

	sameContent := false
	if err := survey.AskOne(&survey.Confirm{
		Message: "Use the same content for both report types?",
		Help:    "If you select no, you'll be asked for the content twice.",
	}, &sameContent); err != nil {
		if errors.Is(err, terminal.InterruptErr) {
			os.Exit(0)
		}
	}

	var lookup answers
	for _, reportType := range reportTypes {
		if reportType == reportTypes[0] || !sameContent {
			lookup = askForAnswers(reportType)
		}

		doc, err := docx.Open("template.docx")
		if err != nil {
			log.Fatalf("error opening template: %v", err)
		}

		if err = doc.ReplaceAll(docx.PlaceholderMap{
			"background": lookup.Background,
			"task":       lookup.Task,
			"result":     lookup.Result,
			"feedback":   lookup.Feedback,
			"week":       lookup.Week,
			"date":       lookup.Date,
			"type":       reportType,
		}); err != nil {
			log.Fatal(err)
		}

		filename := fmt.Sprintf("output/%s-%s.docx", reportType, lookup.Week)
		if err = doc.WriteToFile(filename); err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Created %s.\n", color.HiYellowString(filename))
	}
}

func askForAnswers(reportType string) answers {
	fmt.Printf("%s Filling in report type: %s\n\n", color.YellowString("#"), reportType)

	qs := parseQuestions()

	var lookup answers
	if err := survey.Ask(qs, &lookup); err != nil {
		if errors.Is(err, terminal.InterruptErr) {
			os.Exit(0)
		}
	}

	return lookup
}
