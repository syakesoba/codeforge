"use client";

import { useState } from "react";
import CodeMirror from "@uiw/react-codemirror";
import { go } from "@codemirror/lang-go";
import SlideViewer from "@/components/SlideViewer";
import SuccessModal from "@/components/SuccessModal";
import {
  getAnswer,
  submitSolution,
  type Problem,
  type SubmitResult,
} from "@/lib/api";
import { nextLessonId } from "@/lib/lessons";
import { useAuth } from "@/components/AuthProvider";

export default function LessonWorkspace({ problem }: { problem: Problem }) {
  const { markCompleted } = useAuth();
  const [code, setCode] = useState(problem.starterCode);
  const [result, setResult] = useState<SubmitResult | null>(null);
  const [submitting, setSubmitting] = useState(false);
  const [errorMessage, setErrorMessage] = useState<string | null>(null);
  const [showSuccessModal, setShowSuccessModal] = useState(false);

  const [answerCode, setAnswerCode] = useState<string | null>(null);
  const [showAnswer, setShowAnswer] = useState(false);
  const [loadingAnswer, setLoadingAnswer] = useState(false);
  const [answerError, setAnswerError] = useState<string | null>(null);

  async function handleSubmit() {
    setSubmitting(true);
    setErrorMessage(null);
    setResult(null);
    try {
      const res = await submitSolution(problem.id, code);
      setResult(res);
      if (res.passed) {
        markCompleted(problem.id);
        setShowSuccessModal(true);
      }
    } catch (err) {
      setErrorMessage(
        err instanceof Error ? err.message : "採点中にエラーが発生しました",
      );
    } finally {
      setSubmitting(false);
    }
  }

  async function handleToggleAnswer() {
    if (showAnswer) {
      setShowAnswer(false);
      return;
    }
    if (answerCode !== null) {
      setShowAnswer(true);
      return;
    }
    setLoadingAnswer(true);
    setAnswerError(null);
    try {
      const code = await getAnswer(problem.id);
      setAnswerCode(code);
      setShowAnswer(true);
    } catch (err) {
      setAnswerError(
        err instanceof Error ? err.message : "答えの取得に失敗しました",
      );
    } finally {
      setLoadingAnswer(false);
    }
  }

  return (
    <div className="grid flex-1 grid-cols-1 gap-6 p-6 lg:grid-cols-2">
      <div className="flex h-[70vh] flex-col overflow-hidden rounded-xl border border-black/10 bg-white p-6 dark:border-white/10 dark:bg-neutral-900">
        <SlideViewer markdown={problem.markdown} />
      </div>

      <div className="flex flex-col gap-4">
        <div className="overflow-hidden rounded-xl border border-black/10 dark:border-white/10">
          <CodeMirror
            value={code}
            height="360px"
            extensions={[go()]}
            onChange={setCode}
          />
        </div>

        <div className="flex items-center gap-3">
          <button
            type="button"
            onClick={handleSubmit}
            disabled={submitting}
            className="rounded-full bg-emerald-600 px-6 py-2 font-semibold text-white transition hover:bg-emerald-700 disabled:cursor-not-allowed disabled:opacity-60"
          >
            {submitting ? "採点中..." : "実行して採点"}
          </button>

          <button
            type="button"
            onClick={handleToggleAnswer}
            disabled={loadingAnswer}
            className="rounded-full border border-black/10 px-6 py-2 font-semibold text-neutral-700 transition hover:bg-neutral-50 disabled:cursor-not-allowed disabled:opacity-60 dark:border-white/10 dark:text-neutral-300 dark:hover:bg-neutral-900"
          >
            {loadingAnswer
              ? "読み込み中..."
              : showAnswer
                ? "答えを隠す"
                : "答えを見る"}
          </button>
        </div>

        {answerError && (
          <p className="rounded-lg bg-red-50 p-3 text-sm text-red-700 dark:bg-red-950 dark:text-red-300">
            {answerError}
          </p>
        )}

        {showAnswer && answerCode && (
          <div className="rounded-xl border border-black/10 dark:border-white/10">
            <p className="rounded-t-xl border-b border-black/10 bg-neutral-50 px-4 py-2 text-xs font-bold text-neutral-500 dark:border-white/10 dark:bg-neutral-900 dark:text-neutral-400">
              模範解答
            </p>
            <pre className="overflow-x-auto p-4 text-sm">
              <code>{answerCode}</code>
            </pre>
          </div>
        )}

        {errorMessage && (
          <p className="rounded-lg bg-red-50 p-3 text-sm text-red-700 dark:bg-red-950 dark:text-red-300">
            {errorMessage}
          </p>
        )}

        {result && (
          <div
            className={`rounded-xl border p-4 ${
              result.passed
                ? "border-emerald-300 bg-emerald-50 dark:border-emerald-800 dark:bg-emerald-950"
                : "border-amber-300 bg-amber-50 dark:border-amber-800 dark:bg-amber-950"
            }`}
          >
            <p className="font-bold">
              {result.passed ? "✅ 合格！" : "❌ 不合格"}
              <span className="ml-2 text-xs font-normal opacity-70">
                {result.durationMs}ms
              </span>
            </p>
            <pre className="mt-2 overflow-x-auto whitespace-pre-wrap text-sm">
              {result.output}
            </pre>
          </div>
        )}
      </div>

      {showSuccessModal && (
        <SuccessModal
          nextLessonId={nextLessonId(problem.id)}
          onClose={() => setShowSuccessModal(false)}
        />
      )}
    </div>
  );
}
