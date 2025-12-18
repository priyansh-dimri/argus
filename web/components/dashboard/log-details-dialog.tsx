"use client";

import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogDescription,
} from "@/components/ui/dialog";
import { ThreatLog } from "@/hooks/use-threats";
import { Badge } from "@/components/ui/badge";
import { ScrollArea } from "@/components/ui/scroll-area";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { cn } from "@/lib/utils";

interface LogDetailsDialogProps {
  log: ThreatLog | null;
  open: boolean;
  onOpenChange: (open: boolean) => void;
}

export function LogDetailsDialog({
  log,
  open,
  onOpenChange,
}: LogDetailsDialogProps) {
  if (!log) return null;

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="max-w-2xl bg-zinc-950 border-zinc-800 text-zinc-100 max-h-[85vh] flex flex-col p-0 gap-0">
        <DialogHeader className="p-6 border-b border-zinc-800">
          <div className="flex items-center justify-between">
            <div className="flex items-center gap-3">
              <Badge
                variant={"outline"}
                className={cn(
                  "uppercase",
                  log.is_threat &&
                    "bg-red-500/10 text-red-500 border-red-500/30",
                  !log.is_threat &&
                    "text-neon-green border-neon-green/30 bg-neon-green/10"
                )}
              >
                {log.is_threat ? "Blocked" : "Allowed"}
              </Badge>
              <DialogTitle className="font-mono text-lg">{log.id}</DialogTitle>
            </div>
            <span className="text-xs text-zinc-500 font-mono">
              {new Date(log.timestamp).toLocaleString()}
            </span>
          </div>
          <DialogDescription className="mt-2 text-zinc-400 font-mono text-xs">
            {log.method} {log.route} • {log.ip}
          </DialogDescription>
        </DialogHeader>

        <Tabs defaultValue="payload" className="flex-1 flex flex-col min-h-0">
          <div className="px-6 pt-4">
            <TabsList className="bg-zinc-900 border border-zinc-800">
              <TabsTrigger value="payload">Payload</TabsTrigger>
              <TabsTrigger value="headers">Headers</TabsTrigger>
              <TabsTrigger value="metadata">Metadata</TabsTrigger>
              <TabsTrigger value="analysis">Analysis</TabsTrigger>
            </TabsList>
          </div>

          <ScrollArea className="flex-1 p-6">
            <TabsContent value="payload" className="mt-0">
              <div className="space-y-4">
                <div className="space-y-2">
                  <h4 className="text-sm font-medium text-zinc-400">
                    Raw Payload
                  </h4>
                  <pre className="p-4 rounded-lg bg-zinc-900 border border-zinc-800 font-mono text-xs text-zinc-300 whitespace-pre-wrap break-all">
                    {log.payload || "<Empty Payload>"}
                  </pre>
                </div>
              </div>
            </TabsContent>

            <TabsContent value="headers" className="mt-0">
              <div className="grid gap-2">
                {Object.entries(log.headers || {}).map(([k, v]) => (
                  <div
                    key={k}
                    className="flex flex-col space-y-1 p-2 rounded bg-zinc-900/50 border border-zinc-800/50"
                  >
                    <span className="text-xs font-semibold text-zinc-500">
                      {k}
                    </span>
                    <span className="text-xs font-mono text-zinc-300 break-all">
                      {v as string}
                    </span>
                  </div>
                ))}
                {(!log.headers || Object.keys(log.headers).length === 0) && (
                  <p className="text-sm text-zinc-500 italic">
                    No headers captured.
                  </p>
                )}
              </div>
            </TabsContent>

            <TabsContent value="metadata" className="mt-0">
              <div className="grid gap-2">
                {Object.entries(log.metadata || {}).map(([k, v]) => (
                  <div
                    key={k}
                    className="flex flex-col space-y-1 p-2 rounded bg-zinc-900/50 border border-zinc-800/50"
                  >
                    <span className="text-xs font-semibold text-neon-blue">
                      {k}
                    </span>
                    <span className="text-xs font-mono text-zinc-300 break-all">
                      {v as string}
                    </span>
                  </div>
                ))}
                {(!log.metadata || Object.keys(log.metadata).length === 0) && (
                  <p className="text-sm text-zinc-500 italic">
                    No metadata provided.
                  </p>
                )}
              </div>
            </TabsContent>

            <TabsContent value="analysis" className="mt-0">
              <div className="space-y-4">
                <div className="p-4 rounded-lg bg-zinc-900 border border-zinc-800 space-y-3">
                  <div>
                    <span className="text-xs text-zinc-500 uppercase font-bold">
                      AI Reason
                    </span>
                    <p className="text-sm text-zinc-300 mt-1 leading-relaxed">
                      {log.reason}
                    </p>
                  </div>
                  <div>
                    <span className="text-xs text-zinc-500 uppercase font-bold">
                      Confidence Score
                    </span>
                    <div className="mt-1 flex items-center gap-2">
                      <div className="h-2 flex-1 bg-zinc-800 rounded-full overflow-hidden">
                        <div
                          className={cn(
                            "h-full rounded-full",
                            log.confidence > 0.8
                              ? "bg-red-500"
                              : "bg-yellow-500"
                          )}
                          style={{ width: `${log.confidence * 100}%` }}
                        />
                      </div>
                      <span className="text-xs font-mono">
                        {(log.confidence * 100).toFixed(0)}%
                      </span>
                    </div>
                  </div>
                </div>
              </div>
            </TabsContent>
          </ScrollArea>
        </Tabs>
      </DialogContent>
    </Dialog>
  );
}
