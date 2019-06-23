package publisher

import (
	"context"
	"os"

	"cloud.google.com/go/pubsub"
)

func main() {
	client, err := pubsub.NewClient(context.Background(), os.Getenv("PROJECT_ID"))
	if err != nil {
		panic(err)
	}
	topic := client.Topic(os.Getenv("PUBSUB_TOPIC"))
	message := &pubsub.Message{
		Data: []byte("oc-demo"),
	}
	topic.Publish(context.Background(), message)
}
