export function StatusIndicator() {
  return (
    <div className="inline-flex items-center gap-2 px-3 py-1 rounded-full bg-white/5 border border-white/5 text-xs font-medium text-muted-foreground">
      <span className="relative flex h-2 w-2">
        <span className="animate-ping absolute inline-flex h-full w-full rounded-full bg-neon-green opacity-75"></span>
        <span className="relative inline-flex rounded-full h-2 w-2 bg-neon-green"></span>
      </span>
      System Operational
    </div>
  );
}
