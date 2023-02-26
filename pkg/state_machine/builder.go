package state_machine

import "context"

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

	for _, t := range builder.transitions {
		err = stateMachine.AddTransition(t.from, t.event, t.to)

		if err != nil {
			return nil, err
		}
	}

	return stateMachine, err
}
