# Lesson 3: サービス層でビジネスロジックを分離する

## 今回学ぶこと

- 「データの出し入れ」と「業務ルール」を分けて考える
- サービス層という設計上の位置づけ
- リポジトリ（Lesson 2）とサービス層を組み合わせる

### データアクセスと業務ルールを混ぜない

「口座Aから口座Bへ送金する」処理を考えます。ここには、大きく分けて2種類のことが混ざりがちです。

1. **データの出し入れ**: 口座の情報をどこから読み、どこに書き戻すか（DB、ファイル、インメモリ…）
2. **業務ルール**: 「残高が足りなければ送金できない」「送金元から引いた分だけ送金先に足す」といった、そのアプリケーション固有のロジック

この2つを1つの関数に混ぜて書いてしまうと、ルールをテストしたいだけなのに本物のDB接続が必要になったり、DBの種類を変えるたびにビジネスルールのコードまで触ることになったりします。

### サービス層を導入する

そこで、Lesson 2で学んだ「リポジトリ（データの出し入れ）」の上に、業務ルールだけを持つ**サービス層**を重ねます。

```go
type TransferService struct {
	repo AccountRepository // データの出し入れは repo に任せる
}

func (s *TransferService) Transfer(fromID, toID, amount int) error {
	from, err := s.repo.FindByID(fromID) // データ取得はrepoの仕事
	if err != nil {
		return err
	}
	to, err := s.repo.FindByID(toID)
	if err != nil {
		return err
	}

	if from.Balance < amount { // ここから先が「業務ルール」
		return ErrInsufficientBalance
	}
	from.Balance -= amount
	to.Balance += amount

	if err := s.repo.Save(from); err != nil { // 保存もrepoの仕事
		return err
	}
	return s.repo.Save(to)
}
```

`TransferService` は `AccountRepository` インターフェースにだけ依存しているので、Lesson 2で学んだのと同じ理由で、インメモリ実装でもSQL実装でも差し替え可能です。そして「残高が足りるか」「いくら引いていくら足すか」という**送金というドメインに固有のルール**は、`TransferService` の中に閉じ込められています。

### なぜ途中で失敗した場合の挙動が重要か

`Transfer` の実装では、`from` と `to` の両方の口座を**先に取得してから**残高チェックを行っています。もし送金先の口座が見つからない場合、送金元の残高を減らす**前**にエラーを返す必要があります。そうしないと「送金元からは引かれたのに、送金先に届いていない」という不整合が起きてしまいます。

## 演習

`InMemoryAccountRepository` は実装済みです。`TransferService.Transfer` を実装してください。

**`(s *TransferService) Transfer(fromID, toID, amount int) error`**
1. `s.repo.FindByID(fromID)` で送金元を取得する。エラーならそのまま返す
2. `s.repo.FindByID(toID)` で送金先を取得する。エラーならそのまま返す
3. 送金元の残高が `amount` より小さければ `ErrInsufficientBalance` を返す
4. 送金元の残高から `amount` を引き、送金先の残高に `amount` を足す
5. `s.repo.Save` で両方の口座を保存する。エラーがあればそのまま返す
6. すべて成功したら `nil` を返す

## ヒント

- 2つの `FindByID` 呼び出しは、両方とも残高を変更する**前**に行ってください。片方が見つからない場合に、もう片方の残高だけ変更されてしまうのを防ぐためです
- `s.repo.Save(from)` と `s.repo.Save(to)` は、それぞれの呼び出し直後にエラーチェックしてください
