# A Golang primer

To make the most of this library, it's essential to have a solid understanding of Go's capabilities. Familiarity with certain language and standard library features is crucial for achieving this goal.

It is out of scope of this document providing a complete Go tutorial, it is expected that you have some familiarity already or have other means to cover that requirement.

This document aims to make sure you have a solid understanding of some key aspects.

## Coroutines

Coroutines (aka goroutines) are lightweight threads of execution, managed by the Go runtime instead of the OS.
Therefore, concurrent programming is natively supported. Communication is achieved with channels.

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

Functions `foo` and `bar` are executed in parallel just by utilising the native language syntax `go` statement.

> ⚠️ Please note that using a `time.Sleep` that assumes that both functions had enough time to complete is an improper way of solving the problem.
> It is just intended to avoid including other concepts at the same time.

## Wait Groups

Wait Groups are a mechanism to allow any other goroutine to wait for a group of goroutines to complete before continuing execution.

Take this revised example where the parallel execution of `foo` and `bar` is awaited using this language construct:
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

You will not be often required to used Wait Groups with this library because most of the constructs make your life easier with concurrent or by communicating between constructs running in the background / in parallel.

One use case, however, is if you want to be sure that when publishing an event, into the `eventhub.EventHub` construct, all subscribers got the chance to execute to completion.

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

## Defer

The `defer` statement in Go allows you to schedule a function call to be executed later, when the surrounding function completes. It is often used to ensure that some cleanup code is executed after a function completes, regardless of the path that led to the function's exit.

Take this revised example where the defer is utilised to signal the wait group the function completed its execution:
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

This example is too simple to show the real value of using `difer`, but it illustrates the purpose. With this mechanism you can be rest assured that, when the function terminates the statement will always be  executed. A more complex execution pattern with multiple return points would show off better the value of this statement (since you would not need to make sure on each return block that the wait group was signalled properly).

## Channels

Channels are thread-safe data structures for communicating between goroutines.

## Closures

🚧

## Context

🚧

## Time

🚧

