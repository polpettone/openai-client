package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func BotCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "bot",
		Short: "starts a telegram bot",
		Run: func(cmd *cobra.Command, args []string) {
			stdout, err := handleBotCommand(cmd, args)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Fprint(cmd.OutOrStdout(), stdout)
		},
	}
}

func handleBotCommand(cobraCommand *cobra.Command, args []string) (string, error) {
	err := StartBot()
	if err != nil {
		return "", err
	}
	return "", nil
}

func init() {
	botCmd := BotCmd()
	rootCmd.AddCommand(botCmd)
}
