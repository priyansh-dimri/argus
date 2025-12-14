import { Zap, Brain, Lock, LucideIcon } from "lucide-react";

export type SecurityMode = "LATENCY_FIRST" | "SMART_SHIELD" | "PARANOID";
export type AIPathOption = "never" | "conditional" | "always";
export const MIN_AI_BLOCKING_LATENCY_MS = 100;

export interface ModeConfig {
  id: SecurityMode;
  name: string;
  icon: LucideIcon;
  color: "neon-blue" | "neon-orange" | "neon-green";
  description: string;
  useCase: string;
}

export const MODES: Record<SecurityMode, ModeConfig> = {
  LATENCY_FIRST: {
    id: "LATENCY_FIRST",
    name: "Latency First",
    icon: Zap,
    color: "neon-blue",
    description: "Blocks Regex hits instantly. AI runs async in background.",
    useCase: "Best for Login APIs & High-Frequency Trading",
  },
  SMART_SHIELD: {
    id: "SMART_SHIELD",
    name: "Smart Shield",
    icon: Brain,
    color: "neon-green",
    description:
      "Verifies threats with AI. Fixes false positives automatically.",
    useCase: "Best for CMS, Blogs, & Rich Text Inputs",
  },
  PARANOID: {
    id: "PARANOID",
    name: "Paranoid",
    icon: Lock,
    color: "neon-orange",
    description: "Trusts nothing. Checks every single request with AI.",
    useCase: "Best for Admin Panels & Zero-Trust Zones",
  },
};

export function determineMode({
  latencyBudgetMs,
  aiOnCriticalPath,
}: {
  latencyBudgetMs: number;
  aiOnCriticalPath: AIPathOption;
}): SecurityMode {
  if (latencyBudgetMs < MIN_AI_BLOCKING_LATENCY_MS) {
    return "LATENCY_FIRST";
  }
  if (aiOnCriticalPath === "always") return "PARANOID";
  if (aiOnCriticalPath === "conditional") return "SMART_SHIELD";
  return "LATENCY_FIRST";
}
