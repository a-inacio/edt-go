# Expectable

This construct can be seen as a higher level of abstraction on synchronizing tasks. Using native Go language constructs it is not that difficult to make such orchestration with a `channel` or a `sync.WaitGroup` 

These native constructs can quickly be utilised to produce messy code, though they are undeniably powerful. You get the chance to express this same pattern but taking advantage of `Events` and the `Event Hub`.

In a nutshell you can block your execution until an `Event` is `Published`.

## Usage

```go
expect := NewExpectable(hub, SomeEvent{})

res, err := expect.Go(ctx)
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

res, err := expect.Go(ctx)

...

if res.(SomeEvent).Message == "Hello EDT!" {
	
}
```

### Cancellation

Another way to unblock execution is using context cancellation.

```go
ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
defer cancel()

expect := NewExpectable(hub, SomeEvent{})

res, err := expect.Go(ctx)

if err != nil {
	// Deal with cancellation
}
```