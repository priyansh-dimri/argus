import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { BenchmarkData } from "./benchmark.types";

export function RawDataTable({ data }: { data: BenchmarkData[] }) {
  return (
    <div className="rounded-lg border border-white/10 bg-black/40 overflow-hidden">
      <Table>
        <TableHeader className="bg-white/5">
          <TableRow className="border-white/10 hover:bg-transparent">
            <TableHead className="text-gray-400 font-mono text-xs uppercase">
              Benchmark Name
            </TableHead>
            <TableHead className="text-right text-gray-400 font-mono text-xs uppercase">
              Memory (Bytes)
            </TableHead>
            <TableHead className="text-right text-gray-400 font-mono text-xs uppercase">
              Allocs
            </TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          {data.map((row) => (
            <TableRow
              key={row.name}
              className="border-white/5 hover:bg-white/5"
            >
              <TableCell className="font-medium text-zinc-300 font-mono text-xs">
                {row.name}
              </TableCell>
              <TableCell className="text-right text-zinc-400 font-mono text-xs">
                {(row.bytes_per_op ?? 0).toLocaleString()}
              </TableCell>
              <TableCell className="text-right text-zinc-400 font-mono text-xs">
                {(row.allocs_per_op ?? 0).toLocaleString()}
              </TableCell>
            </TableRow>
          ))}
        </TableBody>
      </Table>
    </div>
  );
}
