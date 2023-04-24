# Welcome

The Event Driven Toolkit (EDT for short) is your swiss army knife for event driven applications, Go edition.
Here you find implementation support for common patterns when developing such applications.

## Motivation

For several years, more than a decade certainly, found myself repeating some patterns for each programming language happened to be using at the time. Go is one of those languages where most recently had to repeat those trusted recipes in.

## Why?

### Designed to improve Code Quality

Having common constructs is a great way to improve your code.

Go has very powerful language capabilities for parallelism and exchanging data, like `coroutines` and `channels`. Such native language constructs reduce immensely the complexity of some tasks but it is equally easy to produce messy code if there is not some effort made into readability.
This library aims for that particular goal: Making code readable, simpler and, consequently, improve the overall Code Quality.

> ⚠️ Disclaimer: There is no hidden claim trying to imply that code produced in this language is prone to get messy. By the contrary, Go is a very elegant language. Newcomers might however, become overwhelmed with some capabilities or use them ineffectively. 

### Event Hub

This is just one of the Constructs you can find within this library, but it is probably the most important one. By definition, an Event Driven application deals with Publishing and Subscribing events, this Construct fulfils just that.

### Consistency

What also makes this library special is that most Constructs work well together.

> 👉 You can control a `State Machine` state changes by publishing events or even have a task on your application that is blocked until an event is dispatched (using the `Expectable` construct).


### Leveraging from Go

This library is intended to leverage from Go, not just hide "features that might look scary for new developers" under useless abstractions.
It is actually important to have a clear understanding of this language capabilities to effectively utilise EDT.
