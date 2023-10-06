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
	"fmt"
	"io"
	"net/http"
	"os"
)

// DownloadFile télécharge un fichier à partir d'une URL `url` et le sauvegarde dans un fichier local `filePath`.
// Si `apiKey` n'est pas vide, il est ajouté à l'en-tête X-Api-Key de la requête.
func DownloadFile(filePath, url, apiKey string) error {
	// Créer le fichier local
	out, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Créer une nouvelle requête HTTP
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	// Ajouter la clé API à l'en-tête de la requête
	if apiKey != "" {
		req.Header.Add("x-api-key", apiKey)
	}

	// Exécuter la requête HTTP et obtenir la réponse
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Vérifier le code de statut de la réponse
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	// Copier le contenu dans le fichier local
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}
