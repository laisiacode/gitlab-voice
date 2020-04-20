FROM golang:1.14 as builder
ADD . /src
RUN cd /src && go build -o app

FROM alpine
WORKDIR /app
COPY --from=builder /src/app /app/
ENTRYPOINT ./app
