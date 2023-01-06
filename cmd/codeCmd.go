package cmd

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
)

func AskCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "ask",
		Short: "",
		Run: func(cmd *cobra.Command, args []string) {
			stdout, err := handleAskCommand(args)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Fprint(cmd.OutOrStdout(), stdout)
		},
	}
}

func handleAskCommand(args []string) (string, error) {

	question := args[0]

	if question == "" {
		return "", errors.New("no question asked")
	}

	result, err := Questioner(question)

	if err != nil {
		return "", err
	}

	return result, nil
}

func init() {
	askCmd := AskCmd()
	rootCmd.AddCommand(askCmd)
}
