package commands

import (
	"fmt"

	"github.com/polpettone/openai-client/cmd/client"

	"github.com/spf13/cobra"
)

func ChatCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "chat",
		Short: "",
		Run: func(cmd *cobra.Command, args []string) {
			stdout, err := handleChatCommand(cmd, args)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Fprint(cmd.OutOrStdout(), stdout)
		},
	}
}

func handleChatCommand(cobraCommand *cobra.Command, args []string) (string, error) {

	client := client.NewLlamaClient()

	var contentFromArg string
	if len(args) > 0 {
		contentFromArg = args[0]
	}

	response, err := client.Complete(contentFromArg)

	if err != nil {
		return "", err
	}

	var output string
	for _, n := range response.Choices {
		output += fmt.Sprintf("%s\n", n.Message.Content)
	}

	return output, nil
}

func init() {
	chatCmd := ChatCmd()
	rootCmd.AddCommand(chatCmd)
}
