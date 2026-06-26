"use client";

import { useEffect, useRef } from "react";
import { cn } from "@/lib/utils";

interface DialogProps {
  open: boolean;
  onClose: () => void;
  title?: string;
  children: React.ReactNode;
  className?: string;
}

export function Dialog({ open, onClose, title, children, className }: DialogProps) {
  const dialogRef = useRef<HTMLDialogElement>(null);

  useEffect(() => {
    const el = dialogRef.current;
    if (!el) return;
    if (open) {
      el.showModal();
    } else {
      el.close();
    }
  }, [open]);

  useEffect(() => {
    const el = dialogRef.current;
    if (!el) return;
    const handler = () => onClose();
    el.addEventListener("close", handler);
    return () => el.removeEventListener("close", handler);
  }, [onClose]);

  if (!open) return null;

  return (
    <dialog
      ref={dialogRef}
      onClick={(e) => {
        if (e.target === e.currentTarget) onClose();
      }}
      className="backdrop:bg-black/50 rounded-xl border bg-background p-0 shadow-xl max-w-lg w-full"
    >
      <div className={cn("p-6", className)}>
        {title && (
          <div className="flex items-center justify-between mb-4">
            <h2 className="text-lg font-semibold">{title}</h2>
            <button
              type="button"
              onClick={onClose}
              aria-label="Close dialog"
              className="rounded-md p-1 hover:bg-muted transition-colors text-muted-foreground"
            >
              ✕
            </button>
          </div>
        )}
        {children}
      </div>
    </dialog>
  );
}
