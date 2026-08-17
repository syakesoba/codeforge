# Lesson 5（道場）: レイヤードアーキテクチャを組み立てる

## 今回学ぶこと

- これまでのレッスンで学んだ要素を組み合わせて、実践的なアプリケーションの一部を書く
- リポジトリ層・サービス層・通知（DI）を組み合わせた「層」の作り方

### これまで学んだことの総仕上げ

このレッスンは「道場」＝総仕上げ問題です。Lesson 1〜3で学んだ以下の要素を、すべて組み合わせて使います。

- インターフェースでの依存性の注入（Lesson 1）: `Notifier` を具体的な通知手段ではなくインターフェースとして受け取る
- リポジトリパターン（Lesson 2）: `BookRepository` が蔵書データの永続化を抽象化する
- サービス層（Lesson 3）: `LibraryService` が「本を借りる」という業務ルールを持つ

`LibraryService` は次の2つの依存を、どちらも**インターフェース**として受け取ります。

```go
type LibraryService struct {
	repo     BookRepository // データアクセス層
	notifier Notifier       // 通知の送信手段
}
```

こうしておくことで、`LibraryService` 自体のコードは「本の貸し出しルール」だけに集中でき、データがどこに保存されるか・通知がどう送られるかの詳細を一切知る必要がありません。テストでは、`InMemoryBookRepository` と、記録するだけの偽の `Notifier` を渡すだけで、外部との通信を一切せずに業務ルールを検証できます。

## 演習

`InMemoryBookRepository` は実装済みです。`LibraryService.Borrow` を実装してください。

**`(s *LibraryService) Borrow(bookID int) error`**
1. `s.repo.FindByID(bookID)` で本を取得する。エラーならそのまま返す
2. 取得した本がすでに貸し出し中（`OnLoan == true`）なら `ErrBookAlreadyOnLoan` を返す
3. `OnLoan` を `true` にする
4. `s.repo.Save(book)` で保存する。エラーならそのまま返す
5. `s.notifier.Notify(...)` で `"「タイトル」を貸し出しました"` という内容の通知を送る。エラーならそのまま返す
6. すべて成功したら `nil` を返す

## ヒント

- 通知メッセージの組み立てには `fmt.Sprintf("「%s」を貸し出しました", book.Title)` が使えます
- 「本が見つからない」チェックと「すでに貸し出し中」チェックは、どちらも `book.OnLoan = true` にする**前**に行う必要があります
- `fmt` のimportが必要になるので、「インポートを自動修正」を活用してください
