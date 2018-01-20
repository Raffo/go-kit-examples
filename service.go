package main

import (
	"fmt"
	"strings"
)

// EchoService is the interface containing the definition of all the function of our service
type EchoService interface {
	Echo(string) string
	Uppercase(string) string
}

type echoService struct {
	base string
}

// Echo this really needs a value receiver
func (e echoService) Echo(s string) string {
	return fmt.Sprintf("%s, %s", e.base, s)
}

func (e echoService) Uppercase(s string) string {
	return strings.ToUpper(s)
}

// ServiceMiddleware is needed only to be able to chain middlewares
type ServiceMiddleware func(EchoService) EchoService
