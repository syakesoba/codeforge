"use client";

import Link from "next/link";
import { useRouter } from "next/navigation";
import { useAuth } from "@/components/AuthProvider";
import ProgressBar from "@/components/ProgressBar";
import { lessons } from "@/lib/lessons";

export default function SiteHeader() {
  const { user, loading, completed, streakDays, logOut } = useAuth();
  const router = useRouter();

  async function handleLogOut() {
    await logOut();
    router.refresh();
  }

  return (
    <header className="flex items-center justify-between gap-4 border-b border-black/10 px-6 py-3 dark:border-white/10">
      <Link href="/" className="shrink-0 font-bold">
        CodeForge
      </Link>

      <div className="hidden max-w-xs flex-1 items-center gap-3 sm:flex">
        <ProgressBar
          completed={completed.size}
          total={lessons.length}
          className="flex-1"
        />
        {streakDays > 0 && (
          <span
            title={`${streakDays}日連続学習中`}
            className="shrink-0 whitespace-nowrap text-xs font-semibold text-amber-600 dark:text-amber-400"
          >
            🔥 {streakDays}日
          </span>
        )}
      </div>

      <div className="flex items-center gap-3 text-sm">
        {loading ? null : user ? (
          <>
            <span className="text-neutral-500">{user.email}</span>
            <button
              type="button"
              onClick={handleLogOut}
              className="rounded-full border border-black/10 px-4 py-1 font-semibold transition hover:bg-neutral-50 dark:border-white/10 dark:hover:bg-neutral-900"
            >
              ログアウト
            </button>
          </>
        ) : (
          <>
            <Link
              href="/login"
              className="font-semibold text-emerald-700 hover:underline dark:text-emerald-400"
            >
              ログイン
            </Link>
            <Link
              href="/signup"
              className="rounded-full bg-emerald-600 px-4 py-1 font-semibold text-white transition hover:bg-emerald-700"
            >
              新規登録
            </Link>
          </>
        )}
      </div>
    </header>
  );
}
