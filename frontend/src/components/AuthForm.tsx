"use client";

import { useState } from "react";
import Link from "next/link";
import { useRouter } from "next/navigation";
import { useAuth } from "@/components/AuthProvider";

export default function AuthForm({ mode }: { mode: "login" | "signup" }) {
  const { logIn, signUp } = useAuth();
  const router = useRouter();

  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState<string | null>(null);
  const [submitting, setSubmitting] = useState(false);

  const isSignUp = mode === "signup";

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault();
    setSubmitting(true);
    setError(null);
    try {
      if (isSignUp) {
        await signUp(email, password);
      } else {
        await logIn(email, password);
      }
      router.push("/");
    } catch (err) {
      setError(
        err instanceof Error ? err.message : "エラーが発生しました",
      );
    } finally {
      setSubmitting(false);
    }
  }

  return (
    <div className="mx-auto flex w-full max-w-sm flex-1 flex-col justify-center gap-6 px-6 py-16">
      <h1 className="text-2xl font-bold">
        {isSignUp ? "新規登録" : "ログイン"}
      </h1>

      <form onSubmit={handleSubmit} className="flex flex-col gap-4">
        <label className="flex flex-col gap-1 text-sm">
          <span className="font-semibold">メールアドレス</span>
          <input
            type="email"
            required
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            className="rounded-lg border border-black/15 px-3 py-2 dark:border-white/15 dark:bg-neutral-900"
          />
        </label>

        <label className="flex flex-col gap-1 text-sm">
          <span className="font-semibold">パスワード</span>
          <input
            type="password"
            required
            minLength={isSignUp ? 8 : undefined}
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            className="rounded-lg border border-black/15 px-3 py-2 dark:border-white/15 dark:bg-neutral-900"
          />
          {isSignUp && (
            <span className="text-xs text-neutral-500">8文字以上</span>
          )}
        </label>

        {error && (
          <p className="rounded-lg bg-red-50 p-3 text-sm text-red-700 dark:bg-red-950 dark:text-red-300">
            {error}
          </p>
        )}

        <button
          type="submit"
          disabled={submitting}
          className="rounded-full bg-emerald-600 px-6 py-2 font-semibold text-white transition hover:bg-emerald-700 disabled:cursor-not-allowed disabled:opacity-60"
        >
          {submitting ? "送信中..." : isSignUp ? "登録する" : "ログイン"}
        </button>
      </form>

      <p className="text-sm text-neutral-600 dark:text-neutral-400">
        {isSignUp ? (
          <>
            すでにアカウントをお持ちですか？{" "}
            <Link
              href="/login"
              className="font-semibold text-emerald-700 hover:underline dark:text-emerald-400"
            >
              ログイン
            </Link>
          </>
        ) : (
          <>
            アカウントをお持ちでないですか？{" "}
            <Link
              href="/signup"
              className="font-semibold text-emerald-700 hover:underline dark:text-emerald-400"
            >
              新規登録
            </Link>
          </>
        )}
      </p>
    </div>
  );
}
