// コースごとに異なるアクセントカラーを持たせ、Progate風の「コースごとに色が違う」
// にぎやかな見た目にする。
//
// Tailwindはソースコード中に完全な形で書かれたクラス名文字列だけを検出するため、
// `bg-${color}-500` のような動的な組み立ては（実行時には正しくても）ビルド時に
// 検出されずスタイルが欠落する。そのため、色ごとの完全なクラス名をこのオブジェクトに
// 直接書き出しておく。
export type CourseColorKey = "violet" | "sky" | "teal" | "amber" | "rose";

export type CourseColorClasses = {
  /** コースカード上部のアクセントバー */
  bar: string;
  /** 「Course」バッジの背景・文字色 */
  badge: string;
  /** 進捗バーの塗りつぶし色 */
  progress: string;
  /** レッスンページのヘッダー上部アクセント */
  border: string;
};

export const courseColors: Record<CourseColorKey, CourseColorClasses> = {
  violet: {
    bar: "bg-violet-500",
    badge: "bg-violet-100 text-violet-700 dark:bg-violet-950 dark:text-violet-300",
    progress: "bg-violet-500",
    border: "border-t-violet-500",
  },
  sky: {
    bar: "bg-sky-500",
    badge: "bg-sky-100 text-sky-700 dark:bg-sky-950 dark:text-sky-300",
    progress: "bg-sky-500",
    border: "border-t-sky-500",
  },
  teal: {
    bar: "bg-teal-500",
    badge: "bg-teal-100 text-teal-700 dark:bg-teal-950 dark:text-teal-300",
    progress: "bg-teal-500",
    border: "border-t-teal-500",
  },
  amber: {
    bar: "bg-amber-500",
    badge: "bg-amber-100 text-amber-700 dark:bg-amber-950 dark:text-amber-300",
    progress: "bg-amber-500",
    border: "border-t-amber-500",
  },
  rose: {
    bar: "bg-rose-500",
    badge: "bg-rose-100 text-rose-700 dark:bg-rose-950 dark:text-rose-300",
    progress: "bg-rose-500",
    border: "border-t-rose-500",
  },
};
