Scenario: Command-Line Argument Parser ⌨️

You're creating a small command-line tool. A core piece of this tool is
a function that parses the user's raw input string into a structured
command and a set of flags.

The input string follows a few simple rules:

    The first word is always the command.

    Words starting with -- (long flag) or - (short flag) are considered flag names.

    A flag can be a boolean flag (it has no value) or it can have a value, which is the very next word after the flag.

    To handle values with spaces, arguments can be enclosed in "double quotes". The quotes should be removed from the final value.

Your Task

Write a Go function ParseCommand(input string) (\*ParsedCommand, error)
that takes a raw input string and returns a ParsedCommand struct.

Here's the struct you should populate: Go

type ParsedCommand struct { Command string // Flags can be bool, string,
etc. Flags map\[string\]interface{} }

Requirements:

    The function should correctly identify the command and all flags.

    Flag names in the map should be stored without the leading - or --.

    If a flag is present but has no value following it (and the next word isn't another flag), it's a boolean flag and its value in the map should be true.

    Your parser must correctly handle quoted values containing spaces.

Example:

    Input String: remote add --name origin "/users/alex/repo.git" -f

    Expected Output:

        Command: "remote"

        Flags: map[string]interface{}{"name": "origin", "f": true}

        Wait, what about add? Good catch. For this problem, let's assume the command can be multi-part. Anything before the first flag is part of the command. So the command is actually "remote add".

Corrected Example:

    Input String: remote add --name origin "/users/alex/repo.git" -f

    Expected Output (ParsedCommand struct):

        Command: "remote add"

        Flags: map[string]interface{}{"name": "origin", "f": true}
