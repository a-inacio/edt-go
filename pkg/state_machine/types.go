package state_machine

import (
	"context"
)

type Event interface{}

type State struct {
	Name     string
	OnBefore func(ctx context.Context, trigger Trigger)
	OnEnter  func(ctx context.Context, trigger Trigger)
	OnAfter  func(ctx context.Context, trigger Trigger)
}

type Trigger struct {
	Event     *Event
	FromState *State
	ToState   *State
}

// Transition is a type that represents a transition from one state to another in response to a specific event.
type Transition struct {
	EventName string
	To        *Node
}

type NodeType int

const (
	InitialNode NodeType = iota
	TerminalNode
	ChildNode
)

type Node struct {
	Type        NodeType
	State       *State
	Transitions map[string]Transition
}

// StateMachine is a type that represents a generic state machine.
type StateMachine struct {
	nodes   map[string]Node
	context context.Context
	current string
	initial string
}

type transitionBuilder struct {
	from  string
	to    string
	event Event
}

type StateMachineBuilder struct {
	initialState *State
	states       []*State
	transitions  []transitionBuilder
	context      context.Context
}
