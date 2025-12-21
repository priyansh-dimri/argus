"use client";

import {
  BarChart,
  Bar,
  LineChart,
  Line,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  ResponsiveContainer,
  Cell,
} from "recharts";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { BenchmarkData } from "./benchmark.types";

const TOOLTIP_STYLE = {
  backgroundColor: "#09090b",
  border: "1px solid #27272a",
  borderRadius: "8px",
  fontSize: "12px",
  color: "#fff",
};

export function LatencyChart({ data }: { data: BenchmarkData[] }) {
  return (
    <Card className="bg-black/20 border-white/10">
      <CardHeader>
        <CardTitle className="text-sm font-mono text-gray-400 tracking-widest">
          Latency Overhead (µs)
        </CardTitle>
      </CardHeader>
      <CardContent className="h-[300px]">
        <ResponsiveContainer width="100%" height="100%">
          <BarChart data={data} layout="vertical" margin={{ left: 40 }}>
            <CartesianGrid
              strokeDasharray="3 3"
              stroke="#ffffff10"
              horizontal={false}
            />
            <XAxis
              type="number"
              stroke="#666"
              fontSize={10}
              tickLine={false}
              axisLine={false}
              unit="µs"
            />
            <YAxis
              dataKey="name"
              type="category"
              stroke="#999"
              fontSize={11}
              tickLine={false}
              axisLine={false}
              width={140}
            />
            <Tooltip
              cursor={{ fill: "#ffffff05" }}
              contentStyle={TOOLTIP_STYLE}
            />
            <Bar
              dataKey="value"
              fill="#3b82f6"
              radius={[0, 4, 4, 0]}
              barSize={32}
            >
              {data.map((_, index) => (
                <Cell key={`cell-${index}`} fill="#3b82f6" />
              ))}
            </Bar>
          </BarChart>
        </ResponsiveContainer>
      </CardContent>
    </Card>
  );
}

export function ScalabilityChart({ data }: { data: BenchmarkData[] }) {
  const transformed = data.map((d) => ({
    cores: d.cores,
    latency: ((d.value ?? 0) / 1000).toFixed(1),
  }));

  return (
    <Card className="bg-black/20 border-white/10">
      <CardHeader>
        <CardTitle className="text-sm font-mono text-gray-400 uppercase tracking-widest">
          Parallel Scaling
        </CardTitle>
      </CardHeader>
      <CardContent className="h-[300px]">
        <ResponsiveContainer width="100%" height="100%">
          <LineChart data={transformed}>
            <CartesianGrid
              strokeDasharray="3 3"
              stroke="#ffffff10"
              vertical={false}
            />
            <XAxis
              dataKey="cores"
              stroke="#666"
              fontSize={10}
              tickLine={false}
              axisLine={false}
              label={{
                value: "CPU Cores",
                position: "insideBottom",
                offset: -5,
                fill: "#666",
                fontSize: 10,
              }}
            />
            <YAxis
              stroke="#666"
              fontSize={10}
              tickLine={false}
              axisLine={false}
              label={{
                value: "Latency (µs)",
                angle: -90,
                position: "insideLeft",
                fill: "#666",
                fontSize: 10,
              }}
            />
            <Tooltip contentStyle={TOOLTIP_STYLE} />
            <Line
              type="monotone"
              dataKey="latency"
              stroke="#10b981"
              strokeWidth={3}
              dot={{ r: 4, fill: "#10b981" }}
              activeDot={{ r: 6 }}
              name="Latency (µs)"
            />
          </LineChart>
        </ResponsiveContainer>
      </CardContent>
    </Card>
  );
}

export function ThroughputChart({ data }: { data: BenchmarkData[] }) {
  return (
    <Card className="bg-black/20 border-white/10">
      <CardHeader>
        <CardTitle className="text-sm font-mono text-gray-400 tracking-widest">
          Request Processing Time (µs)
        </CardTitle>
      </CardHeader>
      <CardContent className="h-[300px]">
        <ResponsiveContainer width="100%" height="100%">
          <BarChart data={data}>
            <CartesianGrid
              strokeDasharray="3 3"
              stroke="#ffffff10"
              vertical={false}
            />
            <XAxis
              dataKey="name"
              stroke="#666"
              fontSize={10}
              tickLine={false}
              axisLine={false}
            />
            <YAxis
              stroke="#666"
              fontSize={10}
              tickLine={false}
              axisLine={false}
            />
            <Tooltip
              cursor={{ fill: "#ffffff05" }}
              contentStyle={TOOLTIP_STYLE}
            />
            <Bar
              dataKey="value"
              fill="#f97316"
              radius={[4, 4, 0, 0]}
              barSize={40}
              name="Latency (µs)"
            />
          </BarChart>
        </ResponsiveContainer>
      </CardContent>
    </Card>
  );
}
