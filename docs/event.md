# Event

This is a very simple Construct that has one very important purposed of claiming that anything can be an Event:

```go
type Event interface{}
```

Like an `Action`, the Event work as a glue for Event oriented Constructs:
 - Event Hub
 - State Machine

## Named Events

Normally you should not care much about the Name, nevertheless a `State Machine` has states names that you might have some trouble getting from a struct name.
This gives you freedom of naming in return of implementing the `NamedEvent` interface.

> 👉 The `State Machine` constructs leverages from named events.

e.g.:

```go
hub.Publish(*event.WithName("SomeEvent"), nil)
```

## Key Values

Another option to define your own specific type, containing fields for passing parameters during the event publishing is to use a `GenericNamedEvent`.

e.g.:
```go
hub.Publish(*event.WithNameAndKeyValues("SomeEvent", "Message", 42), nil)
```