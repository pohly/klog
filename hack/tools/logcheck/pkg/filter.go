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
	"bufio"
	"flag"
	"fmt"
	"os"
	"regexp"
	"strings"
)

// RegexpFilter implements flag.Value by accepting a file name and parsing that
// file.
type RegexpFilter struct {
	filename string
	regexps  []*regexp.Regexp
}

var _ flag.Value = &RegexpFilter{}

func (f *RegexpFilter) String() string {
	return f.filename
}

func (f *RegexpFilter) Set(value string) error {
	file, err := os.Open(value)
	if err != nil {
		return err
	}
	defer file.Close()

	// Reset before parsing.
	f.filename = value
	f.regexps = nil

	// Read line-by-line.
	scanner := bufio.NewScanner(file)
	for line := 0; scanner.Scan(); line++ {
		text := scanner.Text()
		if strings.HasPrefix(text, "#") {
			continue
		}
		text = strings.TrimSpace(text)
		if text == "" {
			continue
		}
		// Must match entire string.
		re, err := regexp.Compile("^" + text + "$")
		if err != nil {
			return fmt.Errorf("%s:%d: %v", value, line, err)
		}
		f.regexps = append(f.regexps, re)
	}

	if err := scanner.Err(); err != nil {
		return err
	}
	return nil
}

// Matches returns true if the text matches any of the regular expressions.
func (f *RegexpFilter) Matches(text string) bool {
	for _, re := range f.regexps {
		if re.MatchString(text) {
			return true
		}
	}
	return false
}
