"use client";

import { motion, AnimatePresence } from "framer-motion";
import { MODES, SecurityMode } from "./mode-config";
import { cn } from "@/lib/utils";

interface PlaygroundCardProps {
  mode: SecurityMode;
}

export function PlaygroundCard({ mode }: PlaygroundCardProps) {
  const config = MODES[mode];
  const Icon = config.icon;

  return (
    <div className="relative h-full min-h-[300px] flex items-center justify-center p-6">
      <motion.div
        animate={{
          backgroundColor:
            mode === "LATENCY_FIRST"
              ? "var(--neon-blue)"
              : mode === "SMART_SHIELD"
              ? "var(--neon-green)"
              : "var(--neon-orange)",
        }}
        className="absolute inset-0 opacity-10 blur-[80px] transition-colors duration-700"
      />

      <AnimatePresence mode="wait">
        <motion.div
          key={mode}
          initial={{ opacity: 0, scale: 0.9, y: 20 }}
          animate={{ opacity: 1, scale: 1, y: 0 }}
          exit={{ opacity: 0, scale: 0.95, y: -20 }}
          transition={{ duration: 0.4, type: "spring", stiffness: 100 }}
          className="relative w-full max-w-sm"
        >
          <div className="bg-black/80 border border-white/10 rounded-2xl p-8 shadow-2xl backdrop-blur-xl relative overflow-hidden group hover:border-white/20 transition-colors">
            <div
              className={cn(
                "w-16 h-16 rounded-2xl flex items-center justify-center mb-6 mx-auto transition-colors duration-500",
                mode === "LATENCY_FIRST" && "bg-neon-blue/20 text-neon-blue",
                mode === "SMART_SHIELD" && "bg-neon-green/20 text-neon-green",
                mode === "PARANOID" && "bg-neon-orange/20 text-neon-orange"
              )}
            >
              <Icon className="w-8 h-8" />
            </div>

            <div className="text-center space-y-4">
              <h3 className="text-2xl font-bold text-white tracking-tight">
                {config.name}
              </h3>

              <div className="h-px w-12 bg-white/10 mx-auto" />

              <p className="text-gray-400 leading-relaxed min-h-[3rem]">
                {config.description}
              </p>

              <div
                className={cn(
                  "inline-block px-3 py-1 rounded-full text-xs font-medium bg-white/5 border border-white/10",
                  mode === "LATENCY_FIRST" &&
                    "text-neon-blue border-neon-blue/20",
                  mode === "SMART_SHIELD" &&
                    "text-neon-green border-neon-green/20",
                  mode === "PARANOID" &&
                    "text-neon-orange border-neon-orange/20"
                )}
              >
                {config.useCase}
              </div>
            </div>

            <div
              className={cn(
                "absolute top-0 right-0 w-20 h-20 bg-gradient-to-bl from-white/5 to-transparent opacity-0 group-hover:opacity-100 transition-opacity duration-500",
                mode === "LATENCY_FIRST" && "from-neon-blue/10",
                mode === "SMART_SHIELD" && "from-neon-green/10",
                mode === "PARANOID" && "from-neon-orange/10"
              )}
            />
          </div>
        </motion.div>
      </AnimatePresence>
    </div>
  );
}
