FROM golang:1.21-alpine AS build

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
WORKDIR /app/cmd
RUN go build -o /meai-server main.go


FROM alpine:latest
WORKDIR /app
COPY --from=build /meai-server .
COPY ../migrations ../migrations
COPY ../.env .env
EXPOSE 8081
CMD ["./meai-server"]