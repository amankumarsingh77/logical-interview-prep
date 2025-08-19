# String Decompression Utility 📦

You're working on a system that receives compressed data over a network to save bandwidth. The data consists of strings that have been compacted using a special format. Your task is to write a utility that decompresses these strings back into their original form.

## Compression Format

The format is `k[string]`, where the string inside the square brackets is repeated `k` times.

### Rules

- `k[string]` means the string is repeated `k` times. `k` will always be a **positive integer**.
- The input can contain regular characters that are not part of a compression block.
- Compression blocks **can be nested**. For example, in `3[a2[c]]`, the `2[c]` block is nested inside the `3[a...]` block.

## Your Task

Write a Go function:

```go
Decompress(s string) (string, error)
```

that takes the compressed string and returns the fully decompressed version.

### Requirements

- Correctly handle simple, sequential, and nested patterns.
- Return an error for malformed input, including but not limited to:
  - Mismatched brackets (`[` without `]` or extra `]`).
  - Invalid or missing number for `k`.
  - Non-positive repeat counts (`k <= 0`).
  - Brackets in the wrong order (`]...[`).
  - Trailing digits not followed by `[`.

## Examples

| Input          | Output        |
|----------------|---------------|
| `3[a]2[bc]`    | `aaabcbc`     |
| `3[a2[c]]`     | `accaccacc`   |
| `2[abc]3[cd]ef`| `abcabccdcdcdef` |
| `a2[b3[c]]d`   | `abcccbcccd`  |

---

## Reference Implementation (Go)

```go
package decompress

import (
	"errors"
	"strings"
	"unicode"
)

// Decompress expands strings in the form k[substr],
// supporting nesting like 3[a2[c]] -> accaccacc.
// Returns an error for malformed inputs.
func Decompress(s string) (string, error) {
	var countStack []int
	var strStack []string

	currCount := 0
	var currStr strings.Builder

	for _, r := range s {
		switch {
		case unicode.IsDigit(r):
			// Accumulate multi-digit numbers
			currCount = currCount*10 + int(r-'0')

		case r == '[':
			// A '[' must be preceded by a positive repeat count
			if currCount <= 0 {
				return "", errors.New("malformed input: '[' must be preceded by positive repeat count")
			}
			countStack = append(countStack, currCount)
			strStack = append(strStack, currStr.String())
			currCount = 0
			currStr.Reset()

		case r == ']':
			// There must be a matching '['
			if len(countStack) == 0 || len(strStack) == 0 {
				return "", errors.New("malformed input: unmatched ']'")
			}
			// Pop stacks
			repeat := countStack[len(countStack)-1]
			countStack = countStack[:len(countStack)-1]

			prev := strStack[len(strStack)-1]
			strStack = strStack[:len(strStack)-1]

			// Build repeated segment
			var seg strings.Builder
			seg.Grow(currStr.Len() * repeat)
			piece := currStr.String()
			for i := 0; i < repeat; i++ {
				seg.WriteString(piece)
			}

			// Append to previous layer
			currStr.Reset()
			currStr.Grow(len(prev) + seg.Len())
			currStr.WriteString(prev)
			currStr.WriteString(seg.String())

		default:
			// Regular character
			currStr.WriteRune(r)
		}
	}

	// Any leftover count means digits not followed by '['
	if currCount != 0 {
		return "", errors.New("malformed input: trailing digits not followed by '['")
	}
	// Unmatched '[' if stacks not empty
	if len(countStack) != 0 || len(strStack) != 0 {
		return "", errors.New("malformed input: unmatched '['")
	}

	return currStr.String(), nil
}
```

### Notes on Error Handling

- **`[` without a preceding positive integer** → error.
- **Extra `]`** → error.
- **Leftover digits** at the end (like `12abc`) → error.
- **Unclosed bracket** (like `3[a`) → error.

### Time & Space Complexity

- Let *n* be the length of the input and *m* the length of the output.
- Time: **O(n + m)** — every input rune is processed once; output writing is proportional to the expanded size.
- Space: **O(n + m)** — stacks up to nesting depth; output buffer of size *m*.

---

## Minimal Test Cases

```go
package decompress_test

import (
	"testing"

	"github.com/your/module/decompress"
)

func TestDecompress_Valid(t *testing.T) {
	cases := map[string]string{
		"3[a]2[bc]":     "aaabcbc",
		"3[a2[c]]":      "accaccacc",
		"2[abc]3[cd]ef": "abcabccdcdcdef",
		"a2[b3[c]]d":    "abcccbcccd",
		"xyz":           "xyz",
	}

	for in, want := range cases {
		got, err := decompress.Decompress(in)
		if err != nil {
			t.Fatalf("unexpected error for %q: %v", in, err)
		}
		if got != want {
			t.Fatalf("for %q: got %q, want %q", in, got, want)
		}
	}
}

func TestDecompress_Errors(t *testing.T) {
	errInputs := []string{
		"3[a",        // unmatched '['
		"3a]",        // unmatched ']'
		"[a]",        // missing count
		"0[a]",       // non-positive count
		"10",         // trailing digits
		"2[abc]x]",   // extra ']'
	}
	for _, in := range errInputs {
		if _, err := decompress.Decompress(in); err == nil {
			t.Fatalf("expected error for %q, got nil", in)
		}
	}
}
```

---

## CLI Snippet (Optional)

```go
// go run ./cmd/decompress '2[abc]3[cd]ef'
package main

import (
	"fmt"
	"log"
	"os"

	"your/module/decompress"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatalf("usage: %s <compressed-string>", os.Args[0])
	}
	out, err := decompress.Decompress(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(out)
}
```

---

Happy decompressing! 🎈
