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
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func main() {
	if len(os.Args) != 4 {
		printUsage()
		return
	}

	inputFilePath := os.Args[1]
	outputPrefix := os.Args[2]
	chunkSizeStr := os.Args[3]

	chunkSize, err := parseSize(chunkSizeStr)
	if err != nil {
		fmt.Println("error:", err)
		printUsage()
		return
	}

	err = splitFile(inputFilePath, outputPrefix, chunkSize)
	if err != nil {
		fmt.Println("error:", err)
	}

}

func printUsage() {
	fmt.Println("=====================================================")
	fmt.Println("                  CHUNKER - USAGE                     ")
	fmt.Println("=====================================================")
	fmt.Println("Usage: chunker <inputFilePath> <outputPrefix> <chunkSize>")
	fmt.Println("-----------------------------------------------------")
	fmt.Println("- inputFilePath: Path to the input file to be split.")
	fmt.Println("- outputPrefix: Prefix for the output chunk files.")
	fmt.Println("- chunkSize: Size of each chunk. Examples: 1KB, 5MB, 1GB.")
	fmt.Println("- Supported units: B, KB, MB, GB. Default unit is bytes.")
	fmt.Println("-----------------------------------------------------")
	fmt.Println("Example: chunker input.txt outputChunk 5MB")
	fmt.Println("-----------------------------------------------------")
	fmt.Println("For more information, visit: https://github.com/itpey/chunker")
	fmt.Println("=====================================================")
}

func parseSize(sizeStr string) (int64, error) {
	re := regexp.MustCompile(`^(\d+)[ \t]*([KkMmGgBb]?)B?$`)
	match := re.FindStringSubmatch(strings.ToUpper(sizeStr))

	if len(match) == 0 {
		return 0, fmt.Errorf("invalid size format")
	}

	value, err := strconv.Atoi(match[1])
	if err != nil {
		return 0, fmt.Errorf("invalid size value: %v", err)
	}

	unit := match[2]
	switch strings.ToUpper(unit) {
	case "K":
		return int64(value) * 1024, nil
	case "M":
		return int64(value) * 1024 * 1024, nil
	case "G":
		return int64(value) * 1024 * 1024 * 1024, nil
	case "B":
		return int64(value), nil
	default:
		return 0, fmt.Errorf("invalid unit: %s", unit)
	}
}

func splitFile(inputFilePath, outputPrefix string, chunkSize int64) error {
	file, err := os.Open(inputFilePath)
	if err != nil {
		return fmt.Errorf("error opening input file: %v", err)
	}
	defer file.Close()

	i := 1
	for {
		outputFileName := fmt.Sprintf("%s_%d", outputPrefix, i)
		outputFile, err := os.Create(outputFileName)
		if err != nil {
			return fmt.Errorf("error creating output file: %v", err)
		}

		written, err := io.CopyN(outputFile, file, chunkSize)
		outputFile.Close()

		if written == 0 {
			os.Remove(outputFileName)
			break
		}

		if err != nil && err != io.EOF {
			return fmt.Errorf("error writing chunk: %v", err)
		}

		fmt.Printf("Chunk %s created.\n", outputFileName)

		i++
	}

	return nil
}
