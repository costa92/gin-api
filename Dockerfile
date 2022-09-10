FROM golang:1.18 as  builder

ENV GO111MODULE=on \
    GOPROXY=https://goproxy.cn,direct \
    GIN_MODE=release
ENV CGO_ENABLED 0
ENV GOOS linux
# 将代码复制到构建镜像中。

WORKDIR /workspace
COPY . .
RUN go mod download
RUN --mount=type=cache,target=/go

RUN go build -buildmode=pie -ldflags "-linkmode external -extldflags -static -w" -o /workspace/goweb
RUN ls -l && cd /workspace && ls -l

## 运行时镜像。
## Alpine兼顾了镜像大小和运维性。
FROM alpine:latest
RUN apk update --no-cache && apk --no-cache add ca-certificates
ENV TZ Asia/Shanghai

WORKDIR /app
## 复制构建产物。
COPY --from=builder /workspace/goweb /app/goweb
COPY --from=builder /workspace/config/config.yaml /app/

EXPOSE 8080

RUN mkdir /app/logs
## 指定默认的启动命令。
CMD ["./goweb", "--config", "/app/config.yaml"]