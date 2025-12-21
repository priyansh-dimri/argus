import { cn } from "@/lib/utils";
import { LucideIcon } from "lucide-react";

interface MetricCardProps {
  title: string;
  value: string;
  unit: string;
  description: string;
  icon: LucideIcon;
  trend?: string;
  trendUp?: boolean;
}

export function MetricCard({
  title,
  value,
  unit,
  description,
  icon: Icon,
  trend,
  trendUp,
}: MetricCardProps) {
  return (
    <div className="relative group overflow-hidden rounded-xl border border-white/10 bg-black/40 p-6 backdrop-blur-md transition-all hover:border-neon-blue/30">
      <div className="absolute inset-0 bg-gradient-to-br from-neon-blue/5 to-transparent opacity-0 transition-opacity group-hover:opacity-100" />

      <div className="relative z-10 flex flex-col justify-between h-full">
        <div className="flex items-start justify-between mb-4">
          <div className="p-2 rounded-lg bg-white/5 text-neon-blue">
            <Icon className="w-5 h-5" />
          </div>
          {trend && (
            <span
              className={cn(
                "text-xs font-mono px-2 py-1 rounded-full border",
                trendUp
                  ? "border-green-500/30 text-green-400 bg-green-500/10"
                  : "border-red-500/30 text-red-400 bg-red-500/10"
              )}
            >
              {trend}
            </span>
          )}
        </div>

        <div>
          <div className="flex items-baseline gap-1">
            <h3 className="text-3xl font-bold text-white tracking-tight">
              {value}
            </h3>
            <span className="text-sm font-mono text-muted-foreground">
              {unit}
            </span>
          </div>
          <p className="text-sm font-medium text-gray-400 mt-1">{title}</p>
          <p className="text-xs text-muted-foreground mt-3 leading-relaxed border-t border-white/5 pt-3">
            {description}
          </p>
        </div>
      </div>
    </div>
  );
}
