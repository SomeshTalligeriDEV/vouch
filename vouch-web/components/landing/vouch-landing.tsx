"use client";

import { useEffect, useRef } from "react";
import Script from "next/script";

import { DESIGN_HTML, DESIGN_KEYFRAMES } from "./design";

// Renders the imported Claude Design (Vouch.dc.html) and wires up its
// interactions: hover micro-effects, FAQ toggles, scroll reveals, the score
// gauge, and the waitlist form. Falls back to IntersectionObserver when GSAP
// has not loaded.
declare global {
  interface Window {
    gsap?: any;
    ScrollTrigger?: any;
  }
}

export function VouchLanding() {
  const ref = useRef<HTMLDivElement>(null);

  useEffect(() => {
    const root = ref.current?.querySelector<HTMLElement>("[data-vouch-root]");
    if (!root) return;

    // --- hover micro-interactions (style-hover="...") ---
    const hoverEls = root.querySelectorAll<HTMLElement>("[style-hover]");
    const cleanups: Array<() => void> = [];
    hoverEls.forEach((el) => {
      const hover = el.getAttribute("style-hover") || "";
      const base = el.getAttribute("style") || "";
      const enter = () => el.setAttribute("style", base + ";" + hover);
      const leave = () => el.setAttribute("style", base);
      el.addEventListener("mouseenter", enter);
      el.addEventListener("mouseleave", leave);
      cleanups.push(() => {
        el.removeEventListener("mouseenter", enter);
        el.removeEventListener("mouseleave", leave);
      });
    });

    // --- FAQ +/– marks ---
    root.querySelectorAll<HTMLDetailsElement>("details[data-faq]").forEach((d) => {
      const handler = () => {
        const m = d.querySelector("[data-faq-mark]");
        if (m) m.textContent = d.open ? "–" : "+";
      };
      d.addEventListener("toggle", handler);
      cleanups.push(() => d.removeEventListener("toggle", handler));
    });

    // --- waitlist form -> toast ---
    const form = root.querySelector<HTMLFormElement>("[data-waitlist-form]");
    const input = root.querySelector<HTMLInputElement>("[data-waitlist-input]");
    const toast = root.querySelector<HTMLElement>("[data-toast]");
    let toastTimer: ReturnType<typeof setTimeout>;
    if (form) {
      const onSubmit = (e: Event) => {
        e.preventDefault();
        if (!input?.value) return;
        try {
          localStorage.setItem("vouch_waitlist_email", input.value);
        } catch {
          /* ignore */
        }
        input.value = "";
        if (toast) {
          toast.hidden = false;
          clearTimeout(toastTimer);
          toastTimer = setTimeout(() => (toast.hidden = true), 3600);
        }
      };
      form.addEventListener("submit", onSubmit);
      cleanups.push(() => form.removeEventListener("submit", onSubmit));
    }

    // --- gauge + rail fill ---
    const runGauge = (sec: HTMLElement) => {
      const ring = sec.querySelector<SVGElement>("[data-gauge-ring]");
      const num = sec.querySelector<HTMLElement>("[data-gauge-num]");
      const fill = sec.querySelector<HTMLElement>("[data-rail-vfill]");
      const gold = sec.querySelector<HTMLElement>("[data-tier-gold]");
      const TARGET = 842;
      if (fill && gold) {
        fill.style.height =
          Math.max(0, gold.offsetTop + gold.offsetHeight / 2 - 14) + "px";
      }
      if (ring) {
        const C = parseFloat(ring.getAttribute("stroke-dasharray") || "578");
        ring.animate(
          [{ strokeDashoffset: String(C) }, { strokeDashoffset: String(C * (1 - 0.842)) }],
          { duration: 1700, easing: "cubic-bezier(.16,1,.3,1)", fill: "forwards" },
        );
      }
      if (num) {
        const start = performance.now();
        const tick = (now: number) => {
          const t = Math.min(1, (now - start) / 1700);
          const eased = 1 - Math.pow(1 - t, 3);
          num.textContent = String(Math.round(TARGET * eased));
          if (t < 1) requestAnimationFrame(tick);
        };
        requestAnimationFrame(tick);
      }
    };

    // --- scroll reveal via IntersectionObserver (works with or without GSAP) ---
    const reveals = root.querySelectorAll<HTMLElement>("[data-reveal]");
    const obs = new IntersectionObserver(
      (entries) => {
        entries.forEach((entry) => {
          if (!entry.isIntersecting) return;
          const el = entry.target as HTMLElement;
          el.style.opacity = "1";
          el.style.transform = "none";
          el.querySelectorAll<HTMLElement>("[data-stagger]").forEach((kid, i) => {
            kid.style.transition =
              "opacity .5s ease, transform .5s cubic-bezier(.34,1.56,.64,1)";
            kid.style.transitionDelay = 150 + i * 90 + "ms";
            kid.style.opacity = "1";
            kid.style.transform = "none";
          });
          root
            .querySelectorAll<SVGPathElement>("[data-conn] path")
            .forEach((p) => {
              const L = p.getTotalLength();
              p.style.strokeDasharray = String(L);
              p.style.strokeDashoffset = String(L);
              p.animate(
                [{ strokeDashoffset: String(L) }, { strokeDashoffset: "0" }],
                { duration: 1000, easing: "ease-in-out", fill: "forwards" },
              );
            });
          if (el.hasAttribute("data-gauge")) runGauge(el);
          obs.unobserve(el);
        });
      },
      { threshold: 0.12, rootMargin: "0px 0px -8% 0px" },
    );
    reveals.forEach((el) => obs.observe(el));

    return () => {
      cleanups.forEach((c) => c());
      obs.disconnect();
      clearTimeout(toastTimer);
    };
  }, []);

  return (
    <>
      <link
        href="https://fonts.googleapis.com/css2?family=Space+Grotesk:wght@400;500;600;700&family=DM+Sans:opsz,wght@9..40,400;9..40,500;9..40,600&family=Space+Mono:wght@400;700&display=swap"
        rel="stylesheet"
      />
      <style dangerouslySetInnerHTML={{ __html: DESIGN_KEYFRAMES }} />
      <Script
        src="https://cdnjs.cloudflare.com/ajax/libs/gsap/3.12.5/gsap.min.js"
        strategy="afterInteractive"
      />
      <div ref={ref} dangerouslySetInnerHTML={{ __html: DESIGN_HTML }} />
    </>
  );
}
