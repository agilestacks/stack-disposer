FROM golang:1.17-alpine as builder

RUN apk update && apk upgrade && \
    apk add --no-cache git

WORKDIR /go/src/stack-disposer
COPY . .
RUN go mod download
RUN go build -o bin/worker main.go

FROM gcr.io/superhub/gcp-toolbox:d5d4d96-985cbfc-b0ea077-fd95688

COPY --from=builder /go/src/stack-disposer/bin/worker /usr/local/bin/
RUN chmod +x /usr/local/bin/worker

ENTRYPOINT [ "worker" ]
