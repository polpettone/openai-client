package cmd

import (
	"errors"
	"fmt"
	"io/ioutil"

	"github.com/spf13/cobra"
)

func AskCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "ask",
		Short: "",
		Run: func(cmd *cobra.Command, args []string) {
			stdout, err := handleAskCommand(cmd, args)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Fprint(cmd.OutOrStdout(), stdout)
		},
	}
}

func handleAskCommand(cobraCommand *cobra.Command, args []string) (string, error) {

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
		return "", errors.New("no question asked")
	}

	result, err := Questioner(query)

	if err != nil {
		return "", err
	}

	return result, nil
}

func init() {
	askCmd := AskCmd()
	rootCmd.AddCommand(askCmd)

	askCmd.Flags().StringP(
		"file",
		"f",
		"",
		"query from file, get appended to question from argument")
}

func readFromFile(path string) (string, error) {
	//Read the file content
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(data), nil

}
