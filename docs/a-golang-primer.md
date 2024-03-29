# A Golang primer

To make the most out of this library, it's essential to have a solid understanding of Go's capabilities. Familiarity with certain language and standard library features is crucial for achieving this goal.

It is out of scope of this document to provide a complete Go tutorial, it is expected that you have some familiarity already or have other means to cover that requirement.

This document aims to bring your attention to key aspects only.

## Functions

Functions in Go have a key part in the language design philosophy:

- Can return multiple values.
- Are first-class citizens (in short, they are treated as any other data type, thus assignable to a variable).
- Can be defined as closures, thus capturing surrounding declared variables that can be accessed and modified even after the surrounding function returned.
- Variable number of arguments with the ellipsis syntax (`...`).
- Return values can be named

### Coroutines

Coroutines (aka goroutines) are lightweight threads of execution, managed by the Go runtime instead of the OS.
Therefore, concurrent programming is natively supported. Communication is commonly achieved with `channels`.

Take the following naive example:
```go
package main

import (
    "fmt"
    "time"
)

func foo() {
    for i := 0; i < 5; i++ {
        fmt.Println("foo", i)
        time.Sleep(100 * time.Millisecond)
    }
}

func bar() {
    for i := 0; i < 5; i++ {
        fmt.Println("bar", i)
        time.Sleep(100 * time.Millisecond)
    }
}

func main() {
    // Start the goroutines in parallel
    go foo()
    go bar()

    // Naive approach for awaiting the goroutines to finish
    time.Sleep(1 * time.Second)

    fmt.Println("Done")
}
```

The functions `foo` and `bar` are executed in parallel, just by utilising the native language `go` statement.

> ⚠️ Please note that using a `time.Sleep`, to assume that both functions had enough time to complete, is an improper way of solving the problem.
> The idea is to avoid including other concepts simultaneously.

### Closures

Combined with the fact that functions in Go are first-class citizens, closures are a language feature that bring a lot of power and flexibility.
With such capability, it becomes possible to create functions that capture and maintain a state, behave like objects, facilitate functional and concurrency logic.

Take the following example, that defines a simple sequencer:
```go
package main

import (
	"fmt"
)

func sequencer() func() int {
    i := 0
    return func () int {
        i++
        return i
    }
}

func main() {
    seq1 := sequencer()

    fmt.Println(seq1()) // prints 1
    fmt.Println(seq1()) // prints 2
    fmt.Println(seq1()) // prints 3

    seq2 := sequencer()
    fmt.Println(seq2()) // prints 1
}
```

With this library you will often be leveraging from this language feature, by passing anonymous functions over. Such approach will require to follow a certain signature but will simplify integrating your code seamlessly without too much boilerplate or obscure code.

### Defer

The `defer` statement in Go allows you to schedule a function call to be executed later, when the surrounding function completes. It is often used to ensure that some cleanup code is executed after a function completes, regardless of the path that led to the function's exit.

Take this revised example where the `defer` statement is utilised to signal the wait group the function completed its execution:

```go
package main

import (
    "fmt"
    "sync"
    "time"
)

func foo(wg *sync.WaitGroup) {
    defer wg.Done()
    for i := 0; i < 5; i++ {
        fmt.Println("Foo:", i)
        time.Sleep(100 * time.Millisecond)
    }
}

func bar(wg *sync.WaitGroup) {
    defer wg.Done()
    for i := 0; i < 5; i++ {
        fmt.Println("Bar:", i)
        time.Sleep(100 * time.Millisecond)
    }
}

func main() {
    // Define a wait group...
    var wg sync.WaitGroup
    // ... with 2 slots
    wg.Add(2)
    
    // Start the operations in parallel
    go foo(&wg)
    go bar(&wg)
    
    // Wait for completion
    wg.Wait()
    
    fmt.Println("Done!")
}
```

This example is too simple to show the real value of using `defer`, but it tries to illustrate it's purpose. With this mechanism you can be rest assured that, when the function terminates, the statement will always be executed. A more complex execution pattern, with multiple return points, would show off better the value of this statement (since you would not be hard-pressed, to make sure in each return block, the wait group is signalled properly).

## Wait Groups

Wait Groups are a mechanism to allow any other goroutine to wait for a group of goroutines to complete before continuing execution.

Take this revised example, where the parallel execution of `foo` and `bar` is awaited using this construct:
```go
package main

import (
    "fmt"
    "sync"
    "time"
)

func foo(wg *sync.WaitGroup) {
    for i := 0; i < 5; i++ {
        fmt.Println("Foo:", i)
        time.Sleep(100 * time.Millisecond)
    }
    wg.Done()
}

func bar(wg *sync.WaitGroup) {
    for i := 0; i < 5; i++ {
        fmt.Println("Bar:", i)
        time.Sleep(100 * time.Millisecond)
    }
    wg.Done()
}

func main() {
    // Define a wait group...
    var wg sync.WaitGroup
    // ... with 2 slots
    wg.Add(2)
    
    // Start the operations in parallel
    go foo(&wg)
    go bar(&wg)
    
    // Wait for completion
    wg.Wait()
    
    fmt.Println("Done!")
}
```

