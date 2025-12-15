"use client";

import { motion } from "framer-motion";
import { AlertTriangle, CheckCircle2 } from "lucide-react";
import { useState, useEffect } from "react";
import AutoToggle from "./auto-toggle";

export function CircuitVisual() {
  const [isOpen, setIsOpen] = useState(false);

  useEffect(() => {
    const interval = setInterval(() => {
      setIsOpen((prev) => !prev);
    }, 3000);
    return () => clearInterval(interval);
  }, []);

  return (
    <div className="flex items-center gap-3">
      <div className="flex flex-col items-center gap-2">
        <div className="text-[10px] text-gray-500 font-mono uppercase">
          AI Status
        </div>
        <motion.div
          animate={{ color: isOpen ? "#ef4444" : "#22c55e" }}
          className="flex items-center gap-1 font-bold text-sm font-mono "
        >
          {isOpen ? (
            <AlertTriangle className="w-4 h-4" />
          ) : (
            <CheckCircle2 className="w-4 h-4" />
          )}
          {isOpen ? "OUTAGE" : "ONLINE"}
        </motion.div>
      </div>

      <AutoToggle />

      <div className="flex flex-col items-center gap-2">
        <div className="text-[10px] text-gray-500 font-mono uppercase">
          Traffic
        </div>
        <motion.div
          animate={{ opacity: [0.5, 1, 0.5] }}
          transition={{ repeat: Infinity, duration: 2 }}
          className="text-white text-xs font-mono bg-white/10 px-2 py-1 rounded"
        >
          ALLOWED
        </motion.div>
      </div>
    </div>
  );
}
