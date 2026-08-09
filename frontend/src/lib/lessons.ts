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
