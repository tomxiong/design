FROM golang:1.19-bullseye as builder

RUN go install github.com/swaggo/swag/cmd/swag@latest

ARG TARGETPLATFORM
ARG TARGETARCH
RUN echo building for "$TARGETPLATFORM"

WORKDIR /workspace

# Copy the Go Modules manifests
COPY go.mod go.mod
COPY go.sum go.sum
# cache deps before building and copying source so that we don't need to re-download as much
# and so that source changes don't invalidate our downloaded layer

# Copy the go source
COPY cmd/ cmd/
COPY internal/ internal/

RUN swag init -d . -g cmd/main.go -o docs

RUN CGO_ENABLED=0 GOOS=linux GOARCH=$TARGETARCH GO111MODULE=on && \
    go mod download && go mod tidy && \
    go build -a -o desgin-server ./cmd/main.go

FROM jgoerzen/debian-base-standard

COPY --from=builder /workspace/design-server /design-server

# Run the binary.
ENTRYPOINT ["/design-server"]