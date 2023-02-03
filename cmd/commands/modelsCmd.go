package commands

import (
	"fmt"
	"github.com/polpettone/openai-client/cmd/provider"

	"github.com/spf13/cobra"
)

func ModelsCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "models",
		Short: "",
		Run: func(cmd *cobra.Command, args []string) {
			stdout, err := handleModelsCommand(cmd, args)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Fprint(cmd.OutOrStdout(), stdout)
		},
	}
}

func handleModelsCommand(cobraCommand *cobra.Command, args []string) (string, error) {

	models, err := provider.ListModels()

	if err != nil {
		return "", err
	}

	result := ""
	for _, m := range models {
		result += fmt.Sprintf("%s   %s\n", m.ID, m.OwnedBy)
	}

	return fmt.Sprintf("%v", result), nil
}

func init() {
	modelsCmd := ModelsCmd()
	rootCmd.AddCommand(modelsCmd)
}
