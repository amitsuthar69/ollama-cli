package cmd

import (
	"strings"

	"github.com/amitsuthar69/ollama-cli/cmd/client"
	"github.com/spf13/cobra"
)

var askCmd = &cobra.Command{
	Use:   "ask <message>",
	Short: "prompt the LLM",
	Args:  cobra.MinimumNArgs(1),
	Run:   askQroq,
}

func askQroq(cmd *cobra.Command, args []string) {
	userInput := strings.Join(args, " ")
	client.ChatCompletion(userInput)
}

func init() {
	rootCmd.AddCommand(askCmd)
}
