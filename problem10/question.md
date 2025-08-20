# Scenario: Structured Query String Parser 🗺️

You're building a backend for a web service. When a user submits an HTML
form or a complex query, the data arrives as a single URL-encoded
string. Your task is to parse this string into a nested Go map that
represents the data's structure.

## The Rules for Keys:

-   Keys are separated from values by `=`, and pairs are separated by
    `&`.
-   A `.` in a key indicates a nested map. For example, `user.name=Alex`
    should result in a map where the key `user` contains another map
    with the key `name`.
-   A `[index]` in a key indicates an array/slice. For example,
    `roles[0]=admin` should result in a map where the key `roles`
    contains a slice. For this problem, you can assume the indices will
    be valid, in order, and without gaps (e.g., 0, 1, 2...).
-   All keys and values are URL-encoded. Your function must decode them.
    For example, `%20` should become a space.

## Your Task

Write a Go function:

``` go
ParseQuery(query string) (map[string]interface{}, error)
```

that takes the encoded query string as input and returns a
`map[string]interface{}` representing the structured data.

### Requirements:

-   **Parsing**: Correctly split the string into key-value pairs.
-   **Decoding**: Decode the URL-encoded keys and values.
-   **Structure Building**: Dynamically create nested maps and slices
    based on the key paths.
-   **Error Handling**: The function should return an error if it
    encounters an invalid structure. For instance, if you process
    `user.name=Alex` and later process `user=notamap`, this is a
    conflict that should produce an error.

### Example:

**Input String:**

    user.name=Alex%20Doe&user.id=123&roles[0]=admin&roles[1]=editor&active=true

**Expected `map[string]interface{}` Output (Go):**

``` go
map[string]interface{}{
    "user": map[string]interface{}{
        "name": "Alex Doe",
        "id":   "123",
    },
    "roles": []interface{}{"admin", "editor"},
    "active": "true",
}
```
