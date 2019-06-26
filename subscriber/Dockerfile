FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /go/

COPY bin/app .

CMD ["./app"]
