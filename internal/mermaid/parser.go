package mermaid

import (
	"fmt"
	"strings"
)

type Node struct {
	From             string
	To               string
	SourceIsInitial  bool
	TargetIsTerminal bool
}

func Parse(graph string) ([]Node, error) {
	lines := strings.Split(graph, "\n")
	var nodes []Node
	for _, line := range lines {
		line = strings.TrimSpace(line)

		if line == "" || strings.HasPrefix(line, "stateDiagram") {
			continue
		}

		matches := strings.Split(line, "-->")

		if len(matches) == 2 {
			fromStateName := strings.TrimSpace(matches[0])
			toStateName := strings.TrimSpace(matches[1])

			if fromStateName == "" || toStateName == "" {
				return nil, fmt.Errorf("invalid transition: %s", line)
			}

			isTerminal := false
			if toStateName == "[*]" {
				isTerminal = true
			}

			isInitial := false
			if fromStateName == "[*]" {
				isInitial = true
			}

			nodes = append(nodes, Node{
				From:             fromStateName,
				To:               toStateName,
				SourceIsInitial:  isInitial,
				TargetIsTerminal: isTerminal,
			})
		} else {
			return nil, fmt.Errorf("invalid transition: %s", line)
		}
	}

	return nodes, nil
}
