import Link from "next/link";
import { getProblem } from "@/lib/api";
import { courseForLesson } from "@/lib/lessons";
import LessonWorkspace from "@/components/LessonWorkspace";

export default async function LessonPage(props: PageProps<"/lessons/[id]">) {
  const { id } = await props.params;
  const problem = await getProblem(id);
  const course = courseForLesson(id);
  const positionInCourse = course
    ? course.lessons.findIndex((l) => l.id === id) + 1
    : 0;

  return (
    <div className="flex min-h-screen flex-col">
      <header className="flex items-center justify-between gap-4 border-b border-black/10 px-6 py-4 dark:border-white/10">
        <div>
          <p className="text-sm text-neutral-500">
            {course?.title}
            {course && (
              <span className="ml-2 text-neutral-400">
                レッスン {positionInCourse}/{course.lessons.length}
              </span>
            )}
          </p>
          <h1 className="text-xl font-bold">{problem.title}</h1>
        </div>
        <Link
          href="/"
          className="shrink-0 text-sm font-semibold text-emerald-700 hover:underline dark:text-emerald-400"
        >
          ← レッスン一覧に戻る
        </Link>
      </header>
      <LessonWorkspace problem={problem} />
    </div>
  );
}
