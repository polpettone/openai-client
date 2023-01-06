package cmd

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
)

func ImageCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "image",
		Short: "",
		Run: func(cmd *cobra.Command, args []string) {
			stdout, err := handleImageCommand(cmd, args)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Fprint(cmd.OutOrStdout(), stdout)
		},
	}
}

func handleImageCommand(cobraCommand *cobra.Command, args []string) (string, error) {

	path, err := cobraCommand.Flags().GetString("file")

	if err != nil {
		return "", err
	}

	var fileContent string
	if path != "" {
		fileContent, err = readFromFile(path)
		if err != nil {
			return "", err
		}
	}

	contentFromArg := args[0]

	query := fmt.Sprintf("%s \n %s", contentFromArg, fileContent)

	if query == "" {
		return "", errors.New("no query")
	}

	imageName, err := ImageGenerator(query)

	if err != nil {
		return "", err
	}

	return imageName, err
}

func init() {
	imageCmd := ImageCmd()
	rootCmd.AddCommand(imageCmd)

	imageCmd.Flags().StringP(
		"file",
		"f",
		"",
		"query from file, get appended to question from argument")
}
