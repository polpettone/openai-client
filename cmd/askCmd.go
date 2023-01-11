package cmd

import (
	"errors"
	"fmt"
	"io/ioutil"

	"github.com/polpettone/labor/openai-client/pkg"
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

	openEditor, err := cobraCommand.Flags().GetBool("queryFromEditor")
	if err != nil {
		return "", err
	}

	path, err := cobraCommand.Flags().GetString("file")
	if err != nil {
		return "", err
	}

	model, err := cobraCommand.Flags().GetString("model")
	if err != nil {
		return "", err
	}

	temperature, err := cobraCommand.Flags().GetFloat64("temperature")
	if err != nil {
		return "", err
	}

	var queryContentFromFile string

	if openEditor {
		queryContentFromFile, err = pkg.CaptureInputFromEditor("")
		if err != nil {
			return "", err
		}
	} else {
		if path != "" {
			queryContentFromFile, err = readFromFile(path)
			if err != nil {
				return "", err
			}
		}
	}

	var contentFromArg string
	if len(args) > 0 {
		contentFromArg = args[0]
	}

	if contentFromArg == "" && queryContentFromFile == "" {
		return "", errors.New("No query. Provide one as argument or via file")
	}

	query := fmt.Sprintf("%s \n %s", contentFromArg, queryContentFromFile)

	if query == "" {
		return "", errors.New("no question asked")
	}

	result, err := Questioner(query, model, temperature)

	if err != nil {
		return "", err
	}

	return result, nil
}

func init() {
	askCmd := AskCmd()
	rootCmd.AddCommand(askCmd)

	askCmd.Flags().BoolP(
		"queryFromEditor",
		"q",
		false,
		"opens a temp file in an editor to write a query")

	askCmd.Flags().StringP(
		"file",
		"f",
		"",
		"query from file, get appended to question from argument")

	askCmd.Flags().StringP(
		"model",
		"m",
		"text-davinci-003",
		"model")

	askCmd.Flags().Float64(
		"temperature",
		0.7,
		"https://beta.openai.com/docs/api-reference/completions/create#completions/create-temperature")

}

func readFromFile(path string) (string, error) {
	//Read the file content
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(data), nil

}
