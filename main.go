package main

import (
	"flag"
	"net/http"
	"os"

	"github.com/go-kit/kit/log"
	httptransport "github.com/go-kit/kit/transport/http"
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

	echoHandler := httptransport.NewServer(
		makeEchoEndpoint(svc),
		decodeEchoRequest,
		encodeResponse,
		httptransport.ServerErrorLogger(logger),
	)
	http.Handle("/echo", echoHandler)
	logger.Log("msg", "HTTP server started", "addr", *listen)
	logger.Log("err", http.ListenAndServe(*listen, nil))
}
