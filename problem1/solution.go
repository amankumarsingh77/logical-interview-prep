package main

import "fmt"

func main() {
	tasks := map[string][]string{
		"compile":       {"download_deps"},
		"run_tests":     {"compile"},
		"create_image":  {"compile"},
		"deploy":        {"run_tests", "create_image"},
		"download_deps": {},
	}
	fmt.Println(ResolveOrder(tasks))
}

func ResolveOrder(tasks map[string][]string) ([]string, error) {
	visited := make(map[string]bool)
	onPath := make(map[string]bool)
	var result []string
	var visit func(task string) error
	visit = func(task string) error {
		if onPath[task] {
			return fmt.Errorf("cycle detected on filed %v", task)
		}
		if visited[task] {
			return nil
		}

		onPath[task] = true

		for _, t := range tasks[task] {
			if _, ok := tasks[t]; !ok {
				return fmt.Errorf("could not find task %v in the given tasks", t)
			}
			if !visited[t] {
				if err := visit(t); err != nil {
					return err
				}
			}
		}

		visited[task] = true
		onPath[task] = false
		result = append(result, task)
		return nil
	}

	for task := range tasks {
		if !visited[task] {
			if err := visit(task); err != nil {
				return nil, err
			}
		}
	}

	return result, nil
}
