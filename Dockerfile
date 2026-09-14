# --- ビルド用ステージ ---
FROM golang:1.26-alpine AS builder
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY cmd/ ./cmd/
COPY internal/ ./internal/
# CGO無効・静的リンクにすることで、後段のイメージでも単体で動くバイナリになる
RUN CGO_ENABLED=0 GOOS=linux go build -o /out/server ./cmd/server

# --- 実行用ステージ ---
# 以前は distroless（シェルもツールも無い最小イメージ）を使っていたが、
# バックエンドは実行時に次の外部コマンドを使うため、それでは採点とエディタ支援が動かなかった。
#   - 採点（internal/judge）: docker CLI で採点用コンテナを起動する
#   - /check・/format（internal/lint）: go コマンドと golang.org/x/tools/imports
# Go のバージョンは採点イメージ（judge-image/Dockerfile）と揃え、
# 「エディタではエラーにならないのに採点ではコンパイルエラー」といった食い違いを防ぐ。
FROM golang:1.25-alpine
RUN apk add --no-cache docker-cli

# 実行時はモジュールをネットワークから取得しない（キャッシュに無ければ即座に失敗させる）。
# cgo はエディタ支援に不要なため無効にする。キャッシュの焼き込みより前に設定しておかないと、
# 設定の違いでビルドキャッシュが使われなくなる。
ENV GOTOOLCHAIN=local CGO_ENABLED=0

# /check・/format で Gin/GORM/gRPC 等の外部モジュールをオフラインで解決できるよう、
# 採点イメージと同じ priming モジュールでモジュールキャッシュとビルドキャッシュを焼き込む。
COPY judge-image/priming /opt/priming
RUN cd /opt/priming && go mod download && go build ./... && go test -vet=off -run='^$' ./...
ENV GOPROXY=off

WORKDIR /app
COPY --from=builder /out/server ./server
# problems/ 以下はレッスンのMarkdown・スターターコード・模範解答で、
# 実行時にファイルシステムから読み込む（バイナリには埋め込んでいない）。
COPY problems/ ./problems/

ENV PORT=8080
ENV DB_PATH=/data/codeforge.db
# 採点ワークスペースの置き場所。採点コンテナはホストのDockerデーモンが起動し、
# `-v` のパスはホスト側で解決されるため、起動時にホストと同じパスでバインドマウントすること
# （docs/05-operations/setup.html の「バックエンドをDockerで動かす」を参照）。
ENV JUDGE_WORK_DIR=/tmp/codeforge-judge
VOLUME ["/data"]

EXPOSE 8080
HEALTHCHECK --interval=10s --timeout=3s CMD ["/app/server", "-healthcheck"]
ENTRYPOINT ["/app/server"]
