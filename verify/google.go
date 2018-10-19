package verify

import (
	bytes "bytes"
	b64 "encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var cred Credentials
var config *oauth2.Config

const (
	InvalidTokenFormat = "10"
)

// Credentials which stores google ids.
type Credentials struct {
	Cid     string `json:"cid"`
	Csecret string `json:"csecret"`
}

type expirationTime int32

func (e *expirationTime) UnmarshalJSON(b []byte) error {
	var n json.Number
	err := json.Unmarshal(b, &n)
	if err != nil {
		return err
	}
	i, err := n.Int64()
	if err != nil {
		return err
	}
	*e = expirationTime(i)
	return nil
}

// tokenJSON is the struct representing the HTTP response from OAuth2
// providers returning a token in JSON form.
type tokenJSON struct {
	AccessToken  string         `json:"access_token"`
	TokenType    string         `json:"token_type"`
	RefreshToken string         `json:"refresh_token"`
	ExpiresIn    expirationTime `json:"expires_in"` // at least PayPal returns string, while most return number
	Expires      expirationTime `json:"expires"`    // broken Facebook spelling of expires_in
}

func (e *tokenJSON) expiry() (t time.Time) {
	if v := e.ExpiresIn; v != 0 {
		return time.Now().Add(time.Duration(v) * time.Second)
	}
	if v := e.Expires; v != 0 {
		return time.Now().Add(time.Duration(v) * time.Second)
	}
	return
}

// Token represents the credentials used to authorize
// the requests to access protected resources on the OAuth 2.0
// provider's backend.
//
// This type is a mirror of oauth2.Token and exists to break
// an otherwise-circular dependency. Other internal packages
// should convert this Token into an oauth2.Token before use.
type Token struct {
	// AccessToken is the token that authorizes and authenticates
	// the requests.
	AccessToken string

	// TokenType is the type of token.
	// The Type method returns either this or "Bearer", the default.
	TokenType string

	// RefreshToken is a token that's used by the application
	// (as opposed to the user) to refresh the access token
	// if it expires.
	RefreshToken string

	// Expiry is the optional expiration time of the access token.
	//
	// If zero, TokenSource implementations will reuse the same
	// token forever and RefreshToken or equivalent
	// mechanisms for that TokenSource will not be used.
	Expiry time.Time

	// Raw optionally contains extra metadata from the server
	// when updating a token.
	Raw interface{}
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

func VerifyToken(allBytes []byte) (*Token, error) {
	var fragments = bytes.Split(allBytes, []byte("."))

	if len(fragments) != 3 {
		return nil, InvalidTokenFormatErr
	}

	var jsonTokenBytes []byte = make([]byte, len(fragments[1]))
	var _, err = b64.StdEncoding.Decode(jsonTokenBytes, fragments[1])

	if err != nil {
		return nil, InvalidTokenFormatErr
	}

	return getToken(jsonTokenBytes)
}

func getToken(tokenBytes []byte) (*Token, error) {
	var tj tokenJSON
	var err error
	if err = json.Unmarshal(tokenBytes, &tj); err != nil {
		return nil, InvalidTokenFormatErr
	}
	var token = &Token{
		AccessToken:  tj.AccessToken,
		TokenType:    tj.TokenType,
		RefreshToken: tj.RefreshToken,
		Expiry:       tj.expiry(),
		Raw:          make(map[string]interface{}),
	}
	err = json.Unmarshal(tokenBytes, &token.Raw) // no error checks for optional fields

	if err != nil {
		return nil, InvalidTokenFormatErr
	}

	fmt.Println("AccessToken: " + token.AccessToken)
	fmt.Println("TokenType: " + token.TokenType)
	fmt.Println("RefreshToken: " + token.RefreshToken)
	fmt.Println("Expiry: " + token.Expiry.Format("2006-01-02T15:04:05.999999-07:00"))

	// if token.RefreshToken == "" {
	// 	token.RefreshToken = v.Get("refresh_token")
	// }
	if token.AccessToken == "" {
		return token, InvalidTokenFormatErr
	}

	return token, nil
}

// func verifyGoogleLogin(AccessToken string) {
// 	var token = oauth2.Token{
// 		AccessToken,
// 		TokenType: "Bearer",
// 	}

// }

// oauthCallbackHandler completes the OAuth flow, retreives the user's profile
// information and stores it in a session.
/*
func oauthCallbackHandler(w http.ResponseWriter, r *http.Request) *appError {
	oauthFlowSession, err := bookshelf.SessionStore.Get(r, r.FormValue("state"))
	if err != nil {
		return appErrorf(err, "invalid state parameter. try logging in again.")
	}

	redirectURL, ok := oauthFlowSession.Values[oauthFlowRedirectKey].(string)
	// Validate this callback request came from the app.
	if !ok {
		return appErrorf(err, "invalid state parameter. try logging in again.")
	}

	code := r.FormValue("code")
	tok, err := bookshelf.OAuthConfig.Exchange(context.Background(), code)
	if err != nil {
		return appErrorf(err, "could not get auth token: %v", err)
	}

	session, err := bookshelf.SessionStore.New(r, defaultSessionID)
	if err != nil {
		return appErrorf(err, "could not get default session: %v", err)
	}

	ctx := context.Background()
	profile, err := fetchProfile(ctx, tok)
	if err != nil {
		return appErrorf(err, "could not fetch Google profile: %v", err)
	}

	session.Values[oauthTokenSessionKey] = tok
	// Strip the profile to only the fields we need. Otherwise the struct is too big.
	session.Values[googleProfileSessionKey] = stripProfile(profile)
	if err := session.Save(r, w); err != nil {
		return appErrorf(err, "could not save session: %v", err)
	}

	http.Redirect(w, r, redirectURL, http.StatusFound)
	return nil
}
*/
