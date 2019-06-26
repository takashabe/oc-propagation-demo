package main

import (
	"context"
	"os"

	"cloud.google.com/go/pubsub"
	"github.com/takashabe/oc-propagation-demo/internal/exporter"
	"github.com/takashabe/oc-propagation-demo/internal/propagation"
	"go.opencensus.io/trace"
)

func main() {
	exporter.InitStackdriver(os.Getenv("PROJECT_ID"))

	ctx, _ := trace.StartSpan(context.Background(), "publish")

	client, err := pubsub.NewClient(ctx, os.Getenv("PROJECT_ID"))
	if err != nil {
		panic(err)
	}

	topic := client.Topic(os.Getenv("PUBSUB_TOPIC"))
	msg := &pubsub.Message{
		Data: []byte("oc-demo"),
	}
	if _, err := topic.Publish(ctx, propagation.WrapMessage(ctx, msg)).Get(ctx); err != nil {
		panic(err)
	}
}
