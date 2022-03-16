package helpers

import (
	"archive/zip"
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func Unzip(b []byte) (*zip.Reader, error) {
	buf := bytes.NewReader(b)
	fileSystem, err := zip.NewReader(buf, buf.Size())
	if err != nil {
		return nil, err
	}

	return fileSystem, nil
}

func ZipFromDirectory(path string) ([]byte, error) {
	buffer := bytes.NewBuffer(nil)
	myZip := zip.NewWriter(buffer)
	err := filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		if err != nil {
			return err
		}
		relPath := strings.TrimPrefix(filePath, filepath.Dir(path))
		zipFile, err := myZip.Create(strings.TrimPrefix(relPath, "/"))
		if err != nil {
			return err
		}
		fsFile, err := os.Open(filePath)
		if err != nil {
			return err
		}
		_, err = io.Copy(zipFile, fsFile)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	err = myZip.Close()
	if err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}
