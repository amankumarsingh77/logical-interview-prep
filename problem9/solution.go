package main

import (
	"errors"
	"log"
	"strings"
)

type ParsedCommand struct {
	Command string
	Flags   map[string]interface{}
}

func main() {
	log.Println(ParseCommand(`git status`))
}
func ParseCommand(input string) (*ParsedCommand, error) {
	input = strings.TrimSpace(input)
	parts, err := tokenize(input)
	if err != nil {
		return nil, err
	}
	if len(parts) == 0 {
		return nil, errors.New("empty input")
	}

	pc := &ParsedCommand{
		Flags: make(map[string]interface{}),
	}

	idx := 0
	for idx < len(parts) && !strings.HasPrefix(parts[idx], "-") {
		if idx > 0 {
			pc.Command += " "
		}
		pc.Command += parts[idx]
		idx++
	}

	for idx < len(parts) {
		token := parts[idx]
		if strings.HasPrefix(token, "-") {
			flag := cleanFlag(token)
			if idx+1 < len(parts) && !strings.HasPrefix(parts[idx+1], "-") {
				pc.Flags[flag] = parts[idx+1]
				idx += 2
			} else {
				pc.Flags[flag] = true
				idx++
			}
		} else {
			return nil, errors.New("unexpected token: " + token)
		}
	}

	return pc, nil
}
func tokenize(input string) ([]string, error) {
	var parts []string
	var current strings.Builder
	inQuotes := false

	for _, r := range input {
		switch {
		case r == '"':
			inQuotes = !inQuotes
		case r == ' ' && !inQuotes:
			if current.Len() > 0 {
				parts = append(parts, current.String())
				current.Reset()
			}
		default:
			current.WriteRune(r)
		}
	}

	if inQuotes {
		return nil, errors.New("unclosed quote")
	}
	if current.Len() > 0 {
		parts = append(parts, current.String())
	}

	return parts, nil
}

func cleanFlag(s string) string {
	if strings.HasPrefix(s, "--") {
		return s[2:]
	}
	if strings.HasPrefix(s, "-") {
		return s[1:]
	}
	return s
}
