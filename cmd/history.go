package cmd

import (
	"fmt"

	"github.com/amitsuthar69/ollama-cli/cmd/client"
	"github.com/spf13/cobra"
)

var historyCmd = &cobra.Command{
	Use:   "history",
	Short: "Show the history of prompts",
	Long: `Usage: ollama history -d 6
this is will show history of past 6 days only.
Default: -1 (all history)
`,
	Run: showHistory,
}

var days int16

func showHistory(cmd *cobra.Command, args []string) {
	if days < -1 {
		fmt.Println("Invalid flag value for --days")
	} else {
		client.DisplayHistory(days)
	}
}

func init() {
	rootCmd.AddCommand(historyCmd)
	historyCmd.Flags().Int16VarP(&days, "days", "d", -1, "Number of days past.")
}
