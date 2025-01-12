package cmd

import (
	"strings"

	"github.com/amitsuthar69/ollama-cli/cmd/client"
	"github.com/spf13/cobra"
)

var askCmd = &cobra.Command{
	Use:          "ask <message>",
	Short:        "prompt the LLM",
	Args:         cobra.MinimumNArgs(1),
	SilenceUsage: true,
	Run:          askQroq,
}

var context bool

func askQroq(cmd *cobra.Command, args []string) {
	userInput := strings.Join(args, " ")

	client.ChatCompletion(userInput, context)
}

func init() {
	rootCmd.AddCommand(askCmd)
	askCmd.PersistentFlags().BoolVarP(&context, "ctx", "c", false, "Use Context Mode. [Context window: 10 mins]")
	rootCmd.TraverseChildren = true
}
