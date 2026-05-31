FROM node:22-alpine AS web-builder
WORKDIR /app
COPY package.json package-lock.json ./
RUN npm ci
COPY index.html vite.config.js ./
COPY src/ src/
RUN npm run build

FROM golang:1.23-alpine AS go-builder
WORKDIR /app
COPY go.mod ./
COPY main.go ./
COPY --from=web-builder /app/dist dist/
RUN CGO_ENABLED=0 go build -o page-viewer .

FROM scratch
COPY --from=go-builder /app/page-viewer /page-viewer
ENTRYPOINT ["/page-viewer"]
