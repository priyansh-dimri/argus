export interface BenchmarkData {
  name?: string;
  value?: number;
  rps?: number;
  latency?: number;
  cores?: number;
  ns_op?: number;
  efficiency?: number;
  bytes?: number;
  allocs?: number;
  [key: string]: string | number | undefined;
}
