package verify

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"os"
	token "votecube-id/verify/token"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var cred Credentials
var config *oauth2.Config

const (
	InvalidTokenFormat = "10"
)

var audience = make([]string, 1)

// Credentials which stores google ids.
type Credentials struct {
	Cid     string `json:"cid"`
	Csecret string `json:"csecret"`
}

var InvalidTokenFormatErr = errors.New(InvalidTokenFormat)

// var (
// 	OAuthConfig *oauth2.Config
// 	Endpoint    = google.Endpoint
// )

func ConfigureOAuthClient(clientID, clientSecret string) *oauth2.Config {
	redirectURL := os.Getenv("OAUTH2_CALLBACK")
	if redirectURL == "" {
		redirectURL = "http://localhost:8080/go/s"
	}

	audience[0] = "246713460433-amkv0f52nfalrhm5s7r2jgoqup6jtt4m.apps.googleusercontent.com"

	return &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURL,
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
		Endpoint:     google.Endpoint,
	}
}

func SetConfig() {
	file, err := ioutil.ReadFile("./creds.json")
	if err != nil {
		log.Printf("File error: %v\n", err)
		os.Exit(1)
	}
	json.Unmarshal(file, &cred)

	config = ConfigureOAuthClient(
		cred.Cid,
		cred.Csecret)
}

func VerifyToken(allBytes []byte) (*token.ClaimSet, error) {
	var base64Token = string(allBytes)

	return token.VerifyIDToken(base64Token, audience)
}
