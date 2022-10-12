# Copyright (c) 2022 EPAM Systems, Inc.
# 
# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at http://mozilla.org/MPL/2.0/.

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
