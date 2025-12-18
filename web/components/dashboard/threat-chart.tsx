"use client";

import { useMemo, useState } from "react";
import {
  AreaChart,
  Area,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  ResponsiveContainer,
} from "recharts";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { format, subHours, isAfter } from "date-fns";

interface ThreatLog {
  timestamp: string;
  is_threat: boolean;
}

interface ThreatChartProps {
  threats: ThreatLog[];
}

export function ThreatChart({ threats }: ThreatChartProps) {
  const [view, setView] = useState<"all" | "threats">("all");

  const chartData = useMemo(() => {
    const now = new Date();
    const sixHoursAgo = subHours(now, 6);
    const bins: Record<
      string,
      { time: string; timestamp: number; threats: number; allowed: number }
    > = {};

    for (let i = 0; i < 72; i++) {
      const binTime = new Date(now.getTime() - i * 5 * 60 * 1000);
      const label = format(binTime, "HH:mm");
      const key = Math.floor(binTime.getTime() / 300000); // equal to 5 minutes
      bins[key] = {
        time: label,
        timestamp: binTime.getTime(),
        threats: 0,
        allowed: 0,
      };
    }

    threats.forEach((threat) => {
      const threatTime = new Date(threat.timestamp);
      if (isAfter(threatTime, sixHoursAgo)) {
        const key = Math.floor(threatTime.getTime() / 300000);
        if (bins[key]) {
          if (threat.is_threat) {
            bins[key].threats += 1;
          } else {
            bins[key].allowed += 1;
          }
        }
      }
    });

    return Object.values(bins).sort((a, b) => a.timestamp - b.timestamp);
  }, [threats]);

  return (
    <Card className="bg-zinc-950/50 border-zinc-800/50 backdrop-blur-xl ring-1 ring-white/5 shadow-2xl overflow-hidden">
      <CardHeader className="flex flex-col sm:flex-row items-start sm:items-center justify-between space-y-4 sm:space-y-0 pb-6">
        <div className="space-y-1">
          <CardTitle className="text-lg font-semibold tracking-tight text-zinc-100">
            Traffic Intelligence
          </CardTitle>
          <p className="text-xs text-zinc-500 font-medium">
            Real-time threat landscape (6h Window)
          </p>
        </div>

        <div className="flex items-center p-1 bg-zinc-900 rounded-lg border border-zinc-800">
          <button
            onClick={() => setView("all")}
            className={`px-3 py-1.5 text-[10px] font-bold uppercase tracking-wider transition-all duration-200 rounded-md ${
              view === "all"
                ? "bg-zinc-800 text-neon-blue shadow-lg"
                : "text-zinc-500 hover:text-zinc-300"
            }`}
          >
            Full Stack
          </button>
          <button
            onClick={() => setView("threats")}
            className={`px-3 py-1.5 text-[10px] font-bold uppercase tracking-wider transition-all duration-200 rounded-md ${
              view === "threats"
                ? "bg-red-500/10 text-red-500 shadow-lg"
                : "text-zinc-500 hover:text-zinc-300"
            }`}
          >
            Threats Only
          </button>
        </div>
      </CardHeader>

      <CardContent>
        <div className="h-[320px] w-full pt-4">
          <ResponsiveContainer width="100%" height="100%">
            <AreaChart
              data={chartData}
              margin={{ top: 0, right: 0, left: -20, bottom: 0 }}
            >
              <defs>
                <linearGradient id="blueGradient" x1="0" y1="0" x2="0" y2="1">
                  <stop offset="5%" stopColor="#4ba3f7" stopOpacity={0.5} />
                  <stop offset="95%" stopColor="#4ba3f7" stopOpacity={0} />
                </linearGradient>
                <linearGradient id="redGradient" x1="0" y1="0" x2="0" y2="1">
                  <stop offset="5%" stopColor="#ef4444" stopOpacity={0.5} />
                  <stop offset="95%" stopColor="#ef4444" stopOpacity={0} />
                </linearGradient>
              </defs>

              <CartesianGrid
                strokeDasharray="3 3"
                stroke="#27272a"
                vertical={false}
              />

              <XAxis
                dataKey="time"
                stroke="#52525b"
                fontSize={10}
                tickLine={false}
                axisLine={false}
                interval={11}
                tick={{ dy: 10 }}
              />

              <YAxis
                stroke="#52525b"
                fontSize={10}
                tickLine={false}
                axisLine={false}
                allowDecimals={false}
              />

              <Tooltip
                cursor={{ stroke: "#3f3f46", strokeWidth: 1 }}
                contentStyle={{
                  backgroundColor: "rgba(9, 9, 11, 0.95)",
                  border: "1px solid rgba(63, 63, 70, 0.5)",
                  borderRadius: "12px",
                  boxShadow: "0 20px 25px -5px rgb(0 0 0 / 0.5)",
                  backdropFilter: "blur(8px)",
                }}
                itemStyle={{
                  fontSize: "11px",
                  fontWeight: "bold",
                  textTransform: "uppercase",
                }}
                labelStyle={{
                  color: "#71717a",
                  marginBottom: "4px",
                  fontSize: "10px",
                }}
              />

              {view === "all" && (
                <Area
                  type="monotone"
                  dataKey="allowed"
                  stroke="#4ba3f7"
                  strokeWidth={2}
                  fillOpacity={1}
                  fill="url(#blueGradient)"
                  animationDuration={1000}
                />
              )}

              <Area
                type="monotone"
                dataKey="threats"
                stroke="#ef4444"
                strokeWidth={2}
                fillOpacity={1}
                fill="url(#redGradient)"
                animationDuration={1000}
              />
            </AreaChart>
          </ResponsiveContainer>
        </div>

        <div className="flex gap-6 mt-6 px-2 border-t border-zinc-900 pt-4">
          <div className="flex flex-col">
            <span className="text-[10px] text-zinc-500 uppercase font-bold tracking-widest">
              Allowed
            </span>
            <span className="text-lg font-mono text-neon-blue">
              {chartData
                .reduce((acc, curr) => acc + curr.allowed, 0)
                .toLocaleString()}
            </span>
          </div>
          <div className="flex flex-col">
            <span className="text-[10px] text-zinc-500 uppercase font-bold tracking-widest">
              Blocked
            </span>
            <span className="text-lg font-mono text-red-500">
              {chartData
                .reduce((acc, curr) => acc + curr.threats, 0)
                .toLocaleString()}
            </span>
          </div>
        </div>
      </CardContent>
    </Card>
  );
}