You will not likely be often required to rely on Wait Groups interacting with this library. Most of the constructs make your life easier with concurrency or by communicating between them, even if running in the background or in parallel.

One use case, however, is if you want to be sure that when publishing an event, into the `eventhub.EventHub` construct, all subscribers got the chance to execute until completion.

> 👉 Normally you just fire and forget events, this is just meant for special scenarios where you need to be sure everything finished before continuing.

e.g.:
```go
func main(){
    ...
    wg := hub.Publish(SomeEvent{}, nil)

    // Wait for all subscribers to complete their execution
    wg.Wait()
}
```

## Channels

Channels are thread-safe data structures for communicating between goroutines.

Take the following revised example:
```go
package main

import (
    "fmt"
    "sync"
    "time"
)

func foo(ch chan<- bool) {
    for i := 0; i < 5; i++ {
        fmt.Println("Foo:", i)
        time.Sleep(100 * time.Millisecond)
    }
    ch <- true
}

func bar(ch chan<- bool) {
    for i := 0; i < 5; i++ {
        fmt.Println("Bar:", i)
        time.Sleep(100 * time.Millisecond)
    }
    ch <- true
}

func main() {
    // Define a channel
    ch := make(chan bool)
    
    // Start the operations in parallel
    go foo(ch)
    go bar(ch)
    
    // Wait for completion
    <-ch
    <-ch
    
    fmt.Println("Done!")
}
```

> ⚠️ This would not be the proper way to solve the problem (actually the Wait Group example would be the more appropriate way). This is just meant to demonstrate the usage of a channel.

This library will likely prevent you to often need to utilise `channels`. It is key, however, to understand this mechanism and its important role on implementing synchronisation. 

### More on channels

Take the following naive example to compute values in parallel, using a channel to communicate back the result:

```go
package main

import (
    "fmt"
    "sync"
)

func main() {
    // a wait group to wait for all results
    var wg sync.WaitGroup

    // a channel to pass by resuld
    ch := make(chan int)

    // launch 5 goroutines for calculating the square route of a number
    for i := 0; i < 5; i++ {
        wg.Add(1) // inform the wait group we have one more execution to await for
        go func(n int) {
            ch <- n * n // send the squared value of n through the channel
            wg.Done() // signal that this goroutine is done
        }(i)
    }

    // wait, in the background, for all calculations to finish before closing the channel
    go func() {
        wg.Wait() // wait for all goroutines to finish
        close(ch) // close the channel
    }()

    // print results
    for val := range ch {
        fmt.Println(val)
    }
}
```

Though still academic, this example should give you a better understanding about the role of channel and of a Wait Group. 

### Switch case statement on steroids

In Go, the switch statement is a versatile control structure, with unique features that familiarity with other languages might leave  one at loss when having a first contact. The keyword is `select` and this is the least odd difference one might spot.

Take the following example, many other languages offer equivalent behavior:

```go
package main

import "fmt"

func main() {
    name := "Gopher"
    switch name {
    case "Gopher":
        fmt.Println("Hello, Gopher!")
    case "World":
        fmt.Println("Hello, World!")
    default:
        fmt.Println("Hello, stranger!")
    }
}
```

What makes Go special is the capability of dealing with channels, so this statement is crucial to work with such data structures.

Take the following example:
```go
package main

import (
	"fmt"
)

func main() {
	ch1 := make(chan string)
	ch2 := make(chan int)
	ch3 := make(chan bool)

	go func() {
		defer close(ch1)
		ch1 <- "Hello"
	}()

	go func() {
		defer close(ch2)
		ch2 <- 42
	}()

	ch1Closed := false
	ch2Closed := false

	defer close(ch3)

	// when the count reach 0, cancel the context
	go func() {
		for {
			if ch1Closed && ch2Closed {
				break
			}
		}

		ch3 <- true
	}()

	for {
		select {
		case msg, ok1 := <-ch1:
			if ok1 {
				fmt.Println("Received number from ch1:", msg)
			} else {
				fmt.Println("Channel 1 closed")
				ch1Closed = true
			}
		case num, ok2 := <-ch2:
			if ok2 {
				fmt.Println("Received number from ch2:", num)
			} else {
				fmt.Println("Channel 2 closed")
				ch2Closed = true
			}
		case <-ch3:
			fmt.Println("All channels closed")
			return
		}
	}
}
```

