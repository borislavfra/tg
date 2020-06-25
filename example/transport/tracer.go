// GENERATED BY 'T'ransport 'G'enerator. DO NOT EDIT.
package transport

import (
	"encoding/json"
	"net/http"
	"os"
	"strings"

	otg "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	zipkinTracer "github.com/openzipkin-contrib/zipkin-go-opentracing"
	"github.com/openzipkin/zipkin-go"
	httpReporter "github.com/openzipkin/zipkin-go/reporter/http"
	"github.com/savsgio/gotils"
	"github.com/sirupsen/logrus"
	"github.com/uber/jaeger-client-go/config"
	"github.com/uber/jaeger-client-go/log"
	"github.com/uber/jaeger-lib/metrics"
	"github.com/valyala/fasthttp"
)

func (srv *Server) TraceJaeger(serviceName string) *Server {

	environment, _ := os.LookupEnv("ENV")

	cfg, err := config.FromEnv()
	ExitOnError(srv.log, err, "jaeger config err")

	if cfg.ServiceName == "" {
		cfg.ServiceName = environment + serviceName
	}

	var trace otg.Tracer
	trace, srv.reporterCloser, err = cfg.NewTracer(config.Logger(log.NullLogger), config.Metrics(metrics.NullFactory))

	ExitOnError(srv.log, err, "could not create jaeger tracer")

	otg.SetGlobalTracer(trace)
	return srv
}

func (srv *Server) TraceZipkin(serviceName string, zipkinUrl string) *Server {

	reporter := httpReporter.NewReporter(zipkinUrl)
	srv.reporterCloser = reporter

	environment, envExists := os.LookupEnv("ENV")

	if envExists {
		serviceName = environment + serviceName
	}

	endpoint, err := zipkin.NewEndpoint(serviceName, "")
	ExitOnError(srv.log, err, "could not create endpoint")

	nativeTracer, err := zipkin.NewTracer(reporter, zipkin.WithLocalEndpoint(endpoint))
	ExitOnError(srv.log, err, "could not create tracer")

	trace := zipkinTracer.Wrap(nativeTracer)
	otg.SetGlobalTracer(trace)

	return srv
}

func injectSpan(log logrus.FieldLogger, span otg.Span, ctx *fasthttp.RequestCtx) {

	headers := make(http.Header)

	if err := otg.GlobalTracer().Inject(span.Context(), otg.HTTPHeaders, otg.HTTPHeadersCarrier(headers)); err != nil {
		log.WithError(err).Warning("inject span to HTTP headers")
	}

	for key, values := range headers {
		ctx.Response.Header.Set(key, strings.Join(values, ";"))
	}
}

func extractSpan(log logrus.FieldLogger, opName string, ctx *fasthttp.RequestCtx) (span otg.Span) {

	headers := make(http.Header)

	ctx.Request.Header.VisitAll(func(key, value []byte) {
		headers.Set(gotils.B2S(key), gotils.B2S(value))
	})

	var opts []otg.StartSpanOption
	wireContext, err := otg.GlobalTracer().Extract(otg.HTTPHeaders, otg.HTTPHeadersCarrier(headers))

	if err != nil {
		log.WithError(err).Warning("extract span from HTTP headers")
	} else {
		opts = append(opts, otg.ChildOf(wireContext))
	}

	span = otg.GlobalTracer().StartSpan(opName, opts...)

	ext.HTTPUrl.Set(span, ctx.URI().String())
	ext.HTTPMethod.Set(span, gotils.B2S(ctx.Method()))

	return
}

func toString(value interface{}) string {
	data, _ := json.Marshal(value)
	return gotils.B2S(data)
}
