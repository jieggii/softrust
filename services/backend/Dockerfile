FROM golang:1.25 AS builder

WORKDIR /build
COPY go.mod go.sum* ./

# install dependencies:
RUN go mod download

# copy service source code:
#COPY pkg ./pkg/
COPY internal ./internal/
COPY main.go ./

# build service binary:
RUN CGO_ENABLED=0 go build -o service main.go

FROM scratch AS runner

WORKDIR /usr/local/bin
COPY --from=builder /build/service ./
ENTRYPOINT ["/usr/local/bin/service"]