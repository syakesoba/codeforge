export default function ProgressBar({
  completed,
  total,
  label,
  className = "",
  fillClassName = "bg-emerald-600",
}: {
  completed: number;
  total: number;
  label?: string;
  className?: string;
  /** 塗りつぶし色のTailwindクラス（コースごとのアクセントカラーに差し替える場合に使う）。 */
  fillClassName?: string;
}) {
  const percent = total > 0 ? Math.round((completed / total) * 100) : 0;

  return (
    <div className={className}>
      {label && (
        <div className="mb-1 flex items-center justify-between text-xs font-semibold text-neutral-500 dark:text-neutral-400">
          <span>{label}</span>
          <span>
            {completed}/{total}
          </span>
        </div>
      )}
      <div
        role="progressbar"
        aria-valuenow={percent}
        aria-valuemin={0}
        aria-valuemax={100}
        className="h-2 w-full overflow-hidden rounded-full bg-neutral-200 dark:bg-neutral-800"
      >
        <div
          className={`h-full rounded-full transition-[width] duration-500 ${fillClassName}`}
          style={{ width: `${percent}%` }}
        />
      </div>
    </div>
  );
}
