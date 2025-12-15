"use client";

import { cn } from "@/lib/utils";
import { motion } from "framer-motion";
import { ReactNode } from "react";

interface BentoCardProps {
  title?: string;
  description?: string;
  children: ReactNode;
  className?: string;
  delay?: number;
  minimal?: boolean;
}

export function BentoCard({
  title,
  description,
  children,
  className,
  delay = 0,
  minimal = false,
}: BentoCardProps) {
  return (
    <motion.div
      initial={{ opacity: 0, y: 20 }}
      whileInView={{ opacity: 1, y: 0 }}
      viewport={{ once: true }}
      transition={{ duration: 0.5, delay }}
      className={cn(
        "group relative overflow-hidden rounded-xl border border-white/10 bg-black/50 backdrop-blur-md transition-all hover:border-white/20",
        className
      )}
    >
      <div className="absolute inset-0 bg-gradient-to-br from-white/5 to-transparent opacity-0 transition-opacity duration-500 group-hover:opacity-100 pointer-events-none" />

      <div
        className={cn(
          "relative h-full flex flex-col",
          minimal ? "p-4 justify-center" : "p-6"
        )}
      >
        {!minimal && (
          <div className="flex-1 w-full min-h-[140px] flex items-center justify-center mb-4 relative z-10">
            {children}
          </div>
        )}

        {minimal ? (
          children
        ) : (
          <div className="relative z-10">
            <h3 className="text-xl font-semibold text-white mb-2">{title}</h3>
            <p className="text-sm text-muted-foreground leading-relaxed">
              {description}
            </p>
          </div>
        )}
      </div>
    </motion.div>
  );
}
