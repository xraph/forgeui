package plugin

import "fmt"

// TopologicalSort returns plugins in dependency order using Kahn's algorithm.
// Returns an error if a circular dependency is detected.
func (r *Registry) TopologicalSort() ([]Plugin, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if len(r.plugins) == 0 {
		return []Plugin{}, nil
	}

	// Build adjacency list and in-degree map
	graph := make(map[string][]string)
	inDegree := make(map[string]int)

	// Initialize all plugins with zero in-degree
	for name := range r.plugins {
		graph[name] = []string{}
		inDegree[name] = 0
	}

	// Build the graph
	for name, p := range r.plugins {
		for _, dep := range p.Dependencies() {
			// Only consider dependencies that are registered
			if _, ok := r.plugins[dep.Name]; ok {
				// dep.Name -> name (dependency points to dependent)
				graph[dep.Name] = append(graph[dep.Name], name)
				inDegree[name]++
			}
		}
	}

	// Find all nodes with no incoming edges
	var queue []string
	for name, degree := range inDegree {
		if degree == 0 {
			queue = append(queue, name)
		}
	}

	// Kahn's algorithm
	var sorted []Plugin
	for len(queue) > 0 {
		// Dequeue
		name := queue[0]
		queue = queue[1:]
		sorted = append(sorted, r.plugins[name])

		// For each dependent of the current plugin
		for _, dependent := range graph[name] {
			inDegree[dependent]--
			if inDegree[dependent] == 0 {
				queue = append(queue, dependent)
			}
		}
	}

	// If we didn't process all plugins, there's a cycle
	if len(sorted) != len(r.plugins) {
		return nil, fmt.Errorf("circular dependency detected")
	}

	return sorted, nil
}

