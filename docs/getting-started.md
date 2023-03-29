# Getting Started

## Importing

``` bash
got get -u github.com/a-inacio/edt-go
```

## Exploring the Constructs

This library implements a bunch of Design Patterns, so you will want to understand their roles and how they fit together in order to leverage from this toolkit.

### Creational
These Constructs have the purpose of bridging your code with other Constructs:

 - Action
 - Event

### Behavioral 

These Constructs will help you on specific scenarios where you want to have your code behaving in a certain way:

 - Actor
 - Director
 - Event Hub
 - State Machine

### Concurrency

There Constructs will help you on specific scenarios where you need to deal with asynchronous operations:

 - Awaitable
 - Expirable

## Construct Builders

**Some Constructs offer a Builder method**, this is the preferable way of instantiating such objects, since it will hide complexity of setting them up (you will learn that some have different options that affect behavior).

Take the following example:

``` go
res, err := expirable.NewBuilder().
    FromOperation(func(ctx context.Context) (action.Result, error) {
        operationCalled = true

        return awaitable.RunAfter(ctx, 5*time.Second, func(ctx context.Context) (action.Result, error) {
            return 42, nil
        })
    }).
    WithTimeout(10 * time.Second).
    Go(ctx)
```

The code above sets up an `Expirable`, using the `Builder` method, that will wait up to `10 seconds` to execute a task.
The task is just an `Awaitable` (a simple Construct) that will return `42` after `5 seconds` elapse. 

> 👉 Simpler Constructs do not offer a Builder method, because there would not be added value with the extra verbosity.