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
}

export function useThreats(limit = 50) {
  const [threats, setThreats] = useState<ThreatLog[]>([]);
  const [loading, setLoading] = useState(true);

  const supabase = useMemo(() => createClient(), []);

  useEffect(() => {
    const fetchInitial = async () => {
      const { data, error } = await supabase
        .from("threat_logs")
        .select("*")
        .order("timestamp", { ascending: false })
        .limit(limit);

      if (!error && data) {
        setThreats(data as ThreatLog[]);
      }
      setLoading(false);
    };

    fetchInitial();

    const channel: RealtimeChannel = supabase
      .channel("realtime-threats")
      .on(
        "postgres_changes",
        { event: "INSERT", schema: "public", table: "threat_logs" },
        (payload) => {
          const newThreat = payload.new as ThreatLog;
          setThreats((prev) => [newThreat, ...prev].slice(0, limit));
        }
      )
      .subscribe();

    return () => {
      supabase.removeChannel(channel);
    };
  }, [limit, supabase]);

  return { threats, loading };
}
