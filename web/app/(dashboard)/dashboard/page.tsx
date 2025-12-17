"use client";

import { useThreats } from "@/hooks/use-threats";
import { StatsCards } from "@/components/dashboard/stats-cards";
import { ThreatTable } from "@/components/dashboard/threat-table";
import { ThreatChart } from "@/components/dashboard/threat-chart";
import { Loader2 } from "lucide-react";

export default function DashboardPage() {
  const { threats, loading } = useThreats();

  if (loading) {
    return (
      <div className="flex h-[50vh] w-full items-center justify-center">
        <Loader2 className="h-8 w-8 animate-spin text-neon-blue" />
      </div>
    );
  }

  return (
    <div className="space-y-6 animate-in fade-in duration-500">
      <div className="flex items-center justify-between">
        <h1 className="text-3xl font-bold tracking-tight text-white">
          Command Center
        </h1>
      </div>

      <StatsCards threats={threats} />

      <div className="grid grid-cols-1 lg:grid-cols-7 gap-6">
        <div className="lg:col-span-4">
          <ThreatChart threats={threats} />
        </div>
        <div className="lg:col-span-3">
          {/* TODO: add a live feed */}
          <div className="h-full min-h-[300px] rounded-xl border border-white/10 bg-black/40 backdrop-blur-md p-6 flex flex-col items-center justify-center text-center">
            <span className="text-neon-blue font-mono text-sm mb-2">
              ● LIVE FEED
            </span>
            <p className="text-muted-foreground text-sm">
              Waiting for incoming traffic...
            </p>
          </div>
        </div>
      </div>

      <ThreatTable threats={threats} />
    </div>
  );
}
