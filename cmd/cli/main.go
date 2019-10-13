package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"

	"github.com/fatih/color"

	"google.golang.org/api/photoslibrary/v1"

	"github.com/dkoshkin/instagram-photos-resizer/pkg/instagram"
	"github.com/dkoshkin/instagram-photos-resizer/pkg/photos"
)

func main() {
	clientID := os.Getenv("GOOGLE_CLIENT_ID")
	clientSecret := os.Getenv("GOOGLE_CLIENT_SECRET")
	if clientID == "" || clientSecret == "" {
		log.Fatal(`Error: GOOGLE_CLIENT_ID and GOOGLE_CLIENT_SECRET must be set.
1. Open https://console.cloud.google.com/apis/credentials
2. Create an OAuth client ID where the application type is other.
3. Set the following environment variables:
export GOOGLE_CLIENT_ID=
export GOOGLE_CLIENT_SECRET=
`)
	}
	if len(os.Args) < 2 {
		log.Fatalf("usage: %s FILES...", os.Args[0])
	}

	ctx := context.Background()
	client, err := photos.NewOAuthClient(ctx, clientID, clientSecret)
	if err != nil {
		log.Fatal(err)
	}

	photos, err := photos.NewUploader(client)
	if err != nil {
		log.Fatal(err)
	}

	files := os.Args[1:]
	fmt.Printf("Uplading %d File(s)...\n", len(files))
	for _, filepath := range files {
		fmt.Printf("  Processing %q\n", filepath)
		b, err := ioutil.ReadFile(filepath)
		if err != nil {
			log.Fatal(err)
		}
		prepared, err := instagram.Prepare(bytes.NewReader(b))
		if err != nil {
			log.Fatal(err)
		}

		// upload in reverse order so it shows up in the correct order in Google Photos
		for n := len(prepared.Reader) - 1; n >= 0; n-- {
			filename := fmt.Sprintf("resized-%d-%s", n, path.Base(filepath))
			err := upload(photos, prepared.Reader[n], filename)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}

func upload(photos *photos.Uploader, file io.Reader, filename string) error {
	uploadToken, err := photos.Upload(file, filename)
	if err != nil {
		return fmt.Errorf("could not get upload token: %v", err)
	}

	batch, err := photos.MediaItems.BatchCreate(&photoslibrary.BatchCreateMediaItemsRequest{
		NewMediaItems: []*photoslibrary.NewMediaItem{
			{
				Description:     filename,
				SimpleMediaItem: &photoslibrary.SimpleMediaItem{UploadToken: uploadToken},
			},
		},
	}).Do()
	if err != nil {
		return fmt.Errorf("could not upload file: %v", err)
	}

	for _, result := range batch.NewMediaItemResults {
		green := color.New(color.FgGreen).SprintFunc()
		fmt.Printf("  - %s: %s\n", result.MediaItem.Description, green("OK"))
	}

	return nil
}
