# Expectable

This construct can be seen as a higher level abstraction for synchronizing tasks. Using native Go language constructs it is not that difficult to make such orchestration with a `channel` or a `sync.WaitGroup` 

These native constructs, though undeniably powerful, can quickly produce messy and/or complex code. You get the chance to express this same pattern but taking advantage of `Events` and the `Event Hub`.

In a nutshell you can block your execution until an `Event` is `Published`.

> 👉 You can set a timeout and a criteria for filtering a specific Event configuration.

## Usage

```go
expectable.NewBuilder().
    On(hub).
    Expect(SomeEvent{}).
    Go(context.Background())
```

The code above will block until `SomeEvent` gets published.

```go
hub.Publish(SomeEvent{}, ctx)
```

### Retrieving event values

The result value, after unblocking is the event type itself. You can use this mechanism to pass parameters further down into your logic.

```go
type SomeEvent struct {
	Message string
}

...

res, err := expectable.NewBuilder().
    On(hub).
    Expect(SomeEvent{}).
    Go(context.Background())

...

if res.(SomeEvent).Message == "Hello EDT!" {
	
}
```

### Filtering specific event values

You can target a specific event configuration that unlocks execution.

```go
type SomeEvent struct {
	Message string
}

...

res, err := expectable.NewBuilder().
    On(hub).
    Expect(SomeEvent{}).
    Where(func(e event.Event) bool {
        if e.(SomeEvent).Message == "Hello EDT!" {
            return true
        }

        return false
    }).
    Go(context.Background())
}
```

The example above will only unlock with an event message containing `Hello EDT!`.

> 👉 You can also combine this scenario with a Timeout.


### Cancellation

Another way to unblock execution is using context cancellation.

```go
ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
defer cancel()

res, err := expectable.NewBuilder().
    On(hub).
    Expect(SomeEvent{}).
    Go(ctx)

if err != nil {
	// Deal with cancellation
}
```

### Timeout

You can directly specify a timeout as a limit for receiving an event.

```go
ctx := context.Background()

res, err := expectable.NewBuilder().
    On(hub).
    Expect(SomeEvent{}).
	WithTimeout(1 * time.Second).
    Go(ctx)

if err != nil {
	// Deal with the timeout
}
```