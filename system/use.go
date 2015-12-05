// Copyright 2014 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package system

import "unsafe"

// use is a no-op, but the compiler cannot see that it is.
// Calling use(p) ensures that p is kept live until that point.
// https://github.com/golang/go/commit/cf622d758cd51cfa09f5b503d323c81ed3a5541e
//go:noescape
func use(p unsafe.Pointer)
