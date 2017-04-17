# Go Kit example 

This repo exists just to play with go-kit, to get a better understanding of what it can offer. 

- First commit contains the most basic example ever that you can build.

- Second commit have additional logging middlewares and corrects some errors from the previous example

- Third commit added instrumentation / metrics


## A "review"

The whole idea is: 

- stuff is modular and pluggable 
- we wire things in the main of the application. This includes logging, instrumentation and definition of routes 
- we create one handler per route. The handlers are initialized with the NewServer function
- All the logic is cleanly separated in `service.go`. This means it's easy to rewrite the rest.
- Being mostly interested in the HTTP part, having logic on error code in `transport.go` sounds a bit weird.
- The middleware are tied to the interface they are wrapping. This means that it is probably better to have only one interface for all the operations of our microservice.