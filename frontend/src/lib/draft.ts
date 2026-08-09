const PREFIX = "codeforge:draft:";

/** 未ログイン時のコード自動保存はlocalStorageのみで完結させる。 */
export function getLocalDraft(problemId: string): string | null {
  if (typeof window === "undefined") {
    return null;
  }
  return window.localStorage.getItem(PREFIX + problemId);
}

export function saveLocalDraft(problemId: string, code: string): void {
  if (typeof window === "undefined") {
    return;
  }
  window.localStorage.setItem(PREFIX + problemId, code);
}
