FROM golang:1.14 as builder
ENV GO111MODULE=on
WORKDIR /go/src/app
COPY go.mod .
COPY go.sum .
# Get dependencies - will also be cached if we won't change mod/sum
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o sidecar-auth-proxy ./cmd/sidecar-auth-proxy

FROM alpine:latest

RUN apk add --no-cache bash && \
    apk add --update tzdata && \
    apk add --no-cache ca-certificates && \
    addgroup -S appgroup && adduser -u 1000 -S appuser -G appgroup

COPY --from=builder /go/src/app/sidecar-auth-proxy /usr/bin/

ENTRYPOINT ["/usr/bin/sidecar-auth-proxy"]
USER appuser
