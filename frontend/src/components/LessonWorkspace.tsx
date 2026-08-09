"use client";

import { useCallback, useEffect, useMemo, useRef, useState } from "react";
import CodeMirror from "@uiw/react-codemirror";
import { go } from "@codemirror/lang-go";
import { indentUnit } from "@codemirror/language";
import {
  linter,
  lintGutter,
  type Diagnostic as CMDiagnostic,
} from "@codemirror/lint";
import type { EditorView } from "@codemirror/view";
import SlideViewer from "@/components/SlideViewer";
import SuccessModal from "@/components/SuccessModal";
import {
  checkCode,
  formatCode,
  getAnswer,
  getDraft,
  saveDraft,
  submitSolution,
  type Diagnostic,
  type Problem,
  type SubmitResult,
} from "@/lib/api";
import { getLocalDraft, saveLocalDraft } from "@/lib/draft";
import { nextLessonId } from "@/lib/lessons";
import { useAuth } from "@/components/AuthProvider";

/** 自動保存のデバウンス間隔（入力停止からこの時間が経ったら保存する）。 */
const AUTOSAVE_DELAY_MS = 1200;

// go build の診断（行・列は1始まり）をCodeMirrorの文字オフセットに変換する。
// 列以降〜行末までを範囲として下線を引く（コンパイラは終了位置を教えてくれないため）。
function toCMDiagnostic(
  d: Diagnostic,
  view: EditorView,
  onFixImports: (view: EditorView) => void,
): CMDiagnostic {
  const lineCount = view.state.doc.lines;
  const lineNumber = Math.min(Math.max(d.line, 1), lineCount);
  const line = view.state.doc.line(lineNumber);
  const from = Math.min(line.from + Math.max(d.column - 1, 0), line.to);

  const isImportIssue = d.message.startsWith("undefined: ");

  return {
    from,
    to: line.to,
    severity: "error",
    message: d.message,
    actions: isImportIssue
      ? [
          {
            name: "インポートを自動修正",
            apply: (v) => onFixImports(v),
          },
        ]
      : undefined,
  };
}

