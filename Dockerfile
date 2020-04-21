FROM golang:1.14 as builder
WORKDIR /src
ADD . /src
RUN CGO_ENABLED=0 go build -o app

FROM alpine
COPY --from=builder /src/app /bin/app
ENTRYPOINT [ "/bin/app" ]
