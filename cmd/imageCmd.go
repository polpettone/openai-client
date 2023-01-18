package cmd

import (
	"errors"
	"fmt"
	"strings"
	"sync"

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

	countOfCreations, err := cobraCommand.Flags().GetInt("countOfCreations")
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

	var contentFromArg string
	if len(args) > 0 {
		contentFromArg = args[0]
	}

	query := fmt.Sprintf("%s \n %s", contentFromArg, fileContent)

	if query == "" {
		return "", errors.New("no query")
	}

	var wg sync.WaitGroup

	for n := 0; n < countOfCreations; n++ {

		imageName := generateName(strings.Split(query, " ")[0])

		wg.Add(1)
		go func() {
			defer wg.Done()
			ImageGenerator(query, imageName)
		}()
	}
	wg.Wait()

	if err != nil {
		return "", err
	}

	return "done", err
}

func init() {
	imageCmd := ImageCmd()
	rootCmd.AddCommand(imageCmd)

	imageCmd.Flags().StringP(
		"file",
		"f",
		"",
		"query from file, get appended to question from argument")

	imageCmd.Flags().IntP(
		"countOfCreations",
		"n",
		1,
		"image get n times created")
}
