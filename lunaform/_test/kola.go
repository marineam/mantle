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

package main

import (
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/coreos/go-semver/semver"

	"github.com/coreos/mantle/kola/cluster"
	"github.com/coreos/mantle/kola/register"
	"github.com/coreos/mantle/lunaform"

	// link it in!
	_ "github.com/coreos/mantle/kola/registry"
)

func filterTests(platform, arch string, version semver.Version) map[string]*register.Test {
	r := make(map[string]*register.Test)

	for name, t := range register.Tests {
		// Tests that assemble clusters on the fly are not supprted
		if t.ClusterSize < 1 {
			continue
		}

		// Blank user data not currently supported because lunaform
		// blocks project level ssh keys to avoid external state.
		// Tests should probably be able to add per-instance ssh
		// keys in gce metadata instead.
		if t.UserData == "" {
			continue
		}

		// Discovery not supported yet
		if strings.Contains(t.UserData, "$discovery") {
			continue
		}

		// Check the test's min and end versions when running more then one test
		if versionOutsideRange(version, t.MinVersion, t.EndVersion) {
			continue
		}

		allowed := true
		for _, p := range t.Platforms {
			if p == platform {
				allowed = true
				break
			} else {
				allowed = false
			}
		}
		if !allowed {
			continue
		}

		for _, a := range t.Architectures {
			if a == arch {
				allowed = true
				break
			} else {
				allowed = false
			}
		}
		if !allowed {
			continue
		}

		r[name] = t
	}

	return r
}

// versionOutsideRange checks to see if version is outside [min, end). If end
// is a zero value, it is ignored and there is no upper bound. If version is a
// zero value, the bounds are ignored.
func versionOutsideRange(version, minVersion, endVersion semver.Version) bool {
	if version == (semver.Version{}) {
		return false
	}

	if version.LessThan(minVersion) {
		return true
	}

	if (endVersion != semver.Version{}) && !version.LessThan(endVersion) {
		return true
	}

	return false
}

// register all kola tests in lunaform
func init() {
	tests := filterTests("gce", "amd64", semver.Version{})

	for _, test := range tests {
		test := test // for closure
		lunaform.Register(lunaform.Test{
			Name:        test.Name,
			ClusterSize: test.ClusterSize,
			UserData:    test.UserData,
			Run: func(c *lunaform.Cluster) {
				runTest(c, test)
			},
		})
	}
}

func runTest(c *lunaform.Cluster, t *register.Test) {
	// pass along all registered native functions
	var names []string
	for k := range t.NativeFuncs {
		names = append(names, k)
	}

	// Cluster -> TestCluster
	tcluster := cluster.TestCluster{
		H:           c.H,
		Cluster:     c,
		NativeFuncs: names,
	}

	// drop kolet binary on machines
	if t.NativeFuncs != nil {
		scpKolet(tcluster, "amd64")
	}

	defer func() {
		// give some time for the remote journal to be flushed so it can be read
		// before we run the deferred machine destruction
		time.Sleep(2 * time.Second)
	}()

	// run test
	t.Run(tcluster)
}

// scpKolet searches for a kolet binary and copies it to the machine.
func scpKolet(c cluster.TestCluster, mArch string) {
	for _, d := range []string{
		".",
		filepath.Dir(os.Args[0]),
		filepath.Join(filepath.Dir(os.Args[0]), mArch),
		filepath.Join("/usr/lib/kola", mArch),
	} {
		kolet := filepath.Join(d, "kolet")
		if _, err := os.Stat(kolet); err == nil {
			if err := c.DropFile(kolet); err != nil {
				c.Fatalf("dropping kolet binary: %v", err)
			}
			return
		}
	}
	c.Fatalf("Unable to locate kolet binary for %s", mArch)
}
