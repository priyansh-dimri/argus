"use client";

import { useCallback, useState } from "react";
import {
  ReactFlow,
  Node,
  Edge,
  Background,
  Controls,
  MiniMap,
  useNodesState,
  useEdgesState,
  MarkerType,
  Panel,
} from "@xyflow/react";
import "@xyflow/react/dist/style.css";
import { Card } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";

interface NodeData extends Record<string, unknown> {
  label: string;
  desc: string;
}

const initialNodes: Node<NodeData>[] = [
  {
    id: "sdk",
    type: "input",
    data: { label: "SDK Client", desc: "Bearer argus_xxx" },
    position: { x: 50, y: 50 },
    style: {
      background: "#4F46E5",
      color: "white",
      border: "none",
      borderRadius: "8px",
      padding: "12px",
    },
  },
  {
    id: "dashboard",
    type: "input",
    data: { label: "Next.js Dashboard", desc: "JWT Auth" },
    position: { x: 250, y: 50 },
    style: {
      background: "#10B981",
      color: "white",
      border: "none",
      borderRadius: "8px",
      padding: "12px",
    },
  },
  {
    id: "http",
    data: { label: "HTTP Server", desc: ":8080 + CORS" },
    position: { x: 150, y: 180 },
    style: {
      background: "#1F2937",
      border: "1px solid #374151",
      borderRadius: "8px",
      padding: "12px",
    },
  },
  {
    id: "auth",
    data: { label: "Auth Middleware", desc: "SDK: API Key\nDashboard: JWT" },
    position: { x: 150, y: 300 },
    style: {
      background: "#1F2937",
      border: "1px solid #60A5FA",
      borderRadius: "8px",
      padding: "12px",
    },
  },
  {
    id: "middleware",
    data: { label: "Protection Middleware", desc: "Body buffer + WAF check" },
    position: { x: 150, y: 420 },
    style: {
      background: "#1F2937",
      border: "1px solid #F59E0B",
      borderRadius: "8px",
      padding: "12px",
    },
  },
  {
    id: "waf",
    data: { label: "WAF Singleton", desc: "Coraza + CRS\n262µs avg" },
    position: { x: 350, y: 420 },
    style: {
      background: "#7C3AED",
      color: "white",
      border: "none",
      borderRadius: "8px",
      padding: "12px",
      fontWeight: "600",
    },
  },
  {
    id: "mode-router",
    data: {
      label: "Mode Router",
      desc: "Config-driven switch:\nLATENCY_FIRST | SMART_SHIELD | PARANOID",
    },
    position: { x: 150, y: 540 },
    style: {
      background: "#F59E0B",
      color: "#000",
      border: "2px solid #FBBF24",
      borderRadius: "8px",
      padding: "12px",
      fontWeight: "600",
    },
  },
  {
    id: "mode-latency",
    data: {
      label: "LATENCY_FIRST",
      desc: "WAF block → async log\nWAF pass → async log",
    },
    position: { x: 50, y: 680 },
    style: {
      background: "#FCD34D",
      color: "#000",
      border: "none",
      borderRadius: "8px",
      padding: "10px",
      fontSize: "12px",
    },
  },
  {
    id: "mode-smart",
    data: {
      label: "SMART_SHIELD",
      desc: "WAF pass → async log\nWAF block → sync AI",
    },
    position: { x: 250, y: 680 },
    style: {
      background: "#60A5FA",
      color: "#000",
      border: "none",
      borderRadius: "8px",
      padding: "10px",
      fontSize: "12px",
    },
  },
  {
    id: "mode-paranoid",
    data: {
      label: "PARANOID",
      desc: "All requests → sync AI\nBlock on threat",
    },
    position: { x: 450, y: 680 },
    style: {
      background: "#EF4444",
      color: "white",
      border: "none",
      borderRadius: "8px",
      padding: "10px",
      fontSize: "12px",
    },
  },
  {
    id: "breaker",
    data: {
      label: "Circuit Breaker",
      desc: "40ns overhead\n3 failures → OPEN",
    },
    position: { x: 250, y: 820 },
    style: {
      background: "#1F2937",
      border: "1px solid #F59E0B",
      borderRadius: "8px",
      padding: "12px",
    },
  },
  {
    id: "ai",
    data: { label: "Gemini AI", desc: "30K token limit\n60s timeout" },
    position: { x: 250, y: 960 },
    style: {
      background: "#8B5CF6",
      color: "white",
      border: "none",
      borderRadius: "8px",
      padding: "12px",
      fontWeight: "600",
    },
  },
  {
    id: "db",
    type: "output",
    data: { label: "Supabase DB", desc: "Async threat storage" },
    position: { x: 250, y: 1100 },
    style: {
      background: "#10B981",
      color: "white",
      border: "none",
      borderRadius: "8px",
      padding: "12px",
    },
  },
];

