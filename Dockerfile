FROM golang:1.19-buster as builder
WORKDIR /build
COPY go.mod ./
COPY go.sum ./
COPY . .
RUN go mod download
RUN go mod download
RUN go build -o server ./cmd/receiver/main.go

FROM gcr.io/distroless/base-debian10
COPY --from=builder /build/server .
ENTRYPOINT ["/server"]