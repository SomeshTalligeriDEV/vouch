import { cn } from "@/lib/utils";

interface TextareaProps extends React.TextareaHTMLAttributes<HTMLTextAreaElement> {
  label?: string;
  error?: string;
  hint?: string;
  maxLength?: number;
  showCount?: boolean;
}

export function Textarea({
  label,
  error,
  hint,
  id,
  className,
  maxLength,
  showCount,
  value,
  ...props
}: TextareaProps) {
  const textareaId = id ?? label?.toLowerCase().replace(/\s+/g, "-");
  const charCount = typeof value === "string" ? value.length : 0;

  return (
    <div className="space-y-1.5">
      {label && (
        <div className="flex items-center justify-between">
          <label htmlFor={textareaId} className="block text-sm font-medium text-foreground">
            {label}
            {props.required && <span className="text-destructive ml-1" aria-hidden>*</span>}
          </label>
          {showCount && maxLength && (
            <span
              className={cn(
                "text-xs tabular-nums",
                charCount > maxLength * 0.9 ? "text-destructive" : "text-muted-foreground",
              )}
            >
              {charCount}/{maxLength}
            </span>
          )}
        </div>
      )}
      <textarea
        id={textareaId}
        maxLength={maxLength}
        value={value}
        aria-invalid={!!error}
        className={cn(
          "w-full rounded-md border bg-background px-3 py-2 text-sm resize-y min-h-[80px]",
          "placeholder:text-muted-foreground",
          "focus:outline-none focus:ring-2 focus:ring-ring",
          "disabled:opacity-50 disabled:cursor-not-allowed",
          error ? "border-destructive" : "border-border",
          className,
        )}
        {...props}
      />
      {error && <p role="alert" className="text-xs text-destructive">{error}</p>}
      {hint && !error && <p className="text-xs text-muted-foreground">{hint}</p>}
    </div>
  );
}
