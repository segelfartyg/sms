FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o k8s-reporter .

FROM scratch
COPY --from=builder /app/k8s-reporter /k8s-reporter
ENTRYPOINT ["/k8s-reporter"]
