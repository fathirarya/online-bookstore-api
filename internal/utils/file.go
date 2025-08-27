package utils

import (
	"encoding/base64"
	"io"
	"mime/multipart"
)

func FileToBase64(file multipart.File) (string, error) {
	defer file.Close()

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(fileBytes), nil
}
