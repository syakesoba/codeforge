// コースごとに異なるアクセントカラーを持たせ、Progate風の「コースごとに色が違う」
// にぎやかな見た目にする。
//
// Tailwindはソースコード中に完全な形で書かれたクラス名文字列だけを検出するため、
// `bg-${color}-500` のような動的な組み立ては（実行時には正しくても）ビルド時に
// 検出されずスタイルが欠落する。そのため、色ごとの完全なクラス名をこのオブジェクトに
// 直接書き出しておく。
export type CourseColorKey =
  | "violet"
  | "sky"
  | "teal"
  | "amber"
  | "rose"
  | "emerald"
  | "orange"
  | "indigo"
  | "cyan";

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
  emerald: {
    bar: "bg-emerald-500",
    badge:
      "bg-emerald-100 text-emerald-700 dark:bg-emerald-950 dark:text-emerald-300",
    progress: "bg-emerald-500",
    border: "border-t-emerald-500",
  },
  orange: {
    bar: "bg-orange-500",
    badge:
      "bg-orange-100 text-orange-700 dark:bg-orange-950 dark:text-orange-300",
    progress: "bg-orange-500",
    border: "border-t-orange-500",
  },
  indigo: {
    bar: "bg-indigo-500",
    badge:
      "bg-indigo-100 text-indigo-700 dark:bg-indigo-950 dark:text-indigo-300",
    progress: "bg-indigo-500",
    border: "border-t-indigo-500",
  },
  cyan: {
    bar: "bg-cyan-500",
    badge: "bg-cyan-100 text-cyan-700 dark:bg-cyan-950 dark:text-cyan-300",
    progress: "bg-cyan-500",
    border: "border-t-cyan-500",
  },
};
