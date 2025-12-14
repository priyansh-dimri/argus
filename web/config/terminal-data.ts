import { TerminalLine } from "@/components/marketing/terminal-window";

export const ATTACKER_LINES: TerminalLine[] = [
  { text: "nmap -p 80,443 target.api.com", delay: 900 },
  { text: "Starting Nmap 7.93...", delay: 500 },
  { text: "Scanning target.api.com (192.0.2.1)", delay: 450 },
  { text: "PORT     STATE SERVICE", delay: 250 },
  { text: "80/tcp   open  http", delay: 250 },
  { text: "443/tcp  open  https", delay: 250 },
  { text: "", delay: 200 },
  {
    text: 'curl -X POST /v1/data/query -d \'{"query": "DROP ALL TABLES;"}\'',
    delay: 400,
  },
  {
    text: "Error: 403 Forbidden. Access Denied by Argus Proxy.",
    delay: 2500,
    color: "neon-orange",
  },
];

export const ARGUS_LINES: TerminalLine[] = [
  { text: "Argus v1.0 initialized...", delay: 150 },
  { text: "Mode: SMART_SHIELD", delay: 200, className: "text-neon-blue" },
  { text: "Listening on port 8080...", delay: 200 },
  { text: "", delay: 2800 },
  {
    text: "[Coraza] Suspicious Pattern Detected (942100)",
    delay: 300,
    className: "text-yellow-400",
  },
  { text: "[Gemini] Analyzing Context...", delay: 400 },
  {
    text: "Context: Login Form (Untrusted)",
    delay: 300,
    className: "text-muted-foreground",
  },
  {
    text: "[Gemini] Verdict: THREAT (High Confidence)",
    delay: 800,
    className: "text-red-400",
  },
  {
    text: "[Shield] BLOCKED 403 (12ms)",
    delay: 200,
    className:
      "text-red-500 font-bold bg-red-500/10 px-2 py-1 inline-block rounded",
  },
];
