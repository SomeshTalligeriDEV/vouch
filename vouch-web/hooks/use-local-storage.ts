"use client";

import { useCallback, useEffect, useState } from "react";

export function useLocalStorage<T>(key: string, initialValue: T) {
  const readValue = useCallback((): T => {
    if (typeof window === "undefined") return initialValue;
    try {
      const item = window.localStorage.getItem(key);
      return item ? (JSON.parse(item) as T) : initialValue;
    } catch {
      return initialValue;
    }
  }, [key, initialValue]);

  const [storedValue, setStoredValue] = useState<T>(readValue);

  const setValue = useCallback(
    (value: T | ((prev: T) => T)) => {
      const next = value instanceof Function ? value(storedValue) : value;
      try {
        window.localStorage.setItem(key, JSON.stringify(next));
      } catch {
        // quota exceeded or private mode — silently ignore
      }
      setStoredValue(next);
    },
    [key, storedValue],
  );

  const removeValue = useCallback(() => {
    try {
      window.localStorage.removeItem(key);
    } catch {
      // ignore
    }
    setStoredValue(initialValue);
  }, [key, initialValue]);

  useEffect(() => {
    const handler = (e: StorageEvent) => {
      if (e.key === key) setStoredValue(readValue());
    };
    window.addEventListener("storage", handler);
    return () => window.removeEventListener("storage", handler);
  }, [key, readValue]);

  return [storedValue, setValue, removeValue] as const;
}
