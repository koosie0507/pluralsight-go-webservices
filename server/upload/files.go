package upload

import (
	"io/ioutil"
	"path/filepath"
	"time"
)

// UploadsDir is the path to the directory where this service keeps uploaded files.
var UploadsDir string = filepath.Join("uploads")

// Upload items store minimal information about files being uploaded to the server.
type Upload struct {
	Name       string    `json:"name"`
	UploadDate time.Time `json:"uploadDate"`
}

// GetUploads retrieves a list of files that were uploaded to the server.
func GetUploads() ([]Upload, error) {
	result := make([]Upload, 0)
	files, err := ioutil.ReadDir(UploadsDir)
	if err != nil {
		return nil, err
	}
	for _, f := range files {
		result = append(result, Upload{Name: f.Name(), UploadDate: f.ModTime()})
	}
	return result, nil
}
