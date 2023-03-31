# Loopable

Use this constructs to keep something running forever.
Execution can be canceled using a `context`.

## Usage

```go
loopable.RunForever(context.Background(), 5*time.Second, func(ctx context.Context) (action.Result, error) {
    fmt.Println("Hello World! See you in a moment again soon 👋")
    return action.Nothing()
})
```

### Cancellation

```go
ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
defer cancel()

loopable.RunForever(ctx, 5*time.Second, func(ctx context.Context) (action.Result, error) {
    fmt.Println("Hello World! Will only be called once 👋")
    return action.Nothing()
})
```