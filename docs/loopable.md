# Loopable

🚧

## Samples

```go
loopable.RunForever(context.Background(), 5*time.Second, func(ctx context.Context) (action.Result, error) {
    fmt.Println("Hello World! See you in a moment again 👋")
    return action.Nothing()
})
```