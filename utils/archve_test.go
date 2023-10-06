package utils

import (
	"os"
	"testing"
)

const SRCPATH = "/home/mordicus/opt/panacee"
const DESTFILE = "/home/mordicus/tmp/panacee.zip"

func TestNewZipArchiveManager(t *testing.T) {
	archiveManager := NewZipArchiveManager()
	if archiveManager == nil {
		t.Error("NewZipArchiveManager() should not return nil")
	}
}

func TestZipArchiveManager_Pack(t *testing.T) {
	// Efface le fichier de destination s'il existe
	err := os.Remove(DESTFILE)
	archiveManager := NewZipArchiveManager()
	err = archiveManager.Pack(SRCPATH, DESTFILE, false)
	if err != nil {
		t.Error(err)
	}
}

func TestZipArchiveManager_Unpack(t *testing.T) {
	// Efface le répertoire de destination s'il existe
	err := os.RemoveAll(SRCPATH + "_unpacked")
	archiveManager := NewZipArchiveManager()
	err = archiveManager.Unpack(DESTFILE, SRCPATH+"_unpacked")
	if err != nil {
		t.Error(err)
	}
}
