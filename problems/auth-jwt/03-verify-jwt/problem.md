# Lesson 3: JWTを検証する

## 今回学ぶこと

- `jwt.Parse` によるトークンの検証
- 署名アルゴリズムを必ず確認すべき理由（`alg: none` 攻撃）
- クレームの取り出し方

### jwt.Parse の基本形

トークンの検証は `jwt.Parse` で行います。第2引数には「秘密鍵を返す関数（キー関数）」を渡すのが特徴的な設計です。

```go
token, err := jwt.Parse(tokenString, func(t *jwt.Token) (any, error) {
	return []byte(secret), nil
})
```

なぜ鍵を直接渡さず関数にするかというと、「トークンのヘッダーを見てから、どの鍵を使うか決めたい」ケースがあるためです（例: 複数の鍵をローテーションしている場合）。

`jwt.Parse` は次をすべて自動で検証してくれます。

- 署名が正しいか（改ざんされていないか）
- `exp`（有効期限）が切れていないか
- `nbf`（有効開始時刻）に達しているか

どれかが失敗すれば `err` が返るので、呼び出し側は `err != nil` を見るだけで済みます。

### 署名アルゴリズムの確認は必須

キー関数の中で、**署名アルゴリズムが想定どおりかを必ず確認してください**。

```go
token, err := jwt.Parse(tokenString, func(t *jwt.Token) (any, error) {
	if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, errors.New("unexpected signing method")
	}
	return []byte(secret), nil
})
```

これを怠ると、攻撃者が「アルゴリズムは `none`（署名なし）です」と書いたトークンを送りつけ、署名検証を素通りさせる**アルゴリズム混同攻撃**が成立する恐れがあります。ライブラリ側でも対策は入っていますが、アプリケーション側で明示的に確認するのが安全な作法です。

`jwt.SigningMethodHS256` は `*jwt.SigningMethodHMAC` 型なので、上の型アサーションで「HMAC系のアルゴリズムであること」を確認できます。

### クレームを取り出す

検証に成功したら、`token.Claims` からクレームを取り出します。`jwt.MapClaims` に型アサーションしてから、専用のゲッターを使うのが v5 の作法です。

```go
claims, ok := token.Claims.(jwt.MapClaims)
if !ok {
	return "", errors.New("invalid claims")
}

sub, err := claims.GetSubject() // "sub" クレームを取り出す
if err != nil {
	return "", err
}
```

`GetSubject()` のほか、`GetExpirationTime()` / `GetIssuedAt()` などが用意されています。`claims["sub"].(string)` と直接書くこともできますが、ゲッターを使うほうが型安全です。

### 完成形

```go
func VerifyToken(secret, tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(secret), nil
	})
	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", errors.New("invalid claims")
	}

	return claims.GetSubject()
}
```

## 演習

`VerifyToken` 関数を実装してください。

**`VerifyToken(secret, tokenString string) (string, error)`**
- `jwt.Parse` でトークンを検証する
- キー関数の中で署名アルゴリズムがHMAC系であることを確認し、そうでなければエラーを返す
- 検証に失敗したら（署名不正・期限切れなど）エラーを返す
- 成功したら `sub` クレーム（ユーザー名）を返す

トークンを発行する `IssueToken` はすでに実装済みです。

## ヒント

- `claims.GetSubject()` の戻り値 `(string, error)` はそのまま返せます
- 「期限切れ」「署名不正」の判定は `jwt.Parse` が自動でやってくれるので、`err` をそのまま返すだけで構いません
