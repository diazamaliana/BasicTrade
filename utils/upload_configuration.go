package utils

import (
	"basictrade/helpers"
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"path"
	"strings"
	"time"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

const (
	// MaxFileSize is the maximum allowed file size in bytes (2 MB in this example)
	MaxFileSize = 2 * 1024 * 1024
)

func UploadFile(fileHeader *multipart.FileHeader, fileName string) (string, error) {
	c, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Check if the file size exceeds the maximum allowed size
	if fileHeader.Size > MaxFileSize {
		return "", errors.New("File size exceeds the maximum allowed size!")
	}

	// Check if the file is an image based on its content type
	if !isImageFile(fileHeader.Header.Get("Content-Type")) {
		return "", errors.New("File is not an image!")
	}

	// Add Cloudinary product environment credentials.
	cld, err := cloudinary.NewFromParams(helpers.EnvCloudName(), helpers.EnvCloudAPIKey(), helpers.EnvCloudAPISecret())
	if err != nil {
		return "", err
	}
	fmt.Println(helpers.EnvCloudName())

	// Convert file
	fileReader, err := convertFile(fileHeader)
	if err != nil {
		return "", err
	}

	// Upload file
	uploadParam, err := cld.Upload.Upload(c, fileReader, uploader.UploadParams{
		PublicID: fileName,
		Folder:   helpers.EnvCloudUploadFolder(),
	})
	if err != nil {
		return "", err
	}

	return uploadParam.SecureURL, nil
}

func convertFile(fileHeader *multipart.FileHeader) (*bytes.Reader, error) {
    if fileHeader == nil {
        return nil, errors.New("File header is nil.")
    }

    file, err := fileHeader.Open()
    if err != nil {
        return nil, err
    }
    defer file.Close()

    // Read the file content into an in-memory buffer
    buffer := new(bytes.Buffer)
    if _, err := io.Copy(buffer, file); err != nil {
        return nil, err
    }

    // Create a bytes.Reader from the buffer
    fileReader := bytes.NewReader(buffer.Bytes())
    return fileReader, nil
}


func RemoveExtension(filename string) string {
	return path.Base(filename[:len(filename)-len(path.Ext(filename))])
}

// isImageFile checks if the given content type corresponds to an image file
func isImageFile(contentType string) bool {
	return strings.HasPrefix(contentType, "image/")
}
