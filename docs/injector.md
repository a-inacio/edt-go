# Injector

This construct allows you to append context information to any `Action` execution, by implementing a `Dependency Injection` pattern.

> 👉 If it is a Global or a Local context, it is up to you.

It only relies on Go's standard `context.Context` so, as long as you keep track of if it (something you should be doing anyway) your contextual information will be easily accessible.

> 👉 Means you are not limited only to these library's constructs. 

## Usage

### Setting Values

There are two distinct patterns at play here:

- Singleton
- Factory

#### Singleton Instance

An instance must be given, and it will be the same everytime it is retrieved.

```go
dependencies := injector.WithContext(nil)
dependencies.SetSingleton(SomeValue{message: "42"})
```

It is possible to define a constructor method (i.e. a function), this function is guaranteed to invoked only once.

```go
dependencies := injector.WithContext(nil)
dependencies.SetSingleton(func() SomeValue {
	return SomeValue{message: "42"}
})
```

The constructor method can have parameters as long as they can be satisfied by the injector:

```go
dependencies := injector.WithContext(nil)
dependencies.SetSingleton(func() SomeValue{
    return SomeValue{message: "42"}
})

dependencies.SetSingleton(func(value SomeValue) AnotherValue {
    return AnotherValue{message: value.message}
})
```

#### Factory

A callback function is to be utilised, and you should be creating a new instance on each invocation.

```go
dependencies := injector.WithContext(nil)

counter := 0
dependencies.SetFactory(func() SomeValue {
    counter++
    return SomeValue{counter: counter}
})
```

### Retrieving Values

Independently of how the setter is defined, you get the values always in the same manner.

```go
dependencies := injector.FromContext(ctx)

value, err := injector.GetValue[SomeValue](dependencies)

if err != nil {
	// do something with `value`
}
```

> ⚠️ Though you can use an interface, currently finding a dependency that implements such interface is not supported.
> Only direct type names mather.

## Inheriting contexts

🚧TODO
