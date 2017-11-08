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

const port = 8080
const path = "hash"
const key = "JAVA_SERVICE_HOST"

var hostName string
var uri string

func SetHostname(name string) {
	hostName = name
}

func lookupHostname() {
	// If the hostName is empty ...
	if hostName == "" {
		// then look it up via environment variable
		var ok bool
		hostName, ok = os.LookupEnv(key)
		if !ok {
			// ... and kill execution if we can't get the variable
			log.Logger.Panicf("The environment variable %s must be set!", key)
		}
	}
	uri = fmt.Sprintf("http://%v:%v/%v", hostName, port, path)
	log.Logger.Infof("Using Java Hash Service at: %s", uri)
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

	lookupHostname()

	req, err := http.NewRequest("POST", uri, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	return req, err
}
