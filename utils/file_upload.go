/*
Copyright 2023 Kavoos Bojnourdi
Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation
files (the “Software”), to deal in the Software without restriction, including without limitation the rights to use, copy,
modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the
Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED “AS IS”, WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE
WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR
OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/
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
