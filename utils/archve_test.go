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
	"os"
	"testing"
)

const SRCPATH = "/bin"
const DESTFILE = "/temp/test.zip"

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
	err := os.RemoveAll("/tmp/test_unpacked")
	archiveManager := NewZipArchiveManager()
	err = archiveManager.Unpack(DESTFILE, "/tmp/test_unpacked")
	if err != nil {
		t.Error(err)
	}
}
