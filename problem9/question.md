# Scenario: Command-Line Argument Parser ⌨️

You're building a command-line interface (CLI) for an application. Your
first task is to write a robust parser that takes the raw user input as
a single string and transforms it into a structured format that your
application can easily use.

## Your Task

Write a Go function,
`ParseCommand(input string) (*ParsedCommand, error)`, that accepts the
raw input string. It should return a pointer to a `ParsedCommand` struct
containing the parsed data.

Here is the struct definition:

``` go
type ParsedCommand struct {
    Command string
    // Flag values can be strings or booleans
    Flags   map[string]interface{}
}
```

## Parsing Rules

Your function must follow these five rules:

1.  **The Command**: Any text from the beginning of the input string
    until the first word that starts with a hyphen (-) is considered the
    command. This means the command can be a single word (e.g.,
    `publish`) or multiple words (e.g., `remote add`).

2.  **Flags and Values**: A word beginning with `-` or `--` is a flag.
    The very next word is its value, unless that word is also a flag.

3.  **Boolean Flags**: A flag is considered a boolean flag (with a value
    of `true`) if it's the last word in the input, or if the word
    immediately following it is also a flag.

4.  **Quoted Values**: To handle values containing spaces, a value can
    be enclosed in double quotes (`"`). Your parser must treat the
    entire quoted segment as a single value, and the quotes themselves
    should be removed from the final stored value.

5.  **Flag Naming**: When storing a flag in the `Flags` map, you must
    remove the leading `-` or `--` from its name. For example, `--path`
    becomes `path`.

## Definitive Example

**Input String:**

    publish --path "/reports/q1 report.pdf" --user "alex doe" -v

**Expected ParsedCommand Output:**

    Command: "publish"

    Flags: map[string]interface{}{
        "path": "/reports/q1 report.pdf",
        "user": "alex doe",
        "v": true,
    }
