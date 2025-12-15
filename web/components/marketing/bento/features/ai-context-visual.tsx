"use client";

import { motion } from "framer-motion";
import { ShieldCheck, User } from "lucide-react";

export function AiContextVisual() {
  return (
    <div className="w-full max-w-[320px] flex flex-col gap-3 text-xs font-mono">
      <motion.div
        initial={{ opacity: 0, x: -10 }}
        whileInView={{ opacity: 1, x: 0 }}
        transition={{ delay: 0.5 }}
        className="self-start bg-neon-orange/10 border border-neon-orange/20 rounded-lg rounded-tl-none p-3 max-w-[90%]"
      >
        <div className="flex font-bold items-center gap-2 mb-1 text-gray-400 text-[11px] uppercase">
          <User className="w-3 h-3" /> User Payload
        </div>
        <div className="text-gray-200">
          &quot;Here is how to delete a table:{" "}
          <span className="text-red-400">DROP TABLE users;</span>&quot;
        </div>
      </motion.div>

      <motion.div
        initial={{ opacity: 0 }}
        whileInView={{ opacity: 1 }}
        transition={{ delay: 1.5, duration: 0.5 }}
        className="self-center flex gap-1"
      >
        <span className="w-1 h-1 bg-neon-blue rounded-full animate-bounce" />
        <span className="w-1 h-1 bg-neon-blue rounded-full animate-bounce delay-75" />
        <span className="w-1 h-1 bg-neon-blue rounded-full animate-bounce delay-150" />
      </motion.div>

      <motion.div
        initial={{ opacity: 0, x: 10 }}
        whileInView={{ opacity: 1, x: 0 }}
        transition={{ delay: 2.2 }}
        className="self-end bg-neon-green/10 border border-neon-green/20 rounded-lg rounded-tr-none p-3 max-w-[90%]"
      >
        <div className="flex items-center gap-2 mb-1 text-neon-green font-bold text-[11px] uppercase font-bold">
          <ShieldCheck className="w-3 h-3" /> Argus Verdict
        </div>
        <div className="text-white">
          Context detected: <span className="text-neon-green">Educational SQL blog</span>.{" "}
          <br />
          Action: <span className="text-neon-green font-bold">ALLOW</span>.
        </div>
      </motion.div>
    </div>
  );
}
