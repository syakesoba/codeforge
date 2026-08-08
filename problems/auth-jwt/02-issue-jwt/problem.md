# Lesson 2: JWTを発行する

## 今回学ぶこと

- JWT（JSON Web Token）とは何か、なぜ使うのか
- クレーム（claims）の意味と代表的な項目
- `golang-jwt/jwt/v5` でトークンを発行する方法

### JWTとは

ログイン後、ユーザーが「自分は認証済みだ」とサーバーに示すための仕組みが必要です。**JWT**は、そのための「署名付きの通行証」だと考えてください。

JWTは3つの部分をドット（`.`）で繋いだ文字列です。

```
eyJhbGciOiJIUzI1NiJ9.eyJzdWIiOiJhbGljZSJ9.4pcPyMD09olPSyXnrXCjTwXyr4BsezdI1AVTmud2fU4
└─── ヘッダー ───┘└─── ペイロード ───┘└──────── 署名 ────────┘
```

- **ヘッダー**: 署名アルゴリズムの情報
- **ペイロード（クレーム）**: ユーザーIDなどの情報
- **署名**: 秘密鍵を使って生成された改ざん検知用のデータ

重要なのは、**ペイロードは暗号化されていない**という点です。誰でもデコードして中身を読めます。したがって、パスワードなどの秘密情報をJWTに入れてはいけません。

一方、**署名があるので改ざんはできません**。攻撃者がペイロードを書き換えても、秘密鍵を知らない限り正しい署名を作れないため、検証時に必ず弾かれます。

### クレーム

ペイロードに入れる情報を「クレーム」と呼びます。よく使われる標準クレームは次のとおりです。

| クレーム | 意味 |
|---|---|
| `sub` (subject) | トークンの主体。ユーザーIDやユーザー名を入れる |
| `exp` (expiration time) | 有効期限（UNIX時刻）。これを過ぎたトークンは無効 |
| `iat` (issued at) | 発行時刻 |

`exp` は必ず設定してください。有効期限のないトークンは、一度漏洩すると永久に悪用されてしまいます。

### トークンを発行する

```go
import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func IssueToken(secret, username string, ttl time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"sub": username,
		"exp": time.Now().Add(ttl).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(secret))
}
```

- `jwt.MapClaims` は `map[string]any` で、クレームを自由に組み立てられます
- `exp` には `time.Time` ではなく **UNIX時刻（`int64`）** を入れます。`time.Now().Add(ttl).Unix()` で得られます
- `jwt.SigningMethodHS256` は、1つの秘密鍵で署名も検証も行う「共通鍵方式」のアルゴリズムです
- `SignedString` は秘密鍵を `[]byte` で受け取り、完成したトークン文字列を返します

`ttl`（Time To Live）は有効期間です。`time.Hour` を渡せば1時間有効なトークンになります。

## 演習

`IssueToken` 関数を実装してください。

**`IssueToken(secret, username string, ttl time.Duration) (string, error)`**
- クレームに `sub`（`username`）と `exp`（現在時刻 + `ttl` のUNIX時刻）を設定する
- 署名アルゴリズムは `jwt.SigningMethodHS256` を使う
- `secret` を鍵として署名した、トークン文字列を返す

## ヒント

- `exp` の値は `time.Now().Add(ttl).Unix()` です
- `token.SignedString([]byte(secret))` の戻り値 `(string, error)` はそのまま返せます
