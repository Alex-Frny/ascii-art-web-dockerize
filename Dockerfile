# Étape 1 : compilation du serveur Go
FROM golang:1.26.5-alpine AS builder

WORKDIR /app

COPY go.mod ./
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o server .

# Étape 2 : création de l'image finale
FROM alpine:3.22

LABEL maintainer="Alexis Fourny"
LABEL version="1.0"
LABEL description="Serveur web ASCII Art développé en Go"

RUN addgroup -S appgroup && adduser -S appuser -G appgroup

WORKDIR /app

COPY --from=builder /app/server .
COPY --from=builder /app/templates ./templates
COPY --from=builder /app/banners ./banners
COPY --from=builder /app/static ./static

USER appuser

EXPOSE 8080

CMD ["./server"]