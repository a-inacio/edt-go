# Action

This simple Construct is the glue that makes most all other Constructs work together. Without being too opinionated, your functions just have to follow a simple convention, in order to have them injected as callbacks.

```go
func(ctx context.Context) (Result, error) {
	return nil, nil
}
```
Can be written as:

```go
func(ctx context.Context) (Result, error) {
	return action.Nothing()
}
```

Which in turn is the same as:

```go
action.DoNothing()
```

Most other Constructs have operations with this very same signature. 

> 👉 That's the magic sauce that allows them to be chained together.