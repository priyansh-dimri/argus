"use client";

import { useEffect, useState, useMemo } from "react";
import { createClient } from "@/lib/supabase/client";
import type { RealtimeChannel } from "@supabase/supabase-js";

export interface ThreatLog {
  id: string;
  ip: string;
  route: string;
  method: string;
  is_threat: boolean;
  reason: string;
  confidence: number;
  timestamp: string;
  project_id: string;
  headers: Record<string, string>;
  metadata: Record<string, string>;
}

export function useThreats(projectId: string | undefined, limit = 50) {
  const [threats, setThreats] = useState<ThreatLog[]>([]);
  const [loading, setLoading] = useState(false);

  const supabase = useMemo(() => createClient(), []);

  useEffect(() => {
    let isCancelled = false;
    let channel: RealtimeChannel | null = null;

    Promise.resolve().then(async () => {
      if (isCancelled) return;

      if (!projectId) {
        setThreats([]);
        setLoading(false);
        return;
      }

      setLoading(true);

      const { data, error } = await supabase
        .from("threat_logs")
        .select("*")
        .eq("project_id", projectId)
        .order("timestamp", { ascending: false })
        .limit(limit);

      if (!isCancelled) {
        if (!error && data) {
          setThreats(data as ThreatLog[]);
        }
        setLoading(false);
      }

      if (isCancelled) return;

      channel = supabase
        .channel(`realtime-threats-${projectId}`)
        .on(
          "postgres_changes",
          {
            event: "INSERT",
            schema: "public",
            table: "threat_logs",
            filter: `project_id=eq.${projectId}`,
          },
          (payload) => {
            const newThreat = payload.new as ThreatLog;
            setThreats((prev) => [newThreat, ...prev].slice(0, limit));
          }
        )
        .subscribe();
    });

    return () => {
      isCancelled = true;
      if (channel) supabase.removeChannel(channel);
    };
  }, [projectId, limit, supabase]);

  return { threats, loading };
}
