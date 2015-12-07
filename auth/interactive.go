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

// Package auth provides Google oauth2 bindings for mantle.
package auth

import (
	"fmt"

	"github.com/coreos/mantle/Godeps/_workspace/src/github.com/coreos/pkg/capnslog"
	"github.com/coreos/mantle/Godeps/_workspace/src/golang.org/x/net/context"
	"github.com/coreos/mantle/Godeps/_workspace/src/golang.org/x/oauth2"
)

var plog = capnslog.NewPackageLogger("github.com/coreos/mantle", "auth")

// InteractiveTokenSource blah blah
func InteractiveTokenSource(ctx context.Context, conf *oauth2.Config, name string) (oauth2.TokenSource, error) {
	tokCache, err := UserTokenCache(name)
	if err != nil {
		return nil, err
	}

	tok, err := tokCache.Token()
	if err == nil {
		plog.Infof("Using cached OAuth token for %s", name)
		return tokCache.TokenSource(tok, conf.TokenSource(ctx, tok)), nil
	}
	plog.Infof("Failed to read cached OAuth token for %s: %v", name, err)

	url := conf.AuthCodeURL("state", oauth2.AccessTypeOffline)
	fmt.Printf("Visit the URL for the auth dialog: %v\n", url)
	fmt.Print("Enter token: ")

	var code string
	if _, err := fmt.Scan(&code); err != nil {
		return nil, err
	}
	tok, err = conf.Exchange(ctx, code)
	if err != nil {
		return nil, err
	}

	err = tokCache.Cache(tok)
	if err != nil {
		plog.Errorf("Failed to cache OAuth token for %s: %v", err)
	}
	return tokCache.TokenSource(tok, conf.TokenSource(ctx, tok)), err
}
