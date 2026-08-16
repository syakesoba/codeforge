# Lesson 1: 環境変数で設定を切り替える

## 今回学ぶこと

- なぜコンテナ化されたアプリの設定は「環境変数」で渡すのが基本なのか
- `os.Getenv` を使った設定の読み込みと、デフォルト値の扱い方
- テストしやすい設定読み込みの書き方

### なぜ設定をコードに埋め込んではいけないのか

```go
const dbPath = "/Users/me/dev/myapp/data.db" // 悪い例
```

このように設定値をコードに直接書いてしまうと、開発環境用の設定のままDockerイメージをビルドしてしまい、本番環境で動かなくなります。Dockerイメージは「1回ビルドしたら、開発・検証・本番のどの環境にも同じものをデプロイする」のが基本的な考え方（Immutable Infrastructure）なので、**環境ごとに変わる値をイメージの外から渡せる**必要があります。

### 環境変数で渡す

そこで使われるのが環境変数です。Dockerでは `docker run -e PORT=3000 myapp` のように、コンテナ起動時に環境変数を注入できます（`docker-compose.yml` の `environment:` や、Kubernetesの `env:` も同じ考え方です）。

```go
import "os"

port := os.Getenv("PORT")
if port == "" {
	port = "8080" // 環境変数が無ければデフォルト値
}
```

`os.Getenv` は、環境変数が設定されていなければ空文字列を返すので、ローカル開発時は何も設定しなくてもデフォルト値で動き、本番環境だけ環境変数で上書きする、という運用ができます。

### テストしやすくする工夫

`os.Getenv` を設定読み込みの関数の中で直接呼んでしまうと、テストのたびに実際の環境変数を書き換える必要が出てきて面倒です（しかも他のテストと干渉する可能性もあります）。

そこで、「環境変数名を受け取って値を返す関数」を引数として受け取るようにします。

```go
func LoadConfig(getenv func(string) string) Config {
	port := getenv("PORT")
	// ...
}

// 本番: LoadConfig(os.Getenv)
// テスト: LoadConfig(fakeGetenv)  ← mapから値を返すだけの偽の関数
```

`os.Getenv` 自身がちょうど `func(string) string` という型を持つ関数なので、本番コードでは `LoadConfig(os.Getenv)` とそのまま渡せます。この「外部依存を関数として注入する」考え方は、Course「実践的なアーキテクチャ」のLesson 1で学んだ依存性の注入と同じ発想です。

## 演習

`LoadConfig(getenv func(string) string) Config` を実装してください。

- `PORT`: 空文字列なら `"8080"` を使う
- `DB_PATH`: 空文字列なら `"data/app.db"` を使う
- `DEBUG`: 値が `"true"` なら `true`、それ以外（未設定を含む）は `false`

## ヒント

- `getenv("PORT")` の戻り値が空文字列かどうかは `if port == "" { ... }` で判定します
- `DEBUG` の判定は比較演算子1つで書けます: `getenv("DEBUG") == "true"`
