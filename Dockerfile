FROM golang:alpine AS builder
LABEL authors="scodfield"
LABEL stage=gobuilder

ENV GOPROXY https://goproxy.cn,direct
RUN apk update --no-cache && apk add --no-cache tzdata

WORKDIR /app

COPY go.mod go.sum ./

RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

COPY . .

ARG VERSION="dev"
ARG BUILD_TIME
ARG GIT_COMMIT

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go build -ldflags="\
    -X config.Version=$VERSION \
    -X config.BuildTime=${BUILD_TIME:-$(date -u +'%Y-%m-%dT%H:%M:%SZ')} \
    -X config.GitCommit=${GIT_COMMIT:-$(git rev-parse HEAD 2>/dev/null || echo 'unknown')} " \
    -o ppt .

FROM alpine:latest
RUN apk update --no-cache && apk add --no-cache ca-certificates

COPY --from=builder /usr/local/bin /usr/local/bin
COPY --from=builder /usr/share/zoneinfo/Asia/Shanghai /usr/share/zoneinfo/Asia/Shanghai

ENV TZ Asia/Shanghai

WORKDIR /app
COPY --from=builder /app/ppt ./
EXPOSE 8081

CMD ["/app/ppt"]