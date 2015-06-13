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
	"go/ast"
	"go/build"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"

	"github.com/coreos/mantle/harness"
)

func ScanFile(filePath string) error {
	fileSet := token.NewFileSet()
	fileAst, err := parser.ParseFile(fileSet, filePath, nil, 0)
	if err != nil {
		return err
	}
	if !ast.FileExports(fileAst) {
		return nil
	}

	for _, d := range fileAst.Decls {
		ast.Print(fileSet, d)
	}

	return nil
}

func ScanPackage(pkg string) error {
	pkgInfo, err := build.Import(pkg, ".", 0)
	if err != nil {
		return err
	}

	for _, f := range pkgInfo.GoFiles {
		p := filepath.Join(pkgInfo.Dir, f)
		if err := ScanFile(p); err != nil {
			return err
		}
	}

	for _, f := range pkgInfo.CgoFiles {
		p := filepath.Join(pkgInfo.Dir, f)
		if err := ScanFile(p); err != nil {
			return err
		}
	}

	return nil
}

func main() {
	_ = harness.TestCase{}
	for _, pkg := range os.Args[1:] {
		if err := ScanPackage(pkg); err != nil {
			fmt.Fprintf(os.Stderr,
				"Inspecting package %s failed: %v\n", pkg, err)
			os.Exit(1)
		}
	}
}
