package cmd

import (
	"fmt"
	"log"

	"github.com/fatih/color"

	"github.com/senkawa/pswi/umi"

	"github.com/spf13/cobra"
)

var generateCmd = &cobra.Command{
	Use:   "template start_date [output_file=output.md]",
	Args:  cobra.RangeArgs(1, 2),
	Short: "Generate a template for a given start date until today's week",
	Long:  `Generates a template for a given start date until today's week. This should be filled in and fed into the parse command.`,
	Run: func(cmd *cobra.Command, args []string) {
		var filename string
		if len(args) <= 1 {
			filename = "output.md"
		} else {
			filename = args[1]
		}

		date, err := umi.NewDateGenerator(filename, args[0])
		if err != nil {
			log.Fatal(err)
		}

		dates, err := date.RangeOfDates()
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Written template to %s\n\n", color.YellowString(filename))
		c := color.New(color.FgBlack).Add(color.Bold).Add(color.BgMagenta)
		c.Printf(" Expected Reports \n\n")

		fmt.Print(dates)

		if err = date.Persist(dates); err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(generateCmd)
}
