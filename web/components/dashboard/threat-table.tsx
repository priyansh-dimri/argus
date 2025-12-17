"use client";

import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { Badge } from "@/components/ui/badge";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { ThreatLog } from "@/hooks/use-threats";
import { cn } from "@/lib/utils";

interface ThreatTableProps {
  threats: ThreatLog[];
}

export function ThreatTable({ threats }: ThreatTableProps) {
  return (
    <Card className="col-span-4 lg:col-span-3 bg-black/40 border-white/10 backdrop-blur-md">
      <CardHeader>
        <CardTitle className="text-white">Recent Activity</CardTitle>
      </CardHeader>
      <CardContent>
        <div className="rounded-md border border-white/5">
          <Table>
            <TableHeader className="bg-white/5">
              <TableRow className="border-white/5 hover:bg-transparent">
                <TableHead className="w-[100px] text-gray-400">
                  Verdict
                </TableHead>
                <TableHead className="text-gray-400">Method</TableHead>
                <TableHead className="text-gray-400">Path</TableHead>
                <TableHead className="text-gray-400">IP Address</TableHead>
                <TableHead className="text-gray-400">Confidence</TableHead>
                <TableHead className="text-right text-gray-400">Time</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {threats.length === 0 ? (
                <TableRow>
                  <TableCell
                    colSpan={6}
                    className="h-24 text-center text-muted-foreground"
                  >
                    No threats detected (yet).
                  </TableCell>
                </TableRow>
              ) : (
                threats.map((log) => (
                  <TableRow
                    key={log.id}
                    className="border-white/5 hover:bg-white/5 transition-colors"
                  >
                    <TableCell>
                      <Badge
                        variant="outline"
                        className={cn(
                          "border-0 font-mono text-xs uppercase",
                          log.is_threat
                            ? "bg-red-500/10 text-red-500 hover:bg-red-500/20"
                            : "bg-green-500/10 text-green-500 hover:bg-green-500/20"
                        )}
                      >
                        {log.is_threat ? "BLOCK" : "ALLOW"}
                      </Badge>
                    </TableCell>
                    <TableCell className="font-mono text-xs text-white">
                      {log.method}
                    </TableCell>
                    <TableCell
                      className="font-mono text-xs text-gray-300 max-w-[200px] truncate"
                      title={log.route}
                    >
                      {log.route}
                    </TableCell>
                    <TableCell className="font-mono text-xs text-gray-400">
                      {log.ip}
                    </TableCell>
                    <TableCell className="font-mono text-xs text-gray-400">
                      {(log.confidence * 100).toFixed(0)}%
                    </TableCell>
                    <TableCell className="text-right text-xs text-gray-500 font-mono">
                      {new Date(log.timestamp).toLocaleTimeString()}
                    </TableCell>
                  </TableRow>
                ))
              )}
            </TableBody>
          </Table>
        </div>
      </CardContent>
    </Card>
  );
}
