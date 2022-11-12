FROM golang:1.17 as builder

WORKDIR /src
ADD . /src
RUN CGO_ENABLED=0 go build -o gitlab-voice

FROM alpine:latest
COPY --from=builder /src/app /bin/gitlab-voice
ENTRYPOINT [ "/bin/gitlab-voice" ]
