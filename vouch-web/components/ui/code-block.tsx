import { cn } from "@/lib/utils";

interface CodeBlockProps {
  children: string;
  language?: string;
  className?: string;
}

export function CodeBlock({ children, language, className }: CodeBlockProps) {
  return (
    <pre
      className={cn(
        "rounded-lg bg-muted p-4 text-sm overflow-x-auto font-mono",
        className,
      )}
    >
      {language && (
        <div className="text-xs text-muted-foreground mb-2 select-none">{language}</div>
      )}
      <code>{children}</code>
    </pre>
  );
}

export function InlineCode({ children }: { children: React.ReactNode }) {
  return (
    <code className="rounded bg-muted px-1.5 py-0.5 text-sm font-mono">
      {children}
    </code>
  );
}
