package telemetry

import (
	"context"
	"runtime"
	"sync"

	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

var mut sync.Mutex

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

	mut.Lock()
	defer mut.Unlock()

	if ctx == nil {
		newBackCtx := context.Background()
		ctx = &newBackCtx
	}

	newCtx, span := tr.Start(*ctx, name, opts...)
	span.SetAttributes(attribute.Key("function").String(GetCallerFuncName()))

	// Add memory usage
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	span.SetAttributes(attribute.Key("m.Alloc").Int64(int64(m.Alloc)))
	span.SetAttributes(attribute.Key("m.TotalAlloc").Int64(int64(m.TotalAlloc)))
	span.SetAttributes(attribute.Key("m.NumGC").Int64(int64(m.NumGC)))

	*ctx = newCtx
	return span
}

// ReInjectContext carries the trace context in a new background context
func ReInjectContext(ctx *context.Context) {
	// carrying the trace context in a new background context
	newCtx := CarryTelemetryContext(*ctx)
	*ctx = newCtx
}

// CarryTelemetryContext carries the trace context in a new background context
func CarryTelemetryContext(ctx context.Context) context.Context {
	// carrying the trace context in a new background context
	carrier := propagation.MapCarrier{}
	traceContext := propagation.TraceContext{}
	traceContext.Inject(ctx, carrier)
	newCtx := traceContext.Extract(context.Background(), carrier)
	return newCtx
}

// GetCallerFuncName returns the name of the function 2 levels higher in the stack.
func GetCallerFuncName() string {
	fpcs := make([]uintptr, 1)

	// Skip 3 levels to get the caller
	runtime.Callers(3, fpcs)

	frame, _ := runtime.CallersFrames(fpcs).Next()

	return frame.Function
}
