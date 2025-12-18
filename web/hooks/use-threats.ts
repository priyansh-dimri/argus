"use client";

import { useEffect, useState, useMemo } from "react";
import { createClient } from "@/lib/supabase/client";
import { RealtimeChannel } from "@supabase/supabase-js";

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
  payload: string;
}

export function useThreats(projectId: string | undefined, limit = 50) {
  const [threats, setThreats] = useState<ThreatLog[]>([]);
  const [loading, setLoading] = useState(!!projectId);
  const supabase = useMemo(() => createClient(), []);

  useEffect(() => {
    if (!projectId) {
      setThreats([]);
      setLoading(false);
      return;
    }

    setLoading(true);

    const fetchInitial = async () => {
      try {
        const { data, error } = await supabase
          .from("threat_logs")
          .select("*")
          .eq("project_id", projectId)
          .order("timestamp", { ascending: false })
          .limit(limit);

        if (error) throw error;
        setThreats(data || []);
      } catch (err) {
        console.error("Fetch error:", err);
      } finally {
        setLoading(false);
      }
    };

    fetchInitial();

    const channel: RealtimeChannel = supabase
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
          setThreats((prev) => {
            if (prev.some((t) => t.id === newThreat.id)) return prev;
            return [newThreat, ...prev].slice(0, limit);
          });
        }
      )
      .subscribe();

    return () => {
      supabase.removeChannel(channel);
    };
  }, [projectId, limit, supabase]);

  return { threats, loading };
}
