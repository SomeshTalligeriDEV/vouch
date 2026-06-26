"use client";

import { useRef, useState } from "react";

import { uploadFile } from "@/lib/api";

export function ImageUpload({
  value,
  onChange,
  label = "Logo",
}: {
  value?: string;
  onChange: (url: string) => void;
  label?: string;
}) {
  const inputRef = useRef<HTMLInputElement>(null);
  const [status, setStatus] = useState<"idle" | "uploading" | "error">("idle");

  const onPick = async (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (!file) return;
    setStatus("uploading");
    try {
      const url = await uploadFile(file);
      onChange(url);
      setStatus("idle");
    } catch {
      setStatus("error");
    }
  };

  return (
    <div>
      <span className="mb-1 block text-sm text-ink/60">{label}</span>
      <div className="flex items-center gap-3">
        <button
          type="button"
          onClick={() => inputRef.current?.click()}
          className="flex h-16 w-16 items-center justify-center overflow-hidden rounded-lg border border-line bg-panel text-xs text-ink/50 hover:border-crimson"
        >
          {value ? (
            // eslint-disable-next-line @next/next/no-img-element
            <img src={value} alt={label} className="h-full w-full object-cover" />
          ) : status === "uploading" ? (
            "…"
          ) : (
            "Upload"
          )}
        </button>
        <div className="text-sm">
          {status === "uploading" && (
            <span className="text-ink/60">Uploading…</span>
          )}
          {status === "error" && (
            <span className="text-red-400">Upload failed — sign in and retry.</span>
          )}
          {status === "idle" && value && (
            <button
              type="button"
              onClick={() => onChange("")}
              className="text-ink/60 hover:text-red-400"
            >
              Remove
            </button>
          )}
        </div>
      </div>
      <input
        ref={inputRef}
        type="file"
        accept="image/png,image/jpeg,image/webp,image/gif,image/svg+xml"
        className="hidden"
        onChange={onPick}
      />
    </div>
  );
}
