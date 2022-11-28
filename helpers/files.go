package helpers

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

/* UploadFileLocal Uploads a file to the local server */
func UploadFileLocal(filepath string, file multipart.File) error {
	f, err := os.OpenFile(filepath, os.O_WRONLY|os.O_CREATE, 0666)

	defer func() {
		if f != nil {
			f.Close()
		}
	}()

	if err != nil {
		return fmt.Errorf("error uploading the file: %s", err.Error())
	}

	_, err = io.Copy(f, file)

	if err != nil {
		return fmt.Errorf("error creating the file: %s", err.Error())
	}

	return nil
}

/* UploadFileRemote Uploads a file to a remote server */
func UploadRemote(file multipart.File, publicId string) (string, error) {
	cld, err := cloudinary.New()

	if err != nil {
		return "", fmt.Errorf("failed to intialize Cloudinary: %s", err.Error())
	}

	ctx := context.Background()
	uploadResult, err := cld.Upload.Upload(ctx, file, uploader.UploadParams{PublicID: publicId})

	if err != nil {
		return "", fmt.Errorf("failed to upload file: %s", err.Error())
	}

	return uploadResult.SecureURL, nil
}

/* DestroyRemote Destroys a file in the remote server */
func DestroyRemote(publicId string) error {
	cld, err := cloudinary.New()

	if err != nil {
		return fmt.Errorf("failed to intialize Cloudinary: %s", err.Error())
	}

	ctx := context.Background()
	_, err = cld.Upload.Destroy(ctx, uploader.DestroyParams{PublicID: publicId})

	if err != nil {
		return fmt.Errorf("failed to destroy file: %s", err.Error())
	}

	return nil
}

/* GetFileLocal Gets a file form the local server */
func GetFileLocal(filepath string) (io.ReadCloser, error) {
	f, err := os.Open(filepath)

	return f, err
}

/* GetFileRemote Gets a file form a remote server */
func GetFileRemote(fileUrl string) (io.ReadCloser, error) {
	response, err := http.Get(fileUrl)

	if err != nil {
		return nil, err
	}

	if response.StatusCode >= 400 {
		return nil, fmt.Errorf("file not found in the remote server")
	}

	return response.Body, nil
}
