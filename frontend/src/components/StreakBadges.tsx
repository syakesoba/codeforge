"use client";

import { useAuth } from "@/components/AuthProvider";
import { computeEarnedBadges } from "@/lib/badges";

export default function StreakBadges() {
  const { completed, streakDays } = useAuth();
  const badges = computeEarnedBadges(completed, streakDays);

  if (completed.size === 0) {
    return null;
  }

  return (
    <div className="flex flex-col gap-3 rounded-xl border border-black/10 p-5 dark:border-white/10">
      <div className="flex items-center justify-between">
        <p className="text-sm text-neutral-500">学習の記録</p>
        {streakDays > 0 && (
          <span className="text-sm font-semibold text-amber-600 dark:text-amber-400">
            🔥 {streakDays}日連続学習中
          </span>
        )}
      </div>

      {badges.length > 0 && (
        <div className="flex flex-wrap gap-2">
          {badges.map((badge) => (
            <span
              key={badge.id}
              title={badge.description}
              className="flex items-center gap-1 rounded-full bg-amber-50 px-3 py-1 text-xs font-semibold text-amber-700 dark:bg-amber-950 dark:text-amber-300"
            >
              <span>{badge.emoji}</span>
              <span>{badge.label}</span>
            </span>
          ))}
        </div>
      )}
    </div>
  );
}
