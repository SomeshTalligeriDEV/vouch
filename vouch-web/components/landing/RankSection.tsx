"use client";

import { useEffect, useRef, useState } from "react";

const RANKS = [
  {
    id: "e",
    label: "E-Rank",
    range: "0 – 99",
    sub: "unranked · lurking",
    color: "#5C6268",
    glow: "rgba(92,98,104,0.4)",
    border: "rgba(92,98,104,0.3)",
  },
  {
    id: "d",
    label: "D-Rank",
    range: "100 – 199",
    sub: "first ships",
    color: "#78a1bb",
    glow: "rgba(120,161,187,0.45)",
    border: "rgba(120,161,187,0.35)",
  },
  {
    id: "c",
    label: "C-Rank",
    range: "200 – 499",
    sub: "gaining users",
    color: "#4caf80",
    glow: "rgba(76,175,128,0.45)",
    border: "rgba(76,175,128,0.35)",
  },
  {
    id: "b",
    label: "B-Rank",
    range: "500 – 999",
    sub: "shipping for revenue",
    color: "#C8F24C",
    glow: "rgba(200,242,76,0.6)",
    border: "rgba(200,242,76,0.5)",
    active: true,
  },
  {
    id: "a",
    label: "A-Rank",
    range: "1,000 – 1,999",
    sub: "proven operator",
    color: "#e8a838",
    glow: "rgba(232,168,56,0.5)",
    border: "rgba(232,168,56,0.4)",
  },
  {
    id: "s",
    label: "S-Rank",
    range: "2,000+",
    sub: "top 1% of builders",
    color: "#ff6b6b",
    glow: "rgba(255,107,107,0.55)",
    border: "rgba(255,107,107,0.45)",
  },
];

const FORMULA = [
  { label: "Ships", key: "ships", value: 12, max: 30, color: "#C8F24C" },
  { label: "Users", key: "users", value: 3400, max: 10000, color: "#78a1bb" },
  { label: "Revenue", key: "revenue", value: 2100, max: 10000, color: "#4caf80" },
  { label: "Velocity", key: "velocity", value: 280, max: 500, color: "#e8a838" },
];

function useInView(threshold = 0.15) {
  const ref = useRef<HTMLDivElement>(null);
  const [visible, setVisible] = useState(false);
  useEffect(() => {
    const el = ref.current;
    if (!el) return;
    const obs = new IntersectionObserver(
      ([e]) => { if (e.isIntersecting) { setVisible(true); obs.disconnect(); } },
      { threshold }
    );
    obs.observe(el);
    return () => obs.disconnect();
  }, [threshold]);
  return { ref, visible };
}

function useCountUp(target: number, active: boolean, duration = 1800) {
  const [value, setValue] = useState(0);
  useEffect(() => {
    if (!active) return;
    const start = performance.now();
    const tick = (now: number) => {
      const t = Math.min(1, (now - start) / duration);
      const eased = 1 - Math.pow(1 - t, 3);
      setValue(Math.round(target * eased));
      if (t < 1) requestAnimationFrame(tick);
    };
    requestAnimationFrame(tick);
  }, [active, target, duration]);
  return value;
}

function ScoreBar({ item, active, delay }: { item: typeof FORMULA[0]; active: boolean; delay: number }) {
  const pct = Math.min(100, (item.value / item.max) * 100);
  return (
    <div style={{ display: "flex", alignItems: "center", gap: 12, opacity: active ? 1 : 0, transform: active ? "none" : "translateX(-16px)", transition: `opacity .5s ${delay}ms, transform .5s cubic-bezier(.34,1.56,.64,1) ${delay}ms` }}>
      <span style={{ width: 64, fontSize: 11, fontFamily: "'Space Mono',monospace", color: "rgba(244,246,240,.5)", textTransform: "uppercase", letterSpacing: ".04em", flexShrink: 0 }}>
        {item.label}
      </span>
      <div style={{ flex: 1, height: 4, background: "rgba(244,246,240,.08)", borderRadius: 2, overflow: "hidden" }}>
        <div style={{ height: "100%", width: active ? `${pct}%` : "0%", background: item.color, borderRadius: 2, boxShadow: `0 0 6px ${item.color}`, transition: `width 1.4s cubic-bezier(.16,1,.3,1) ${delay + 200}ms` }} />
      </div>
      <span style={{ width: 40, fontSize: 11, fontFamily: "'Space Mono',monospace", color: item.color, textAlign: "right", flexShrink: 0 }}>
        {item.value >= 1000 ? `${(item.value / 1000).toFixed(1)}k` : item.value}
      </span>
    </div>
  );
}

