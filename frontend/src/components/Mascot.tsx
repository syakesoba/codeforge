// CodeForgeのオリジナルマスコット（既存キャラクターを模したものではない、独自デザイン）。
// pose="wave" は歓迎ポーズ（トップページ・ログイン/新規登録ページ）、
// pose="celebrate" は合格祝いポーズ（SuccessModal）で使う。
type MascotPose = "wave" | "celebrate";

function Sparkle({ cx, cy, size }: { cx: number; cy: number; size: number }) {
  return (
    <g transform={`translate(${cx} ${cy}) rotate(45)`}>
      <rect
        x={-size / 2}
        y={-size / 8}
        width={size}
        height={size / 4}
        rx={size / 8}
        fill="#FFC94D"
      />
      <rect
        x={-size / 8}
        y={-size / 2}
        width={size / 4}
        height={size}
        rx={size / 8}
        fill="#FFC94D"
      />
    </g>
  );
}

export default function Mascot({
  pose = "wave",
  className = "",
}: {
  pose?: MascotPose;
  className?: string;
}) {
  return (
    <svg
      viewBox="0 0 200 200"
      className={className}
      role="img"
      aria-label="CodeForgeのマスコット"
    >
      <ellipse cx="100" cy="186" rx="52" ry="8" className="fill-black/10 dark:fill-black/40" />

      {pose === "celebrate" && (
        <>
          <Sparkle cx={28} cy={38} size={14} />
          <Sparkle cx={172} cy={32} size={10} />
          <Sparkle cx={100} cy={12} size={12} />
        </>
      )}

      {/* 腕（後ろ側） */}
      {pose === "wave" ? (
        <path
          d="M45 122 Q22 132 28 158"
          stroke="#FF7A3D"
          strokeWidth="11"
          strokeLinecap="round"
          fill="none"
        />
      ) : (
        <path
          d="M48 112 Q16 92 22 54"
          stroke="#FF7A3D"
          strokeWidth="11"
          strokeLinecap="round"
          fill="none"
        />
      )}

      {/* 本体 */}
      <path
        d="M100 22 C147 22 172 58 172 102 C172 150 141 180 100 180 C59 180 28 150 28 102 C28 58 53 22 100 22 Z"
        fill="#FF7A3D"
      />

      {/* おなか */}
      <ellipse cx="100" cy="118" rx="44" ry="38" fill="#FFE4D1" />

      {/* ほっぺ */}
      <circle cx="66" cy="106" r="7" fill="#FFB199" opacity="0.7" />
      <circle cx="134" cy="106" r="7" fill="#FFB199" opacity="0.7" />

      {/* 目 */}
      <circle cx="82" cy="92" r="8.5" fill="#2A3541" />
      <circle cx="118" cy="92" r="8.5" fill="#2A3541" />
      <circle cx="85" cy="89" r="2.5" fill="#fff" />
      <circle cx="121" cy="89" r="2.5" fill="#fff" />

      {/* 口 */}
      {pose === "celebrate" ? (
        <path
          d="M83 109 Q100 132 117 109"
          stroke="#2A3541"
          strokeWidth="4.5"
          fill="none"
          strokeLinecap="round"
        />
      ) : (
        <path
          d="M87 110 Q100 120 113 110"
          stroke="#2A3541"
          strokeWidth="4"
          fill="none"
          strokeLinecap="round"
        />
      )}

      {/* アンテナ */}
      <line x1="100" y1="22" x2="100" y2="6" stroke="#FF7A3D" strokeWidth="4.5" strokeLinecap="round" />
      <circle cx="100" cy="6" r="6.5" fill="#FFC94D" />

      {/* 腕（前側） */}
      {pose === "wave" ? (
        <path
          d="M155 112 Q186 96 178 62"
          stroke="#FF9A5C"
          strokeWidth="11"
          strokeLinecap="round"
          fill="none"
        />
      ) : (
        <path
          d="M152 112 Q184 92 178 54"
          stroke="#FF9A5C"
          strokeWidth="11"
          strokeLinecap="round"
          fill="none"
        />
      )}
    </svg>
  );
}
