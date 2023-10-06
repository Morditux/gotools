package utils

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
)

// UploadFile envoie un fichier spécifié par `filePath` à une URL `url` via HTTPS.
// Si `key` n'est pas vide, il est ajouté à l'en-tête X-Api-Key de la requête.
func UploadFile(url, filePath string, key string) error {
	// Ouvrir le fichier
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Créer un tampon pour stocker la partie multipart de la requête
	var requestBody bytes.Buffer
	multiPartWriter := multipart.NewWriter(&requestBody)

	// Ajouter le fichier à la partie multipart
	fileWriter, err := multiPartWriter.CreateFormFile("file", file.Name())
	if err != nil {
		return err
	}
	_, err = io.Copy(fileWriter, file)
	if err != nil {
		return err
	}

	// Fermer la partie multipart
	err = multiPartWriter.Close()
	if err != nil {
		return err
	}

	// Créer une requête HTTP POST avec le fichier en tant que corps
	req, err := http.NewRequest("POST", url, &requestBody)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", multiPartWriter.FormDataContentType())
	if key != "" {
		req.Header.Set("X-Api-Key", key)
	}
	// Envoyer la requête
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Vérifier la réponse
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	return nil
}
