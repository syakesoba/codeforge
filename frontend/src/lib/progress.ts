const STORAGE_KEY = "codeforge:completed-lessons";
const EMPTY_SET: ReadonlySet<string> = new Set();
const EMPTY_MAP: ReadonlyMap<string, string> = new Map();

/**
 * 学習進捗はログイン中はサーバー、未ログイン中はブラウザのlocalStorageで保持する。
 * 保存形式は `{ [problemId]: passedAtISOString }`。ストリーク計算のため合格日時を持つ。
 *
 * 旧形式（IDの配列のみ）のデータが残っている場合は、読み込み時に現在時刻を
 * 合格日時として補完する（正確な過去の日付は失われるが、ストリークは今後の
 * 学習から数え直せば十分なため許容する）。
 *
 * useSyncExternalStore から使うため、中身が変化していない限り同じ参照を
 * 返すようにキャッシュしている（参照が毎回変わるとhydrationミスマッチや
 * 無限再レンダリングの原因になるため）。
 */
let cachedRaw: string | null = null;
let cachedIds: ReadonlySet<string> = EMPTY_SET;
let cachedPassedAt: ReadonlyMap<string, string> = EMPTY_MAP;

function parse(raw: string | null): Record<string, string> {
  if (!raw) {
    return {};
  }
  const data: unknown = JSON.parse(raw);
  if (Array.isArray(data)) {
    // 旧形式（string[]）からの移行
    const now = new Date().toISOString();
    const migrated: Record<string, string> = {};
    for (const id of data) {
      migrated[id] = now;
    }
    return migrated;
  }
  return data as Record<string, string>;
}

function ensureSynced(): void {
  if (typeof window === "undefined") {
    return;
  }
  const raw = window.localStorage.getItem(STORAGE_KEY);
  if (raw === cachedRaw) {
    return;
  }
  cachedRaw = raw;
  const record = parse(raw);
  cachedIds = new Set(Object.keys(record));
  cachedPassedAt = new Map(Object.entries(record));
}

export function getCompletedLessonsSnapshot(): ReadonlySet<string> {
  if (typeof window === "undefined") {
    return EMPTY_SET;
  }
  ensureSynced();
  return cachedIds;
}

export function getPassedAtSnapshot(): ReadonlyMap<string, string> {
  if (typeof window === "undefined") {
    return EMPTY_MAP;
  }
  ensureSynced();
  return cachedPassedAt;
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
  const raw = window.localStorage.getItem(STORAGE_KEY);
  const record = parse(raw);
  if (!(id in record)) {
    record[id] = new Date().toISOString();
    window.localStorage.setItem(STORAGE_KEY, JSON.stringify(record));
    // 同じタブ内ではstorageイベントが発火しないため、変更を即座に反映させる
    window.dispatchEvent(new Event("storage"));
  }
}
