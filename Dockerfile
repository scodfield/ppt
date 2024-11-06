FROM golang:alpine AS builder
LABEL authors="scodfield"
LABEL stage=gobuilder

ENV GOPROXY https://goproxy.cn,direct
RUN apk update --no-cache && apk add --no-cache tzdata

WORKDIR /app
ADD go.mod .
RUN go mod download
COPY . .
RUN go build -o ppt

FROM alpine:latest
RUN apk update --no-cache && apk add --no-cache ca-certificates
COPY --from=builder /usr/share/zoneinfo/Asia/Shanghai /usr/share/zoneinfo/Asia/Shanghai
ENV TZ Asia/Shanghai

WORKDIR /app
COPY --from=builder /app/ppt /app/ppt
EXPOSE 8081

CMD ["/app/ppt"]