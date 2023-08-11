FROM docker.io/library/golang:latest AS builder

WORKDIR /app

COPY . .

RUN CGO_ENABLED=0 go build -o fake-rest .

FROM docker.io/library/alpine:latest

WORKDIR /app

COPY --from=builder /app/fake-rest .

ENV DELAY_SECONDS=0 ERROR_PERCENTAGE=0 HTTP_STATUS_CODE=200 HTTP_RESPONSE_MESSAGE="Success"

EXPOSE 8080

USER 1001

CMD ["/app/fake-rest"]

