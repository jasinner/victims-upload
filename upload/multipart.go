package upload


import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"fmt"
	"github.com/victims/victims-common/log"
)

var uri string
const port = 8080
const path = "hash"

func New(hostname string){
	uri = fmt.Sprintf("http://%v:%v/%v", hostname, port, path)
	log.Logger.Infof("Using Java Hash Service at: %v", uri)
}

// Creates a new file upload http request from file path
func UploadRequest(paramName, path string) (*http.Request, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile(paramName, filepath.Base(path))
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(part, file)

	err = writer.Close()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", uri, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	return req, err
}
