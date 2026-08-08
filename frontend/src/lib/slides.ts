export type Slide = {
  title: string;
  body: string;
};

/**
 * レッスンのMarkdownを `##` / `###` 見出しの境界でスライドに分割する。
 * 先頭の `#` （レッスンタイトル）はページヘッダー側で表示済みのため読み飛ばす。
 * フェンス付きコードブロック（```）内の `#` 始まりの行は見出しとして扱わない。
 */
export function splitIntoSlides(markdown: string): Slide[] {
  const lines = markdown.split("\n");
  const slides: Slide[] = [];

  let title = "";
  let bodyLines: string[] = [];
  let inFence = false;

  const flush = () => {
    const body = bodyLines.join("\n").trim();
    if (title || body) {
      slides.push({ title, body });
    }
  };

  for (const line of lines) {
    if (/^\s*```/.test(line)) {
      inFence = !inFence;
      bodyLines.push(line);
      continue;
    }

    const heading = !inFence && /^(#{1,3})\s+(.*)$/.exec(line);
    if (heading) {
      const level = heading[1].length;
      if (level === 1) {
        // レッスンタイトル行はページヘッダーで表示済みなのでスキップ
        continue;
      }
      flush();
      title = heading[2].trim();
      bodyLines = [];
      continue;
    }

    bodyLines.push(line);
  }
  flush();

  return slides;
}
