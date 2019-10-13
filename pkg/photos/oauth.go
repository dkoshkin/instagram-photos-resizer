package photos

import (
	"context"
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"net/http"
	"os/exec"
	"runtime"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/photoslibrary/v1"
)

const (
	redirectURL = "urn:ietf:wg:oauth:2.0:oob"
)

// NewOAuthClient creates a new http.Client with a bearer access token
func NewOAuthClient(ctx context.Context, clientID string, clientSecret string) (*http.Client, error) {
	config := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Endpoint:     google.Endpoint,
		Scopes:       []string{photoslibrary.PhotoslibraryScope},
		RedirectURL:  redirectURL,
	}
	state, err := generateOAuthState()
	if err != nil {
		return nil, err
	}
	url := config.AuthCodeURL(state)

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		return nil, fmt.Errorf("could not open your browser: %v", err)
	}

	fmt.Print("Enter Code: ")

	var authCode string
	if _, err := fmt.Scanln(&authCode); err != nil {
		return nil, fmt.Errorf("could not read input: %v", err)
	}
	accessToken, err := config.Exchange(ctx, authCode)
	if err != nil {
		return nil, fmt.Errorf("could not get toke: %v", err)
	}
	return config.Client(ctx, accessToken), nil
}

func generateOAuthState() (string, error) {
	var n uint64
	if err := binary.Read(rand.Reader, binary.LittleEndian, &n); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", n), nil
}