export function RankSection() {
  const { ref, visible } = useInView(0.1);
  const score = useCountUp(842, visible);

  return (
    <section
      ref={ref}
      style={{
        background: "#0C0E10",
        padding: "100px 28px",
        position: "relative",
        overflow: "hidden",
      }}
    >
      {/* background grid */}
      <div
        aria-hidden
        style={{
          position: "absolute",
          inset: 0,
          backgroundImage:
            "linear-gradient(rgba(200,242,76,.04) 1px,transparent 1px),linear-gradient(90deg,rgba(200,242,76,.04) 1px,transparent 1px)",
          backgroundSize: "48px 48px",
          maskImage: "radial-gradient(ellipse 80% 60% at 50% 40%, black 20%, transparent 100%)",
        }}
      />
      {/* ambient glow */}
      <div
        aria-hidden
        style={{
          position: "absolute",
          top: "20%",
          left: "50%",
          transform: "translateX(-50%)",
          width: 600,
          height: 300,
          background: "radial-gradient(ellipse,rgba(200,242,76,.07) 0%,transparent 70%)",
          pointerEvents: "none",
        }}
      />

      <div style={{ maxWidth: 1200, margin: "0 auto", position: "relative" }}>

        {/* formula + live score */}
        <div style={{ display: "grid", gridTemplateColumns: "1fr auto 1fr", gap: 48, alignItems: "center", marginBottom: 80 }}>

          {/* formula bars */}
          <div style={{ display: "flex", flexDirection: "column", gap: 16 }}>
            {FORMULA.map((item, i) => (
              <ScoreBar key={item.key} item={item} active={visible} delay={i * 80} />
            ))}
          </div>

          {/* equals + total */}
          <div style={{ display: "flex", flexDirection: "column", alignItems: "center", gap: 8, opacity: visible ? 1 : 0, transition: "opacity .6s .4s" }}>
            <div style={{ fontFamily: "'Space Mono',monospace", fontSize: 32, color: "rgba(244,246,240,.25)", lineHeight: 1 }}>×</div>
            <div style={{ fontSize: 11, fontFamily: "'Space Mono',monospace", color: "rgba(244,246,240,.3)", letterSpacing: ".1em", textTransform: "uppercase" }}>stripe</div>
          </div>

          {/* score orb */}
          <div style={{ display: "flex", flexDirection: "column", alignItems: "center", gap: 12, opacity: visible ? 1 : 0, transform: visible ? "none" : "scale(.9)", transition: "opacity .7s .2s, transform .7s cubic-bezier(.34,1.56,.64,1) .2s" }}>
            <div style={{ position: "relative", width: 160, height: 160, display: "flex", alignItems: "center", justifyContent: "center" }}>
              {/* rings */}
              {[160, 130, 100].map((s, i) => (
                <div key={s} style={{ position: "absolute", width: s, height: s, border: `${i === 0 ? 1 : i === 1 ? 1.5 : 2}px solid rgba(200,242,76,${.08 + i * .06})`, borderRadius: "50%", animation: `spin ${18 + i * 7}s linear infinite ${i % 2 ? "reverse" : ""}` }} />
              ))}
              <div style={{ position: "absolute", inset: 0, background: "radial-gradient(circle, rgba(200,242,76,.12) 0%, transparent 70%)", borderRadius: "50%" }} />
              <div style={{ textAlign: "center", position: "relative" }}>
                <div style={{ fontFamily: "'Space Grotesk',sans-serif", fontWeight: 700, fontSize: 48, color: "#C8F24C", lineHeight: 1, textShadow: "0 0 30px rgba(200,242,76,.6)" }}>
                  {score}
                </div>
                <div style={{ fontFamily: "'Space Mono',monospace", fontSize: 10, color: "rgba(200,242,76,.7)", letterSpacing: ".12em", textTransform: "uppercase", marginTop: 4 }}>
                  Builder Score
                </div>
              </div>
            </div>
            <div style={{ padding: "5px 14px", background: "rgba(200,242,76,.12)", border: "1px solid rgba(200,242,76,.4)", borderRadius: 999, fontFamily: "'Space Mono',monospace", fontSize: 11, color: "#C8F24C", letterSpacing: ".06em" }}>
              ★ B-RANK
            </div>
          </div>
        </div>

        {/* rank ladder */}
        <div style={{ display: "grid", gridTemplateColumns: "repeat(6,1fr)", gap: 12 }}>
          {RANKS.map((rank, i) => (
            <div
              key={rank.id}
              style={{
                position: "relative",
                padding: "20px 16px",
                borderRadius: 12,
                border: `1.5px solid ${rank.border}`,
                background: rank.active
                  ? `linear-gradient(135deg, rgba(200,242,76,.12) 0%, rgba(200,242,76,.04) 100%)`
                  : "rgba(255,255,255,.02)",
                boxShadow: rank.active ? `0 0 32px ${rank.glow}, inset 0 0 32px rgba(200,242,76,.04)` : "none",
                opacity: visible ? 1 : 0,
                transform: visible ? "none" : "translateY(20px)",
                transition: `opacity .5s ${i * 70}ms, transform .5s cubic-bezier(.34,1.56,.64,1) ${i * 70}ms`,
                cursor: "default",
              }}
            >
              {rank.active && (
                <div style={{ position: "absolute", top: -1, left: "50%", transform: "translateX(-50%)", padding: "2px 10px", background: "#C8F24C", borderRadius: "0 0 8px 8px", fontFamily: "'Space Mono',monospace", fontSize: 9, fontWeight: 700, color: "#14181B", letterSpacing: ".08em", textTransform: "uppercase" }}>
                  YOU
                </div>
              )}
              <div style={{ fontFamily: "'Space Grotesk',sans-serif", fontWeight: 700, fontSize: 18, color: rank.color, textShadow: rank.active ? `0 0 20px ${rank.glow}` : "none", marginBottom: 8 }}>
                {rank.label}
              </div>
              <div style={{ fontFamily: "'Space Mono',monospace", fontSize: 10, color: rank.color, opacity: .7, marginBottom: 6, letterSpacing: ".04em" }}>
                {rank.range}
              </div>
              <div style={{ fontFamily: "'DM Sans',sans-serif", fontSize: 11, color: "rgba(244,246,240,.45)", lineHeight: 1.4 }}>
                {rank.sub}
              </div>
              {rank.active && (
                <div style={{ marginTop: 12, fontFamily: "'Space Mono',monospace", fontSize: 13, fontWeight: 700, color: "#C8F24C" }}>
                  842
                </div>
              )}
            </div>
          ))}
        </div>

      </div>

      <style>{`
        @keyframes spin { to { transform: rotate(360deg); } }
      `}</style>
    </section>
  );
}
