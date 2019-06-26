package main

import (
	"fmt"
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

	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		ctx, span := trace.StartSpan(req.Context(), "publish")

		client, err := pubsub.NewClient(ctx, os.Getenv("PROJECT_ID"))
		if err != nil {
			panic(err)
		}
		topic := client.Topic(os.Getenv("PUBSUB_TOPIC"))

		msg := &pubsub.Message{
			Data: []byte(fmt.Sprintf("oc-demo: %s", req.URL)),
		}
		if _, err := topic.Publish(ctx, propagation.WrapMessage(ctx, msg)).Get(ctx); err != nil {
			panic(err)
		}
		span.End()
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
