package event

import "testing"

type SomeEvent struct{}
type SomeOtherEventNamed struct{}
type SomeOtherEventNamedWithPointer struct{}

func (e SomeOtherEventNamed) EventName() string {
	return "SomethingElse"
}

func (e *SomeOtherEventNamedWithPointer) EventName() string {
	return "SomethingElsePointy"
}

func TestEventGetName(t *testing.T) {
	type args struct {
		event Event
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "anything can be an event",
			args: args{
				event: SomeEvent{},
			},
			want: "SomeEvent",
		},
		{
			name: "a pointer for a struct should work too",
			args: args{
				event: &SomeEvent{},
			},
			want: "SomeEvent",
		},
		{
			name: "devs should be able to name the events",
			args: args{
				event: SomeOtherEventNamed{},
			},
			want: "SomethingElse",
		},
		{
			name: "devs should be able to name the events, even with pointers",
			args: args{
				event: &SomeOtherEventNamed{},
			},
			want: "SomethingElse",
		},
		{
			// Looks yet another simple ... but this use case is what caused the "Here be dragons" reflection
			// madness...
			name: "devs should be able to name the events, even with pointers on the declaration",
			args: args{
				event: SomeOtherEventNamedWithPointer{},
			},
			want: "SomethingElsePointy",
		},
		{
			name: "devs should be able to name the events, even with pointers inception",
			args: args{
				event: &SomeOtherEventNamedWithPointer{},
			},
			want: "SomethingElsePointy",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetName(tt.args.event); got != tt.want {
				t.Errorf("GetName() = %v, want %v", got, tt.want)
			}
		})
	}
}
