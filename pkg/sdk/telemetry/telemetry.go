package telemetry

import (
	"context"
	"runtime"

	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// StartSpan creates a new Open Telemetry tracer and starts a span with the given name and returns it.
// the context is updated with the new span.
func StartSpan(ctx *context.Context, tracerName string, name string) trace.Span {

	if otel.GetTracerProvider() == nil {
		logrus.Warn("otel tracer provider not set")
		return nil
	}
	tr := otel.Tracer(tracerName)
	if tr == nil {
		logrus.Warn("otel tracer not set")
		return nil
	}

	if ctx == nil {
		newCtx := context.Background()
		ctx = &newCtx
	}

	newCtx, span := tr.Start(*ctx, name)
	span.SetAttributes(attribute.Key("function").String(GetCallerFuncName()))
	*ctx = newCtx
	return span
}

// GetCallerFuncName returns the name of the function 2 levels higher in the stack.
func GetCallerFuncName() string {
	fpcs := make([]uintptr, 1)

	// Skip 3 levels to get the caller
	runtime.Callers(3, fpcs)

	frame, _ := runtime.CallersFrames(fpcs).Next()

	return frame.Function
}
