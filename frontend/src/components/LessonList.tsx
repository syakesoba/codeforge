"use client";

import Link from "next/link";
import { courses } from "@/lib/lessons";
import { courseColors } from "@/lib/courseColors";
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
        const colors = courseColors[course.accent];
        return (
          <div
            key={course.id}
            className={`overflow-hidden rounded-xl border-x border-b border-t-4 border-x-black/10 border-b-black/10 bg-white shadow-sm dark:border-x-white/10 dark:border-b-white/10 dark:bg-neutral-900 ${colors.border}`}
          >
            <div className="p-5">
              <span
                className={`inline-block rounded-full px-3 py-0.5 text-xs font-bold tracking-wide ${colors.badge}`}
              >
                COURSE
              </span>
              <h2 className="mt-2 mb-3 text-lg font-bold">{course.title}</h2>
              <ProgressBar
                completed={completedInCourse}
                total={course.lessons.length}
                fillClassName={colors.progress}
                className="mb-3"
              />
              <ul className="flex flex-col gap-2">
                {course.lessons.map((lesson) => (
                  <li key={lesson.id}>
                    <Link
                      href={`/lessons/${lesson.id}`}
                      className="flex items-center justify-between rounded-lg px-3 py-2 font-medium text-neutral-700 transition hover:bg-neutral-50 dark:text-neutral-300 dark:hover:bg-neutral-800"
                    >
                      <span>{lesson.title}</span>
                      {completed.has(lesson.id) && (
                        <span
                          aria-label="合格済み"
                          title="合格済み"
                          className="flex h-5 w-5 shrink-0 items-center justify-center rounded-full bg-emerald-600 text-xs font-bold text-white"
                        >
                          ✓
                        </span>
                      )}
                    </Link>
                  </li>
                ))}
              </ul>
            </div>
          </div>
        );
      })}
    </div>
  );
}
