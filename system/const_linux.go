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

// +build ignore

package system

/*
#define _GNU_SOURCE
#include <fcntl.h>
*/
import "C"

// Various values missing from os, syscall, and sys/unix.
const (
	// https://github.com/golang/go/issues/7830
	O_PATH    = C.O_PATH
	O_TMPFILE = C.O_TMPFILE

	// For linkat(2)
	AT_EMPTY_PATH     = C.AT_EMPTY_PATH
	AT_SYMLINK_FOLLOW = C.AT_SYMLINK_FOLLOW
	AT_FDCWD          = C.AT_FDCWD
)
