package state_machine

import (
	"context"
	"errors"
	"fmt"
	"reflect"
)

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

func (sm *StateMachine) AddTransition(fromStateName string, event Event, toStateName string) error {
	eventName := reflect.TypeOf(event).Name()

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

func (sm *StateMachine) TriggerEvent(event Event) error {
	if !sm.IsRunning() {
		return fmt.Errorf("state machine not started")
	}

	currentNode := sm.currentNode()

	eventName := reflect.TypeOf(event).Name()

	transition, exists := currentNode.Transitions[eventName]
	if !exists {
		return fmt.Errorf("current state %s has no transition named: %s", sm.current, eventName)
	}

	trigger := Trigger{
		FromState: currentNode.State,
		ToState:   transition.To.State,
		Event:     &event,
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
