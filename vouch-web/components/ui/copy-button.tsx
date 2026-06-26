"use client";

import { useState } from "react";
import { cn } from "@/lib/utils";

interface CopyButtonProps {
  text: string;
  className?: string;
}

export function CopyButton({ text, className }: CopyButtonProps) {
  const [copied, setCopied] = useState(false);

  const copy = async () => {
    try {
      await navigator.clipboard.writeText(text);
      setCopied(true);
      setTimeout(() => setCopied(false), 2000);
    } catch {
      // Clipboard API not available (e.g., non-HTTPS context)
    }
  };

  return (
    <button
      type="button"
      onClick={copy}
      aria-label={copied ? "Copied!" : "Copy to clipboard"}
      className={cn(
        "inline-flex items-center gap-1.5 rounded px-2 py-1 text-xs font-medium transition-colors",
        copied
          ? "bg-green-100 text-green-700"
          : "bg-muted text-muted-foreground hover:text-foreground",
        className,
      )}
    >
      {copied ? "✓ Copied" : "Copy"}
    </button>
  );
}
