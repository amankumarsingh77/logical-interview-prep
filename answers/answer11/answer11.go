package main

import (
	"errors"
	"log"
	"strings"
	"unicode"
)

func main() {
	log.Println(Decompress("2[abcasdas]2[cd5[a]]ef"))
}

// without nesting approach
//func Decompress(s string) (string, error) {
//	var finalString strings.Builder
//	start := 0
//	for start < len(s) {
//		if unicode.IsDigit(rune(s[start])) {
//			end := start
//			for end < len(s) && unicode.IsDigit(rune(s[end])) {
//				end++
//			}
//			repeated, _ := strconv.Atoi(s[start:end])
//			start = end
//			if start >= len(s) || s[start] != '[' {
//				return "", errors.New("malformed string")
//			}
//			end = strings.Index(s[start:], "]")
//			repeatingStr := s[start+1 : start+end]
//			log.Println(repeatingStr)
//			for i := 0; i < repeated; i++ {
//				finalString.WriteString(repeatingStr)
//			}
//			start += end + 1
//		} else {
//			finalString.WriteByte(s[start])
//			start++
//		}
//	}
//	return finalString.String(), nil
//}

// with nesting approach
func Decompress(s string) (string, error) {
	var countStack []int
	var strStack []string
	currCount := 0
	currStr := ""

	for _, r := range s {
		switch {
		case unicode.IsDigit(r):
			currCount = currCount*10 + int(r-'0')
		case r == '[':
			countStack = append(countStack, currCount)
			strStack = append(strStack, currStr)
			currCount = 0
			currStr = ""
		case r == ']':
			if len(countStack) == 0 || len(strStack) == 0 {
				return "", errors.New("malformed string: unmatched ']'")
			}
			repeat := countStack[len(countStack)-1]
			countStack = countStack[:len(countStack)-1]
			prevStr := strStack[len(strStack)-1]
			strStack = strStack[:len(strStack)-1]
			expanded := strings.Repeat(currStr, repeat)
			currStr = prevStr + expanded
		default:
			currStr += string(r)
		}
	}
	if len(strStack) != 0 {
		return "", errors.New("malformed string: unmatched '['")
	}
	return currStr, nil
}
