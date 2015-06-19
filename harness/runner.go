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
	"path/filepath"
	"sort"
	"sync"
)

var RunnerFailed = errors.New("Failure while running tests!")

/*
type Runner interface {
	List() []string
	Run(match string) error
}*/

type Local struct {
	Tests []interface{}

	parallel chan struct{}  // Limits number of concurrent tests.
	pending  sync.WaitGroup // Wait outstanding tests.
	failed   int
	passed   int
	skipped  int
}

func (s *Local) List() []string {
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

func (s *Local) Run(match string, parallel int) error {
	s.parallel = make(chan struct{}, parallel)
	for _, value := range s.Tests {
		s.pending.Add(1)
		switch t := value.(type) {
		case Group:
			go s.runGroup(t, match)
		case TestFunc:
			go s.runTestFunc(t, match)
		default:
			panic("Unexpected test type")
		}
	}

	s.pending.Wait()
	plog.Noticef("FAILED:  %d", s.failed)
	plog.Noticef("PASSED:  %d", s.passed)
	plog.Noticef("SKIPPED: %d", s.skipped)

	if s.failed != 0 {
		return RunnerFailed
	}
	return nil
}

func (s *Local) runGroup(tc Group, match string) {
	defer s.pending.Done()

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

	s.parallel <- struct{}{}
	defer func() { <-s.parallel }()

	prepare := NewGroupHarness(tc, "Prepare")
	prepare.Run()

	if !prepare.Failed() && !prepare.Skipped() {
		cleanup := NewGroupHarness(tc, "Cleanup")
		defer cleanup.Run()
	}

	for _, test := range tests {
		if prepare.Failed() {
			test.Fail()
		} else if prepare.Skipped() {
			test.Skip()
		}

		test.Run()
		if test.Failed() {
			s.failed++
		} else if test.Skipped() {
			s.skipped++
		} else {
			s.passed++
		}
	}
}

func (s *Local) runTestFunc(tf TestFunc, match string) {
	defer s.pending.Done()

	name := TestName(tf)
	if ok, _ := filepath.Match(match, name); !ok {
		return
	}

	s.parallel <- struct{}{}
	defer func() { <-s.parallel }()

	test := NewTestFuncHarness(tf)
	test.Run()
	if test.Failed() {
		s.failed++
	} else if test.Skipped() {
		s.skipped++
	} else {
		s.passed++
	}
}
