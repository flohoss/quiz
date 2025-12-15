ARG V_GOLANG=1.25.3
ARG V_NODE=25
ARG V_ALPINE=3
FROM golang:${V_GOLANG}-alpine AS final
RUN apk update && apk add --no-cache tzdata dumb-init && \
    rm -rf /tmp/* /var/tmp/* /usr/share/man /var/cache/apk/*

RUN go install github.com/air-verse/air@latest

WORKDIR /app

COPY ./go.mod ./go.sum ./
RUN go mod download > /dev/null 2>&1

COPY ./docker/mono-12.txt /mono-12.txt

ENTRYPOINT ["dumb-init", "--", "/app/docker/dev.entrypoint.sh"]
