FROM golang:1.25-alpine

RUN adduser -D -u 1000 judge

ENV GOCACHE=/home/judge/.cache/go-build
RUN mkdir -p /home/judge/.cache && chown -R judge:judge /home/judge

USER judge

# 標準ライブラリのビルドキャッシュをイメージに焼き込んでおくことで、
# 実行時（コンテナ起動のたびにファイルシステムがリセットされる）に
# runtime/net/http 等をゼロからコンパイルし直す必要がなくなる。
RUN go build std

# Gin/GORM/SQLite/JWT/bcrypt など、外部モジュールに依存するレッスン用に
# モジュールキャッシュ（$HOME/go/pkg/mod）とビルドキャッシュを事前に焼き込む。
# 採点実行時は --network none でオフラインになるため、ここで一度だけ
# ネットワーク越しにダウンロード・ビルドしておく必要がある。
# 各レッスンの go.mod/go.sum はここと同じバージョンを参照すること。
COPY --chown=judge:judge priming /home/judge/priming
# go build だけでなく go test まで走らせることで、testing/httptest を含む
# 「採点時に実際に使われるビルド経路」のキャッシュも焼き込む。
# これが無いと、外部モジュールを使うレッスンの初回提出だけ
# テストバイナリのフルコンパイルが走ってタイムアウトしてしまう。
RUN cd /home/judge/priming && go mod tidy && go build ./... && go test -vet=off ./...

# 実行時は完全オフラインであることを明示し、万一バージョン不一致等があっても
# ネットワーク越しの解決を試みず即座に失敗させる。
ENV GOPROXY=off

WORKDIR /workspace
