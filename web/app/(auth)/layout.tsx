import Link from "next/link";
import { AuthVisual } from "@/components/auth/auth-visual";
import { ArgusLogo } from "@/components/shared/argus-logo";

export default function AuthLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <div className="min-h-screen w-full lg:grid lg:grid-cols-2">
      <div className="relative flex min-h-screen flex-col items-center justify-center p-8 md:p-12 lg:p-16 bg-background">
        <Link
          href="/"
          className="absolute left-8 top-8 flex items-center gap-2 group"
        >
          <ArgusLogo />
        </Link>

        <div className="w-full max-w-sm space-y-8 animate-in fade-in slide-in-from-bottom-4 duration-700">
          {children}
        </div>
      </div>

      {/* right auth visual column which is hidden in mobile screens */}
      <div className="relative hidden h-full flex-col bg-muted p-10 text-white lg:flex dark:border-l border-white/5">
        <div className="absolute inset-0 bg-zinc-900" />
        <AuthVisual />
      </div>
    </div>
  );
}
