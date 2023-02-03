package cmd

import (
	"errors"
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

	if len(args) < 2 {
		return "", errors.New("2 Arguments required: telegramBotToken, memoryID")
	}

	telegramBotToken := args[0]
	if telegramBotToken == "" {
		return "", errors.New("telegramBotToken must not be empty")
	}

	memoryID := args[1]
	if memoryID == "" {
		return "", errors.New("memoryID must not be empty")
	}

	err := StartBot(telegramBotToken, memoryID)

	if err != nil {
		return "", err
	}
	return "", nil
}

func init() {
	botCmd := BotCmd()
	rootCmd.AddCommand(botCmd)
}
