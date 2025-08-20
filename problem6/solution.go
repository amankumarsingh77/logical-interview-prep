package main

import (
	"bufio"
	"fmt"
	"io"
)

type Options struct {
	Count      bool
	Duplicates bool
	Unique     bool
}

func Uniq(reader io.Reader, writer io.Writer, opts Options) error {
	scanner := bufio.NewScanner(reader)

	if !opts.Count && !opts.Duplicates && !opts.Unique {
		if !scanner.Scan() {
			return nil
		}
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
		type lineInfo struct {
			line  string
			count int
		}
		var orderedLines []lineInfo
		lineMap := make(map[string]int)

		for scanner.Scan() {
			currLine := scanner.Text()
			if _, exists := lineMap[currLine]; !exists {
				orderedLines = append(orderedLines, lineInfo{line: currLine, count: 0})
				lineMap[currLine] = len(orderedLines) - 1
			}
			orderedLines[lineMap[currLine]].count++
		}

		for _, info := range orderedLines {
			write := false
			if opts.Duplicates && info.count > 1 {
				write = true
			} else if opts.Unique && info.count == 1 {
				write = true
			} else if opts.Count && !opts.Duplicates && !opts.Unique {
				write = true
			}

			if write {
				var out string
				if opts.Count {
					out = fmt.Sprintf("%d %s\n", info.count, info.line)
				} else {
					out = info.line + "\n"
				}
				if _, err := writer.Write([]byte(out)); err != nil {
					return err
				}
			}
		}
	}

	return scanner.Err()
}
