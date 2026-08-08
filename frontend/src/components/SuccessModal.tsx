"use client";

import Link from "next/link";

export default function SuccessModal({
  nextLessonId,
  onClose,
}: {
  nextLessonId: string | null;
  onClose: () => void;
}) {
  return (
    <div
      className="fixed inset-0 z-50 flex items-center justify-center bg-black/50 p-4"
      onClick={onClose}
    >
      <div
        className="w-full max-w-sm rounded-2xl bg-white p-8 text-center shadow-xl dark:bg-neutral-900"
        onClick={(e) => e.stopPropagation()}
      >
        <p className="text-4xl">🎉</p>
        <h2 className="mt-3 text-xl font-bold">おめでとうございます！</h2>
        <p className="mt-2 text-sm text-neutral-600 dark:text-neutral-400">
          このレッスンに合格しました。
        </p>

        <div className="mt-6 flex flex-col gap-2">
          {nextLessonId ? (
            <Link
              href={`/lessons/${nextLessonId}`}
              className="rounded-full bg-emerald-600 px-6 py-2 font-semibold text-white transition hover:bg-emerald-700"
            >
              次のレッスンへ進む →
            </Link>
          ) : (
            <Link
              href="/"
              className="rounded-full bg-emerald-600 px-6 py-2 font-semibold text-white transition hover:bg-emerald-700"
            >
              🎓 コース制覇！ レッスン一覧へ
            </Link>
          )}
          <button
            type="button"
            onClick={onClose}
            className="rounded-full px-6 py-2 text-sm font-semibold text-neutral-500 transition hover:bg-neutral-100 dark:text-neutral-400 dark:hover:bg-neutral-800"
          >
            このまま続ける
          </button>
        </div>
      </div>
    </div>
  );
}
