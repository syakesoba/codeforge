# CodeForge

コードを書いて実行しながら学ぶ、Progate風のプログラミング学習サイト。

ブラウザ上でスライド形式の解説を読みながらコードを書き、その場で自動採点されます。
第一弾としてGoの応用編（並行処理・Web API・Gin・DB連携・認証）を全5コース27レッスンで提供します。

## 特徴

- **ブラウザで完結する学習体験**: スライド解説 → コード記述 → 自動採点
- **安全なコード実行**: ユーザーのコードはDockerサンドボックス内で隔離実行（ネットワーク遮断・リソース制限・非root）
- **実践的な自動採点**: `go test` による検証。HTTPハンドラやDB操作、並行処理の安全性まで判定できる
- **学習の記録**: 合格したレッスンの進捗を保存（ログイン時はサーバー、未ログイン時はブラウザに保存）

## コース一覧

| # | コース | レッスン数 | 内容 |
|---|---|---|---|
| 1 | ゴルーチン・並行処理 | 6 | goroutine、channel、Mutex、select、context、ワーカープール |
| 2 | Web APIの基礎 | 6 | `net/http` によるルーティング・リクエスト/レスポンス・ミドルウェア |
| 3 | フレームワークで作るAPI | 5 | Ginのルーティング・バリデーション・ミドルウェア |
| 4 | データベース連携 | 5 | `database/sql` とGORMによるCRUD |
| 5 | 認証・JWT | 5 | bcrypt、JWT発行/検証、認証ミドルウェア |

各コースの最終レッスンは「道場」として、そのコースの内容を統合するまとめ問題になっています。

## 技術スタック

| 領域 | 技術 |
|---|---|
| バックエンド | Go（標準ライブラリ中心） |
| フロントエンド | Next.js 16 (App Router) / TypeScript / Tailwind CSS v4 |
| コードエディタ | CodeMirror 6 |
| コード実行 | Docker サンドボックス |
| データ保存 | SQLite |

## セットアップ

前提: Go / Docker Desktop（Linuxコンテナモード）/ Node.js

```bash
# 1. 採点用Dockerイメージをビルド（初回は数分かかります）
docker build -t codeforge-judge:latest ./judge-image

# 2. フロントエンドの依存をインストール
cd frontend && npm install && cd ..
```

## 起動

```bash
# バックエンド（:8080）
go run ./cmd/server
```

```bash
# フロントエンド（:3000）別ターミナルで
cd frontend && npm run dev
```

ブラウザで http://localhost:3000 を開きます。

停止方法やトラブルシューティングを含む詳細は [起動方法](docs/codeforge/setup.html) を参照してください。

## ドキュメント

`docs/codeforge/` 配下のHTMLをブラウザで開いてください。

| ドキュメント | 内容 |
|---|---|
| [概要](docs/codeforge/index.html) | プロジェクト全体の概要とコース一覧 |
| [要件定義](docs/codeforge/requirements.html) | 対象ユーザー、コア機能、技術選定の意思決定 |
| [実装プラン](docs/codeforge/plan.html) | コース構成、進捗、今後のロードマップ |
| [詳細設計](docs/codeforge/design.html) | サンドボックス実行の仕組み、API仕様、ディレクトリ構成 |
| [起動方法](docs/codeforge/setup.html) | セットアップ・起動・停止・環境変数 |

## ディレクトリ構成

```
cmd/server/          HTTPサーバー（採点API・認証API）
internal/
  judge/             Dockerサンドボックスでの採点実行
  problems/          問題定義の許可リストとロード
  auth/              サインアップ・ログイン・セッション管理
  store/             SQLite（users / sessions / progress）
problems/            レッスンのコンテンツ（コース別）
judge-image/         採点用Dockerイメージの定義
frontend/            Next.jsフロントエンド
docs/codeforge/      設計ドキュメント
```

## ライセンス

未定
