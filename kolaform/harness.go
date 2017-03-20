// Copyright 2017 CoreOS, Inc.
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

package kolaform

import (
	"flag"
	"fmt"
	"os"

	"github.com/coreos/mantle/harness"
)

type Test struct {
	Name        string
	Run         func(c *Cluster)
	ClusterSize int
}

var (
	tests harness.Tests
	opts  = harness.Options{
		OutputDir: "_kola_temp",
	}
)

func Register(test Test) {
	if test.Name == "" {
		panic(fmt.Errorf("Missing Name: %#v", test))
	}
	if test.ClusterSize < 1 {
		panic(fmt.Errorf("Invalid ClusterSize: %#v", test))
	}
	tests.Add(test.Name, test.run)
}

func FlagSet(prefix string, errorHandling flag.ErrorHandling) *flag.FlagSet {
	return opts.FlagSet(prefix, errorHandling)
}

func Run() {
	suite := harness.NewSuite(opts, tests)
	if err := suite.Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		fmt.Println("FAIL")
		os.Exit(1)
	}
	fmt.Println("PASS")
	os.Exit(0)
}

func (test Test) run(h *harness.H) {
	h.Parallel()

	c := newCluster(h, test)

	// setup may fail with an incomplete state so schedule
	// the destroy to cleanup anything first.
	defer c.destroy()
	// also need to clean up in the event of any signals.
	ch := onSignal(c.destroy)
	defer offSignal(ch)
	// BUG(marineam): All this cleanup is awesome *but* doesn't
	// actually mean anything when interrupting terraform before
	// it writes out state. So orphaning things is easy.
	// Also, we never get test logs in the event of a signal. :(

	c.setup()

	test.Run(c)
}
