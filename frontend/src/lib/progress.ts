const STORAGE_KEY = "codeforge:completed-lessons";
const EMPTY_SET: ReadonlySet<string> = new Set();

/**
 * 学習進捗はユーザー認証が未実装のため、ひとまずブラウザのlocalStorageで保持する。
 * 将来ログイン機能を追加したら、サーバー側の保存に置き換える想定。
 *
 * useSyncExternalStore から使うため、中身が変化していない限り同じSetの参照を
 * 返すようにキャッシュしている（参照が毎回変わるとhydrationミスマッチや
 * 無限再レンダリングの原因になるため）。
 */
let cachedRaw: string | null = null;
let cachedSnapshot: ReadonlySet<string> = EMPTY_SET;

export function getCompletedLessonsSnapshot(): ReadonlySet<string> {
  if (typeof window === "undefined") {
    return EMPTY_SET;
  }
  const raw = window.localStorage.getItem(STORAGE_KEY);
  if (raw !== cachedRaw) {
    cachedRaw = raw;
    cachedSnapshot = new Set(raw ? (JSON.parse(raw) as string[]) : []);
  }
  return cachedSnapshot;
}

export function getServerCompletedLessonsSnapshot(): ReadonlySet<string> {
  return EMPTY_SET;
}

export function subscribeToProgress(onChange: () => void): () => void {
  window.addEventListener("storage", onChange);
  return () => window.removeEventListener("storage", onChange);
}

export function markLessonCompleted(id: string): void {
  if (typeof window === "undefined") {
    return;
  }
  const completed = new Set(getCompletedLessonsSnapshot());
  completed.add(id);
  window.localStorage.setItem(
    STORAGE_KEY,
    JSON.stringify(Array.from(completed)),
  );
  // 同じタブ内ではstorageイベントが発火しないため、変更を即座に反映させる
  window.dispatchEvent(new Event("storage"));
}
