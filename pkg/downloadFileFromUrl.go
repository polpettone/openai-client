package pkg

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
)

func DownloadFileFromUrl(url string, imageName string) error {

	fmt.Println("Download Image ...")
	client := &http.Client{}

	req, _ := http.NewRequest("GET", url, nil)

	res, err := client.Do(req)

	if err != nil {
		return err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return errors.New(fmt.Sprintf("Fehler beim Herunterladen des Bildes: %s", res.Status))
	}

	file, err := os.Create(imageName)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, res.Body)
	if err != nil {
		return err
	}

	return nil
}
