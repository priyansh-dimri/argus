import { Metadata } from "next";
import Link from "next/link";
import { ArrowLeft } from "lucide-react";
import { ArchitectureFlow } from "@/components/marketing/architecture/architecture-flow";
import { ADRAccordion } from "@/components/marketing/architecture/adr-accordion";
import { TechStackCard } from "@/components/marketing/architecture/tech-stack-card";

export const metadata: Metadata = {
  title: "Architecture | Argus",
  description:
    "Deep dive into Argus hybrid AI WAF architecture and design decisions",
};

export default function ArchitecturePage() {
  return (
    <div className="min-h-screen bg-background pt-32 pb-20">
      <div className="container mx-auto px-4 max-w-7xl">
        <div className="mb-12">
          <Link
            href="/"
            className="inline-flex items-center text-sm text-muted-foreground hover:text-neon-blue mb-6 transition-colors"
          >
            <ArrowLeft className="w-4 h-4 mr-2" /> Back to Home
          </Link>

          <div className="flex items-start justify-between">
            <div>
              <h1 className="text-4xl md:text-6xl font-bold tracking-tight mb-4 bg-clip-text text-transparent bg-gradient-to-b from-white to-white/50">
                System Architecture
              </h1>
              <p className="text-xl text-muted-foreground max-w-2xl">
                Argus combines deterministic rule-based detection with
                probabilistic AI analysis through a three-layer architecture:
                WAF engine for microsecond blocking, circuit-protected API
                calls, and async threat storage.
              </p>
            </div>
          </div>
        </div>

        <section className="mb-16">
          <h2 className="text-2xl font-semibold mb-6">Request Flow</h2>
          <p className="text-muted-foreground mb-6">
            Interactive diagram showing request processing through middleware,
            WAF, Argus protection modes, and AI analysis. Click nodes for
            implementation details.
          </p>
          <ArchitectureFlow />
        </section>

        <section className="mb-16">
          <h2 className="text-2xl font-semibold mb-6">
            Architecture Decisions
          </h2>
          <p className="text-muted-foreground mb-6">
            Key design choices that define Argus&apos;s performance, resilience,
            and security characteristics.
          </p>
          <ADRAccordion />
        </section>

        <section className="mb-16">
          <h2 className="text-2xl font-semibold mb-6">Technology Stack</h2>
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
            <TechStackCard
              category="Backend"
              stack={[
                "Go 1.22+",
                "net/http (stdlib)",
                "Coraza WAF v3",
                "pgx/v5",
                "Google Genai SDK",
              ]}
            />
            <TechStackCard
              category="Frontend"
              stack={[
                "Next.js 16",
                "React 19",
                "TypeScript 5",
                "Tailwind CSS 4",
                "Radix UI (shadcn)",
              ]}
            />
            <TechStackCard
              category="Security"
              stack={[
                "OWASP CRS 4.0",
                "JWT (golang-jwt)",
                "gobreaker/v2",
                "Supabase Auth",
              ]}
            />
            <TechStackCard
              category="Infrastructure"
              stack={[
                "Supabase PostgreSQL",
                "Render (backend)",
                "Gemini 2.5 Flash",
                "Docker (sidecar)",
              ]}
            />
          </div>
        </section>
      </div>
    </div>
  );
}
