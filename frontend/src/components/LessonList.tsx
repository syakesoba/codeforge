"use client";

import Link from "next/link";
import { courses } from "@/lib/lessons";
import { useAuth } from "@/components/AuthProvider";
import ProgressBar from "@/components/ProgressBar";

export default function LessonList() {
  const { completed } = useAuth();

  return (
    <div className="flex flex-col gap-6">
      {courses.map((course) => {
        const completedInCourse = course.lessons.filter((l) =>
          completed.has(l.id),
        ).length;
        return (
        <div
          key={course.id}
          className="rounded-xl border border-black/10 p-5 dark:border-white/10"
        >
          <p className="text-sm text-neutral-500">Course</p>
          <h2 className="mb-2 text-lg font-semibold">{course.title}</h2>
          <ProgressBar
            completed={completedInCourse}
            total={course.lessons.length}
            className="mb-3"
          />
          <ul className="flex flex-col gap-2">
            {course.lessons.map((lesson) => (
              <li key={lesson.id}>
                <Link
                  href={`/lessons/${lesson.id}`}
                  className="flex items-center justify-between rounded-lg px-3 py-2 font-medium text-emerald-700 transition hover:bg-emerald-50 dark:text-emerald-400 dark:hover:bg-emerald-950"
                >
                  <span>{lesson.title}</span>
                  {completed.has(lesson.id) && (
                    <span
                      aria-label="合格済み"
                      title="合格済み"
                      className="flex h-5 w-5 items-center justify-center rounded-full bg-emerald-600 text-xs font-bold text-white"
                    >
                      ✓
                    </span>
                  )}
                </Link>
              </li>
            ))}
          </ul>
        </div>
        );
      })}
    </div>
  );
}
