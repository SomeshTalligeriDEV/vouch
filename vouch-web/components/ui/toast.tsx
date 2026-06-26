"use client";

import { useEffect, useState } from "react";
import { cn } from "@/lib/utils";

export type ToastVariant = "default" | "success" | "error";

interface Toast {
  id: string;
  message: string;
  variant: ToastVariant;
}

let toastQueue: ((t: Omit<Toast, "id">) => void) | null = null;

export function toast(message: string, variant: ToastVariant = "default") {
  toastQueue?.({ message, variant });
}

export function Toaster() {
  const [toasts, setToasts] = useState<Toast[]>([]);

  useEffect(() => {
    toastQueue = (t) => {
      const id = crypto.randomUUID();
      setToasts((prev) => [...prev, { ...t, id }]);
      setTimeout(() => {
        setToasts((prev) => prev.filter((x) => x.id !== id));
      }, 4000);
    };
    return () => {
      toastQueue = null;
    };
  }, []);

  if (!toasts.length) return null;

  return (
    <div className="fixed bottom-4 right-4 z-50 flex flex-col gap-2">
      {toasts.map((t) => (
        <div
          key={t.id}
          role="alert"
          className={cn(
            "rounded-lg px-4 py-3 text-sm font-medium shadow-md min-w-[240px] max-w-sm transition-all",
            t.variant === "success" && "bg-green-600 text-white",
            t.variant === "error" && "bg-destructive text-destructive-foreground",
            t.variant === "default" && "bg-foreground text-background",
          )}
        >
          {t.message}
        </div>
      ))}
    </div>
  );
}
