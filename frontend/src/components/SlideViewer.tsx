"use client";

import { useState } from "react";
import ReactMarkdown from "react-markdown";
import { splitIntoSlides } from "@/lib/slides";

export default function SlideViewer({ markdown }: { markdown: string }) {
  const [slides] = useState(() => splitIntoSlides(markdown));
  const [index, setIndex] = useState(0);

  const slide = slides[index];
  const isFirst = index === 0;
  const isLast = index === slides.length - 1;

  return (
    <div className="flex h-full flex-col">
      <div className="flex-1 overflow-y-auto">
        <p className="mb-2 text-xs font-bold tracking-wide text-emerald-600 dark:text-emerald-400">
          スライド {index + 1} / {slides.length}
        </p>
        {slide.title && (
          <h2 className="mb-3 text-lg font-bold">{slide.title}</h2>
        )}
        <div className="prose prose-neutral dark:prose-invert max-w-none prose-pre:bg-neutral-900 prose-pre:text-neutral-100">
          <ReactMarkdown>{slide.body}</ReactMarkdown>
        </div>
      </div>

      <div className="mt-4 flex items-center justify-between border-t border-black/10 pt-4 dark:border-white/10">
        <button
          type="button"
          onClick={() => setIndex((i) => Math.max(0, i - 1))}
          disabled={isFirst}
          className="rounded-full border border-black/10 px-4 py-1.5 text-sm font-semibold text-neutral-700 transition disabled:cursor-not-allowed disabled:opacity-40 dark:border-white/10 dark:text-neutral-300"
        >
          ← 前へ
        </button>

        <div className="flex gap-1.5">
          {slides.map((s, i) => (
            <button
              key={i}
              type="button"
              aria-label={`スライド ${i + 1}: ${s.title || ""}`}
              onClick={() => setIndex(i)}
              className={`h-1.5 w-1.5 rounded-full transition ${
                i === index
                  ? "bg-emerald-600"
                  : "bg-neutral-300 dark:bg-neutral-700"
              }`}
            />
          ))}
        </div>

        <button
          type="button"
          onClick={() => setIndex((i) => Math.min(slides.length - 1, i + 1))}
          disabled={isLast}
          className="rounded-full bg-emerald-600 px-4 py-1.5 text-sm font-semibold text-white transition disabled:cursor-not-allowed disabled:opacity-40"
        >
          次へ →
        </button>
      </div>
    </div>
  );
}
