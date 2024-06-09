# syntax=docker/dockerfile:1

FROM golang:1.22.3

WORKDIR /app

COPY . ./

RUN chmod +x /app/run_app.sh

ENTRYPOINT ["/app/run_app.sh"]