import { courses, lessons } from "@/lib/lessons";

export type Badge = {
  id: string;
  emoji: string;
  label: string;
  description: string;
};

/**
 * バッジはサーバーに保存せず、合格済みレッスンとストリーク日数から
 * その場で導出する（シンプルに保つため、バッジ専用の永続化は持たない）。
 */
export function computeEarnedBadges(
  completed: ReadonlySet<string>,
  streakDays: number,
): Badge[] {
  const badges: Badge[] = [];

  if (completed.size >= 1) {
    badges.push({
      id: "first-lesson",
      emoji: "🎉",
      label: "はじめの一歩",
      description: "最初のレッスンに合格した",
    });
  }

  for (const course of courses) {
    if (course.lessons.every((l) => completed.has(l.id))) {
      badges.push({
        id: `course-clear-${course.id}`,
        emoji: "🏆",
        label: `${course.title} 制覇`,
        description: `「${course.title}」の全レッスンに合格した`,
      });
    }
  }

  if (completed.size >= lessons.length) {
    badges.push({
      id: "all-clear",
      emoji: "👑",
      label: "全コース制覇",
      description: "全レッスンに合格した",
    });
  }

  if (streakDays >= 3) {
    badges.push({
      id: "streak-3",
      emoji: "🔥",
      label: "3日連続学習",
      description: "3日連続でレッスンに合格した",
    });
  }

  if (streakDays >= 7) {
    badges.push({
      id: "streak-7",
      emoji: "⚡",
      label: "7日連続学習",
      description: "7日連続でレッスンに合格した",
    });
  }

  return badges;
}
