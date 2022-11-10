package telemetry

import (
	"context"
	"runtime"

	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

// StartSpan creates a new Open Telemetry tracer and starts a span with the given name and returns it.
// the context is updated with the new span.
func StartSpan(ctx *context.Context, tracerName string, name string, opts ...trace.SpanStartOption) trace.Span {

	if otel.GetTracerProvider() == nil {
		logrus.Warn("otel tracer provider not set")
	}
	tr := otel.Tracer(tracerName)
	if tr == nil {
		logrus.Warn("otel tracer not set")
	}

	if ctx == nil {
		newBackCtx := context.Background()
		ctx = &newBackCtx
	}

	newCtx, span := tr.Start(*ctx, name, opts...)
	span.SetAttributes(attribute.Key("function").String(GetCallerFuncName()))
	*ctx = newCtx
	return span
}

// ReInjectContext carries the trace context in a new background context
func ReInjectContext(ctx *context.Context) {
	// carrying the trace context in a new background context
	carrier := propagation.MapCarrier{}
	traceContext := propagation.TraceContext{}
	traceContext.Inject(*ctx, carrier)
	newCtx := traceContext.Extract(context.Background(), carrier)
	*ctx = newCtx
}

// GetCallerFuncName returns the name of the function 2 levels higher in the stack.
func GetCallerFuncName() string {
	fpcs := make([]uintptr, 1)

	// Skip 3 levels to get the caller
	runtime.Callers(3, fpcs)

	frame, _ := runtime.CallersFrames(fpcs).Next()

	return frame.Function
}
