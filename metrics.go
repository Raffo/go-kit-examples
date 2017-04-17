package main

import (
	"time"

	"github.com/go-kit/kit/metrics"
)

func instrumentingMiddleware(
	requestCount metrics.Counter,
	requestLatency metrics.Histogram,
	countResult metrics.Histogram,
) ServiceMiddleware {
	return func(next EchoService) EchoService {
		return instrmw{requestCount, requestLatency, countResult, next}
	}
}

type instrmw struct {
	requestCount   metrics.Counter
	requestLatency metrics.Histogram
	countResult    metrics.Histogram
	EchoService
}

func (mw instrmw) Echo(s string) (output string) {
	defer func(begin time.Time) {
		lvs := []string{"method", "echo", "error"}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	output = mw.EchoService.Echo(s)
	return
}
