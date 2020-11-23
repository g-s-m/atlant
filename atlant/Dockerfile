
FROM golang:1.13 AS builder
RUN apt-get update && \
    apt-get install curl make unzip
RUN curl -OL https://github.com/protocolbuffers/protobuf/releases/download/v3.14.0/protoc-3.14.0-linux-x86_64.zip
RUN unzip -o protoc-3.14.0-linux-x86_64.zip -d /usr/local bin/protoc
RUN go get google.golang.org/protobuf/cmd/protoc-gen-go
RUN go get google.golang.org/grpc/cmd/protoc-gen-go-grpc
WORKDIR /atlant
COPY . /atlant
RUN make

FROM alpine
COPY --from=builder /atlant/bin/atlant /app/atlant
ENTRYPOINT ["/app/atlant"]

