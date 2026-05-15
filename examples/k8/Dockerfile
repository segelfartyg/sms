FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o k8s-exporter .

FROM scratch
COPY --from=builder /app/k8s-exporter /k8s-exporter
ENTRYPOINT ["/k8s-exporter"]
