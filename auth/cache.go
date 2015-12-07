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

package auth

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"

	"github.com/coreos/mantle/Godeps/_workspace/src/golang.org/x/oauth2"
)

type TokenCache string

type cachedTokenSource struct {
	cache  TokenCache
	source oauth2.TokenSource
	token  *oauth2.Token
}

// Token returns the token from the wrapped TokenSource, writing it to
// disk if it has changed.
func (c *cachedTokenSource) Token() (tok *oauth2.Token, err error) {
	if tok, err = c.source.Token(); err != nil {
		return
	}
	if c.token == nil || *c.token != *tok {
		if err = c.cache.Cache(tok); err != nil {
			return
		}
		c.token = tok
	}
	return
}

// Cache saves the given token to disk. Update is atomic.
func (t TokenCache) Cache(tok *oauth2.Token) error {
	file, err := ioutil.TempFile(
		filepath.Dir(string(t)), filepath.Base(string(t)))
	if err != nil {
		return err
	}

	if err := json.NewEncoder(file).Encode(tok); err != nil {
		os.Remove(file.Name())
		file.Close()
		return err
	}

	if err := file.Close(); err != nil {
		os.Remove(file.Name())
		return err
	}

	if err := os.Rename(file.Name(), string(t)); err != nil {
		os.Remove(file.Name())
		return err
	}

	return nil
}

// Token reads the cached token from disk.
func (t TokenCache) Token() (*oauth2.Token, error) {
	file, err := os.Open(string(t))
	if err != nil {
		return nil, err
	}
	defer file.Close()

	tok := &oauth2.Token{}
	if err := json.NewDecoder(file).Decode(tok); err != nil {
		return nil, err
	}

	return tok, nil
}

// TokenSource wraps the provides source and caches refreshed tokens to disk.
func (t TokenCache) TokenSource(tok *oauth2.Token, src oauth2.TokenSource) oauth2.TokenSource {
	return &cachedTokenSource{
		cache:  t,
		source: src,
		token:  tok,
	}
}

// UserTokenCache is a simple TokenSource that only reads and writes to disk.
func UserTokenCache(name string) (TokenCache, error) {
	userInfo, err := user.Current()
	if err != nil {
		return TokenCache(""), err
	}
	return TokenCache(filepath.Join(userInfo.HomeDir,
		fmt.Sprintf(".mantle-cache-%s.json", name))), nil
}
