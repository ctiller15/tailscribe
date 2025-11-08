FROM golang:1.24-alpine AS build

WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o /app/tailscribe

FROM alpine:latest

WORKDIR /app
COPY --from=build /app/tailscribe .
COPY --from=build /app/ui/ ./ui/
COPY --from=build /app/assets ./assets

EXPOSE 8080

CMD ["./tailscribe"]