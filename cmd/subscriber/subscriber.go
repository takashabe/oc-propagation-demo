package main

import (
	"context"
	"os"
	"time"

	"cloud.google.com/go/pubsub"
	"github.com/takashabe/oc-propagation-demo/internal/exporter"
	"github.com/takashabe/oc-propagation-demo/internal/propagation"
	"go.opencensus.io/trace"
)

func main() {
	exporter.InitStackdriver(os.Getenv("PROJECT_ID"))

	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, os.Getenv("PROJECT_ID"))
	if err != nil {
		panic(err)
	}

	subscription := client.Subscription(os.Getenv("PUBSUB_SUBSCRIPTION"))
	subscription.Receive(ctx, func(ctx, msg *pubsub.Message) error {
		sc := propagation.SpanContextFromMessage(msg)
		ctx, span := trace.StartSpanWithRemoteParent(ctx, "receive", sc)

		time.Sleep(100 * time.Millisecond)

		span.End()
		return nil
	})
}
