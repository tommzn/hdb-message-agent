FROM alpine:latest

RUN apk add --no-cache ca-certificates

WORKDIR /go

COPY build_artifact_bin hdb-bin

RUN chmod 755 /go/hdb-bin
ENTRYPOINT ["/go/hdb-bin"]