export default function LessonWorkspace({ problem }: { problem: Problem }) {
  const { user, loading: authLoading, markCompleted } = useAuth();
  const [code, setCode] = useState(problem.starterCode);
  const [result, setResult] = useState<SubmitResult | null>(null);
  const [submitting, setSubmitting] = useState(false);
  const [errorMessage, setErrorMessage] = useState<string | null>(null);
  const [showSuccessModal, setShowSuccessModal] = useState(false);
  const [formatting, setFormatting] = useState(false);
  const [saveStatus, setSaveStatus] = useState<"idle" | "saving" | "saved">(
    "idle",
  );

  const [answerCode, setAnswerCode] = useState<string | null>(null);
  const [showAnswer, setShowAnswer] = useState(false);
  const [loadingAnswer, setLoadingAnswer] = useState(false);
  const [answerError, setAnswerError] = useState<string | null>(null);

  // CodeMirrorの診断ホバーツールチップからの「インポートを自動修正」は
  // マウスを正確に赤い波線へ重ねないと出てこず気づきにくいため、
  // 未解決のimportがある間は常設のボタンでも同じ修正を実行できるようにする。
  const [unresolvedImports, setUnresolvedImports] = useState<string[]>([]);

  // 自動保存の初期読み込みが終わるまでは、読み込み前のコードを保存で
  // 上書きしてしまわないようにガードする。
  const draftLoadedRef = useRef(false);
  const lastSavedCodeRef = useRef<string | null>(null);
  const saveTimerRef = useRef<ReturnType<typeof setTimeout> | null>(null);

  // レッスンを開いたら、保存済みの続きがあれば復元する
  // （ログイン中はサーバー優先、無ければlocalStorage、それも無ければ雛形コード）。
  useEffect(() => {
    if (authLoading) return;

    draftLoadedRef.current = false;
    let cancelled = false;

    (async () => {
      let draft: string | null = null;
      if (user) {
        try {
          draft = await getDraft(problem.id);
        } catch {
          draft = null;
        }
      }
      if (draft === null) {
        draft = getLocalDraft(problem.id);
      }

      if (cancelled) return;
      if (draft !== null) {
        setCode(draft);
        lastSavedCodeRef.current = draft;
      } else {
        lastSavedCodeRef.current = problem.starterCode;
      }
      setUnresolvedImports([]);
      draftLoadedRef.current = true;
    })();

    return () => {
      cancelled = true;
    };
  }, [problem.id, problem.starterCode, user, authLoading]);

  // 入力停止後、一定時間経ったら自動保存する。
  useEffect(() => {
    if (!draftLoadedRef.current || code === lastSavedCodeRef.current) {
      return;
    }

    setSaveStatus("saving");
    if (saveTimerRef.current) {
      clearTimeout(saveTimerRef.current);
    }
    saveTimerRef.current = setTimeout(() => {
      saveLocalDraft(problem.id, code);
      lastSavedCodeRef.current = code;
      setSaveStatus("saved");
      if (user) {
        saveDraft(problem.id, code).catch(() => {
          // 自動保存の失敗はユーザー操作を妨げないよう無視する
        });
      }
    }, AUTOSAVE_DELAY_MS);

    return () => {
      if (saveTimerRef.current) {
        clearTimeout(saveTimerRef.current);
      }
    };
  }, [code, problem.id, user]);

  const applyFixImports = useCallback(
    async (view: EditorView) => {
      try {
        const current = view.state.doc.toString();
        const fixed = await formatCode(problem.id, current, unresolvedImports);
        view.dispatch({
          changes: { from: 0, to: view.state.doc.length, insert: fixed },
        });
        setUnresolvedImports([]);
      } catch {
        // 失敗しても何もしない（次の入力で再チェックされる）
      }
    },
    [problem.id, unresolvedImports],
  );

  const goLinter = useMemo(
    () =>
      linter(
        async (view) => {
          const src = view.state.doc.toString();

          let diagnostics: Diagnostic[];
          try {
            diagnostics = await checkCode(problem.id, src);
          } catch {
            return [];
          }

          // 応答が届くまでの間にさらに入力が進んでいた場合、
          // 古い結果で最新の入力の診断を上書きしないよう破棄する。
          if (view.state.doc.toString() !== src) {
            return [];
          }

          const missingImports = Array.from(
            new Set(
              diagnostics
                .filter((d) => d.message.startsWith("undefined: "))
                .map((d) => d.message.slice("undefined: ".length)),
            ),
          );
          setUnresolvedImports(missingImports);

          return diagnostics.map((d) =>
            toCMDiagnostic(d, view, applyFixImports),
          );
        },
        { delay: 800 },
      ),
    [problem.id, applyFixImports],
  );

  async function handleFormat() {
    setFormatting(true);
    setErrorMessage(null);
    try {
      const fixed = await formatCode(problem.id, code, unresolvedImports);
      setCode(fixed);
      setUnresolvedImports([]);
    } catch (err) {
      setErrorMessage(
        err instanceof Error ? err.message : "フォーマットに失敗しました",
      );
    } finally {
      setFormatting(false);
    }
  }

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
            // GoはgofmtでタブインデントするためindentUnitをタブにする。
            // 未指定だとCodeMirrorのデフォルト（スペース2つ）が使われ、
            // 改行時に挿入される字下げが既存のタブ行と揃わなくなる。
            extensions={[go(), indentUnit.of("\t"), lintGutter(), goLinter]}
            onChange={setCode}
          />
        </div>

        {unresolvedImports.length > 0 && (
          <div className="flex flex-wrap items-center justify-between gap-3 rounded-lg border border-sky-200 bg-sky-50 px-4 py-2 text-sm text-sky-800 dark:border-sky-900 dark:bg-sky-950 dark:text-sky-200">
            <span>
              未解決のimportがあります:{" "}
              <code className="font-mono">{unresolvedImports.join(", ")}</code>
            </span>
            <button
              type="button"
              onClick={handleFormat}
              disabled={formatting}
              className="shrink-0 rounded-full bg-sky-600 px-3 py-1 text-xs font-semibold text-white transition hover:bg-sky-700 disabled:cursor-not-allowed disabled:opacity-60"
            >
              {formatting ? "修正中..." : "インポートを自動修正"}
            </button>
          </div>
        )}

        <div className="flex items-center gap-3">
          <button
            type="button"
            onClick={handleSubmit}
            disabled={submitting}
            className="rounded-full bg-orange-500 px-6 py-2 font-semibold text-white transition hover:bg-orange-600 disabled:cursor-not-allowed disabled:opacity-60"
          >
            {submitting ? "採点中..." : "実行して採点"}
          </button>

          <button
            type="button"
            onClick={handleFormat}
            disabled={formatting}
            title="未使用・不足しているimportの整理とgofmt整形を行います"
            className="rounded-full border border-black/10 px-6 py-2 font-semibold text-neutral-700 transition hover:bg-neutral-50 disabled:cursor-not-allowed disabled:opacity-60 dark:border-white/10 dark:text-neutral-300 dark:hover:bg-neutral-900"
          >
            {formatting ? "整形中..." : "フォーマット"}
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

          {saveStatus !== "idle" && (
            <span className="ml-auto text-xs text-neutral-400" aria-live="polite">
              {saveStatus === "saving" ? "保存中..." : "自動保存済み"}
            </span>
          )}
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
            <p className="mt-2 text-xs font-semibold tracking-wide text-neutral-500 dark:text-neutral-400">
              実行結果ログ（fmt.Printなどの出力もここに表示されます）
            </p>
            <pre className="mt-1 overflow-x-auto whitespace-pre-wrap text-sm">
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
