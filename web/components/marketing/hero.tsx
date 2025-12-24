"use client";

import { useRef, useState } from "react";
import Link from "next/link";
import { motion } from "framer-motion";
import {
  ArrowRight,
  Zap,
  Terminal,
  Shield,
  RotateCcw,
  MonitorPlay,
} from "lucide-react";

import { cn } from "@/lib/utils";
import { Button } from "@/components/ui/button";
import { TerminalWindow } from "./terminal-window";
import { useTerminalTimer } from "@/hooks/use-terminal-timer";
import { ATTACKER_LINES, ARGUS_LINES } from "@/config/terminal-data";

export function Hero() {
  const terminalRef = useRef<HTMLDivElement>(null);
  const { elapsedTime, restartKey, restart } = useTerminalTimer(terminalRef);
  const [activeView, setActiveView] = useState<"attacker" | "argus">(
    "attacker"
  );

  return (
    <section className="relative min-h-screen pt-32 pb-20 overflow-hidden bg-background">
      <div className="absolute inset-0 bg-grid-small-white opacity-20 pointer-events-none" />
      <div className="absolute top-0 left-1/2 -translate-x-1/2 w-[800px] h-[500px] bg-neon-blue/10 rounded-full blur-[120px] pointer-events-none" />
      <div className="absolute bottom-0 right-0 w-[600px] h-[600px] bg-neon-green/5 rounded-full blur-[100px] pointer-events-none" />

      <div className="container relative z-10 flex flex-col items-center text-center max-w-6xl px-4 mx-auto">
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.5 }}
          className="mb-8"
        >
          <div className="inline-flex items-center rounded-full border border-neon-blue/30 bg-neon-blue/5 px-3 py-1 text-sm text-neon-blue backdrop-blur-sm">
            <span className="flex h-2 w-2 rounded-full bg-neon-blue mr-2 animate-pulse" />
            v1.0 • Powered by Gemini 2.5
          </div>
        </motion.div>

        <motion.h1
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.5, delay: 0.1 }}
          className="text-5xl md:text-7xl font-bold tracking-tight mb-6 bg-clip-text text-transparent bg-gradient-to-b from-white to-white/50"
        >
          The Cognitive Shield <br /> for High-Scale APIs.
        </motion.h1>

        <motion.p
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.5, delay: 0.2 }}
          className="text-lg md:text-xl text-muted-foreground mb-10 max-w-2xl mx-auto"
        >
          Hybrid WAF. Sub-millisecond blocking. AI-verified edge cases. Three
          risk profiles for{" "}
          <a href="#playground" className="text-neon-blue hover:underline">
            different routes
          </a>
          .
        </motion.p>

        <motion.div
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.5, delay: 0.3 }}
          className="flex flex-col sm:flex-row gap-4 mb-24"
        >
          <Link href="/signup">
            <Button
              size="lg"
              className="rounded-full h-12 px-8 text-base font-semibold bg-white text-black hover:bg-white/90 shadow-[0_0_20px_rgba(255,255,255,0.3)] transition-all"
            >
              Start Protecting <Zap className="ml-2 h-4 w-4" />
            </Button>
          </Link>
          <Link href="/architecture">
            <Button
              variant="outline"
              size="lg"
              className="rounded-full h-12 px-8 text-base border-white/10 hover:bg-white/5 hover:text-white transition-colors"
            >
              View Architecture <ArrowRight className="ml-2 h-4 w-4" />
            </Button>
          </Link>
        </motion.div>

        {/* TERMINAL SECTION */}
        <motion.div
          ref={terminalRef}
          initial={{ opacity: 0, scale: 0.95 }}
          whileInView={{ opacity: 1, scale: 1 }}
          viewport={{ once: true, margin: "-100px" }}
          transition={{ duration: 0.7, delay: 0.2 }}
          className="w-full relative"
        >
          <div className="absolute -inset-1 bg-gradient-to-r from-neon-blue/20 to-neon-green/20 rounded-xl blur opacity-30" />

          <div className="relative rounded-xl border border-white/10 bg-black/80 backdrop-blur-xl shadow-2xl overflow-hidden">
            <div className="flex items-center justify-between px-4 py-3 border-b border-white/10 bg-white/5 flex-wrap gap-2">
              <div className="flex items-center gap-4">
                <div className="flex gap-2">
                  <div className="w-3 h-3 rounded-full bg-red-500/80" />
                  <div className="w-3 h-3 rounded-full bg-yellow-500/80" />
                  <div className="w-3 h-3 rounded-full bg-green-500/80" />
                </div>
                <div className="text-xs font-mono text-muted-foreground hidden sm:block">
                  argus-defense-console
                </div>
              </div>

              <div className="flex items-center gap-3">
                <button
                  onClick={() =>
                    setActiveView((v) =>
                      v === "attacker" ? "argus" : "attacker"
                    )
                  }
                  className="flex md:hidden items-center gap-1.5 px-2 py-1 rounded bg-white/5 hover:bg-white/10 text-xs font-medium text-muted-foreground hover:text-white transition-colors border border-white/10"
                >
                  <MonitorPlay className="w-3 h-3" />
                  <span>
                    {activeView === "attacker" ? "Show Argus" : "Show Attacker"}
                  </span>
                </button>

                <button
                  onClick={restart}
                  className="flex items-center gap-1.5 px-2 py-1 rounded hover:bg-white/10 text-xs font-medium text-muted-foreground hover:text-white transition-colors"
                >
                  <RotateCcw className="w-3 h-3" />
                  <span>Restart</span>
                </button>
              </div>
            </div>

            <div
              key={restartKey}
              className="grid md:grid-cols-2 divide-y md:divide-y-0 md:divide-x divide-white/10 h-[500px] text-left"
            >
              {/* ATTACKED PANEL */}
              <div
                className={cn(
                  "h-full overflow-hidden transition-all duration-300",
                  activeView === "attacker" ? "block" : "hidden md:block"
                )}
              >
                <TerminalWindow
                  title="attacker@kali:~$"
                  icon={Terminal}
                  color="neon-orange"
                  time={elapsedTime}
                  lines={ATTACKER_LINES}
                  className="border-r border-border/50"
                />
              </div>

              {/* ARGUS PANEL */}
              <div
                className={cn(
                  "h-full overflow-hidden bg-white/5 transition-all duration-300",
                  activeView === "argus" ? "block" : "hidden md:block"
                )}
              >
                <TerminalWindow
                  title="argus@proxy:~$"
                  icon={Shield}
                  color="neon-green"
                  time={elapsedTime}
                  lines={ARGUS_LINES}
                />
              </div>
            </div>
          </div>
        </motion.div>
      </div>
    </section>
  );
}
