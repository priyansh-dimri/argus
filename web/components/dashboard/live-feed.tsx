"use client";

import { ThreatLog } from "@/hooks/use-threats";
import { cn } from "@/lib/utils";
import { Activity } from "lucide-react";

interface LiveFeedProps {
  threats: ThreatLog[];
  isConnected: boolean;
}

export function LiveFeed({ threats, isConnected = true }: LiveFeedProps) {
  return (
    <div className="h-full min-h-[300px] rounded-xl border border-white/10 bg-black/40 backdrop-blur-md flex flex-col overflow-hidden">
      <div className="flex items-center justify-between p-4 border-b border-white/5 bg-white/5">
        <div className="flex items-center gap-2">
          <div className="relative flex h-2 w-2">
            {isConnected && (
              <span className="animate-ping absolute inline-flex h-full w-full rounded-full bg-neon-blue opacity-75"></span>
            )}
            <span
              className={cn(
                "relative inline-flex rounded-full h-2 w-2",
                isConnected ? "bg-neon-blue" : "bg-gray-500"
              )}
            ></span>
          </div>
          <span className="text-xs font-mono text-neon-blue font-bold tracking-wider">
            LIVE FEED
          </span>
        </div>
        <Activity className="h-4 w-4 text-muted-foreground/50" />
      </div>

      <div className="flex-1 p-4 overflow-y-auto space-y-3">
        {threats.length === 0 ? (
          <div className="h-full flex items-center justify-center text-xs text-muted-foreground">
            Waiting for incoming traffic...
          </div>
        ) : (
          threats.slice(0, 5).map((t) => (
            <div
              key={t.id}
              className="group flex flex-col gap-1 border-l-2 border-white/10 pl-3 py-1 hover:border-neon-blue/50 transition-colors"
            >
              <div className="flex items-center justify-between text-[10px] font-mono text-gray-500">
                <span>{new Date(t.timestamp).toLocaleTimeString()}</span>
                <span
                  className={cn(
                    "px-1.5 py-0.5 rounded-[2px]",
                    t.is_threat
                      ? "bg-red-500/20 text-red-400"
                      : "bg-green-500/20 text-green-400"
                  )}
                >
                  {t.is_threat ? "BLOCK" : "ALLOW"}
                </span>
              </div>
              <div
                className="text-xs text-gray-300 font-mono truncate"
                title={t.reason}
              >
                {t.method} {t.route}
              </div>
              <div className="text-[10px] text-gray-600 font-mono">{t.ip}</div>
            </div>
          ))
        )}
      </div>
    </div>
  );
}
