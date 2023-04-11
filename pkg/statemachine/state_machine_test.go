package statemachine

import (
	"context"
	"github.com/a-inacio/edt-go/pkg/action"
	"github.com/a-inacio/edt-go/pkg/awaitable"
	"testing"
	"time"
)

func TestNewStateMachine(t *testing.T) {
	onBeforeCalled := false
	onEnterCalled := false

	initialState := &State{
		Name: "Initial",
		OnBefore: func(ctx context.Context, trigger Trigger) {
			onBeforeCalled = true
		},
		OnEnter: func(ctx context.Context, trigger Trigger) {
			onEnterCalled = true
		},
	}
	sm, err := NewStateMachine(initialState, context.Background())

	if err != nil {
		t.Error("Creating the state machine should not have failed")
	}

	sm.Start()

	if !onBeforeCalled {
		t.Error("Initial State, OnBefore should have been called")
	}

	if !onEnterCalled {
		t.Error("Initial State, OnEnter should have been called")
	}
}

func TestStateMachine_TriggerEvent(t *testing.T) {
	onAfterCalledA := false
	onBeforeCalledB := false
	onEnterCalledB := false

	initialState := &State{
		Name: "A",
		OnAfter: func(ctx context.Context, trigger Trigger) {
			onAfterCalledA = true
		},
	}
	sm, _ := NewStateMachine(initialState, context.Background())

	sm.AddState(&State{
		Name: "B",
		OnBefore: func(ctx context.Context, trigger Trigger) {
			onBeforeCalledB = true
		},
		OnEnter: func(ctx context.Context, trigger Trigger) {
			onEnterCalledB = true
		},
	})

	sm.AddState(&State{
		Name: "C",
	})

	sm.AddState(&State{
		Name: "D",
	})

	type GoToB struct {
	}

	type GoToC struct {
	}

	type GoToD struct {
	}

	sm.AddTransition("A", GoToB{}, "B")
	sm.AddTransition("B", GoToC{}, "C")
	sm.AddTransition("C", GoToD{}, "D")

	sm.Start()
	sm.TriggerEvent(GoToB{})

	if !onAfterCalledA {
		t.Error("A State, OnAfter should have been called")
	}

	if !onBeforeCalledB {
		t.Error("B State, OnBefore should have been called")
	}

	if !onEnterCalledB {
		t.Error("B State, OnEnter should have been called")
	}

	err := sm.TriggerEvent(GoToD{})
	if err == nil {
		t.Error("B State, should not be possible to transition to D")
	}

	err = sm.TriggerEvent(GoToC{})
	if err != nil {
		t.Error("B State, should be possible to transition to C")
	}

	err = sm.TriggerEvent(GoToD{})
	if err != nil {
		t.Error("C State, should be possible to transition to D")
	}
}

func TestNewStateMachine_Cancellation(t *testing.T) {
	onBeforeCalled := false
	onEnterCalled := false

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	initialState := &State{
		Name: "Initial",
		OnBefore: func(ctx context.Context, trigger Trigger) {
			onBeforeCalled = true
		},
		OnEnter: func(ctx context.Context, trigger Trigger) {
			awaitable.RunAfter(ctx, 5*time.Second, func(ctx context.Context) (action.Result, error) {
				onEnterCalled = true
				return action.Nothing()
			})
		},
	}
	sm, err := NewStateMachine(initialState, ctx)

	if err != nil {
		t.Error("Creating the state machine should not have failed")
	}

	sm.Start()

	if !onBeforeCalled {
		t.Error("Initial State, OnBefore should have been called")
	}

	if onEnterCalled {
		t.Error("Initial State, OnEnter should have not been called")
	}
}
