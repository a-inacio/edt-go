package state_machine

import (
	"context"
	"testing"
)

func TestNewStateMachine_FromBuilder(t *testing.T) {
	onBeforeCalled := false
	onEnterCalled := false

	sm, err := NewBuilder().
		WithInitialState(&State{
			Name: "Initial",
			OnBefore: func(ctx context.Context, trigger Trigger) {
				onBeforeCalled = true
			},
			OnEnter: func(ctx context.Context, trigger Trigger) {
				onEnterCalled = true
			},
		}).
		WithContext(context.Background()).
		Build()

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

func TestStateMachine_TriggerEvent_FromBuilder(t *testing.T) {
	onAfterCalledA := false
	onBeforeCalledB := false
	onEnterCalledB := false

	type GoToB struct {
	}

	type GoToC struct {
	}

	type GoToD struct {
	}

	sm, _ := NewBuilder().
		WithInitialState(&State{
			Name: "A",
			OnAfter: func(ctx context.Context, trigger Trigger) {
				onAfterCalledA = true
			},
		}).
		WithContext(context.Background()).
		AddState(&State{
			Name: "B",
			OnBefore: func(ctx context.Context, trigger Trigger) {
				onBeforeCalledB = true
			},
			OnEnter: func(ctx context.Context, trigger Trigger) {
				onEnterCalledB = true
			},
		}).
		AddState(&State{
			Name: "C",
		}).
		AddState(&State{
			Name: "D",
		}).
		AddTransition("A", GoToB{}, "B").
		AddTransition("B", GoToC{}, "C").
		AddTransition("C", GoToD{}, "D").
		Build()

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