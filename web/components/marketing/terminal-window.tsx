"use client";

import { cn } from "@/lib/utils";
import { useMemo } from "react";

export interface TerminalLine {
  text: string;
  delay: number;
  color?: "neon-orange" | "neon-green";
  className?: string;
}

interface TerminalWindowProps {
  title: string;
  icon: React.ElementType;
  color: "neon-blue" | "neon-orange" | "neon-green";
  className?: string;
  time: number;
  lines: {
    text: string;
    delay: number;
    color?: "neon-orange" | "neon-green";
    className?: string;
  }[];
}

export function TerminalWindow({
  title,
  icon: Icon,
  color,
  className,
  lines,
  time,
}: TerminalWindowProps) {
  const absoluteLines = useMemo(() => {
    let currentTime = 0;
    const result = [];
    for (const line of lines) {
      currentTime += line.delay;
      result.push({ ...line, showAt: currentTime });
    }
    return result;
  }, [lines]);

  const visibleLines = absoluteLines.filter((line) => line.showAt <= time);

  const totalDuration = absoluteLines[absoluteLines.length - 1]?.showAt || 0;
  const isTyping = time < totalDuration + 500;

  return (
    <div
      className={cn(
        "pt-8 p-4 font-mono text-sm h-full flex flex-col",
        className
      )}
    >
      <div className="flex items-center gap-2 mb-2 text-gray-400 select-none">
        <Icon
          className={cn(
            "h-4 w-4",
            color === "neon-orange" && "text-neon-orange",
            color === "neon-green" && "text-neon-green"
          )}
        />
        <span className="font-semibold">{title}</span>
      </div>

      <div className="flex-1 overflow-y-auto pt-2 scrollbar-hide font-mono">
        {visibleLines.map((line, index) => (
          <div
            key={index}
            className={cn(
              "py-0.5 break-all animate-in fade-in duration-300 slide-in-from-left-2",
              line.color === "neon-orange" &&
                "text-neon-orange drop-shadow-neon-orange",
              line.color === "neon-green" &&
                "text-neon-green drop-shadow-neon-green",
              !line.color && !line.className && "text-gray-300",
              line.className
            )}
          >
            <span className="opacity-30 mr-2 select-none">$</span>
            {line.text}
          </div>
        ))}

        {/* typing cursor simulation */}
        {isTyping && (
          <div className="py-0.5 text-gray-300">
            <span className="opacity-30 mr-2 select-none">$</span>
            <span className="animate-pulse bg-neon-blue/80 inline-block h-4 w-2 align-middle" />
          </div>
        )}
      </div>

      <style jsx global>{`
        .drop-shadow-neon-orange {
          filter: drop-shadow(0 0 4px var(--neon-orange));
        }
        .drop-shadow-neon-green {
          filter: drop-shadow(0 0 4px var(--neon-green));
        }
        .scrollbar-hide::-webkit-scrollbar {
          display: none;
        }
        .scrollbar-hide {
          -ms-overflow-style: none;
          scrollbar-width: none;
        }
      `}</style>
    </div>
  );
}