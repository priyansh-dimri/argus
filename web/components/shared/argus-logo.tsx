import { ShieldCheck } from "lucide-react";

export function ArgusLogo() {
  return (
    <>
      <div className="relative flex h-8 w-8 items-center justify-center rounded-lg bg-primary/10 transition-colors group-hover:bg-primary/20">
        <ShieldCheck className="h-5 w-5 text-neon-blue transition-transform group-hover:scale-110" />
      </div>
      <span className="font-bold text-lg tracking-tight">Argus</span>
    </>
  );
}
