/** ISO日時の配列から、カレンダー日（ローカルタイムゾーン）ごとの学習日数セットを作る。 */
function toLocalDateKeys(passedAtValues: Iterable<string>): Set<string> {
  const keys = new Set<string>();
  for (const iso of passedAtValues) {
    const d = new Date(iso);
    if (Number.isNaN(d.getTime())) continue;
    keys.add(d.toDateString());
  }
  return keys;
}

/**
 * 現在の連続学習日数を計算する。
 * 今日まだ学習していなくても、昨日までの連続記録が続いていれば0にしない
 * （今日中に学習すれば継続扱いになる、一般的なストリーク方式）。
 */
export function computeStreakDays(passedAtValues: Iterable<string>): number {
  const dateKeys = toLocalDateKeys(passedAtValues);
  if (dateKeys.size === 0) {
    return 0;
  }

  const cursor = new Date();
  if (!dateKeys.has(cursor.toDateString())) {
    cursor.setDate(cursor.getDate() - 1);
    if (!dateKeys.has(cursor.toDateString())) {
      return 0;
    }
  }

  let streak = 0;
  while (dateKeys.has(cursor.toDateString())) {
    streak += 1;
    cursor.setDate(cursor.getDate() - 1);
  }
  return streak;
}
