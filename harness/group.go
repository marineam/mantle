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
	"reflect"
	"strings"
)

// Group is a set of related tests that share common set-up and tear-down code.
// Similar to stand-alone test functions, test methods in a Group are detected
// by using a Test prefix.
type Group interface {
	// Prepare executes before all test methods in this group. If Prepare
	// calls H.Skip or H.Fail the rest of the group is skipped.
	Prepare(*H)
	// Test methods may have arbitrary names beginning with the string
	// returned by Prefix, e.g. ``Test''.
	//   Test(*H)
	// Cleanup executes after all test methods.
	Cleanup(*H)
	// Prefix returns identifier prefix shared by all test methods.
	Prefix() string
}

// BaseGroup provides a trivial Group that does nothing. Can be embedded to
// avoid defining empty Prepare and/or Cleanup methods. The default prefix
// is ``Test''.
type BaseGroup struct{}

func (b BaseGroup) Prepare(h *H)   {}
func (b BaseGroup) Cleanup(h *H)   {}
func (b BaseGroup) Prefix() string { return "Test" }

// isTestMethod checks if m has the right name and signature to be a test.
func isTestMethod(m *reflect.Method, prefix string) bool {
	if !HasIdentifierPrefix(m.Name, prefix) {
		return false
	}

	if m.Type.NumIn() != 2 || m.Type.NumOut() != 0 {
		return false
	}

	tcType := reflect.TypeOf((*Group)(nil)).Elem()
	if !m.Type.In(0).Implements(tcType) {
		return false
	}

	hType := reflect.TypeOf(&H{})
	if !m.Type.In(1).AssignableTo(hType) {
		return false
	}

	return true
}

// ListGroup returns a list of tests with a given prefix it contains.
func ListGroup(g Group) []string {
	group := TestName(g) + "."
	typ := reflect.TypeOf(g)
	tests := make([]string, 0, typ.NumMethod())
	for i := 0; i < typ.NumMethod(); i++ {
		m := typ.Method(i)
		if isTestMethod(&m, g.Prefix()) {
			tests = append(tests, group+m.Name)
		}
	}
	return tests
}

// NewGroupHarness returns a harness for running the given test name.
func NewGroupHarness(tc Group, name string) *H {
	var methodName string
	i := strings.Index(name, ".")
	// A test name without a '.' is assumed to just be the method name.
	if i < 0 {
		methodName = name
		name = TestName(tc) + "." + methodName
	} else {
		methodName = name[i+1:]
	}

	// Create a test function that can call tc.methodName(*H)
	m := reflect.ValueOf(tc).MethodByName(methodName)
	wrap := func(h *H) {
		m.Call([]reflect.Value{reflect.ValueOf(h)})
	}

	return &H{
		name: name,
		test: wrap,
	}
}
