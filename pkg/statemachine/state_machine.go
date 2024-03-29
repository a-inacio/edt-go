package statemachine

import (
	"context"
	"errors"
	"fmt"
	"github.com/a-inacio/edt-go/pkg/event"
)

type State struct {
	Name     string
	OnBefore func(ctx context.Context, trigger Trigger)
	OnEnter  func(ctx context.Context, trigger Trigger)
	OnAfter  func(ctx context.Context, trigger Trigger)
}

type Trigger struct {
	Event     *event.Event
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

func NewStateMachine(initialState *State, ctx context.Context) (*StateMachine, error) {
	if initialState.Name == "" {
		return nil, errors.New("state name cannot be empty")
	}

	return &StateMachine{
		nodes: map[string]Node{
			initialState.Name: {
				Type:        InitialNode,
				State:       initialState,
				Transitions: map[string]Transition{},
			},
		},
		context: ctx,
		initial: initialState.Name,
	}, nil
}

func (sm *StateMachine) IsRunning() bool {
	return sm.current != ""
}

func (sm *StateMachine) AddState(state *State) error {
	if state.Name == "" {
		return errors.New("state name cannot be empty")
	}

	if _, alreadyAdded := sm.nodes[state.Name]; alreadyAdded {
		return fmt.Errorf("state already added %s", state.Name)
	}

	sm.nodes[state.Name] = Node{
		Type:        ChildNode,
		State:       state,
		Transitions: map[string]Transition{},
	}

	return nil
}

func (sm *StateMachine) AddTransition(fromStateName string, e event.Event, toStateName string) error {
	eventName := event.GetName(e)

	fromNode, ok := sm.nodes[fromStateName]
	if !ok {
		return fmt.Errorf("unknown source state: %s", fromStateName)
	}

	if _, alreadyAdded := fromNode.Transitions[eventName]; alreadyAdded {
		return fmt.Errorf("transition already added to %s: %s", fromStateName, eventName)
	}

	toNode, ok := sm.nodes[toStateName]
	if !ok {
		return fmt.Errorf("unknown destination state: %s", toStateName)
	}

	fromNode.Transitions[eventName] = Transition{
		EventName: eventName,
		To:        &toNode,
	}

	return nil
}

func (sm *StateMachine) TriggerEvent(e event.Event) error {
	if !sm.IsRunning() {
		return fmt.Errorf("state machine not started")
	}

	currentNode := sm.currentNode()

	eventName := event.GetName(e)

	transition, exists := currentNode.Transitions[eventName]
	if !exists {
		return fmt.Errorf("current state %s has no transition named: %s", sm.current, eventName)
	}

	trigger := Trigger{
		FromState: currentNode.State,
		ToState:   transition.To.State,
		Event:     &e,
	}

	sm.executeTransition(&currentNode, &trigger, &transition)

	return nil
}

func (sm *StateMachine) Start() error {
	if sm.IsRunning() {
		return fmt.Errorf("state machine already runnig at state: %s", sm.current)
	}

	initialNode := sm.nodes[sm.initial]
	trigger := Trigger{
		ToState: initialNode.State,
	}

	sm.executeTransition(nil, &trigger, &Transition{
		To:        &initialNode,
		EventName: "__start__",
	})

	return nil
}

func (sm *StateMachine) currentNode() Node {
	return sm.nodes[sm.current]
}

func (sm *StateMachine) executeTransition(currentNode *Node, trigger *Trigger, transition *Transition) {
	if currentNode != nil {
		if currentNode.State.OnAfter != nil {
			currentNode.State.OnAfter(sm.context, *trigger)
		}
	}

	sm.current = transition.To.State.Name

	if transition.To.State.OnBefore != nil {
		transition.To.State.OnBefore(sm.context, *trigger)
	}

	if transition.To.State.OnEnter != nil {
		transition.To.State.OnEnter(sm.context, *trigger)
	}
}
