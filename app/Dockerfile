FROM golang:1.25 AS builder 



COPY ./app/config /app/config
COPY ./app/discord /app/discord
COPY ./app/healthcheck /app/healthcheck
COPY ./app/notification /app/notification
COPY ./app/main.go /app/main.go

WORKDIR /app

RUN go mod download
ARG VERSION="dev"

RUN CGO_ENABLED=0 GOOS=linux go build  -ldflags "-X main.Version=${VERSION}" -a -installsuffix cgo -o main .


FROM gcr.io/distroless/static:nonroot

WORKDIR /app

COPY --from=builder /app/main /app/healthchecker

USER nonroot:nonroot

EXPOSE 8080

CMD ["./healthchecker"]