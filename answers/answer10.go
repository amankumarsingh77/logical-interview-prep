package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strings"
)

type Options struct {
	Count      bool // -c: Prefix each line with its consecutive count.
	Duplicates bool // -d: Only print lines that are repeated.
	Unique     bool // -u: Only print lines that are unique (not repeated).
}

func Uniq(reader io.Reader, writer io.Writer, opts Options) error {
	scanner := bufio.NewScanner(reader)
	if !scanner.Scan() {
		return fmt.Errorf("no lines found")
	}
	freq := make(map[string]int)
	if !opts.Count && !opts.Duplicates && !opts.Unique {
		previous := scanner.Text()
		if _, err := writer.Write([]byte(previous + "\n")); err != nil {
			return fmt.Errorf("failed to write the data to stream :%v", err)
		}
		for scanner.Scan() {
			currText := scanner.Text()
			if currText != previous {
				if _, err := writer.Write([]byte(currText + "\n")); err != nil {
					return fmt.Errorf("failed to write the data to stream :%v", err)
				}
				previous = currText
			}
		}
	} else {
		var lines []string
		for scanner.Scan() {
			currLine := scanner.Text()
			freq[currLine]++
			if freq[currLine] == 1 {
				lines = append(lines, currLine)
			}
		}
		for _, line := range lines {
			count := freq[line]
			write := false
			if opts.Duplicates && count > 1 {
				write = true
			} else if opts.Unique && count == 1 {
				write = true
			} else if opts.Count && !opts.Duplicates && !opts.Unique {
				write = true
			}
			if write {
				var out string
				if opts.Count {
					out = fmt.Sprintf("%d %s\n", count, line)
				} else {
					out = line + "\n"
				}
				if _, err := writer.Write([]byte(out)); err != nil {
					return err
				}
			}

		}
	}

	return nil
}

func main() {
	// Sample input data to test against.
	const sampleInput = `apple
apple
banana
cherry
cherry
mango
cherry
apple
banana
banana`

	// Define all test cases.
	testCases := []struct {
		name    string
		options Options
	}{
		{"No Options", Options{}},
		{"Count (-c)", Options{Count: true}},
		{"Duplicates (-d)", Options{Duplicates: true, Count: true}},
		{"Unique (-u)", Options{Unique: true}},
	}

	// Run each test case.
	for _, tc := range testCases {
		// Use strings.NewReader to create an io.Reader from our sample string.
		reader := strings.NewReader(sampleInput)

		// Use bytes.Buffer to capture the output as an io.Writer.
		var writer bytes.Buffer

		fmt.Printf("--- Testing: %s ---\n", tc.name)

		// Call your function with the test case's options.
		if err := Uniq(reader, &writer, tc.options); err != nil {
			fmt.Printf("Error: %v\n", err)
			continue
		}

		// Print the captured output.
		fmt.Print(writer.String())
		fmt.Println("--------------------\n")
	}
}
