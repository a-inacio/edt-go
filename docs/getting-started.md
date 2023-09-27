# Getting Started

## Importing

``` bash
go get -u github.com/a-inacio/edt-go
```

## Exploring the Constructs

This library implements a bunch of "Design Patterns", so you will want to understand their roles and how they fit together in order to leverage from this toolkit.
Maybe using the term "Design Pattern" is a bit abusive, so we will be referring to "Constructs" instead from now on.

### Creational
These Constructs have the primary purpose of bridging your code with other Constructs:

 - Action
 - Event
 - Injector

> 👉 You can also see them as a glue the chain Constructs together. 

### Behavioral 

These Constructs will help you on specific scenarios where you want to have your code behaving in a certain way:

 - Actor
 - Executor
 - Delayable
 - Director
 - Event Hub
 - State Machine

### Concurrency

These Constructs will help you on specific scenarios where you need to deal with asynchronous operations:

 - Promise
 - Expectable
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
    Do(ctx)
```

The code above sets up an `Expirable`, using the `Builder` method, that will wait up to `10 seconds` to execute a task.
The task is just an `Awaitable` (a simple Construct) that will return `42` after `5 seconds` elapse. 

> 👉 Simpler Constructs do not offer a Builder method, because there would not be added value with the extra verbosity.

## Conventions

### If it has a `Do`, then it is implicitly an action

Some of the constructs allow combinations like: chaining, referencing or wrapping `actions`. Since this library has a philosophy around this loose concept (it is just a function with a specific, yet simple signature) all constructs that can result of an execution of "something", have a method named `Do` that respects this convention. It becomes, then, very simple to implement the aforementioned combinations.

### If it has `Must` in the name, it means it can `panic`

This library has a very strict usage of `panic`, so the developer can use it with peace of mind, not worrying about schizophrenic behavior resulting from misusage.

Nevertheless, in some application designs, it is desirable to fully halt execution if certain fatal conditions occur.
Some examples, where crashing is often preferable:
- Dealing with unsatisfied dependencies.
- Missing mandatory configuration step (e.g. initialize a config file with credentials that can then be used as a dependency by the application).

> ⚡️ The developer can detect early mistakes or the application prevented to cause damage due to missing or wrong configurations. If using such a strategy, it is often best if fatal conditions are checked at early startup times. Evidently, each application design will have its own perks, you might want to allow it to run long enough to set up logging or leave a dead man letter behind.

To avoid the typical `Go` code that tests for an error and then `panics`, this library conventions that any method that is prefixed with `Must` will terminate program execution in case of an error.