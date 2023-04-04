# Injector

This construct allows you to append context information to any `Action` execution, by implementing a `Dependency Injection` pattern.

> 👉 If it is a Global or a Local context, it is up to you.

It only relies on Go's standard `context.Context` so, as long as you keep track of if it (something you should be doing anyway) your contextual information will be easily accessible.

> 👉 Means you are not limited only to these library's constructs. 

## Dependency Injection

### Setting Values

There are two distinct patterns at play here:

- Singleton
- Factory

#### Singleton Instance

An instance must be given, and it will be the same everytime it is retrieved.

```go
injector := WithContext(nil)
injector.SetSingleton(SomeValue{message: "42"})
```
#### Factory

A callback function is to be utilised, you should be creating a new instance on each invocation.

```go
injector := WithContext(nil)

counter := 0
injector.SetFactory(func() SomeValue {
    counter++
    return SomeValue{counter: counter}
})
```

### Retrieving Values

Independently of how the setter is defined, you get the values always in the same manner.

```go
injector := FromContext(ctx)

value, err := GetValue(injector, SomeValue{})

if err != nil {
	// do something with `value`
}
```

> ⚠️ Though you can use an interface, currently finding a dependency that implements such interface is not supported.
> Only direct type names mather.

## Inheriting contexts

🚧TODO
