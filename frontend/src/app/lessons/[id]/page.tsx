import Link from "next/link";
import { getProblem } from "@/lib/api";
import { courseForLesson } from "@/lib/lessons";
import { courseColors } from "@/lib/courseColors";
import LessonWorkspace from "@/components/LessonWorkspace";

export default async function LessonPage(props: PageProps<"/lessons/[id]">) {
  const { id } = await props.params;
  const problem = await getProblem(id);
  const course = courseForLesson(id);
  const positionInCourse = course
    ? course.lessons.findIndex((l) => l.id === id) + 1
    : 0;
  const borderClass = course ? courseColors[course.accent].border : "";

  return (
    <div className="flex min-h-screen flex-col">
      <header
        className={`flex items-center justify-between gap-4 border-b border-t-4 border-b-black/10 px-6 py-4 dark:border-b-white/10 ${borderClass}`}
      >
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
          className="shrink-0 text-sm font-semibold text-orange-600 hover:underline dark:text-orange-400"
        >
          ← レッスン一覧に戻る
        </Link>
      </header>
      <LessonWorkspace problem={problem} />
    </div>
  );
}
