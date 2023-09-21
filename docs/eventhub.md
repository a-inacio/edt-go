# Event Hub

This is the construct that alone justifies the name for this library. It allows the definition of events, where a `struct` represents a unique type. An instance of this `struct` can then be `Published` into the `EventHub`, which in turn will be notifying all `Subscribers` listening to that event.

There are two ways of getting notified:
- As an `Action`
- As a `Handler`

The `EventHub` deals only with `Handlers` which will be the most optimal choice since `Actions` will be wrapped as the former. However, you are encouraged to use `Actions` because they will make the code easier to understand and more maintainable. Consider `Handlers` for advanced use cases.

Events can contain data that will be passed down to the `Subscribers`.
## Usage

Assuming an event:
```go
type SomeEvent struct {  
	SomeValue string  
}
```

This code prints `42`:
```go
hub := eventhub.NewEventHub(nil)  
  
  
hub.Subscribe(SomeEvent{}, func(ctx context.Context) (action.Result, error) {  
	ev, _ := event.Get[SomeEvent](ctx)  
	result = ev.SomeValue

	fmt.Print(result)
	
	return action.Nothing()  
})  
  
hub.Publish(SomeEvent{SomeValue: "42"}, nil)
```