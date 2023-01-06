package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

type ImageCreatingPayload struct {
	Prompt string `json:"prompt"`
	N      int    `json:"n"`
	Size   string `json:"size"`
}

type Response struct {
	Created int           `json:"created"`
	Data    []ResponseURL `json:"data"`
}

type ResponseURL struct {
	URL string `json:"url"`
}

func callImageCreating(imageDescription string) error {

	payload := ImageCreatingPayload{
		Prompt: imageDescription,
		N:      1,
		Size:   "1024x1024",
	}

	// Setze die URL und den Auth-Header des Requests
	url := "https://api.openai.com/v1/images/generations"
	authHeader := "Bearer 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"

	// Setze den Request-Body
	requestBody, err := json.Marshal(payload)

	if err != nil {
		return err
	}

	// Erstelle einen neuen HTTP-Client
	client := &http.Client{}

	// Erstelle einen neuen Request
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))

	// Setze den Auth-Header des Requests
	req.Header.Set("Authorization", authHeader)
	req.Header.Set("Content-Type", "application/json")

	// Führe den Request aus
	res, err := client.Do(req)

	if err != nil {
		return err
	}

	defer res.Body.Close()

	// Lies die Antwort aus
	responseBody, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return err
	}

	fmt.Println(string(responseBody))

	var response Response
	err = json.Unmarshal([]byte(responseBody), &response)

	if err != nil {
		return err
	}

	err = downloadImageFromUrl(response.Data[0].URL, generateName(imageDescription))

	if err != nil {
		return err
	}

	return nil

}

func downloadImageFromUrl(url string, imageName string) error {
	// Setze die URL des Bildes

	// Erstelle einen neuen HTTP-Client
	client := &http.Client{}

	// Erstelle einen neuen Request
	req, _ := http.NewRequest("GET", url, nil)

	// Führe den Request aus
	res, err := client.Do(req)

	if err != nil {
		return err
	}

	defer res.Body.Close()

	// Prüfe den Statuscode des Response
	if res.StatusCode != http.StatusOK {
		fmt.Println("Fehler beim Herunterladen des Bildes:", res.Status)
		return nil
	}

	// Erstelle oder öffne eine Datei, um das Bild zu speichern
	file, err := os.Create(imageName)
	if err != nil {
		fmt.Println("Fehler beim Erstellen der Datei:", err)
		return err
	}
	defer file.Close()

	// Kopiere den Response-Body in die Datei
	_, err = io.Copy(file, res.Body)
	if err != nil {
		fmt.Println("Fehler beim Schreiben des Bildes:", err)
		return err
	}

	fmt.Println("Bild erfolgreich heruntergeladen.")
	return nil
}

func generateName(value string) string {
	// Hol die aktuelle Zeit in Millisekunden
	milliseconds := time.Now().UnixNano() / int64(time.Millisecond)

	return fmt.Sprintf("%s-%d", value, milliseconds)
}
