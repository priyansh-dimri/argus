"use client";

import { motion } from "framer-motion";

export function SpeedVisual() {
  return (
    <div className="w-full flex flex-col gap-4 px-4">
      <div className="space-y-1">
        <div className="flex justify-between text-xs text-gray-400">
          <span>Argus Security Layer</span>
          <span className="text-neon-blue">133µs</span>
        </div>
        <div className="h-2 w-full bg-white/5 rounded-full overflow-hidden">
          <motion.div
            initial={{ width: 0 }}
            whileInView={{ width: "1.3%" }}
            transition={{ duration: 1, delay: 0.5 }}
            className="h-full bg-neon-blue shadow-[0_0_10px_var(--neon-blue)]"
          />
        </div>
      </div>

      <div className="space-y-1">
        <div className="flex justify-between text-xs text-gray-400">
          <span>Typical Database Query</span>
          <span className="text-red-400/70">~10ms</span>
        </div>
        <div className="h-2 w-full bg-white/5 rounded-full overflow-hidden">
          <motion.div
            initial={{ width: 0 }}
            whileInView={{ width: "100%" }}
            transition={{ duration: 1, delay: 0.7 }}
            className="h-full bg-red-500/40"
          />
        </div>
      </div>
    </div>
  );
}
