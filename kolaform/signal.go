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
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

// onSignal runs the given function if a fatal signal is received.
// The signal will be resent after the function runs or after the first
// signal is sent three more times.
func onSignal(f func()) chan<- os.Signal {
	ch := make(chan os.Signal, 1)

	go func() {
		done := make(chan struct{})
		sig, ok := <-ch
		if !ok {
			return
		}

		fmt.Fprintf(os.Stderr, "Received %s, cleaning up...\n", sig)

		go func() {
			defer close(done)
			f()
		}()

		for i := 0; i < 3; i++ {
			select {
			case <-done:
				os.Stderr.WriteString("Done.\n")
				break
			case sig2 := <-ch:
				if sig2 != sig {
					break
				}
			}
		}
		signal.Stop(ch)
		syscall.Kill(os.Getpid(), sig.(syscall.Signal))
	}()

	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	return ch
}

func offSignal(ch chan<- os.Signal) {
	signal.Stop(ch)
	close(ch)
}
