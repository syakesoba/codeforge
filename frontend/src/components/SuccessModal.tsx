"use client";

import Link from "next/link";
import Mascot from "@/components/Mascot";

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
        className="w-full max-w-sm overflow-hidden rounded-2xl bg-white text-center shadow-xl dark:bg-neutral-900"
        onClick={(e) => e.stopPropagation()}
      >
        <div className="bg-gradient-to-br from-orange-100 via-pink-50 to-amber-100 px-8 pt-8 pb-4 dark:from-orange-950 dark:via-neutral-900 dark:to-amber-950">
          <Mascot pose="celebrate" className="mx-auto h-28 w-28" />
        </div>

        <div className="px-8 pb-8">
          <h2 className="mt-2 text-xl font-bold">おめでとうございます！</h2>
          <p className="mt-2 text-sm text-neutral-600 dark:text-neutral-400">
            このレッスンに合格しました。
          </p>

          <div className="mt-6 flex flex-col gap-2">
            {nextLessonId ? (
              <Link
                href={`/lessons/${nextLessonId}`}
                className="rounded-full bg-orange-500 px-6 py-2 font-semibold text-white transition hover:bg-orange-600"
              >
                次のレッスンへ進む →
              </Link>
            ) : (
              <Link
                href="/"
                className="rounded-full bg-orange-500 px-6 py-2 font-semibold text-white transition hover:bg-orange-600"
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
    </div>
  );
}
