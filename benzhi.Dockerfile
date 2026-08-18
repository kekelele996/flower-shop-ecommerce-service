# 评测用镜像：构建并启动 backend 中的真实 HTTP 服务。
FROM golang:1.22
WORKDIR /app/backend
COPY backend/go.mod backend/go.sum ./
RUN go mod download
COPY backend/ ./
RUN go build -o /usr/local/bin/server ./cmd/server
EXPOSE 8080
CMD ["/usr/local/bin/server"]
