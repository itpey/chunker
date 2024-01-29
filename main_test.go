// Copyright 2024 itpey
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"testing"
)

func TestParseSize(t *testing.T) {
	testCases := []struct {
		sizeStr   string
		expected  int64
		expectErr bool
	}{
		{"10B", 10, false},
		{"5KB", 5 * 1024, false},
		{"3MB", 3 * 1024 * 1024, false},
		{"2GB", 2 * 1024 * 1024 * 1024, false},
		{"invalid", 0, true},
	}

	for _, tc := range testCases {
		t.Run(tc.sizeStr, func(t *testing.T) {
			size, err := parseSize(tc.sizeStr)
			if err != nil && !tc.expectErr {
				t.Errorf("unexpected error: %v", err)
			}
			if size != tc.expected {
				t.Errorf("expected: %d, got: %d", tc.expected, size)
			}
		})
	}
}

func TestSplitFile(t *testing.T) {
	inputContent := "Hello, itpey!"
	inputFile, err := os.CreateTemp("", "input*.txt")
	if err != nil {
		t.Fatalf("failed to create temporary input file: %v", err)
	}
	defer os.Remove(inputFile.Name())
	defer inputFile.Close()

	if _, err := inputFile.WriteString(inputContent); err != nil {
		t.Fatalf("failed to write to temporary input file: %v", err)
	}

	outputPrefix := "outputChunk"
	chunkSize := int64(5)

	if err := splitFile(inputFile.Name(), outputPrefix, chunkSize); err != nil {
		t.Fatalf("failed to split file: %v", err)
	}

	for i := 1; ; i++ {
		outputFileName := fmt.Sprintf("%s_%d", outputPrefix, i)
		outputFile, err := os.Open(outputFileName)
		if os.IsNotExist(err) {
			break
		} else if err != nil {
			t.Fatalf("error opening output file %s: %v", outputFileName, err)
		}
		defer outputFile.Close()

		content, err := io.ReadAll(outputFile)
		if err != nil {
			t.Fatalf("error reading content from output file %s: %v", outputFileName, err)
		}

		start := (i - 1) * int(chunkSize)
		end := start + int(chunkSize)
		if end > len(inputContent) {
			end = len(inputContent)
		}
		expectedContent := []byte(inputContent)[start:end]

		if !bytes.Equal(content, expectedContent) {
			t.Errorf("unexpected content in output file %s: expected %q, got %q", outputFileName, expectedContent, content)
		}
	}
}
