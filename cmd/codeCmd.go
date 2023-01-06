package cmd

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
)

func CodeCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "code",
		Short: "",
		Run: func(cmd *cobra.Command, args []string) {
			stdout, err := handleCodeCommand(args)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Fprint(cmd.OutOrStdout(), stdout)
		},
	}
}

func handleCodeCommand(args []string) (string, error) {

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
	codeCmd := CodeCmd()
	rootCmd.AddCommand(codeCmd)
}
