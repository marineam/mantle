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
	"fmt"
	"os"

	"github.com/coreos/mantle/cli"
	"github.com/coreos/mantle/harness"
	"github.com/coreos/mantle/kola"
)

var cmdSingleNode = &cli.Command{
	Name:        "single-node",
	Summary:     "Run run kolet single-node tests",
	Description: "blah blash",
	Run:         runSingleNode,
}

func init() {
	cli.Register(cmdSingleNode)
}

func runSingleNode(args []string) int {
	if len(args) > 0 {
		fmt.Fprintf(os.Stderr, "kaboom\n")
		return 2
	}

	runner := harness.Local{
		Tests: []interface{}{&kola.KolaSingleNode{}},
	}

	err := runner.Run("*", 1)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return 1
	}
	return 0
}
