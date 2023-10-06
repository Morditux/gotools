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
	"archive/zip"
	"compress/flate"
	"io"
	"os"
	"path/filepath"
)

type ArchiveManager interface {
	Pack(srcPath string, destFile string, baseDir bool) error
	Unpack(srcFile string, destPath string) error
}

type ZipArchiveManager struct {
}

func NewZipArchiveManager() *ZipArchiveManager {
	return &ZipArchiveManager{}
}

// Pack creates a zip archive from the given source path and saves it to the given destination file. If baseDir is false, the source path is not included in the archive.
func (z *ZipArchiveManager) Pack(src string, dest string, baseDir bool) error {
	zipFile, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	zipWriter.RegisterCompressor(zip.Deflate, func(out io.Writer) (io.WriteCloser, error) {
		return flate.NewWriter(out, flate.BestSpeed)
	})

	defer zipWriter.Close()

	err = filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Pour éviter d'inclure le répertoire source lui-même
		if path == src {
			return nil
		}

		// Si c'est un répertoire, on ne fait rien
		if info.IsDir() {
			return nil
		}

		// Si c'est un lien symbolique, on ne fait rien
		if info.Mode()&os.ModeSymlink != 0 {
			return nil
		}
		// Ouvrir le fichier à ajouter dans l'archive
		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		// Créer une entrée dans l'archive zip
		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		// Ajuster le chemin de l'entrée dans l'archive
		if baseDir {
			header.Name = filepath.Join(filepath.Base(src), path[len(src):])
		} else {
			header.Name = path[len(src):]
		}
		header.Method = zip.Deflate
		// Ajouter l'entrée dans l'archive zip
		writer, err := zipWriter.CreateHeader(header)
		if err != nil {
			return err
		}

		// Copier le contenu du fichier dans l'entrée de l'archive
		_, err = io.Copy(writer, file)
		return err
	})

	return err
}

// Unpack extracts the given zip archive to the given destination path.
func (z *ZipArchiveManager) Unpack(src string, dest string) error {
	// Ouvrir l'archive ZIP
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer r.Close()

	// Parcourir chaque fichier dans l'archive ZIP
	for _, f := range r.File {
		// Construire le chemin complet pour le fichier/dossier de sortie
		fpath := filepath.Join(dest, f.Name)

		// Si c'est un dossier, le créer et passer au fichier suivant
		if f.FileInfo().IsDir() {
			os.MkdirAll(fpath, os.ModePerm)
			continue
		}

		// S'assurer que le répertoire parent du fichier existe
		if err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			return err
		}

		// Extraire le fichier
		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}

		rc, err := f.Open()
		if err != nil {
			outFile.Close()
			return err
		}

		_, err = io.Copy(outFile, rc)

		// Fermer les fichiers sans ignorer les erreurs
		rc.Close()
		outFile.Close()

		if err != nil {
			return err
		}
	}
	return nil
}
