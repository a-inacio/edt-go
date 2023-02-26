package state_machine

import (
	"context"
	"fmt"
	"github.com/a-inacio/edt-go/pkg/state_machine/mermaid"
	"reflect"
	"strings"
)

func NewBuilder() *StateMachineBuilder {
	return &StateMachineBuilder{}
}

func (builder *StateMachineBuilder) WithInitialState(initialState *State) *StateMachineBuilder {
	builder.initialState = initialState
	return builder
}

func (builder *StateMachineBuilder) WithContext(ctx context.Context) *StateMachineBuilder {
	builder.context = ctx
	return builder
}

func (builder *StateMachineBuilder) WithEvents(events ...Event) *StateMachineBuilder {
	for _, e := range events {
		builder.events = append(builder.events, eventBuilder{
			event: e,
		})
	}

	return builder
}

func (builder *StateMachineBuilder) WithEventForEntering(state string, event Event) *StateMachineBuilder {
	builder.events = append(builder.events, eventBuilder{
		state: state,
		event: event,
	})

	return builder
}

func (builder *StateMachineBuilder) FromGraph(graph string) *StateMachineBuilder {
	builder.graph = graph

	return builder
}

func (builder *StateMachineBuilder) AddState(state *State) *StateMachineBuilder {
	builder.states = append(builder.states, state)
	return builder
}

func (builder *StateMachineBuilder) AddTransition(from string, event Event, to string) *StateMachineBuilder {
	builder.transitions = append(builder.transitions, transitionBuilder{
		from:  from,
		event: event,
		to:    to,
	})
	return builder
}

func (builder *StateMachineBuilder) Build() (*StateMachine, error) {
	stateMachine, err := NewStateMachine(builder.initialState, builder.context)

	if err != nil {
		return nil, err
	}

	for _, s := range builder.states {
		err = stateMachine.AddState(s)

		if err != nil {
			return nil, err
		}
	}

	if builder.graph != "" {
		mNodes, err := mermaid.Parse(builder.graph)

		if err != nil {
			return nil, err
		}

		refTable, err := builder.eventReferenceTable()

		if err != nil {
			return nil, err
		}

		for _, node := range mNodes {
			if node.SourceIsInitial || node.TargetIsTerminal {
				continue
			}

			event, ok := refTable[node.To]
			if !ok {
				return nil, fmt.Errorf("no transition event defined for %s --> %s", node.From, node.To)
			}

			err = stateMachine.AddTransition(node.From, event, node.To)

			if err != nil {
				return nil, err
			}
		}
	}

	for _, t := range builder.transitions {
		err = stateMachine.AddTransition(t.from, t.event, t.to)

		if err != nil {
			return nil, err
		}
	}

	return stateMachine, err
}

func (builder *StateMachineBuilder) eventReferenceTable() (map[string]Event, error) {
	table := map[string]Event{}

	for _, e := range builder.events {
		state := e.state
		eventName := reflect.TypeOf(e.event).Name()

		if state == "" {

			if !strings.HasPrefix(eventName, "GoTo") {
				return nil, fmt.Errorf("a conventional event must have a name starting by GoTo, got: %s", eventName)
			}

			state = strings.TrimPrefix(eventName, "GoTo")
		}

		if _, alreadyAdded := table[state]; alreadyAdded {
			return nil, fmt.Errorf("state %s already has a transition event %s", state, eventName)
		}

		table[state] = e.event
	}

	return table, nil
}
