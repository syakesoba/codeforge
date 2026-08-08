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

/** ログイン中ユーザーの合格済み問題ID。未ログインなら null。 */
export async function getServerProgress(): Promise<string[] | null> {
  const res = await fetch(`${API_BASE_URL}/api/progress`, withCredentials);
  if (res.status === 401) {
    return null;
  }
  if (!res.ok) {
    throw new Error(`failed to load progress: ${res.status}`);
  }
  const data: { completedProblemIds: string[] } = await res.json();
  return data.completedProblemIds;
}
