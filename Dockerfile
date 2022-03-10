# Build the manager binary
FROM golang:1.16 as builder

WORKDIR /workspace
# Copy the Go Modules manifests
COPY go.mod go.mod
COPY go.sum go.sum
# cache deps before building and copying source so that we don't need to re-download as much
# and so that source changes don't invalidate our downloaded layer
RUN go mod download

# Copy the go source
COPY gotracing.go gotracing.go

# Build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on go build -a -o gotracing gotracing.go

FROM alpine:3.14.2
LABEL org.opencontainers.image.source https://github.com/avadhut123pisal/github-actions
RUN apk --no-cache add ca-certificates
COPY --from=builder /workspace/gotracing /usr/local/bin/gotracing/
WORKDIR /usr/local/bin/gotracing
RUN chmod +x gotracing
ENTRYPOINT ["./gotracing"]
