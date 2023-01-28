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

	outputFile, err := cobraCommand.Flags().GetString("outputFile")
	if err != nil {
		return "", err
	}

	maxTokens, err := cobraCommand.Flags().GetInt("maxTokens")
	if err != nil {
		return "", err
	}

	contextMemoryID, err := cobraCommand.Flags().GetString("contextMemoryID")
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

	provider := NewProvider(1000, true, contextMemoryID)

	result, err := provider.Prompt(query, model, temperature, maxTokens)

	if err != nil {
		return "", err
	}

	if outputFile != "" {
		err := pkg.WriteToFile(outputFile, result)
		if err != nil {
			return "", err
		}
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

	askCmd.Flags().StringP(
		"outputFile",
		"o",
		"",
		"write response also to a given file")

	askCmd.Flags().IntP(
		"maxTokens",
		"t",
		3000,
		"")

	askCmd.Flags().StringP(
		"contextMemoryID",
		"c",
		"context-1",
		"id of context memory")
}

func readFromFile(path string) (string, error) {
	//Read the file content
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(data), nil

}
