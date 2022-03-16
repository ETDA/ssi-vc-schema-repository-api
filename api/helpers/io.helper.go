package helpers

import (
	"bytes"
	core "ssi-gitlab.teda.th/ssi/core"
	"io"
	"mime/multipart"
)

func MultiPartFileToIFileHeader(fileHeader *multipart.FileHeader) (core.IFile, error) {
	rawFile, err := fileHeader.Open()
	if err != nil {
		return nil, err
	}

	buffer := new(bytes.Buffer)
	_, err = io.Copy(buffer, rawFile)
	if err != nil {
		return nil, err
	}

	return core.NewFile(fileHeader.Filename, buffer.Bytes()), nil
}

func IOReaderToBytes(reader io.Reader) ([]byte, error) {
	buffer := new(bytes.Buffer)
	_, err := io.Copy(buffer, reader)
	if err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}
