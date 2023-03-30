# Welcome

The Event Driven Toolkit (EDT for short) is your swiss army knife for event driven applications, Go edition.
Here you find implementation support for common patterns when developing such applications.

## Motivation

For several years, more than a decade certainly, found my self repeating some patterns for each programming language happened to be using at the time. Go is one of those languages where most recently had to repeat those trusted recipes in.

> 🚧 This is also learning experience with Go, some rough edges are to be expected.

## Designed to improve Code Quality

Having common constructs is a great way to improve your code. Specially so if, by design, they aim for improving the code itself.

Go has very particular and very powerful language capabilities for parallelism and exchanging data like `coroutines` and `channels`. Language constructs like these reduce immensely the complexity of some tasks but the implementation may result in messy code if not properly refactored and some effort is put into readability.
This library aims for that particular goal: Making Code Readable and, by consequence, improve Code Quality.

## Event Hub

This is just one of the Constructs you can find within this library, but it is probably the most important one. By definition, an Event Driven application deals with Publishing and Subscribing events, this Construct fulfils just that.

What also makes this library special is that most Constructs work well together.

> 👉 You can control a `State Machine` state changes by publishing events or even have a task on your application that is blocked until an event is dispatched (using the `Expectable` construct).