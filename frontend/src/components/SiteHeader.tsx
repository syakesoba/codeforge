"use client";

import Link from "next/link";
import { useRouter } from "next/navigation";
import { useAuth } from "@/components/AuthProvider";

export default function SiteHeader() {
  const { user, loading, logOut } = useAuth();
  const router = useRouter();

  async function handleLogOut() {
    await logOut();
    router.refresh();
  }

  return (
    <header className="flex items-center justify-between border-b border-black/10 px-6 py-3 dark:border-white/10">
      <Link href="/" className="font-bold">
        CodeForge
      </Link>

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
