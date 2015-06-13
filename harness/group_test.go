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

type test1 struct{ BaseGroup }

func (t test1) TestFunc1(h *H)  {}
func (t *test1) TestFunc2(h *H) {}
func (t test1) NotATest()       {}

func TestListGroup(t *testing.T) {
	expect := []string{"test1.TestFunc1"}
	funcs := ListGroup(test1{})
	if diff := pretty.Compare(expect, funcs); diff != "" {
		t.Errorf("ListTests failed:\n%s", diff)
	}
	expect = []string{"test1.TestFunc1", "test1.TestFunc2"}
	funcs = ListGroup(&test1{})
	if diff := pretty.Compare(expect, funcs); diff != "" {
		t.Errorf("ListTests failed:\n%s", diff)
	}
}

type test2 struct {
	BaseGroup
	TestCalled bool
}

func (t *test2) Test(h *H) { t.TestCalled = true }

func TestRunGroup(t *testing.T) {
	tc := &test2{TestCalled: false}
	h := NewGroupHarness(tc, "test2.Test")
	<-h.Start()
	if h.Failed() || h.Skipped() || !tc.TestCalled {
		t.Error("Running Group.Test failed")
	}
}
