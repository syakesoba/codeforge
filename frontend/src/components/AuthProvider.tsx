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
  markLessonCompleted as markLocal,
  subscribeToProgress,
} from "@/lib/progress";

type AuthContextValue = {
  user: User | null;
  loading: boolean;
  /** 合格済みレッスンID。ログイン中はサーバー、未ログインはlocalStorageが情報源。 */
  completed: ReadonlySet<string>;
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
  // localStorage側の進捗は外部ストアなので、変更を購読して再描画する
  const [localCompleted, setLocalCompleted] = useState<ReadonlySet<string>>(
    () => new Set(),
  );

  useEffect(() => {
    const sync = () => setLocalCompleted(getCompletedLessonsSnapshot());
    sync();
    return subscribeToProgress(sync);
  }, []);

  const loadServerProgress = useCallback(async () => {
    const ids = await getServerProgress();
    setServerCompleted(ids ? new Set(ids) : null);
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
  }, []);

  const markCompleted = useCallback(
    (lessonId: string) => {
      // 未ログインでも進捗が見えるよう、localStorageには常に記録する。
      // ログイン中はサーバー側にも submit のレスポンス経路で記録済み。
      markLocal(lessonId);
      if (user) {
        setServerCompleted((prev) => {
          const next = new Set(prev ?? []);
          next.add(lessonId);
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

  const value = useMemo<AuthContextValue>(
    () => ({ user, loading, completed, signUp, logIn, logOut, markCompleted }),
    [user, loading, completed, signUp, logIn, logOut, markCompleted],
  );

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
}
