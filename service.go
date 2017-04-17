package main

import "fmt"

// Echoer is the interface containing the definition of all the function of our service
type Echoer interface {
	Echo(string) string
}

type echoService struct {
	base string
}

func (e *echoService) Echo(s string) string {
	return fmt.Sprintf("%s, %s", e.base, s)
}

// ServiceMiddleware is needed only to be able to chain middlewares
type ServiceMiddleware func(echoService) echoService
