package statemachine

import (
	"context"
	"fmt"
	"github.com/a-inacio/edt-go/internal/mermaid"
	"github.com/a-inacio/edt-go/pkg/event"
	"github.com/a-inacio/edt-go/pkg/eventhub"
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

func (builder *StateMachineBuilder) SubscribeFrom(hub *eventhub.EventHub) *StateMachineBuilder {
	builder.hub = hub
	return builder
}

func (builder *StateMachineBuilder) WithEvents(events ...event.Event) *StateMachineBuilder {
	for _, e := range events {
		builder.events = append(builder.events, eventBuilder{
			event: e,
		})
	}

	return builder
}

func (builder *StateMachineBuilder) WithEventNames(names ...string) *StateMachineBuilder {
	for _, n := range names {
		builder.events = append(builder.events, eventBuilder{
			name: n,
		})
	}

	return builder
}

func (builder *StateMachineBuilder) WithEventForEntering(state string, event event.Event) *StateMachineBuilder {
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

func (builder *StateMachineBuilder) AddTransition(from string, event event.Event, to string) *StateMachineBuilder {
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

			e, ok := refTable[node.To]
			if !ok {
				return nil, fmt.Errorf("no transition event defined for %s --> %s", node.From, node.To)
			}

			err = stateMachine.AddTransition(node.From, e, node.To)

			builder.trySubscribeFromHub(e, stateMachine)

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

		builder.trySubscribeFromHub(t.event, stateMachine)
	}

	return stateMachine, err
}

func (builder *StateMachineBuilder) eventReferenceTable() (map[string]event.Event, error) {
	table := map[string]event.Event{}

	for _, e := range builder.events {
		state := e.state
		var eventName string
		var eventInstance event.Event

		if e.name == "" {
			eventName = event.GetName(e.event)
			eventInstance = e.event
		} else {
			eventName = e.name
			eventInstance = event.WithName(eventName)
		}

		if state == "" {

			if !strings.HasPrefix(eventName, "GoTo") {
				return nil, fmt.Errorf("a conventional event must have a name starting by GoTo, got: %s", eventName)
			}

			state = strings.TrimPrefix(eventName, "GoTo")
		}

		if _, alreadyAdded := table[state]; alreadyAdded {
			return nil, fmt.Errorf("state %s already has a transition event %s", state, eventName)
		}

		table[state] = eventInstance
	}

	return table, nil
}

func (builder *StateMachineBuilder) trySubscribeFromHub(e event.Event, sm *StateMachine) {
	if builder.hub == nil {
		return
	}

	builder.hub.RegisterHandler(e, &stateMachineHubHandler{sm: sm})
}
