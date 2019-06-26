SERVICE = oc-demo-publisher

IMAGE        = 
PUBSUB_TOPIC =
PROJECT_ID   =

all: build docker-build docker-push deploy

build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o bin/app .

docker-build:
	docker build -t $(IMAGE) .

docker-push:
	docker push $(IMAGE)

deploy:
	gcloud beta run deploy $(SERVICE) --image $(IMAGE) --region us-central1 --set-env-vars PUBSUB_TOPIC=$(PUBSUB_TOPIC),PROJECT_ID=$(PROJECT_ID)
