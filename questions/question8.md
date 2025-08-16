# Logical interview prep
## Custom Gem

Of course. Let's tackle a more complex and dynamic problem that involves not just splitting strings, but interpreting them.

### Scenario: Template String Resolver ⚙️

You're working on a feature for a new microservice that sends customized user notifications. The notification templates are stored as strings with placeholders, and your job is to create a function that resolves these templates using a given data structure.

The placeholders follow the format `${path.to.value}`. The path represents a series of keys used to look up a value in a nested `map[string]interface{}`.

#### Example

**Template String:**
```
Hello, ${user.name}! You have an upcoming appointment for ${appointment.service}.
```

**Data Map (in Go):**
```go
data := map[string]interface{}{
    "user": map[string]interface{}{
        "name": "Alex",
        "id":   123,
    },
    "appointment": map[string]interface{}{
        "service": "Dental Check-up",
        "time":    "2025-08-16T14:00:00Z",
    },
}
```

**Expected Output String:**
```
Hello, Alex! You have an upcoming appointment for Dental Check-up.
```

---

### Your Task

Write a Go function:

```go
ResolveTemplate(template string, data map[string]interface{}) (string, error)
```

It must satisfy these requirements:

1. **Successful Resolution:** It should correctly replace all valid placeholders with their corresponding values from the data map.

2. **Nested Lookups:** It must be able to handle nested paths (e.g., `user.name`).

3. **Missing Paths:** If a placeholder refers to a path that does not exist in the data map (e.g., `${user.email}`), it should be replaced with an empty string.

4. **Error Handling:** The function should return an error if it encounters a malformed placeholder.  
   A malformed placeholder is one that starts with `${` but does not have a closing brace `}`.  
   Example: `Hello, ${user.name.`
