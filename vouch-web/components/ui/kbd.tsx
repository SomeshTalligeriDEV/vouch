import { cn } from "@/lib/utils";

interface KbdProps {
  children: React.ReactNode;
  className?: string;
}

export function Kbd({ children, className }: KbdProps) {
  return (
    <kbd
      className={cn(
        "inline-flex items-center justify-center rounded border border-border bg-muted px-1.5 py-0.5",
        "font-mono text-xs text-muted-foreground shadow-sm",
        className,
      )}
    >
      {children}
    </kbd>
  );
}
