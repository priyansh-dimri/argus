"use client";

import { motion } from "framer-motion";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { CodeBlock } from "./code-block";
import { Terminal, Box } from "lucide-react";

const GO_CODE = `package main

import (
    "net/http"
    "time"
    "github.com/priyansh-dimri/argus"
)

func main() {
    // Initialize WAF and client
    waf := argus.NewDefaultWAF()
    client := argus.NewClient(
        "https://api.argus.io",
        "your-api-key",
        5*time.Second,
    )
    
    // Configure security mode
    config := argus.Config{
        AppID:  "your-app-id",
        APIKey: "your-api-key",
        Mode:   argus.SmartShield,
    }
    
    // Create and apply middleware
    shield := argus.NewMiddleware(client, waf, config)
    
    http.Handle("/api", shield.Protect(yourHandler))
}`;

const SIDECAR_CODE = `# !!TO BE UPDATED!!
# Run Argus as a sidecar for Node/Python/Ruby apps
$ docker run -d \\
  -p 8080:8080 \\
  -e TARGET_URL=http://localhost:3000 \\
  -e ARGUS_MODE=SMART_SHIELD \\
  ghcr.io/priyansh-dimri/argus:latest

# Traffic flow: User -> Argus (8080) -> Your App (3000)`;

const CONFIG_CODE = `
!!TO BE UPDATED!!
{
  "app_id": "app_xy7_229",
  "mode": "SMART_SHIELD",
  "resilience": {
    "circuit_breaker": {
      "timeout_ms": 300,
      "fail_open": true
    }
  },
  "rules": {
    "block_sqli": true,
    "block_xss": true,
    "ai_verification": true
  }
}`;

export function IntegrationSection() {
  return (
    <section className="py-24 relative overflow-hidden bg-black/50 border-y border-white/5">
      <div className="container mx-auto px-4 max-w-5xl">
        <div className="grid lg:grid-cols-2 gap-12">
          <motion.div
            initial={{ opacity: 0, x: -20 }}
            whileInView={{ opacity: 1, x: 0 }}
            viewport={{ once: true }}
            className="space-y-6"
          >
            <h2 className="text-3xl md:text-4xl font-bold tracking-tight text-white">
              Drop-in Middleware. <br />
              <span className="text-muted-foreground">Zero Friction.</span>
            </h2>
            <p className="text-lg text-muted-foreground leading-relaxed">
              Whether you are building a Go monolith or a polyglot microservices
              mesh, Argus fits into your stack in minutes.
            </p>

            <ul className="space-y-4 pt-4">
              <li className="flex items-start gap-3">
                <div className="p-2 rounded-lg bg-neon-blue/10 text-neon-blue">
                  <Box className="w-5 h-5" />
                </div>
                <div>
                  <h3 className="font-semibold text-white">Native Go SDK</h3>
                  <p className="text-sm text-muted-foreground">
                    WAF + AI + Circuit Breaker. 5 lines. Done.
                  </p>
                </div>
              </li>
              <li className="flex items-start gap-3">
                <div className="p-2 rounded-lg bg-neon-orange/10 text-neon-orange">
                  <Terminal className="w-5 h-5" />
                </div>
                <div>
                  <h3 className="font-semibold text-white">Sidecar Proxy</h3>
                  <p className="text-sm text-muted-foreground">
                    Protect Node, Python, or Ruby apps without changing code.
                  </p>
                </div>
              </li>
            </ul>
          </motion.div>

          <motion.div
            initial={{ opacity: 0, x: 20 }}
            whileInView={{ opacity: 1, x: 0 }}
            viewport={{ once: true }}
          >
            <Tabs defaultValue="go" className="w-full">
              <TabsList className="grid w-full grid-cols-3 bg-white/5 border border-white/10">
                <TabsTrigger
                  value="go"
                  className="data-[state=active]:bg-white/10"
                >
                  Go SDK
                </TabsTrigger>
                <TabsTrigger
                  value="sidecar"
                  className="data-[state=active]:bg-white/10"
                >
                  Sidecar
                </TabsTrigger>
                <TabsTrigger
                  value="config"
                  className="data-[state=active]:bg-white/10"
                >
                  Config
                </TabsTrigger>
              </TabsList>

              <div className="mt-4 relative">
                <div className="absolute -inset-1 bg-gradient-to-r from-neon-blue/20 to-purple-500/20 rounded-xl blur opacity-20" />

                <TabsContent value="go" className="relative">
                  <CodeBlock code={GO_CODE} language="go" />
                </TabsContent>
                <TabsContent value="sidecar" className="relative">
                  <CodeBlock code={SIDECAR_CODE} language="bash" />
                </TabsContent>
                <TabsContent value="config" className="relative">
                  <CodeBlock code={CONFIG_CODE} language="json" />
                </TabsContent>
              </div>
            </Tabs>
          </motion.div>
        </div>
      </div>
    </section>
  );
}
