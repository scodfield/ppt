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

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go build -o ppt .

FROM alpine:latest
RUN apk update --no-cache && apk add --no-cache ca-certificates

COPY --from=builder /usr/local/bin /usr/local/bin
COPY --from=builder /usr/share/zoneinfo/Asia/Shanghai /usr/share/zoneinfo/Asia/Shanghai

ENV TZ Asia/Shanghai

WORKDIR /app
COPY --from=builder /app/ppt ./
EXPOSE 8081

CMD ["/app/ppt"]