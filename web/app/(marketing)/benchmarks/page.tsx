import { MetricCard } from "@/components/marketing/benchmarks/metric-card";
import {
  LatencyChart,
  ScalabilityChart,
  ThroughputChart,
} from "@/components/marketing/benchmarks/benchmark-charts";
import { RawDataTable } from "@/components/marketing/benchmarks/raw-data-table";
import { Zap, Cpu, Server, Activity, ArrowLeft } from "lucide-react";
import Link from "next/link";
import fs from "fs";
import path from "path";

async function getBenchmarkData() {
  const filePath = path.join(process.cwd(), "public/data/benchmarks.json");
  const jsonData = fs.readFileSync(filePath, "utf8");
  return JSON.parse(jsonData);
}

export default async function BenchmarksPage() {
  const data = await getBenchmarkData();

  return (
    <div className="min-h-screen bg-background pt-32 pb-20">
      <div className="container mx-auto px-4 max-w-6xl">
        <div className="mb-12">
          <Link
            href="/"
            className="inline-flex items-center text-sm text-muted-foreground hover:text-neon-blue mb-6 transition-colors"
          >
            <ArrowLeft className="w-4 h-4 mr-2" /> Back to Home
          </Link>
          <h1 className="text-4xl md:text-6xl font-bold tracking-tight mb-4 bg-clip-text text-transparent bg-gradient-to-b from-white to-white/50">
            Performance Analysis
          </h1>
          <p className="text-xl text-muted-foreground max-w-2xl">
            Comprehensive benchmarks from the Argus WAF engine. All tests run on
            consumer-grade hardware to demonstrate real-world performance
            characteristics and optimization trade-offs.
          </p>
          <div className="mt-4 flex items-center gap-2 text-xs font-mono text-zinc-500">
            <span>Hardware: Intel i5-10210U @ 1.6GHz</span>
            <span>•</span>
            <span>Environment: Go 1.25.5, Linux amd64</span>
          </div>
        </div>

        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4 mb-16">
          <MetricCard
            title="Request Analysis"
            value="2.7"
            unit="µs"
            description="Custom request analyzer pipeline for telemetry extraction (header parsing, body inspection)"
            icon={Server}
          />

          <MetricCard
            title="Circuit Breaker"
            value="150"
            unit="ns"
            description="Per-call overhead from custom resilience state machine implementation"
            icon={Activity}
          />

          <MetricCard
            title="Middleware Overhead"
            value="2.8"
            unit="µs"
            description="Added latency by security layer (WAF integration + context building + body handling)"
            icon={Zap}
          />

          <MetricCard
            title="4-Core Efficiency"
            value="56%"
            unit=""
            description="Parallel efficiency with shared WAF instance. Lock-free transaction model."
            icon={Cpu}
          />
        </div>

        <div className="space-y-16">
          <div className="grid lg:grid-cols-2 gap-8">
            <div className="space-y-4">
              <h3 className="text-2xl font-bold text-white">
                Request Analysis Pipeline
              </h3>
              <p className="text-zinc-400 text-sm leading-relaxed">
                Custom analyzer latency across request types. Minimal requests
                process in 2.5µs, while complex header parsing reaches 9.1µs.
                Shows impact of request complexity on telemetry extraction.
              </p>
              <ThroughputChart data={data.analyzer_pipeline} />
            </div>
            <div className="space-y-4">
              <h3 className="text-2xl font-bold text-white">
                Latency Distribution
              </h3>
              <p className="text-zinc-400 text-sm leading-relaxed">
                Middleware latency across security modes. Baseline processing
                takes 3.5µs, while LatencyFirst adds 2.0µs overhead (5.4µs
                total). Paranoid mode achieves 4.6µs through optimized rule
                selection.
              </p>
              <LatencyChart data={data.middleware_modes} />{" "}
            </div>
          </div>

          <div className="grid lg:grid-cols-3 gap-8">
            <div className="lg:col-span-2 space-y-4">
              <h3 className="text-2xl font-bold text-white">
                Linear Scalability
              </h3>
              <p className="text-zinc-400 text-sm leading-relaxed">
                Parallel WAF processing with per-request transaction isolation.
                Performance scales from 309µs (1 core) to 133µs (8 cores),
                achieving 2.3x speedup with lock-free architecture. Efficiency
                decreases at higher core counts due to shared rule engine
                contention.
              </p>

              <ScalabilityChart data={data.scalability} />
            </div>
            <div className="space-y-4">
              <h3 className="text-2xl font-bold text-white">
                Memory Efficiency
              </h3>
              <p className="text-zinc-400 text-sm leading-relaxed">
                Memory allocation profile per operation. Log truncation achieves
                0 B/op for payloads under 4KB through in-place string
                manipulation. WAF middleware allocates ~6.8KB per request for
                transaction state and rule matching.
              </p>

              <RawDataTable data={data.memory_allocs} />
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
