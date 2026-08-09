"use client";

import {
  createContext,
  useCallback,
  useContext,
  useEffect,
  useMemo,
  useState,
} from "react";
import {
  getMe,
  getServerProgress,
  logIn as apiLogIn,
  logOut as apiLogOut,
  signUp as apiSignUp,
  type User,
} from "@/lib/api";
import {
  getCompletedLessonsSnapshot,
  getPassedAtSnapshot,
  markLessonCompleted as markLocal,
  subscribeToProgress,
} from "@/lib/progress";
import { computeStreakDays } from "@/lib/streak";

type AuthContextValue = {
  user: User | null;
  loading: boolean;
  /** 合格済みレッスンID。ログイン中はサーバー、未ログインはlocalStorageが情報源。 */
  completed: ReadonlySet<string>;
  /** 問題ID -> 合格日時（ISO文字列）。 */
  passedAt: ReadonlyMap<string, string>;
  /** 現在の連続学習日数。 */
  streakDays: number;
  signUp: (email: string, password: string) => Promise<void>;
  logIn: (email: string, password: string) => Promise<void>;
  logOut: () => Promise<void>;
  markCompleted: (lessonId: string) => void;
};

const AuthContext = createContext<AuthContextValue | null>(null);

export function useAuth(): AuthContextValue {
  const ctx = useContext(AuthContext);
  if (!ctx) {
    throw new Error("useAuth must be used within AuthProvider");
  }
  return ctx;
}

export default function AuthProvider({
  children,
}: {
  children: React.ReactNode;
}) {
  const [user, setUser] = useState<User | null>(null);
  const [loading, setLoading] = useState(true);
  const [serverCompleted, setServerCompleted] = useState<Set<string> | null>(
    null,
  );
  const [serverPassedAt, setServerPassedAt] = useState<Map<
    string,
    string
  > | null>(null);
  // localStorage側の進捗は外部ストアなので、変更を購読して再描画する
  const [localCompleted, setLocalCompleted] = useState<ReadonlySet<string>>(
    () => new Set(),
  );
  const [localPassedAt, setLocalPassedAt] = useState<
    ReadonlyMap<string, string>
  >(() => new Map());

  useEffect(() => {
    const sync = () => {
      setLocalCompleted(getCompletedLessonsSnapshot());
      setLocalPassedAt(getPassedAtSnapshot());
    };
    sync();
    return subscribeToProgress(sync);
  }, []);

  const loadServerProgress = useCallback(async () => {
    const progress = await getServerProgress();
    setServerCompleted(progress ? new Set(progress.completedIds) : null);
    setServerPassedAt(
      progress ? new Map(Object.entries(progress.passedAt)) : null,
    );
  }, []);

  useEffect(() => {
    let cancelled = false;

    (async () => {
      try {
        const me = await getMe();
        if (cancelled) return;
        setUser(me);
        if (me) {
          await loadServerProgress();
        }
      } catch {
        // バックエンドが起動していない場合などは未ログイン扱いにする
        if (!cancelled) setUser(null);
      } finally {
        if (!cancelled) setLoading(false);
      }
    })();

    return () => {
      cancelled = true;
    };
  }, [loadServerProgress]);

  const signUp = useCallback(
    async (email: string, password: string) => {
      const u = await apiSignUp(email, password);
      setUser(u);
      await loadServerProgress();
    },
    [loadServerProgress],
  );

  const logIn = useCallback(
    async (email: string, password: string) => {
      const u = await apiLogIn(email, password);
      setUser(u);
      await loadServerProgress();
    },
    [loadServerProgress],
  );

  const logOut = useCallback(async () => {
    await apiLogOut();
    setUser(null);
    setServerCompleted(null);
    setServerPassedAt(null);
  }, []);

  const markCompleted = useCallback(
    (lessonId: string) => {
      // 未ログインでも進捗が見えるよう、localStorageには常に記録する。
      // ログイン中はサーバー側にも submit のレスポンス経路で記録済み。
      markLocal(lessonId);
      if (user) {
        const now = new Date().toISOString();
        setServerCompleted((prev) => {
          const next = new Set(prev ?? []);
          next.add(lessonId);
          return next;
        });
        setServerPassedAt((prev) => {
          const next = new Map(prev ?? []);
          if (!next.has(lessonId)) {
            next.set(lessonId, now);
          }
          return next;
        });
      }
    },
    [user],
  );

  const completed = useMemo<ReadonlySet<string>>(
    () => (user && serverCompleted ? serverCompleted : localCompleted),
    [user, serverCompleted, localCompleted],
  );

  const passedAt = useMemo<ReadonlyMap<string, string>>(
    () => (user && serverPassedAt ? serverPassedAt : localPassedAt),
    [user, serverPassedAt, localPassedAt],
  );

  const streakDays = useMemo(
    () => computeStreakDays(passedAt.values()),
    [passedAt],
  );

  const value = useMemo<AuthContextValue>(
    () => ({
      user,
      loading,
      completed,
      passedAt,
      streakDays,
      signUp,
      logIn,
      logOut,
      markCompleted,
    }),
    [
      user,
      loading,
      completed,
      passedAt,
      streakDays,
      signUp,
      logIn,
      logOut,
      markCompleted,
    ],
  );

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
}
