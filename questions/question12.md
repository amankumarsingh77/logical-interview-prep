## Scenario: Semantic Versioning (SemVer) Parser and Comparator 🧐

You are building a dependency management tool for a new programming
language. A core requirement of this tool is to accurately parse and
compare software versions to resolve dependencies correctly. The
versions follow the popular Semantic Versioning (SemVer) standard.

A SemVer string has the format MAJOR.MINOR.PATCH-PRERELEASE+METADATA.

    Core Version: MAJOR.MINOR.PATCH (e.g., 2.1.10). These are non-negative integers.

    Pre-release Identifier (Optional): A hyphen (-) followed by a series of dot-separated identifiers (e.g., -alpha.1, -beta).

    Build Metadata (Optional): A plus sign (+) followed by a series of dot-separated identifiers (e.g., +build.123, +001).

## Your Task

Your task is to implement a SemVer parser and a comparison function.
This will involve two main parts:

1.  Parsing the Version String

First, you need a function that parses a version string into a
structured format. It's recommended to define a struct to hold these
components.

    Function Signature: ParseVersion(versionString string) (*Version, error)

    Version Struct (Suggestion):
    Go

    type Version struct {
        Major      int
        Minor      int
        Patch      int
        PreRelease string
        Metadata   string
    }

2.  Comparing Two Versions

Next, you need a function or method that compares two parsed Version
structs.

    Function Signature: Compare(v1 *Version, v2 *Version) int

    Return Value:

        -1 if v1 is less than v2

        0 if v1 is equal to v2

        1 if v1 is greater than v2

## Comparison Rules (Precedence)

The comparison logic must follow the official SemVer rules:

    Core Version Comparison: Precedence is determined by comparing Major, Minor, and Patch versions numerically in that order.

        Example: 1.2.0 < 2.0.0

        Example: 1.9.0 < 1.10.0

        Example: 1.2.3 < 1.2.4

    Pre-release Identifier Comparison:

        A version with a pre-release tag has lower precedence than a version without one.

            Example: 1.0.0-alpha < 1.0.0

        If both versions have pre-release tags, they are compared by splitting them by . and comparing each identifier from left to right.

            Numeric identifiers are compared numerically.

            String identifiers are compared lexicographically (alphabetically).

            Example: 1.0.0-alpha < 1.0.0-beta

            Example: 1.0.0-rc.1 < 1.0.0-rc.2

    Build Metadata:

        Build metadata (+...) must be completely ignored when determining version precedence.

        Example: 1.0.0+build.1 is equal to 1.0.0+build.2 in terms of precedence.

## Examples

Version 1 Comparison Version 2 Reason 1.0.0 \< 2.0.0 Major version 1 is
less than 2 1.9.0 \< 1.10.0 Minor version 9 is less than 10 1.0.0 \>
1.0.0-alpha Normal version has higher precedence 1.0.0-beta \>
1.0.0-alpha beta comes after alpha alphabetically 1.0.0-rc.1 \<
1.0.0-rc.2 Numeric identifier 1 is less than 2 1.2.3+001 = 1.2.3+002
Build metadata is ignored
