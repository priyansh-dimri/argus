"use client";

import { Check, Copy } from "lucide-react";
import { useState } from "react";

interface CodeBlockProps {
  code: string;
  language: string;
}

export function CodeBlock({ code, language }: CodeBlockProps) {
  const [copied, setCopied] = useState(false);

  const copyToClipboard = () => {
    navigator.clipboard.writeText(code);
    setCopied(true);
    setTimeout(() => setCopied(false), 2000);
  };

  const lines = code.trim().split("\n");

  return (
    <div className="relative group rounded-lg bg-black/50 border border-white/10 backdrop-blur-sm overflow-hidden font-mono text-sm">
      <div className="flex items-center justify-between px-4 py-2 border-b border-white/5 bg-white/5">
        <span className="text-xs text-muted-foreground uppercase">
          {language}
        </span>
        <button
          onClick={copyToClipboard}
          className="text-muted-foreground hover:text-white transition-colors"
          title="Copy code"
        >
          {copied ? (
            <Check className="w-4 h-4 text-neon-green" />
          ) : (
            <Copy className="w-4 h-4" />
          )}
        </button>
      </div>

      <div className="p-4 overflow-x-auto">
        <pre>
          <code className="grid gap-1">
            {lines.map((line, i) => (
              <div key={i} className="table-row">
                <span className="table-cell text-right w-8 pr-4 text-muted-foreground/30 select-none">
                  {i + 1}
                </span>
                <span className="table-cell whitespace-pre text-gray-300">
                  {highlightSyntax(line, language)}
                </span>
              </div>
            ))}
          </code>
        </pre>
      </div>
    </div>
  );
}

// Simple syntax highlighter helper for visual flair without heavy libs
function highlightSyntax(line: string, lang: string): React.ReactNode {
  // Very basic regex-based highlighting for Go/JSON/Bash
  // Real highlighting requires a parser, but this adds the "vibe"

  if (lang === "go") {
    if (
      line.includes("import") ||
      line.includes("package") ||
      line.includes("func") ||
      line.includes("return")
    ) {
      return <span className="text-neon-blue">{line}</span>;
    }
    if (line.includes("//")) {
      return <span className="text-muted-foreground">{line}</span>;
    }
    return line;
  }

  if (lang === "json") {
    if (line.includes(":")) {
      const parts = line.split(":");
      return (
        <>
          <span className="text-neon-blue">{parts[0]}:</span>
          <span className="text-neon-green">{parts.slice(1).join(":")}</span>
        </>
      );
    }
  }

  if (lang === "bash") {
    if (line.startsWith("$")) {
      return <span className="text-neon-green">{line}</span>;
    }
  }

  return line;
}
