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
