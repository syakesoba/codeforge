# --- ビルド用ステージ ---
FROM golang:1.26-alpine AS builder
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY cmd/ ./cmd/
COPY internal/ ./internal/
# CGO無効・静的リンクにすることで、後段の軽量イメージでも単体で動くバイナリになる
RUN CGO_ENABLED=0 GOOS=linux go build -o /out/server ./cmd/server

# --- 実行用ステージ ---
FROM gcr.io/distroless/static-debian12
WORKDIR /app
COPY --from=builder /out/server ./server
# problems/ 以下はレッスンのMarkdown・スターターコード・模範解答で、
# 実行時にファイルシステムから読み込む（バイナリには埋め込んでいない）。
COPY problems/ ./problems/

ENV PORT=8080
ENV DB_PATH=/data/codeforge.db
VOLUME ["/data"]

EXPOSE 8080
HEALTHCHECK --interval=10s --timeout=3s CMD ["/app/server", "-healthcheck"]
ENTRYPOINT ["/app/server"]
