import LessonList from "@/components/LessonList";
import StreakBadges from "@/components/StreakBadges";
import Mascot from "@/components/Mascot";

export default function Home() {
  return (
    <div className="mx-auto flex w-full max-w-2xl flex-1 flex-col gap-6 px-6 py-12">
      <div className="flex items-center gap-4 rounded-2xl bg-gradient-to-br from-orange-50 via-white to-pink-50 p-6 dark:from-orange-950 dark:via-neutral-900 dark:to-pink-950">
        <Mascot pose="wave" className="h-24 w-24 shrink-0 sm:h-28 sm:w-28" />
        <div>
          <h1 className="bg-gradient-to-r from-orange-500 to-pink-500 bg-clip-text text-3xl font-extrabold text-transparent">
            CodeForge
          </h1>
          <p className="mt-1 text-neutral-600 dark:text-neutral-400">
            コードを書いて実行しながらプログラミングを学べるサイトです。
          </p>
        </div>
      </div>

      <StreakBadges />
      <LessonList />
    </div>
  );
}
