# Scenario: Simplified uniq Command-Line Tool

The Unix command-line tool `uniq` is a utility that filters adjacent matching lines from an input file or stream. We're going to implement a simplified version of it.

---

## Your Task

Write a Go function that reads text from an `io.Reader`, processes it to handle adjacent duplicate lines according to a set of options, and writes the result to an `io.Writer`.

---

## Core Logic

The basic behavior is to read line by line. If a line is identical to the one immediately preceding it, it's considered a duplicate and should be handled according to the options.

---

## Function and Options

Your implementation should be based on this function signature:

```go
// Options struct to hold the command-line flags.
type Options struct {
    Count      bool // -c: Prefix each line with its consecutive count.
    Duplicates bool // -d: Only print lines that are repeated.
    Unique     bool // -u: Only print lines that are unique (not repeated).
}

func Uniq(reader io.Reader, writer io.Writer, opts Options) error {
    // Your implementation goes here.
}
```

---

## Behavior Based on Options

### Given the following input:

```
apple
apple
banana
cherry
cherry
cherry
apple
```

---

### No Options (`Options{}`)

Print each unique **adjacent** line.

```
apple
banana
cherry
apple
```

---

### Count (`-c`)

Prefix each output line with the number of times it appeared consecutively.

```
2 apple
1 banana
3 cherry
1 apple
```

---

### Duplicates (`-d`)

Only print the lines that appeared more than once consecutively.

```
apple
cherry
```

---

### Unique (`-u`)

Only print the lines that were **not** repeated.

```
banana
apple
```

---

**Note:** The `-d` and `-u` options are mutually exclusive and you can assume they won't be true at the same time.

---

## Implementation Hint

Structure the logic inside the `Uniq` function to:

* Read input line by line.
* Track the previous line and its repeat count.
* Depending on the options:

    * Print all deduplicated lines (`default`).
    * Print lines with counts (`-c`).
    * Print only repeated lines (`-d`).
    * Print only non-repeated lines (`-u`).
