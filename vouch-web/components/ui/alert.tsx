import { cn } from "@/lib/utils";

type AlertVariant = "info" | "success" | "warning" | "error";

interface AlertProps {
  variant?: AlertVariant;
  title?: string;
  children: React.ReactNode;
  className?: string;
}

const variantStyles: Record<AlertVariant, string> = {
  info: "bg-blue-50 border-blue-200 text-blue-900 dark:bg-blue-950 dark:border-blue-800 dark:text-blue-100",
  success: "bg-green-50 border-green-200 text-green-900 dark:bg-green-950 dark:border-green-800 dark:text-green-100",
  warning: "bg-yellow-50 border-yellow-200 text-yellow-900 dark:bg-yellow-950 dark:border-yellow-800 dark:text-yellow-100",
  error: "bg-red-50 border-red-200 text-red-900 dark:bg-red-950 dark:border-red-800 dark:text-red-100",
};

const icons: Record<AlertVariant, string> = {
  info: "ℹ️",
  success: "✅",
  warning: "⚠️",
  error: "❌",
};

export function Alert({ variant = "info", title, children, className }: AlertProps) {
  return (
    <div
      role="alert"
      className={cn(
        "rounded-lg border p-4 text-sm",
        variantStyles[variant],
        className,
      )}
    >
      <div className="flex gap-3">
        <span aria-hidden className="shrink-0 mt-0.5">{icons[variant]}</span>
        <div>
          {title && <p className="font-semibold mb-1">{title}</p>}
          <div>{children}</div>
        </div>
      </div>
    </div>
  );
}
