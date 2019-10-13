package photos

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"google.golang.org/api/photoslibrary/v1"
)

const (
	apiVersion = "v1"
	basePath   = "https://photoslibrary.googleapis.com/"
	header     = "X-Goog-Upload-File-Name"
)

// Uploader is a client for uploading a media and extends the photoslibrary service.
// photoslibrary does not provide `/v1/uploads` API so we implement here.
type Uploader struct {
	*photoslibrary.Service

	client *http.Client
}

// NewUploader creates a new client.
func NewUploader(client *http.Client) (*Uploader, error) {
	photos, err := photoslibrary.New(client)
	if err != nil {
		return nil, err
	}

	return &Uploader{photos, client}, nil
}

// Upload sends the media and returns the UploadToken.
func (c *Uploader) Upload(file io.Reader, filename string) (string, error) {
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/%s/uploads", basePath, apiVersion), file)
	if err != nil {
		return "", err
	}
	req.Header.Add(header, filename)

	res, err := c.client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	uploadToken := string(b)
	return uploadToken, nil
}
