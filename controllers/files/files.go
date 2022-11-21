package files

import (
	"fmt"
	"helpers"
	"io"
	"mime/multipart"
	"net/http"
)

/* GetRequestFile gets a file from the request */
func GetRequestFile(field string, r *http.Request) (multipart.File, *multipart.FileHeader, error) {
	file, header, err := r.FormFile(field)

	if err != nil {
		return file, header, err
	}

	return file, header, nil
}

/* SetFileToResponse sets a file to the response */
func SetFileToResponse(filepath string, w http.ResponseWriter, isRemote bool) error {
	var file io.ReadCloser
	var err error

	if isRemote {
		file, err = helpers.GetFileRemote(filepath)
	} else {
		file, err = helpers.GetFileLocal(filepath)
	}

	defer func() {
		if err == nil {
			file.Close()
		}
	}()

	if err != nil {
		return fmt.Errorf("file not found")
	}

	_, err = io.Copy(w, file)

	if err != nil {
		return fmt.Errorf("error trying to get the file: %s", err.Error())
	}

	return nil
}
