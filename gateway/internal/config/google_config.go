package config

import (
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type GoogleAuth struct {
	Config oauth2.Config
}

var Google GoogleAuth

func GoogleConfig() oauth2.Config {
	Google.Config = oauth2.Config{
		RedirectURL:  "http://127.0.0.1:5050/user/google_callback",
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		Scopes: []string{"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile"},
		Endpoint: google.Endpoint,
	}

	return Google.Config
}
