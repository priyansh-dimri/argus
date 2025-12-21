"use client";

import { motion } from "framer-motion";
import { BentoCard } from "./bento-card";
import { AiContextVisual } from "./features/ai-context-visual";
import { SpeedVisual } from "./features/speed-visual";
import { CircuitVisual } from "./features/circuit-visual";

export function BentoGrid() {
  return (
    <section className="py-32 relative overflow-hidden" id="features">
      <div className="container mx-auto px-4 max-w-6xl">
        <div className="mb-16 text-center">
          <motion.h2
            initial={{ opacity: 0, y: 20 }}
            whileInView={{ opacity: 1, y: 0 }}
            viewport={{ once: true }}
            className="text-3xl md:text-5xl font-bold tracking-tight mb-4 bg-clip-text text-transparent bg-gradient-to-b from-white to-white/50"
          >
            How Argus Works
          </motion.h2>
          <p className="text-muted-foreground text-lg">
            Three layers of defense. Five lines of code.
          </p>
        </div>

        <div className="grid grid-cols-1 md:grid-cols-3 md:grid-rows-4 gap-4 h-full">
          <BentoCard
            title="Context Aware Intelligence"
            description="Gemini 2.5 analyzes intent, not just syntax. It understands that a SQL tutorial is not an SQL injection."
            className="md:col-span-2 md:row-span-3"
            delay={0.1}
          >
            <AiContextVisual />
          </BentoCard>

          <BentoCard
            title="Sub-Millisecond Protection"
            description="WAF threat detection completes in 133µs with parallel processing; 75x faster than typical database queries."
            className="md:col-span-1 md:row-span-2"
            delay={0.2}
          >
            <SpeedVisual />
          </BentoCard>

          <BentoCard
            title="Fail-Open Circuit Breaker"
            description="If AI is down, traffic flows. Your uptime is our priority."
            className="md:col-span-1 md:row-span-2"
            delay={0.3}
          >
            <CircuitVisual />
          </BentoCard>

          <BentoCard
            className="md:col-span-1 md:row-span-1"
            minimal
            delay={0.4}
          >
            <div className="flex items-center gap-4 justify-center">
              <div className="p-3 bg-blue-500/40 rounded-full text-blue-400 flex-shrink-0">
                <svg
                  version="1.1"
                  id="Layer_1"
                  xmlns="http://www.w3.org/2000/svg"
                  x="0px"
                  y="0px"
                  viewBox="0 0 205.4 76.7"
                  className="h-6 w-auto"
                >
                  <style type="text/css"></style>
                  <g>
                    <g>
                      <g>
                        <g>
                          <path d="M15.5,23.2c-0.4,0-0.5-0.2-0.3-0.5l2.1-2.7c0.2-0.3,0.7-0.5,1.1-0.5h35.7c0.4,0,0.5,0.3,0.3,0.6l-1.7,2.6      c-0.2,0.3-0.7,0.6-1,0.6L15.5,23.2z" />
                        </g>
                      </g>
                    </g>
                    <g>
                      <g>
                        <g>
                          <path d="M0.4,32.4c-0.4,0-0.5-0.2-0.3-0.5l2.1-2.7c0.2-0.3,0.7-0.5,1.1-0.5h45.6c0.4,0,0.6,0.3,0.5,0.6l-0.8,2.4      c-0.1,0.4-0.5,0.6-0.9,0.6L0.4,32.4z" />
                        </g>
                      </g>
                    </g>
                    <g>
                      <g>
                        <g>
                          <path d="M24.6,41.6c-0.4,0-0.5-0.3-0.3-0.6l1.4-2.5c0.2-0.3,0.6-0.6,1-0.6h20c0.4,0,0.6,0.3,0.6,0.7L47.1,41      c0,0.4-0.4,0.7-0.7,0.7L24.6,41.6z" />
                        </g>
                      </g>
                    </g>
                    <g>
                      <g id="CXHf1q_3_">
                        <g>
                          <g>
                            <path d="M128.4,21.4c-6.3,1.6-10.6,2.8-16.8,4.4c-1.5,0.4-1.6,0.5-2.9-1c-1.5-1.7-2.6-2.8-4.7-3.8       c-6.3-3.1-12.4-2.2-18.1,1.5c-6.8,4.4-10.3,10.9-10.2,19c0.1,8,5.6,14.6,13.5,15.7c6.8,0.9,12.5-1.5,17-6.6       c0.9-1.1,1.7-2.3,2.7-3.7c-3.6,0-8.1,0-19.3,0c-2.1,0-2.6-1.3-1.9-3c1.3-3.1,3.7-8.3,5.1-10.9c0.3-0.6,1-1.6,2.5-1.6       c5.1,0,23.9,0,36.4,0c-0.2,2.7-0.2,5.4-0.6,8.1c-1.1,7.2-3.8,13.8-8.2,19.6c-7.2,9.5-16.6,15.4-28.5,17       c-9.8,1.3-18.9-0.6-26.9-6.6c-7.4-5.6-11.6-13-12.7-22.2c-1.3-10.9,1.9-20.7,8.5-29.3c7.1-9.3,16.5-15.2,28-17.3       c9.4-1.7,18.4-0.6,26.5,4.9c5.3,3.5,9.1,8.3,11.6,14.1C130,20.6,129.6,21.1,128.4,21.4z" />
                          </g>
                          <g>
                            <path d="M161.5,76.7c-9.1-0.2-17.4-2.8-24.4-8.8c-5.9-5.1-9.6-11.6-10.8-19.3c-1.8-11.3,1.3-21.3,8.1-30.2       c7.3-9.6,16.1-14.6,28-16.7c10.2-1.8,19.8-0.8,28.5,5.1c7.9,5.4,12.8,12.7,14.1,22.3c1.7,13.5-2.2,24.5-11.5,33.9       c-6.6,6.7-14.7,10.9-24,12.8C166.8,76.3,164.1,76.4,161.5,76.7z M185.3,36.3c-0.1-1.3-0.1-2.3-0.3-3.3       c-1.8-9.9-10.9-15.5-20.4-13.3c-9.3,2.1-15.3,8-17.5,17.4c-1.8,7.8,2,15.7,9.2,18.9c5.5,2.4,11,2.1,16.3-0.6       C180.5,51.3,184.8,44.9,185.3,36.3z" />
                          </g>
                        </g>
                      </g>
                    </g>
                  </g>
                </svg>
              </div>
              <div>
                <div className="font-bold text-white">Go 1.22+</div>
                <div className="text-xs text-muted-foreground">
                  Standard Library Middleware
                </div>
              </div>
            </div>
          </BentoCard>

          <BentoCard
            className="md:col-span-1 md:row-span-1"
            minimal
            delay={0.5}
          >
            <div className="flex items-center gap-4 justify-center">
              <div className="p-3 bg-blue-500/40 rounded-full text-blue-400 flex-shrink-0">
                <svg
                  width="109"
                  height="113"
                  viewBox="0 0 109 113"
                  fill="none"
                  xmlns="http://www.w3.org/2000/svg"
                  className="h-6 w-auto"
                >
                  <path
                    d="M63.7076 110.284C60.8481 113.885 55.0502 111.912 54.9813 107.314L53.9738 40.0627L99.1935 40.0627C107.384 40.0627 111.952 49.5228 106.859 55.9374L63.7076 110.284Z"
                    fill="url(#paint0_linear)"
                  />
                  <path
                    d="M63.7076 110.284C60.8481 113.885 55.0502 111.912 54.9813 107.314L53.9738 40.0627L99.1935 40.0627C107.384 40.0627 111.952 49.5228 106.859 55.9374L63.7076 110.284Z"
                    fill="url(#paint1_linear)"
                    fillOpacity="0.2"
                  />
                  <path
                    d="M45.317 2.07103C48.1765 -1.53037 53.9745 0.442937 54.0434 5.041L54.4849 72.2922H9.83113C1.64038 72.2922 -2.92775 62.8321 2.1655 56.4175L45.317 2.07103Z"
                    fill="#3ECF8E"
                  />
                  <defs>
                    <linearGradient
                      id="paint0_linear"
                      x1="53.9738"
                      y1="54.974"
                      x2="94.1635"
                      y2="71.8295"
                      gradientUnits="userSpaceOnUse"
                    >
                      <stop stopColor="#249361" />
                      <stop offset="1" stopColor="#3ECF8E" />
                    </linearGradient>
                    <linearGradient
                      id="paint1_linear"
                      x1="36.1558"
                      y1="30.578"
                      x2="54.4844"
                      y2="65.0806"
                      gradientUnits="userSpaceOnUse"
                    >
                      <stop />
                      <stop offset="1" stopOpacity="0" />
                    </linearGradient>
                  </defs>
                </svg>
              </div>
              <div>
                <div className="font-bold text-white">Supabase Realtime</div>
                <div className="text-xs text-muted-foreground">
                  Live Threat Streaming
                </div>
              </div>
            </div>
          </BentoCard>
        </div>
      </div>
    </section>
  );
}
