"use client";

import { useState } from "react";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { RefreshCw, Search, ChevronLeft, ChevronRight } from "lucide-react";
import { cn } from "@/lib/utils";
import { useThreatTable } from "@/hooks/use-threat-table";
import { LogDetailsDialog } from "./log-details-dialog";
import { ThreatLog } from "@/hooks/use-threats";

interface ThreatTableProps {
  projectId?: string;
}

export function ThreatTable({ projectId }: ThreatTableProps) {
  const { logs, totalCount, page, setPage, loading, hasNewData, refresh } =
    useThreatTable(projectId);
  const [selectedLog, setSelectedLog] = useState<ThreatLog | null>(null);
  const [detailsOpen, setDetailsOpen] = useState(false);

  const totalPages = Math.ceil(totalCount / 10);

  const handleRowClick = (log: ThreatLog) => {
    setSelectedLog(log);
    setDetailsOpen(true);
  };

  return (
    <Card className="col-span-4 lg:col-span-3 bg-black/40 border-white/10 backdrop-blur-md flex flex-col h-full">
      <CardHeader className="flex flex-row items-center justify-between pb-4">
        <CardTitle className="text-white text-lg">Detailed Logs</CardTitle>
        <div className="flex items-center gap-2">
          {hasNewData && (
            <Button
              variant="secondary"
              size="sm"
              onClick={refresh}
              className="h-8 bg-neon-blue/10 text-neon-blue hover:bg-neon-blue/20 border border-neon-blue/20 animate-pulse"
            >
              <RefreshCw className="mr-2 h-3 w-3" />
              New Logs Available
            </Button>
          )}
          <Button
            variant="ghost"
            size="icon"
            onClick={refresh}
            disabled={loading}
            className={cn(
              "h-8 w-8 text-muted-foreground",
              loading && "animate-spin"
            )}
          >
            <RefreshCw className="h-4 w-4" />
          </Button>
        </div>
      </CardHeader>

      <CardContent className="flex-1 flex flex-col">
        <div className="rounded-md border border-white/5 flex-1 relative min-h-[400px]">
          <Table>
            <TableHeader className="bg-white/5 sticky top-0 z-10 backdrop-blur-sm">
              <TableRow className="border-white/5 hover:bg-transparent">
                <TableHead className="w-[80px] text-gray-400">
                  Verdict
                </TableHead>
                <TableHead className="text-gray-400 w-[80px]">Method</TableHead>
                <TableHead className="text-gray-400">Path</TableHead>
                <TableHead className="text-gray-400 w-[120px]">IP</TableHead>
                <TableHead className="text-right text-gray-400 w-[100px]">
                  Time
                </TableHead>
                <TableHead className="w-[50px]"></TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {logs.length === 0 ? (
                <TableRow>
                  <TableCell
                    colSpan={6}
                    className="h-24 text-center text-muted-foreground"
                  >
                    {loading ? "Loading..." : "No threats found."}
                  </TableCell>
                </TableRow>
              ) : (
                logs.map((log) => (
                  <TableRow
                    key={log.id}
                    className="border-white/5 hover:bg-white/5 transition-colors cursor-pointer group"
                    onClick={() => handleRowClick(log)}
                  >
                    <TableCell>
                      <Badge
                        variant="outline"
                        className={cn(
                          "border-0 font-mono text-[10px] uppercase px-1.5 py-0.5",
                          log.is_threat
                            ? "bg-red-500/10 text-red-500"
                            : "bg-green-500/10 text-green-500"
                        )}
                      >
                        {log.is_threat ? "BLOCK" : "ALLOW"}
                      </Badge>
                    </TableCell>
                    <TableCell className="font-mono text-xs text-zinc-400">
                      {log.method}
                    </TableCell>
                    <TableCell
                      className="font-mono text-xs text-zinc-300 max-w-[150px] truncate"
                      title={log.route}
                    >
                      {log.route}
                    </TableCell>
                    <TableCell className="font-mono text-xs text-zinc-400">
                      {log.ip}
                    </TableCell>
                    <TableCell className="text-right text-xs text-zinc-500 font-mono whitespace-nowrap">
                      {new Date(log.timestamp).toLocaleTimeString([], {
                        hour: "2-digit",
                        minute: "2-digit",
                      })}
                    </TableCell>
                    <TableCell>
                      <Search className="h-3 w-3 text-zinc-600 group-hover:text-neon-blue transition-colors opacity-0 group-hover:opacity-100" />
                    </TableCell>
                  </TableRow>
                ))
              )}
            </TableBody>
          </Table>
        </div>

        <div className="flex items-center justify-between pt-4 border-t border-white/5 mt-4">
          <div className="text-xs text-muted-foreground">
            Showing {logs.length} of {totalCount} events
          </div>
          <div className="flex items-center gap-2">
            <Button
              variant="outline"
              size="icon"
              className="h-7 w-7 border-white/10 bg-transparent disabled:opacity-30"
              disabled={page <= 1}
              onClick={() => setPage((p) => p - 1)}
            >
              <ChevronLeft className="h-3 w-3" />
            </Button>
            <span className="text-xs font-mono text-zinc-400">
              Page {page} of {totalPages || 1}
            </span>
            <Button
              variant="outline"
              size="icon"
              className="h-7 w-7 border-white/10 bg-transparent disabled:opacity-30"
              disabled={page >= totalPages}
              onClick={() => setPage((p) => p + 1)}
            >
              <ChevronRight className="h-3 w-3" />
            </Button>
          </div>
        </div>
      </CardContent>

      <LogDetailsDialog
        log={selectedLog}
        open={detailsOpen}
        onOpenChange={setDetailsOpen}
      />
    </Card>
  );
}
