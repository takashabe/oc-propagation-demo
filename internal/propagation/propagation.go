package propagation

import (
	"context"
	"encoding/hex"
	"fmt"

	"cloud.google.com/go/pubsub"
	"go.opencensus.io/trace"
)

// distributed tracing meta fields
const (
	TraceIDField = "X-Pubsub-TraceId"
	SpanIDField  = "X-Pubsub-SpanId"
	SampledField = "X-Pubsub-Sampled"
)

// SpanContextFromMessage extracts a pubsub span context from message.
func SpanContextFromMessage(m *pubsub.Message) (trace.SpanContext, bool) {
	if m.Attributes == nil {
		return trace.SpanContext{}, false
	}

	tid, ok := parseTraceID(m.Attributes[TraceIDField])
	if !ok {
		return trace.SpanContext{}, false
	}
	sid, ok := parseSpanID(m.Attributes[SpanIDField])
	if !ok {
		return trace.SpanContext{}, false
	}
	sampled := parseSampled(m.Attributes[SampledField])

	return trace.SpanContext{
		TraceID:      tid,
		SpanID:       sid,
		TraceOptions: sampled,
	}, true
}

func parseTraceID(t string) (trace.TraceID, bool) {
	if t == "" {
		return trace.TraceID{}, false
	}
	b, err := hex.DecodeString(t)
	if err != nil {
		return trace.TraceID{}, false
	}
	tid := trace.TraceID{}
	copy(tid[:], b)
	return tid, true
}

func parseSpanID(s string) (trace.SpanID, bool) {
	if s == "" {
		return trace.SpanID{}, false
	}
	b, err := hex.DecodeString(s)
	if err != nil {
		return trace.SpanID{}, false
	}
	sid := trace.SpanID{}
	copy(sid[:], b)
	return sid, true
}

func parseSampled(sampled string) trace.TraceOptions {
	if sampled == "true" {
		return trace.TraceOptions(1)
	}
	return trace.TraceOptions(0)
}

// WrapMessage embed a span context to message.
func WrapMessage(ctx context.Context, m *pubsub.Message) *pubsub.Message {
	if m.Attributes != nil {
		if m.Attributes[TraceIDField] != "" || m.Attributes[SpanIDField] != "" {
			return m
		}
	} else {
		m.Attributes = make(map[string]string, 3)
	}
	sc := trace.FromContext(ctx).SpanContext()
	m.Attributes[TraceIDField] = sc.TraceID.String()
	m.Attributes[SpanIDField] = sc.SpanID.String()
	m.Attributes[SampledField] = fmt.Sprintf("%t", sc.IsSampled())
	return m
}
