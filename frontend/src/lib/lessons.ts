import type { CourseColorKey } from "@/lib/courseColors";

export type LessonMeta = {
  id: string;
  title: string;
};

export type Course = {
  id: string;
  title: string;
  /** コースカードやレッスンページで使うアクセントカラー。 */
  accent: CourseColorKey;
  lessons: LessonMeta[];
};

// カリキュラム上の意図した順序（ゴルーチン・並行処理 → Web APIの基礎 → ...）で並べる。
// 技術検証の都合でWeb APIの基礎を先に実装したが、表示・進行順はカリキュラム通りにする。
export const courses: Course[] = [
  {
    id: "concurrency",
    title: "ゴルーチン・並行処理",
    accent: "violet",
    lessons: [
      { id: "concurrency-01", title: "Lesson 1: ゴルーチンとWaitGroup" },
      { id: "concurrency-02", title: "Lesson 2: channelでやり取りする" },
      { id: "concurrency-03", title: "Lesson 3: sync.Mutexで排他制御する" },
      { id: "concurrency-04", title: "Lesson 4: selectで複数のチャネルを扱う" },
      { id: "concurrency-05", title: "Lesson 5: contextでタイムアウト制御する" },
      { id: "concurrency-06", title: "Lesson 6（道場）: ワーカープールを作ろう" },
    ],
  },
  {
    id: "web-api-basics",
    title: "Web APIの基礎",
    accent: "sky",
    lessons: [
      { id: "web-api-basics-01", title: "Lesson 1: HTTPサーバーを立てる" },
      { id: "web-api-basics-02", title: "Lesson 2: ルーティングとパスパラメータ" },
      { id: "web-api-basics-03", title: "Lesson 3: リクエストを読み取る" },
      { id: "web-api-basics-04", title: "Lesson 4: レスポンスを返す" },
      { id: "web-api-basics-05", title: "Lesson 5: 共通処理をまとめる" },
      { id: "web-api-basics-06", title: "Lesson 6（道場）: TODO管理APIを作ろう" },
    ],
  },
  {
    id: "gin-api",
    title: "フレームワークで作るAPI（Gin）",
    accent: "teal",
    lessons: [
      { id: "gin-api-01", title: "Lesson 1: Ginの基本" },
      {
        id: "gin-api-02",
        title: "Lesson 2: ルーティンググループとパスパラメータ",
      },
      {
        id: "gin-api-03",
        title: "Lesson 3: リクエストのバインディングとバリデーション",
      },
      { id: "gin-api-04", title: "Lesson 4: Ginのミドルウェア" },
      { id: "gin-api-05", title: "Lesson 5（道場）: GinでCRUD APIを作ろう" },
    ],
  },
  {
    id: "database",
    title: "データベース連携",
    accent: "amber",
    lessons: [
      {
        id: "database-01",
        title: "Lesson 1: データベースに接続してテーブルを作る",
      },
      { id: "database-02", title: "Lesson 2: データを登録・取得する" },
      { id: "database-03", title: "Lesson 3: データを更新・削除する" },
      { id: "database-04", title: "Lesson 4: GORMの基本" },
      {
        id: "database-05",
        title: "Lesson 5（道場）: GORMでCRUDを完成させよう",
      },
    ],
  },
  {
    id: "auth-jwt",
    title: "認証・JWT",
    accent: "rose",
    lessons: [
      { id: "auth-jwt-01", title: "Lesson 1: パスワードをハッシュ化する" },
      { id: "auth-jwt-02", title: "Lesson 2: JWTを発行する" },
      { id: "auth-jwt-03", title: "Lesson 3: JWTを検証する" },
      {
        id: "auth-jwt-04",
        title: "Lesson 4: 認証ミドルウェアでAPIを保護する",
      },
      { id: "auth-jwt-05", title: "Lesson 5（道場）: ログインAPIを作ろう" },
    ],
  },
  {
    id: "testing",
    title: "テストの書き方",
    accent: "emerald",
    lessons: [
      { id: "testing-01", title: "Lesson 1: テーブル駆動テストを書く" },
      { id: "testing-02", title: "Lesson 2: サブテストを使いこなす" },
      { id: "testing-03", title: "Lesson 3: インターフェースでモックする" },
      { id: "testing-04", title: "Lesson 4: エラーケースをテストする" },
      { id: "testing-05", title: "Lesson 5（道場）: HTTPハンドラーをテストする" },
    ],
  },
  {
    id: "error-handling",
    title: "エラーハンドリング",
    accent: "orange",
    lessons: [
      { id: "error-handling-01", title: "Lesson 1: カスタムエラー型を定義する" },
      { id: "error-handling-02", title: "Lesson 2: エラーをラップする" },
      { id: "error-handling-03", title: "Lesson 3: errors.Is / errors.As で判定する" },
      { id: "error-handling-04", title: "Lesson 4: 複数のエラーをまとめる" },
      { id: "error-handling-05", title: "Lesson 5（道場）: 実践的なエラーハンドリング" },
    ],
  },
  {
    id: "architecture",
    title: "実践的なアーキテクチャ",
    accent: "indigo",
    lessons: [
      { id: "architecture-01", title: "Lesson 1: インターフェースで依存を注入する" },
      { id: "architecture-02", title: "Lesson 2: リポジトリパターンでデータアクセスを抽象化する" },
      { id: "architecture-03", title: "Lesson 3: サービス層でビジネスロジックを分離する" },
      { id: "architecture-04", title: "Lesson 4: 関数オプションパターンで柔軟な初期化をする" },
      { id: "architecture-05", title: "Lesson 5（道場）: レイヤードアーキテクチャを組み立てる" },
    ],
  },
  {
    id: "deploy",
    title: "デプロイ・Docker化",
    accent: "cyan",
    lessons: [
      { id: "deploy-01", title: "Lesson 1: 環境変数で設定を切り替える" },
      { id: "deploy-02", title: "Lesson 2: ヘルスチェックエンドポイントを実装する" },
      { id: "deploy-03", title: "Lesson 3: グレースフルシャットダウンを実装する" },
      { id: "deploy-04", title: "Lesson 4: 構造化ロギングでコンテナ環境に対応する" },
      { id: "deploy-05", title: "Lesson 5（道場）: 本番向けのHTTPサーバーを組み立てる" },
    ],
  },
];

export const lessons: LessonMeta[] = courses.flatMap((c) => c.lessons);

export function courseForLesson(id: string): Course | undefined {
  return courses.find((c) => c.lessons.some((l) => l.id === id));
}

export function nextLessonId(id: string): string | null {
  const index = lessons.findIndex((l) => l.id === id);
  if (index === -1 || index === lessons.length - 1) {
    return null;
  }
  return lessons[index + 1].id;
}
