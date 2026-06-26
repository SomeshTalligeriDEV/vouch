"use client";

import { useEffect, useRef, useState } from "react";

interface ScoreUpdate {
  builder_id: string;
  total_score: number;
  tier: number;
}

interface UseScoreSSEOptions {
  onUpdate?: (update: ScoreUpdate) => void;
}

export function useScoreSSE(apiBase: string, opts?: UseScoreSSEOptions) {
  const [connected, setConnected] = useState(false);
  const [lastUpdate, setLastUpdate] = useState<ScoreUpdate | null>(null);
  const esRef = useRef<EventSource | null>(null);

  useEffect(() => {
    const es = new EventSource(`${apiBase}/scores/live`);
    esRef.current = es;

    es.onopen = () => setConnected(true);
    es.onerror = () => setConnected(false);

    es.addEventListener("score_update", (e) => {
      try {
        const data = JSON.parse(e.data) as ScoreUpdate;
        setLastUpdate(data);
        opts?.onUpdate?.(data);
      } catch {
        // malformed event — ignore
      }
    });

    return () => {
      es.close();
      setConnected(false);
    };
  }, [apiBase]); // eslint-disable-line react-hooks/exhaustive-deps

  return { connected, lastUpdate };
}
