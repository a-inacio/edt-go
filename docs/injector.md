# Injector

This construct allows you to append context information to any `Action` execution, by implementing a `Dependency Injection` pattern.

> 👉 If it is a Global or a Local context, it is up to you.

It only relies on Go's standard `context.Context` so, as long as you keep track of if it (something you should be doing anyway) your contextual information will be easily accessible.

> 👉 Actually it means you are not limited only to these library's constructs it can be utilised standalone.

## Usage

### Setting Values

There are two distinct creational patterns at play here:

- Singleton
- Factory

#### Singleton Instance

An instance must be given, and it will be the same every time it is retrieved.

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

##### Functions limitations

> ⚠️ Though functions can have 0 or more arguments, **they must have one and only one return type.**

### Retrieving Values

Independently of the injected dependency's creational design pattern type (Singleton or Factory), you get the values in the same manner.

```go
value, err := injector.GetValue[SomeValue](dependencies)

if err != nil {
	// do something with `value`
}
```

#### Satisfying Interfaces

Not withstanding the type creational pattern, when retrieving a value by interface the injector will either look for an explicit definition of an interface dependency or it will look for all the known types and find the ones that implement such interface.

This process only succeeds if there is exactly one resolution for the target type 
interface.

```go
value, err := injector.GetValue[SomeInterface](dependencies)

if err != nil {
	// do something with `value`
}
```

### Resolving Functions

As an alternative, to explicitly set a factory for every single entity of your application, you can manually control this flow by resolving said functions manually. The same constraints of #Factory methods apply. 

```go
value, err := injector.Resolve[AnotherValue](
	dependencies, 
	func(value SomeValue) AnotherValue {  
		return AnotherValue{message: value.message}  
	})

if err != nil {
	// do something with `value`
}
```