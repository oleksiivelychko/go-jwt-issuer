FROM golang:1.19-alpine as builder

# Cgo enables the creation of Go packages that call C code.
ENV CGO_ENABLED 0

# Allow Go to retreive the dependencies for the build step.
RUN apk add --no-cache git

WORKDIR /build

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY ./ ./

# Disable compiler optimizations.
RUN go build -gcflags="all=-N -l" -o /build/app main.go

# Get Delve from a GOPATH not from a Go Modules project.
WORKDIR /go/src/
RUN go install github.com/go-delve/delve/cmd/dlv@latest

FROM alpine:3

# Mark this container as using the Go language runtime.
ENV GOTRACEBACK=all

WORKDIR /build
COPY --from=builder /build/app /
COPY --from=builder /go/bin/dlv /

EXPOSE 8080 56268

CMD ["/dlv", "--listen=:56268", "--headless=true", "--log", "--accept-multiclient", "--api-version=2", "exec", "/app"]
