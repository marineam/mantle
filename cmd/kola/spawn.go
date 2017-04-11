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
	"io/ioutil"
	"os"

	"github.com/spf13/cobra"

	"github.com/coreos/mantle/kola"
	"github.com/coreos/mantle/platform"
	"github.com/coreos/mantle/platform/machine/aws"
	"github.com/coreos/mantle/platform/machine/gcloud"
	"github.com/coreos/mantle/platform/machine/qemu"
)

var (
	cmdSpawn = &cobra.Command{
		Run:    runSpawn,
		PreRun: preRun,
		Use:    "spawn",
		Short:  "spawn a CoreOS instance",
	}

	spawnNodeCount int
	spawnUserData  string
	spawnShell     bool
	spawnRemove    bool
	spawnVerbose   bool
)

func init() {
	cmdSpawn.Flags().IntVarP(&spawnNodeCount, "nodecount", "c", 1, "number of nodes to spawn")
	cmdSpawn.Flags().StringVarP(&spawnUserData, "userdata", "u", "", "userdata to pass to the instances")
	cmdSpawn.Flags().BoolVarP(&spawnShell, "shell", "s", false, "spawn a shell in an instance before exiting")
	cmdSpawn.Flags().BoolVarP(&spawnRemove, "remove", "r", true, "remove instances after shell exits")
	cmdSpawn.Flags().BoolVarP(&spawnVerbose, "verbose", "v", false, "output information about spawned instances")
	root.AddCommand(cmdSpawn)
}

func runSpawn(cmd *cobra.Command, args []string) {
	var userdata string
	var err error
	var cluster platform.Cluster

	if spawnNodeCount <= 0 {
		die("Cluster Failed: nodecount must be one or more")
	}

	if spawnUserData != "" {
		userbytes, err := ioutil.ReadFile(spawnUserData)
		if err != nil {
			die("Reading userdata failed: %v", err)
		}
		userdata = string(userbytes)
	} else {
		// ensure a key is injected
		userdata = "#cloud-config"
	}

	outputDir, err = kola.CleanOutputDir(outputDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Setup failed: %v\n", err)
		os.Exit(1)
	}

	switch kolaPlatform {
	case "qemu":
		cluster, err = qemu.NewCluster(&kola.QEMUOptions, outputDir)
	case "gce":
		cluster, err = gcloud.NewCluster(&kola.GCEOptions, outputDir)
	case "aws":
		cluster, err = aws.NewCluster(&kola.AWSOptions, outputDir)
	default:
		err = fmt.Errorf("invalid platform %q", kolaPlatform)
	}

	if err != nil {
		die("Cluster failed: %v", err)
	}

	var someMach platform.Machine
	for i := 0; i < spawnNodeCount; i++ {
		mach, err := cluster.NewMachine(userdata)
		if err != nil {
			die("Spawning instance failed: %v", err)
		}

		if spawnVerbose {
			fmt.Printf("Machine spawned at %v\n", mach.IP())
		}

		if spawnRemove {
			defer mach.Destroy()
		}

		someMach = mach
	}

	if spawnShell {
		if err := platform.Manhole(someMach); err != nil {
			die("Manhole failed: %v", err)
		}
	}
}

func die(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}
