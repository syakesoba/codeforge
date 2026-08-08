const API_BASE_URL =
  process.env.NEXT_PUBLIC_API_BASE_URL ?? "http://localhost:8080";

export type Problem = {
  id: string;
  title: string;
  markdown: string;
  starterCode: string;
};

export type SubmitResult = {
  passed: boolean;
  output: string;
  durationMs: number;
  error?: string;
};

export type User = {
  id: number;
  email: string;
};

export type Diagnostic = {
  line: number;
  column: number;
  message: string;
};

/** セッションCookieを送受信するため、認証が絡むリクエストには credentials が必要。 */
const withCredentials: RequestInit = { credentials: "include" };

async function errorMessage(res: Response, fallback: string): Promise<string> {
  try {
    const data = await res.json();
    if (typeof data?.error === "string") {
      return data.error;
    }
  } catch {
    // JSONでないレスポンスは無視してフォールバックを使う
  }
  return fallback;
}

export async function getProblem(id: string): Promise<Problem> {
  const res = await fetch(`${API_BASE_URL}/api/problems/${id}`);
  if (!res.ok) {
    throw new Error(`failed to load problem ${id}: ${res.status}`);
  }
  return res.json();
}

export async function getAnswer(id: string): Promise<string> {
  const res = await fetch(`${API_BASE_URL}/api/problems/${id}/answer`);
  if (!res.ok) {
    throw new Error(`failed to load answer for ${id}: ${res.status}`);
  }
  const data: { code: string } = await res.json();
  return data.code;
}

export async function submitSolution(
  id: string,
  code: string,
): Promise<SubmitResult> {
  const res = await fetch(`${API_BASE_URL}/api/problems/${id}/submit`, {
    ...withCredentials,
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ code }),
  });
  if (!res.ok) {
    throw new Error(`failed to submit solution: ${res.status}`);
  }
  return res.json();
}

/**
 * コンパイルのみを行い、未インポート・未使用変数などのエラーを検知する。
 * ユーザーコードは実行されない。
 */
export async function checkCode(
  id: string,
  code: string,
  signal?: AbortSignal,
): Promise<Diagnostic[]> {
  const res = await fetch(`${API_BASE_URL}/api/problems/${id}/check`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ code }),
    signal,
  });
  if (!res.ok) {
    throw new Error(`failed to check code: ${res.status}`);
  }
  const data: { diagnostics: Diagnostic[] } = await res.json();
  return data.diagnostics;
}

/** 不足importの追加・未使用importの削除・gofmt整形を行う（goimports相当）。 */
export async function formatCode(id: string, code: string): Promise<string> {
  const res = await fetch(`${API_BASE_URL}/api/problems/${id}/format`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ code }),
  });
  if (!res.ok) {
    throw new Error(await errorMessage(res, "フォーマットに失敗しました"));
  }
  const data: { code: string } = await res.json();
  return data.code;
}

export async function signUp(email: string, password: string): Promise<User> {
  const res = await fetch(`${API_BASE_URL}/api/auth/signup`, {
    ...withCredentials,
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ email, password }),
  });
  if (!res.ok) {
    throw new Error(await errorMessage(res, "登録に失敗しました"));
  }
  return res.json();
}

export async function logIn(email: string, password: string): Promise<User> {
  const res = await fetch(`${API_BASE_URL}/api/auth/login`, {
    ...withCredentials,
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ email, password }),
  });
  if (!res.ok) {
    throw new Error(await errorMessage(res, "ログインに失敗しました"));
  }
  return res.json();
}

export async function logOut(): Promise<void> {
  await fetch(`${API_BASE_URL}/api/auth/logout`, {
    ...withCredentials,
    method: "POST",
  });
}

/** ログイン中のユーザーを返す。未ログインなら null。 */
export async function getMe(): Promise<User | null> {
  const res = await fetch(`${API_BASE_URL}/api/me`, withCredentials);
  if (res.status === 401) {
    return null;
  }
  if (!res.ok) {
    throw new Error(`failed to load current user: ${res.status}`);
  }
  return res.json();
}

export type ServerProgress = {
  completedIds: string[];
  /** 問題ID -> 合格日時（ISO文字列）。ストリーク計算に使う。 */
  passedAt: Record<string, string>;
};

/** ログイン中ユーザーの合格済み問題情報。未ログインなら null。 */
export async function getServerProgress(): Promise<ServerProgress | null> {
  const res = await fetch(`${API_BASE_URL}/api/progress`, withCredentials);
  if (res.status === 401) {
    return null;
  }
  if (!res.ok) {
    throw new Error(`failed to load progress: ${res.status}`);
  }
  const data: { completedProblemIds: string[]; passedAt: Record<string, string> } =
    await res.json();
  return { completedIds: data.completedProblemIds, passedAt: data.passedAt };
}

/** ログイン中ユーザーの自動保存済みコードを取得する。保存が無い/未ログインなら null。 */
export async function getDraft(id: string): Promise<string | null> {
  const res = await fetch(`${API_BASE_URL}/api/problems/${id}/draft`, withCredentials);
  if (res.status === 204 || res.status === 401) {
    return null;
  }
  if (!res.ok) {
    throw new Error(`failed to load draft: ${res.status}`);
  }
  const data: { code: string } = await res.json();
  return data.code;
}

/** エディタの内容を自動保存する（未ログインの場合は呼び出し側でスキップすること）。 */
export async function saveDraft(id: string, code: string): Promise<void> {
  const res = await fetch(`${API_BASE_URL}/api/problems/${id}/draft`, {
    ...withCredentials,
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ code }),
  });
  if (!res.ok && res.status !== 401) {
    throw new Error(`failed to save draft: ${res.status}`);
  }
}
