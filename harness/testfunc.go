// Copyright 2015 CoreOS, Inc.
// Copyright 2011 The Go Authors.
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

package harness

import (
	"reflect"
	"strings"
	"unicode"
	"unicode/utf8"
)

// TestFunc is a stand alone test, similar to normal Go tests.
type TestFunc func(*H)

// Creates a test harness for running the given function.
func NewTestFuncHarness(tf TestFunc) *H {
	return &H{
		name: TestName(tf),
		test: tf,
	}
}

// TestName returns the name of the given test function or group.
func TestName(t interface{}) string {
	typ := reflect.TypeOf(t)
	if typ.Kind() == reflect.Ptr {
		return typ.Elem().Name()
	} else {
		return typ.Name()
	}
}

// HasIdentifierPrefix tells whether name has a given word prefix, e.g.
// ``Test''. It is a Test if there is a character after Test that is not
// a lower-case letter. We don't want TesticularCancer.
func HasIdentifierPrefix(name, prefix string) bool {
	// Pulled directly from Go's src/cmd/go/test.go
	if !strings.HasPrefix(name, prefix) {
		return false
	}
	if len(name) == len(prefix) { // "Test" is ok
		return true
	}
	c, _ := utf8.DecodeRuneInString(name[len(prefix):])
	return !unicode.IsLower(c)
}
