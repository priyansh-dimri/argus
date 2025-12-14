"use client";

import { Label } from "@/components/ui/label";
import { Slider } from "@/components/ui/slider";
import { cn } from "@/lib/utils";
import { Zap, Brain, Lock, type LucideIcon } from "lucide-react";
import { AIPathOption, MIN_AI_BLOCKING_LATENCY_MS } from "./mode-config";

interface PlaygroundControlsProps {
  latency: number;
  setLatency: (val: number) => void;
  aiPath: AIPathOption;
  setAiPath: (val: AIPathOption) => void;
}

export function PlaygroundControls({
  latency,
  setLatency,
  aiPath,
  setAiPath,
}: PlaygroundControlsProps) {
  const aiBlockingAllowed = latency >= MIN_AI_BLOCKING_LATENCY_MS;
  const options: {
    value: AIPathOption;
    label: string;
    icon: LucideIcon;
    desc: string;
  }[] = [
    {
      value: "never",
      label: "Async Only",
      icon: Zap,
      desc: "AI never blocks the user. Purely for logging.",
    },
    {
      value: "conditional",
      label: "Conditional",
      icon: Brain,
      desc: "AI blocks only when Regex flags a threat.",
    },
    {
      value: "always",
      label: "Always",
      icon: Lock,
      desc: "AI verifies every single request.",
    },
  ];

  return (
    <div className="space-y-8 p-6 md:p-8 bg-white/5 rounded-xl border border-white/10 backdrop-blur-sm">
      <div className="space-y-4">
        <div className="flex items-center justify-between">
          <Label className="text-base font-medium text-gray-200">
            Latency Budget
          </Label>
          <span
            className={cn(
              "font-mono text-sm px-2 py-0.5 rounded",
              latency < 100
                ? "bg-neon-blue/10 text-neon-blue"
                : "bg-white/10 text-gray-400"
            )}
          >
            {latency}ms
          </span>
        </div>
        <Slider
          defaultValue={[500]}
          max={1000}
          step={10}
          value={[latency]}
          onValueChange={(vals) => setLatency(vals[0])}
          className="cursor-pointer"
        />
        <p className="text-xs text-muted-foreground">
          Simulated network overhead allowed per request.
        </p>
      </div>

      <div className="space-y-4 pt-4 border-t border-white/10">
        <Label className="text-base font-medium text-gray-200">
          AI on Critical Path?
        </Label>

        <div className="grid gap-3">
          {options.map((opt) => {
            const Icon = opt.icon;
            const isSelected = aiPath === opt.value;
            const requiresBlocking = opt.value !== "never";
            const disabled = requiresBlocking && !aiBlockingAllowed;

            return (
              <div
                key={opt.value}
                onClick={() => {
                  if (disabled) return;
                  setAiPath(opt.value);
                }}
                className={cn(
                  "relative flex items-start gap-3 p-3 rounded-lg border cursor-pointer transition-all duration-200",
                  isSelected
                    ? "bg-white/10 border-neon-blue/50 ring-1 ring-neon-blue/50"
                    : "bg-black/20 border-white/5 hover:border-white/20 hover:bg-white/5"
                )}
              >
                <div
                  className={cn(
                    "mt-0.5 p-1.5 rounded-md",
                    isSelected
                      ? "bg-neon-blue/20 text-neon-blue"
                      : "bg-white/5 text-gray-400"
                  )}
                >
                  <Icon className="w-4 h-4" />
                </div>
                <div className="space-y-1">
                  <p
                    className={cn(
                      "text-sm font-medium",
                      isSelected ? "text-white" : "text-gray-300"
                    )}
                  >
                    {opt.label}
                  </p>
                  <p className="text-xs text-muted-foreground leading-relaxed">
                    {opt.desc}
                  </p>
                </div>
              </div>
            );
          })}
        </div>
      </div>
    </div>
  );
}
