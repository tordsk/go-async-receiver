FROM golang:1.19-buster as builder
WORKDIR /build
COPY go.mod ./
COPY go.sum ./
COPY . .
RUN go mod download
RUN go mod download
RUN GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o server ./cmd/receiver/main.go
#RUN go build -o server ./cmd/receiver/main.go

FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /build/server .
ENTRYPOINT ["/server"]