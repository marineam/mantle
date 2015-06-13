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
	"errors"
	"fmt"
	"path/filepath"
	"sort"
)

var RunnerFailed = errors.New("Failure while running tests!")

type Runner interface {
	List() []string
	Run(match string) error
}

type Serial struct {
	Tests   []interface{}
	failed  int
	passed  int
	skipped int
}

func (s *Serial) List() []string {
	tests := make([]string, 0, len(s.Tests))
	for _, value := range s.Tests {
		switch t := value.(type) {
		case Group:
			tests = append(tests, ListGroup(t)...)
		case TestFunc:
			tests = append(tests, TestName(t))
		default:
			panic("Unexpected test type")
		}
	}
	sort.Strings(tests)
	return tests
}

func (s *Serial) Run(match string) error {
	for _, value := range s.Tests {
		switch t := value.(type) {
		case Group:
			s.runGroup(t, match)
		case TestFunc:
			s.runTestFunc(t, match)
		default:
			panic("Unexpected test type")
		}
	}

	fmt.Printf("FAILED:  %d\n", s.failed)
	fmt.Printf("PASSED:  %d\n", s.passed)
	fmt.Printf("SKIPPED: %d\n", s.skipped)

	if s.failed != 0 {
		return RunnerFailed
	}
	return nil
}

func (s *Serial) runGroup(tc Group, match string) {
	names := ListGroup(tc)
	tests := make([]*H, 0, len(names))
	for _, name := range names {
		if ok, _ := filepath.Match(match, name); ok {
			tests = append(tests, NewGroupHarness(tc, name))
		}
	}

	if len(tests) == 0 {
		return
	}

	prepare := NewGroupHarness(tc, "Prepare")
	<-prepare.Start()
	prepare.report()

	if !prepare.Failed() && !prepare.Skipped() {
		cleanup := NewGroupHarness(tc, "Cleanup")
		defer func() {
			<-cleanup.Start()
			cleanup.report()
		}()
	}

	for _, test := range tests {
		if prepare.Failed() {
			test.Fail()
		} else if prepare.Skipped() {
			test.Skip()
		} else {
			<-test.Start()
		}
		s.report(test)
	}
}

func (s *Serial) runTestFunc(tf TestFunc, match string) {
	name := TestName(tf)
	if ok, _ := filepath.Match(match, name); !ok {
		return
	}

	test := NewTestFuncHarness(tf)
	<-test.Start()
	s.report(test)
}

func (s *Serial) report(h *H) {
	if h.Failed() {
		s.failed++
	} else if h.Skipped() {
		s.skipped++
	} else {
		s.passed++
	}
	h.report()
}
