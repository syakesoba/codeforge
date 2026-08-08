import LessonList from "@/components/LessonList";
import StreakBadges from "@/components/StreakBadges";

export default function Home() {
  return (
    <div className="mx-auto flex max-w-2xl flex-1 flex-col justify-center gap-6 px-6 py-16">
      <h1 className="text-3xl font-bold">CodeForge</h1>
      <p className="text-neutral-600 dark:text-neutral-400">
        コードを書いて実行しながらプログラミングを学べるサイトです。
      </p>

      <StreakBadges />
      <LessonList />
    </div>
  );
}
