"use client";

import { useState, useEffect, useCallback, useMemo } from "react";
import { createClient } from "@/lib/supabase/client";
import { ThreatLog } from "./use-threats";

export function useThreatTable(projectId: string | undefined, pageSize = 10) {
  const [logs, setLogs] = useState<ThreatLog[]>([]);
  const [totalCount, setTotalCount] = useState(0);
  const [page, setPage] = useState(1);
  const [loading, setLoading] = useState(false);
  const [hasNewData, setHasNewData] = useState(false);

  const supabase = useMemo(() => createClient(), []);

  const fetchLogs = useCallback(
    async (targetPage: number) => {
      if (!projectId) return;

      setLoading(true);
      const from = (targetPage - 1) * pageSize;
      const to = from + pageSize - 1;

      try {
        const { data, error, count } = await supabase
          .from("threat_logs")
          .select("*", { count: "exact" })
          .eq("project_id", projectId)
          .order("timestamp", { ascending: false })
          .range(from, to);

        if (!error && data) {
          setLogs(data as ThreatLog[]);
          setTotalCount(count || 0);
          setHasNewData(false);
        }
      } catch (err) {
        console.error("[useThreatTable] Fetch error:", err);
      } finally {
        setLoading(false);
      }
    },
    [projectId, pageSize, supabase]
  );

  useEffect(() => {
    if (projectId) {
      fetchLogs(page);
    }
  }, [projectId, page, fetchLogs]);

  useEffect(() => {
    setPage(1);
  }, [projectId]);

  useEffect(() => {
    if (!projectId) return;

    const channel = supabase
      .channel(`table-updates-${projectId}`)
      .on(
        "postgres_changes",
        {
          event: "INSERT",
          schema: "public",
          table: "threat_logs",
          filter: `project_id=eq.${projectId}`,
        },
        () => setHasNewData(true)
      )
      .subscribe();

    return () => {
      supabase.removeChannel(channel);
    };
  }, [projectId, supabase]);

  return {
    logs,
    totalCount,
    page,
    setPage,
    loading,
    hasNewData,
    refresh: () => {
      if (page === 1) fetchLogs(1);
      else setPage(1);
    },
  };
}
