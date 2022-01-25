/*
Copyright 2022 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package pkg

import (
	"os"
	"path"
	"testing"
)

func TestMatch(t *testing.T) {
	temp := t.TempDir()
	filename := path.Join(temp, "expressions")
	if err := os.WriteFile(filename, []byte(`# Example file
hello
abc
x.*y
`), 0666); err != nil {
		t.Fatalf("writing file: %v", err)
	}

	filter := &RegexpFilter{}
	if err := filter.Set(filename); err != nil {
		t.Fatalf("reading file: %v", err)
	}

	for text, expectMatch := range map[string]bool{
		"hello":       true,
		"hello/world": false, // no sub-matches
		"abc":         true,
		"x1y":         true,
		"x2y":         true,
	} {
		actualMatch := filter.Matches(text)
		if actualMatch != expectMatch {
			t.Errorf("%s: expected match %v, got %v", text, expectMatch, actualMatch)
		}
	}
}

func TestSetNoFile(t *testing.T) {
	filter := &RegexpFilter{}
	if err := filter.Set("no such file"); err == nil {
		t.Errorf("did not get expected error")
	}
}

func TestSetInvalid(t *testing.T) {
	temp := t.TempDir()
	filename := path.Join(temp, "expressions")
	if err := os.WriteFile(filename, []byte(`# Example file
[
`), 0666); err != nil {
		t.Fatalf("writing file: %v", err)
	}

	filter := &RegexpFilter{}
	err := filter.Set(filename)
	if err == nil {
		t.Fatalf("reading file: did not get expected error")
	}
	expected := filename + ":1: error parsing regexp: missing closing ]: `[$`"
	if err.Error() != expected {
		t.Errorf("error mismatch\nexpected: %q\n     got: %q", expected, err.Error())
	}
}
