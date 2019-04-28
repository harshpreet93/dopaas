package do_auth

import (
	"context"
	"github.com/digitalocean/godo"
	"golang.org/x/oauth2"
)

//const (
//	pat = "mytoken"
//)

type TokenSource struct {
	AccessToken string
}

func (t *TokenSource) Token() (*oauth2.Token, error) {
	token := &oauth2.Token{
		AccessToken: t.AccessToken,
	}
	return token, nil
}

func getTokenFromFS() string {
	// get token from ~/.dopaas.yaml
	// file format is:
	// DIGITALOCEAN_ACCESS_TOKEN: "blahhhh"
	return ""
}

func auth() *godo.Client {
	pat := getTokenFromFS()
	tokenSource := &TokenSource{
		AccessToken: pat,
	}

	oauthClient := oauth2.NewClient(context.Background(), tokenSource)
	client := godo.NewClient(oauthClient)
	return client
}
