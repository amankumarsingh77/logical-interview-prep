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

## Hints

- Consider using a stack-based approach to handle nested patterns
- Think about how to parse numbers that can have multiple digits
- Error handling is crucial - validate brackets are properly matched
- Consider edge cases like empty strings, invalid numbers, etc.

---

Happy decompressing! 🎈