In a nutshell, the previous example uses two channels to send messages and when both are closed we send another message to stop the execution.

> ⚠️ Notice how it can be verified when a channel becomes closed and how he are using it to toggle the flags.
> 
> One might consider using a counter, keeping track of the open channels, when reaching 0 it could be utilised to signal termination. You might be surprised that the messages stating that a channel got called can be called multiple times.
> Be aware of this behavior!
> 
> Also, please understand this is just for illustrating purposes, it would not be a good implementation of a Parallel Aggregation pattern, neither the proper way to deal with cancellation.

## Context

In a nutshell, the standard data structure `context.Context` is essential to manage the lifecycle of long-running operations.
It is very versatility and heavily used by this library.

This section does not intend to be exhaustive in the topic, only bring awareness to key aspects.

### Cancellation

Dealing with cancellation is paramount to EDT and the most relevant use cases are:
- Timeouts / Deadline enforcement
- Graceful shutdown

#### Timeouts / Deadline enforcement

The distinction between a timeout and deadline is that the latter is defined by an exact point in time for it to be enforced, rather than setting a duration that must elapse for the cancellation condition to be met. 
For time based cancellation this library offers constructs that assist on this concept (e.g. `Expirable` and `Delayable`).

A timeout based cancellation:
```go
parent := context.Background()
ctx, cancel := context.WithTimeout(parent, 10*time.Second)
defer cancel()
```

A deadline based cancellation:
```go
parent := context.Background()
ctx, cancel := context.WithDeadline(parent, time.Now().Add(10*time.Second))
defer cancel()
```

Effectively both samples have equivalent behavior (it is just a matter of expressiveness that might best fit a use case or another).

#### Graceful shutdown

This behavior can be achieved by hooking the `os.Interrupt`, the standard library provides a convenient decorator that does exactly that:

```go
ctx, cancel := context.WithSignal(context.Background(), os.Interrupt)
defer cancel()
```
This library provides you a construct (`Director`) that assists with this responsibility, but it is important to be aware about this recipe in case you have a strong reason for not propagating the scope.

### Data Propagation

In Go, the `context.Context` fulfils an important role by providing a standard way to propagate request-scoped values (or just call it data) across different parts of a program.
There are plenty of use cases where this mechanism can be of use to avoid passing down explicit parameters with specific types across multiple layers so a parameter can be of utilised only on a limited or even single pipeline execution stage (e.g. image getting an authenticate user on a HTTP request, passing it down as a context so a certain authorization verification can be made at a certain Service).

```go
package main

import (
	"context"
	"log"
)

func main() {
	ctx := context.WithValue(context.Background(), "userID", "42")

	doSomething(ctx)
}

func doSomething(ctx context.Context) {
    userID, ok := ctx.Value("userID").(string)
	
    if !ok {
        log.Println("No user defined")
        return
    }

    log.Println("Doing something for user", userID)
}
```

> 👉 This library offers you a construct, the `Injector`, that rely on this very same mechanism. You will get, however, a more convenient and powerful way to retrieve values from your context, by satisfying dependencies. 


## Handling Time

There are certainly a lot to say about handling "Time" in Go, but we will only focus on the key aspects that matter most for EDT.

### Duration

A duration (or `time.Duration`) is utilised to express a length of time that can perhaps represent a delay or a timeout (depending on the context it is applied).
This standard data structure has a key part in most of EDT's constructs.

If passing down hardcoded values you will be doing expressions like:
 - `100 * time.Millisecond` (100ms)
 - `5 * time.Decond` (5s)

What might be something you need to do, that might not be easy to remember is doing the same but with a variable (lets say you refactor your code and instead of the hardcoded value you start using an environment variable for it).

Take the following snippet:
```go
import "time"

var durationInSeconds int = 5
duration := time.Duration(durationInSeconds) * time.Second
```

### Delay

If enforcing a delay, on a certain behavior you are better avoiding, most of the time, the usage of `time.Sleep` due to its blocking nature of the goroutine.
Prefer the usage of `time.After` instead, since it relies on a channel that you can combine with other channels and, most importantly, a cancellation context.

This technique is relevant for dealing properly with graceful (and quick) shutdowns of your application.

Take the following example:
```go
package main

import (
    "context"
    "fmt"
    "time"
)

func main() {
    // Create a cancellation context with a timeout of 2 seconds
    ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
    defer cancel()

    fmt.Println("Starting the timer")

    // Wait for either the context to be cancelled or the timer to expire
    select {
    case <-ctx.Done():
        fmt.Println("Cancelled:", ctx.Err())
    case <-time.After(5 * time.Second):
        fmt.Println("Timer expired")
    }
}
```

> 👉 This library will assist you in common patterns dealing with time delays facilitating the usage of context propagation.