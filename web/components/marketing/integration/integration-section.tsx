"use client";

import { motion } from "framer-motion";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { CodeBlock } from "./code-block";
import { Terminal, Box } from "lucide-react";

const GO_CODE = `package main

import (
    "log"
    "net/http"
    "time"
    
    "github.com/priyansh-dimri/argus/pkg/argus"
)

func main() {
    // Initialize WAF
    waf, _ := argus.NewWAF()
    
    // Initialize Argus client
    client := argus.NewClient(
        "https://api.example.com", // add backend URL
        "api-key",
        20*time.Second,
    )
    
    // Configure security mode
    config := argus.Config{
        Mode: argus.SmartShield,
    }
    
    // Create and apply middleware
    shield := argus.NewMiddleware(client, waf, config)
    
    http.Handle("/api/", shield.Protect(yourHandler))
    http.ListenAndServe(":8080", nil)
}`;

const SIDECAR_CODE = `
docker run -d \\
  --name argus-sidecar \\
  -p 8000:8000 \\
  -e TARGET_URL=http://host.docker.internal:3000 \\
  -e ARGUS_API_KEY=your-api-key \\
  -e ARGUS_API_URL=https://api.argus-security.com \\
  ghcr.io/priyansh-dimri/argus-sidecar:latest

# Route traffic through protection modes:
# Smart Shield:   http://localhost:8000/smart-shield/<route>
# Latency First:  http://localhost:8000/latency-first/<route>
# Paranoid:       http://localhost:8000/paranoid/<route>`;

export function IntegrationSection() {
  return (
    <section
      className="py-24 relative overflow-hidden bg-black/50 border-y border-white/5"
      id="integration"
    >
      <div className="container mx-auto px-4 max-w-5xl">
        <div className="flex flex-col lg:grid lg:grid-cols-2 gap-8 lg:gap-12">
          <motion.div
            initial={{ opacity: 0, x: -20 }}
            whileInView={{ opacity: 1, x: 0 }}
            viewport={{ once: true }}
            className="space-y-6"
          >
            <h2 className="text-2xl md:text-3xl lg:text-4xl font-bold tracking-tight text-white">
              Drop-in Middleware. <br />
              <span className="text-muted-foreground">Zero Friction.</span>
            </h2>
            <p className="text-base md:text-lg text-muted-foreground leading-relaxed">
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
              <TabsList className="grid w-full grid-cols-2 bg-white/5 border border-white/10">
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
              </TabsList>

              <div className="mt-4 relative">
                <div className="absolute -inset-1 bg-gradient-to-r from-neon-blue/20 to-purple-500/20 rounded-xl blur opacity-20" />

                <TabsContent value="go" className="relative">
                  <CodeBlock code={GO_CODE} language="go" />
                </TabsContent>
                <TabsContent value="sidecar" className="relative">
                  <CodeBlock code={SIDECAR_CODE} language="bash" />
                </TabsContent>
              </div>
            </Tabs>
          </motion.div>
        </div>
      </div>
    </section>
  );
}
