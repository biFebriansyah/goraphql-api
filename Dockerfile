FROM golang:1.23.4-alpine AS build

WORKDIR /app

COPY .env.docker /app/.env
COPY . .

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -v -o goapp .

FROM alpine:latest AS release

WORKDIR /app

COPY --from=build /app/.env .
COPY --from=build /app/goapp .

ENV PATH="/app:${PATH}"

USER 405

EXPOSE 8081

CMD [ "goapp" ]

