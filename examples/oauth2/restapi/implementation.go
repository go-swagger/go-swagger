package restapi

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	oidc "github.com/coreos/go-oidc"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
)

var (
	state = "foobar" // Don't do this in production.

	clientID     = ""
	clientSecret = ""
	issuer       = "https://accounts.google.com"
	authURL      = "https://accounts.google.com/o/oauth2/v2/auth"
	tokenURL     = "https://www.googleapis.com/oauth2/v4/token"
	userInfoURL  = "https://www.googleapis.com/oauth2/v3/userinfo"
	callbackURL  = "http://127.0.0.1:12345/api/auth/callback"

	endpoint = oauth2.Endpoint{
		AuthURL:  authURL,
		TokenURL: tokenURL,
	}
	config = oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Endpoint:     endpoint,
		RedirectURL:  callbackURL,
		Scopes:       []string{oidc.ScopeOpenID, "profile", "email"},
	}
)

func login(r *http.Request) string {
	var accessToken string
	http.Redirect(wG, r, config.AuthCodeURL(state), http.StatusFound)
	return accessToken
}

func callback(r *http.Request) (string, error) {
	if r.URL.Query().Get("state") != state {
		log.Println("state did not match")
		return "", fmt.Errorf("state did not match")
	}

	myClient := &http.Client{}

	parentContext := context.Background()
	ctx := oidc.ClientContext(parentContext, myClient)

	authCode := r.URL.Query().Get("code")
	log.Printf("Authorization code: %v\n", authCode)
	oauth2Token, err := config.Exchange(ctx, authCode)
	if err != nil {
		log.Println("failed to exchange token", err.Error())
		return "", fmt.Errorf("failed to exchange token")
	}
	log.Println("Raw token data:", oauth2Token)
	return oauth2Token.AccessToken, nil
}

func authenticated(token string) (bool, error) {
	// validate the token
	bearToken := "Bearer " + token
	req, err := http.NewRequest("GET", userInfoURL, nil)
	if err != nil {
		return false, fmt.Errorf("http request: %v", err)
	}

	req.Header.Add("Authorization", bearToken)
	cli := &http.Client{}
	resp, err := cli.Do(req)
	if err != nil {
		return false, fmt.Errorf("http request: %v", err)
	}
	defer resp.Body.Close()

	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return false, fmt.Errorf("fail to get response: %v", err)
	}
	if resp.StatusCode != 200 {
		return false, nil
	}
	return true, nil
}
