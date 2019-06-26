package main

import (
	"context"
	"log"
	"net/http"
	"os"

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

	go func() {
		err := subscription.Receive(ctx, func(ctx context.Context, msg *pubsub.Message) {
			defer msg.Ack()
			sc := propagation.SpanContextFromMessage(msg)
			_, span := trace.StartSpanWithRemoteParent(ctx, "receive", sc)
			span.AddAttributes(
				trace.StringAttribute("TraceID", sc.TraceID.String()),
				trace.StringAttribute("SpanID", sc.SpanID.String()),
			)
			span.End()
		})
		if err != nil {
			panic(err)
		}
	}()

	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handler(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(200)
}
