# Scenario: Task Dependency Resolver

Imagine you're building a script for an automated build system. The system needs to run a series of tasks, but some tasks must be completed before others can begin. For example, you must compile the code before you can `run_tests`.

## Your Task

Write a function in Go that takes a set of tasks and their dependencies and returns a valid sequence in which to execute them.

### Requirements

- The input will be a map where the key is the task name and the value is a slice of its dependencies.
- Language: Go

```go
// Example Input
tasks := map[string][]string{
    "compile":         {"download_deps"},
    "run_tests":       {"compile"},
    "create_image":    {"compile"},
    "deploy":          {"run_tests", "create_image"},
    "download_deps":   {},
}
```

- Your function signature should be:

```go
func ResolveOrder(tasks map[string][]string) ([]string, error)
```

- The function should return a slice of strings (`[]string`) representing a valid execution order.  
  For the example above, a valid output would be:  
  `["download_deps", "compile", "run_tests", "create_image", "deploy"]`  
  (Note: `run_tests` and `create_image` could be swapped.)

- Your function **must detect circular dependencies**.  
  If Task A depends on Task B, and Task B depends on Task A, it's impossible to resolve.  
  In this case, your function should return an error.

- If a task lists a dependency that isn't defined as a task itself, that should also result in an error.

---

### Hints

- How would you keep track of the tasks you've already added to your execution plan?
- How would you detect a cycle?

> This is a classic problem that can be solved using **topological sorting** (e.g., using DFS). 