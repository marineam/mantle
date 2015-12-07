// Copyright 2014 CoreOS, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package google provides support for Google's various OAuth mechanisms.
package google

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/coreos/mantle/Godeps/_workspace/src/github.com/spf13/pflag"
	"github.com/coreos/mantle/Godeps/_workspace/src/golang.org/x/net/context"
	"github.com/coreos/mantle/Godeps/_workspace/src/golang.org/x/oauth2"
	"github.com/coreos/mantle/Godeps/_workspace/src/golang.org/x/oauth2/google"

	"github.com/coreos/mantle/auth"
)

// FlagSet must be added to any cli tools that access Google services.
var (
	FlagSet     *pflag.FlagSet
	authSource  string // --google-auth-source=
	jsonPath    string // --google-service-json=
	compAccount string // --google-compute-account=

	jsonMissing = errors.New("Google auth JSON file unspecified, use --google-service-json=")

	scopes = []string{
		"https://www.googleapis.com/auth/compute",
		"https://www.googleapis.com/auth/devstorage.full_control",
	}

	// client registered under 'marineam-tools'
	interactive = oauth2.Config{
		ClientID:     "937427706989-nbndmfkp0knqardoagk6lbcamrsh828i.apps.googleusercontent.com",
		ClientSecret: "F6Xs5wGHZzGw-QFXl3aylLUT",
		Endpoint:     google.Endpoint,
		RedirectURL:  "urn:ietf:wg:oauth:2.0:oob",
		Scopes:       scopes,
	}
)

func init() {
	FlagSet = pflag.NewFlagSet("google", pflag.ExitOnError)
	FlagSet.StringVar(&authSource, "google-auth-source",
		"interactive",
		"OAuth authentication type: interactive json compute")
	FlagSet.StringVar(&jsonPath, "google-service-json",
		os.Getenv("GOOGLE_APPLICATION_CREDENTIALS"),
		"Path to a JSON file defining a Google service account")
	FlagSet.StringVar(&compAccount, "google-compute-account",
		"default",
		"Name of a Google service account available on GCE")
}

func jsonTokenSource(ctx context.Context) (oauth2.TokenSource, error) {
	if jsonPath == "" {
		return nil, jsonMissing
	}
	jsonData, err := ioutil.ReadFile(jsonPath)
	if err != nil {
		return nil, err
	}
	jwtConf, err := google.JWTConfigFromJSON(jsonData, scopes...)
	if err != nil {
		return nil, err
	}
	return jwtConf.TokenSource(ctx), nil
}

func TokenSource(ctx context.Context) (oauth2.TokenSource, error) {
	switch authSource {
	case "interactive":
		return auth.InteractiveTokenSource(ctx, &interactive, "google")
	case "json":
		return jsonTokenSource(ctx)
	case "compute":
		return google.ComputeTokenSource(compAccount), nil
	default:
		return nil, fmt.Errorf("Invalid Google authentication source %q", authSource)
	}
}

func NewClient(ctx context.Context) (*http.Client, error) {
	src, err := TokenSource(ctx)
	if err != nil {
		return nil, err
	}
	return oauth2.NewClient(ctx, src), nil
}
