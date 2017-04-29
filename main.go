package main

import (
	"flag"
	"net/http"
	"os"

	"github.com/go-kit/kit/log"
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	// we put args always in the main so that they cannot be used by accident as global vars
	var (
		listen = flag.String("listen", ":8080", "HTTP listen address")
	)
	flag.Parse()

	var logger log.Logger
	logger = log.NewLogfmtLogger(os.Stderr)
	logger = log.With(logger, "listen", *listen, "caller", log.DefaultCaller)

	var svc EchoService
	svc = echoService{"Hello"}
	// middlewares
	svc = loggingMiddleware(logger)(svc)

	// TODO understand, is this really the right place to plug the middlewares?
	// adding instrumentation
	fieldKeys := []string{"method", "error"}
	requestCount := kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
		Namespace: "group",
		Subsystem: "echo_service",
		Name:      "request_count",
		Help:      "Number of requests received.",
	}, fieldKeys)
	requestLatency := kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
		Namespace: "group",
		Subsystem: "echo_service",
		Name:      "request_latency_microseconds",
		Help:      "Total duration of requests in microseconds.",
	}, fieldKeys)
	countResult := kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
		Namespace: "group",
		Subsystem: "echo_service",
		Name:      "count_result",
		Help:      "The result of each count method.",
	}, []string{})

	svc = instrumentingMiddleware(requestCount, requestLatency, countResult)(svc)

	var h http.Handler
	{
		h = MakeHTTPHandler(svc, log.With(logger, "component", "HTTP"))
	}

	errs := make(chan error)

	go func() {
		http.Handle("/metrics", promhttp.Handler())
		http.ListenAndServe(":7979", nil)
	}()

	go func() {
		logger.Log("transport", "HTTP", "addr", ":8080")
		errs <- http.ListenAndServe(":8080", h)
	}()

	logger.Log("exit", <-errs)
}
