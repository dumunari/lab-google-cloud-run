FROM golang:1.24-alpine AS build

WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o cep-temperature /app/cmd/main.go

FROM scratch
WORKDIR /app

COPY --from=build /app/cep-temperature .
EXPOSE 8080
CMD ["./cep-temperature"]