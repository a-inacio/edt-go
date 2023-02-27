package mermaid

type Node struct {
	From             string
	To               string
	SourceIsInitial  bool
	TargetIsTerminal bool
}
