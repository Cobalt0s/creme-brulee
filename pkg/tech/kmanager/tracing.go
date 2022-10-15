package kmanager

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Cobalt0s/creme-brulee/pkg/stateful/logging"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus/ctxlogrus"
	"go.opentelemetry.io/otel/trace"
	"strconv"
)

type Trace struct {
	TraceID    string `json:"traceId"`
	SpanID     string `json:"spanId"`
	TraceFlags string `json:"traceFlags"`
	TraceState string `json:"traceState"`
	Remote     string `json:"remote"`
}

func ExtractTrace(ctx context.Context) *Trace {
	spanCtx := trace.SpanContextFromContext(ctx)

	return &Trace{
		TraceID:    spanCtx.TraceID().String(),
		SpanID:     spanCtx.SpanID().String(),
		TraceFlags: spanCtx.TraceFlags().String(),
		TraceState: spanCtx.TraceState().String(),
		Remote:     "true",
	}
}

func (t *Trace) toSpanContext() (*trace.SpanContext, error) {
	traceID, err := hex.DecodeString(t.TraceID)
	if err != nil {
		return nil, createUnmarshalErr("TraceID", err)
	}
	spanID, err := hex.DecodeString(t.SpanID)
	if err != nil {
		return nil, createUnmarshalErr("SpanID", err)
	}
	traceFlags, err := hex.DecodeString(t.TraceFlags)
	if err != nil {
		return nil, createUnmarshalErr("TraceFlags", err)
	}
	traceState, err := trace.ParseTraceState(t.TraceState)
	if err != nil {
		return nil, createUnmarshalErr("TraceState", err)
	}
	remote, err := strconv.ParseBool(t.Remote)
	if err != nil {
		return nil, createUnmarshalErr("Remote", err)
	}

	spanCtx := trace.NewSpanContext(trace.SpanContextConfig{
		TraceID:    convertToByte16(traceID),
		SpanID:     convertToByte8(spanID),
		TraceFlags: trace.TraceFlags(convertToByte1(traceFlags)),
		TraceState: traceState,
		Remote:     remote,
	})
	return &spanCtx, nil
}

// TraceFromEvent deprecated
func TraceFromEvent(ctx context.Context, event []byte) (context.Context, trace.Span) {
	return TraceFromEventNamed(ctx, event, "EventReceived")
}

func TraceFromEventNamed(ctx context.Context, event []byte, name string) (context.Context, trace.Span) {
	spanCtx, err := traceFromEventInternal(event)
	if err != nil {
		log := ctxlogrus.Extract(ctx)
		log.Errorf("failed extracting trace from event %v", err)
		return logging.StartSpan(ctx, name)
	}
	ctx = trace.ContextWithRemoteSpanContext(ctx, *spanCtx)
	return logging.StartSpan(ctx, name)
}

func traceFromEventInternal(event []byte) (*trace.SpanContext, error) {
	t := &struct {
		Trace *Trace `json:"trace"`
	}{}
	if err := json.Unmarshal(event, t); err != nil {
		return nil, err
	}
	if t.Trace == nil {
		return nil, errors.New("trace is missing in event")
	}
	return t.Trace.toSpanContext()
}

func convertToByte16(source []byte) [16]byte {
	var target [16]byte
	copy(target[:], source)
	return target
}
func convertToByte8(source []byte) [8]byte {
	var target [8]byte
	copy(target[:], source)
	return target
}
func convertToByte1(source []byte) byte {
	var target [1]byte
	copy(target[:], source)
	return target[0]
}

func createUnmarshalErr(target string, err error) error {
	return fmt.Errorf("failed converting %v due %v", target, err)
}
