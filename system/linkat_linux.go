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

package system

import (
	"syscall"
	"unsafe"
)

// Linkat wraps the linkat(2) system call.
func Linkat(olddirfd int, oldpath string, newdirfd int, newpath string, flags int) error {
	p0, err := syscall.BytePtrFromString(oldpath)
	if err != nil {
		return err
	}
	p1, err := syscall.BytePtrFromString(newpath)
	if err != nil {
		return err
	}
	_, _, e1 := syscall.Syscall6(syscall.SYS_LINKAT, uintptr(olddirfd), uintptr(unsafe.Pointer(p0)), uintptr(newdirfd), uintptr(unsafe.Pointer(p1)), uintptr(flags), 0)
	use(unsafe.Pointer(p0))
	use(unsafe.Pointer(p1))
	if e1 != 0 {
		return e1
	}
	return nil
}
