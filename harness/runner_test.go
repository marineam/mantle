// Copyright 2015 CoreOS, Inc.
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
	"testing"

	"github.com/coreos/mantle/Godeps/_workspace/src/github.com/kylelemons/godebug/pretty"
)

func TestListRunner(t *testing.T) {
	// test1 is defined in group_test.go
	ts := &Serial{Tests: []interface{}{test1{}}}
	expect := []string{"test1.TestFunc1"}
	funcs := ts.List()
	if diff := pretty.Compare(expect, funcs); diff != "" {
		t.Errorf("List failed:\n%s", diff)
	}
	ts = &Serial{Tests: []interface{}{&test1{}}}
	expect = []string{"test1.TestFunc1", "test1.TestFunc2"}
	funcs = ts.List()
	if diff := pretty.Compare(expect, funcs); diff != "" {
		t.Errorf("List failed:\n%s", diff)
	}
}

func TestRunRunner(t *testing.T) {
	// test2 is defined in group_test.go
	tc := &test2{TestCalled: false}
	ts := &Serial{Tests: []interface{}{tc}}
	if err := ts.Run("*"); err != nil {
		t.Errorf("Run failed: %v", err)
	}
	if !tc.TestCalled {
		t.Error("test2.Test never called")
	}
}
