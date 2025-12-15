"use client";

import { useState } from "react";
import { motion } from "framer-motion";
import {
  determineMode,
  AIPathOption,
  MIN_AI_BLOCKING_LATENCY_MS,
} from "./mode-config";
import { PlaygroundControls } from "./playground-controls";
import { PlaygroundCard } from "./playground-card";

export function PlaygroundSection() {
  const [latency, setLatency] = useState(500);
  const [aiPath, setAiPath] = useState<AIPathOption>("conditional");

  const mode = determineMode({
    latencyBudgetMs: latency,
    aiOnCriticalPath: aiPath,
  });

  return (
    <section
      className="py-32 relative overflow-hidden bg-black"
      id="playground"
    >
      <div className="absolute top-1/2 left-0 w-[500px] h-[500px] bg-neon-blue/5 rounded-full blur-[120px] -translate-y-1/2 pointer-events-none" />

      <div className="container mx-auto px-4 mb-16 text-center relative z-10">
        <motion.h2
          initial={{ opacity: 0, y: 20 }}
          whileInView={{ opacity: 1, y: 0 }}
          viewport={{ once: true }}
          className="text-3xl md:text-5xl font-bold tracking-tight mb-4 bg-clip-text text-transparent bg-gradient-to-b from-white to-white/50"
        >
          Choose Your Trade-off
        </motion.h2>
        <motion.p
          initial={{ opacity: 0, y: 20 }}
          whileInView={{ opacity: 1, y: 0 }}
          viewport={{ once: true }}
          transition={{ delay: 0.1 }}
          className="text-lg text-muted-foreground max-w-2xl mx-auto"
        >
          Security isn&apos;t one-size-fits-all. Configure Argus to match your
          application&apos;s specific latency and risk profile.
        </motion.p>
      </div>

      <div className="container mx-auto px-4 max-w-6xl relative z-10">
        <div className="grid lg:grid-cols-2 gap-12 lg:gap-24 items-center">
          <motion.div
            initial={{ opacity: 0, x: -20 }}
            whileInView={{ opacity: 1, x: 0 }}
            viewport={{ once: true }}
            transition={{ delay: 0.2 }}
          >
            <PlaygroundControls
              latency={latency}
              setLatency={(val) => {
                setLatency(val);
                if (val < MIN_AI_BLOCKING_LATENCY_MS) {
                  setAiPath("never");
                }
              }}
              aiPath={aiPath}
              setAiPath={setAiPath}
            />
          </motion.div>

          <motion.div
            initial={{ opacity: 0, x: 20 }}
            whileInView={{ opacity: 1, x: 0 }}
            viewport={{ once: true }}
            transition={{ delay: 0.3 }}
          >
            <PlaygroundCard mode={mode} />
          </motion.div>
        </div>
      </div>
    </section>
  );
}
