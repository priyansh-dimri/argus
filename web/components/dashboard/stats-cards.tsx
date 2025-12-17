"use client";

import { ShieldAlert, Clock, BarChart3 } from "lucide-react";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { ThreatLog } from "@/hooks/use-threats";
import { cn } from "@/lib/utils";

interface StatsCardsProps {
  threats: ThreatLog[];
}

export function StatsCards({ threats }: StatsCardsProps) {
  const totalBlocked = threats.filter((t) => t.is_threat).length;
  const totalTraffic = threats.length;

  const blockRate =
    totalTraffic > 0 ? ((totalBlocked / totalTraffic) * 100).toFixed(1) : "0";

  const stats = [
    {
      title: "Threats Blocked",
      value: totalBlocked,
      icon: ShieldAlert,
      color: "text-red-500",
      trend: "Total detections",
    },
    {
      title: "Block Rate",
      value: `${blockRate}%`,
      icon: BarChart3,
      color: "text-orange-500",
      trend: "Of analyzed traffic",
    },
    {
      title: "Monitor Mode",
      value: "Active",
      icon: Clock,
      color: "text-neon-blue",
      trend: "Real-time analysis",
    },
  ];

  return (
    <div className="grid gap-4 md:grid-cols-3">
      {stats.map((stat) => (
        <Card
          key={stat.title}
          className="bg-black/40 border-white/10 backdrop-blur-md"
        >
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium text-muted-foreground">
              {stat.title}
            </CardTitle>
            <stat.icon className={cn("h-4 w-4", stat.color)} />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-white">{stat.value}</div>
            <p className="text-xs text-muted-foreground mt-1">{stat.trend}</p>
          </CardContent>
        </Card>
      ))}
    </div>
  );
}