const initialEdges: Edge[] = [
  {
    id: "e1",
    source: "sdk",
    target: "http",
    animated: true,
    markerEnd: { type: MarkerType.ArrowClosed },
    style: { stroke: "#60A5FA" },
  },
  {
    id: "e2",
    source: "dashboard",
    target: "http",
    animated: true,
    markerEnd: { type: MarkerType.ArrowClosed },
    style: { stroke: "#10B981" },
  },
  {
    id: "e3",
    source: "http",
    target: "auth",
    markerEnd: { type: MarkerType.ArrowClosed },
  },
  {
    id: "e4",
    source: "auth",
    target: "middleware",
    markerEnd: { type: MarkerType.ArrowClosed },
  },
  {
    id: "e5",
    source: "waf",
    target: "middleware",
    label: "WAF.Check()",
    markerEnd: { type: MarkerType.ArrowClosed },
    style: { stroke: "#7C3AED", strokeWidth: 2 },
  },
  {
    id: "e6",
    source: "middleware",
    target: "mode-router",
    label: "WAF result",
    markerEnd: { type: MarkerType.ArrowClosed },
    style: { stroke: "#F59E0B", strokeWidth: 2 },
  },
  {
    id: "e7",
    source: "mode-router",
    target: "mode-latency",
    animated: true,
    markerEnd: { type: MarkerType.ArrowClosed },
    style: { stroke: "#FCD34D", strokeWidth: 2 },
  },
  {
    id: "e8",
    source: "mode-router",
    target: "mode-smart",
    animated: true,
    markerEnd: { type: MarkerType.ArrowClosed },
    style: { stroke: "#60A5FA", strokeWidth: 2 },
  },
  {
    id: "e9",
    source: "mode-router",
    target: "mode-paranoid",
    animated: true,
    markerEnd: { type: MarkerType.ArrowClosed },
    style: { stroke: "#EF4444", strokeWidth: 2 },
  },
  {
    id: "e10",
    source: "mode-latency",
    target: "breaker",
    label: "Async only",
    markerEnd: { type: MarkerType.ArrowClosed },
    style: { stroke: "#FCD34D", strokeWidth: 2 },
  },
  {
    id: "e11",
    source: "mode-smart",
    target: "breaker",
    label: "If WAF blocked",
    markerEnd: { type: MarkerType.ArrowClosed },
    style: { stroke: "#60A5FA" },
  },
  {
    id: "e12",
    source: "mode-paranoid",
    target: "breaker",
    label: "All requests",
    markerEnd: { type: MarkerType.ArrowClosed },
    style: { stroke: "#EF4444", strokeWidth: 2 },
  },
  {
    id: "e13",
    source: "breaker",
    target: "ai",
    markerEnd: { type: MarkerType.ArrowClosed },
    style: { stroke: "#8B5CF6", strokeWidth: 2 },
  },
  {
    id: "e14",
    source: "ai",
    target: "db",
    label: "Async goroutine",
    markerEnd: { type: MarkerType.ArrowClosed },
    animated: true,
    style: { stroke: "#10B981" },
  },
];

export function ArchitectureFlow() {
  const [nodes, , onNodesChange] = useNodesState(initialNodes);
  const [edges, , onEdgesChange] = useEdgesState(initialEdges);
  const [selectedNode, setSelectedNode] = useState<Node<NodeData> | null>(null);

  const onNodeClick = useCallback(
    (_: React.MouseEvent, node: Node<NodeData>) => {
      setSelectedNode(node);
    },
    []
  );

  return (
    <Card className="h-[700px] bg-background/60 backdrop-blur-md border border-white/10 relative">
      <ReactFlow
        nodes={nodes}
        edges={edges}
        onNodesChange={onNodesChange}
        onEdgesChange={onEdgesChange}
        onNodeClick={onNodeClick}
        fitView
        className="rounded-xl"
      >
        <Background color="#374151" gap={16} />
        <Controls
          className="bg-background border border-white/10"
          style={{
            backgroundColor: "#222",
            color: "#fff",
          }}
        />
        <MiniMap
          className="bg-background/80 border border-white/10"
          maskColor="rgba(0,0,0,0.6)"
        />

        {selectedNode && (
          <Panel
            position="top-right"
            className="bg-background/95 backdrop-blur-md border border-white/10 rounded-lg p-4 max-w-xs"
          >
            <div className="text-sm">
              <Badge className="mb-2">{selectedNode.id}</Badge>
              <h3 className="font-semibold mb-1">{selectedNode.data.label}</h3>
              <p className="text-xs text-muted-foreground whitespace-pre-line">
                {selectedNode.data.desc}
              </p>
            </div>
          </Panel>
        )}
      </ReactFlow>
    </Card>
  );
}
