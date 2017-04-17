package main

import (
	"time"

	"github.com/go-kit/kit/log"
)

func loggingMiddleware(logger log.Logger) ServiceMiddleware {
	return func(next EchoService) EchoService {
		return logMiddleware{logger, next}
	}
}

// logmw wraps our original middleware
type logMiddleware struct {
	logger log.Logger
	EchoService
}

func (mw logMiddleware) Echo(s string) (output string) {
	defer func(begin time.Time) {
		_ = mw.logger.Log(
			"method", "echo",
			"input", s,
			"output", output,
			"took", time.Since(begin),
		)
	}(time.Now())

	output = mw.EchoService.Echo(s)
	return
}
