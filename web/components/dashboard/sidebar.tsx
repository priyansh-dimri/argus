"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";
import { cn } from "@/lib/utils";
import {
  ShieldCheck,
  LayoutDashboard,
  Settings,
  LogOut,
  ChevronLeft,
  ChevronRight,
} from "lucide-react";
import { Button } from "@/components/ui/button";
import { createClient } from "@/lib/supabase/client";
import { useRouter } from "next/navigation";
import { motion, AnimatePresence } from "framer-motion";

interface SidebarProps {
  collapsed: boolean;
  setCollapsed: (collapsed: boolean) => void;
}

export function Sidebar({ collapsed, setCollapsed }: SidebarProps) {
  const pathname = usePathname();
  const router = useRouter();
  const supabase = createClient();

  const handleSignOut = async () => {
    await supabase.auth.signOut();
    router.push("/login");
  };

  const navItems = [
    { name: "Dashboard", href: "/dashboard", icon: LayoutDashboard },
    { name: "Settings", href: "/settings", icon: Settings },
  ];

  return (
    <aside
      className={cn(
        "fixed left-0 top-0 z-40 h-screen border-r border-white/5 bg-black/95 backdrop-blur-xl transition-all duration-300 ease-in-out",
        collapsed ? "w-16" : "w-64"
      )}
    >
      <div className="flex h-full flex-col justify-between py-4">
        <div className="px-4 mb-8">
          <div className="flex items-center gap-3 overflow-hidden">
            <div className="flex h-8 w-8 shrink-0 items-center justify-center rounded-lg bg-primary/10">
              <ShieldCheck className="h-5 w-5 text-neon-blue" />
            </div>
            <AnimatePresence>
              {!collapsed && (
                <motion.span
                  initial={{ opacity: 0, x: -10 }}
                  animate={{ opacity: 1, x: 0 }}
                  exit={{ opacity: 0, x: -10 }}
                  className="font-bold text-lg tracking-tight whitespace-nowrap"
                >
                  Argus
                </motion.span>
              )}
            </AnimatePresence>
          </div>
        </div>

        <nav className="flex-1 px-2 space-y-1">
          {navItems.map((item) => {
            const isActive = pathname === item.href;
            return (
              <Link
                key={item.href}
                href={item.href}
                className={cn(
                  "flex items-center gap-3 px-3 py-2 rounded-md transition-all duration-200 group relative overflow-hidden",
                  isActive
                    ? "bg-white/5 text-white"
                    : "text-muted-foreground hover:text-white hover:bg-white/5"
                )}
              >
                {isActive && (
                  <motion.div
                    layoutId="activeNav"
                    className="absolute left-0 top-0 bottom-0 w-0.5 bg-neon-blue"
                  />
                )}

                <item.icon
                  className={cn(
                    "h-5 w-5 shrink-0 transition-colors",
                    isActive ? "text-neon-blue" : "group-hover:text-white"
                  )}
                />

                {!collapsed && (
                  <motion.span
                    initial={{ opacity: 0 }}
                    animate={{ opacity: 1 }}
                    exit={{ opacity: 0 }}
                    className="truncate text-sm font-medium"
                  >
                    {item.name}
                  </motion.span>
                )}

                <div className="absolute inset-0 bg-gradient-to-r from-neon-blue/5 to-transparent opacity-0 group-hover:opacity-100 transition-opacity" />
              </Link>
            );
          })}
        </nav>

        <div className="px-2 space-y-2">
          <Button
            variant="ghost"
            size="sm"
            onClick={() => setCollapsed(!collapsed)}
            className="hidden lg:flex w-full justify-start px-3 text-muted-foreground hover:text-white"
          >
            {collapsed ? (
              <ChevronRight className="h-4 w-4" />
            ) : (
              <ChevronLeft className="h-4 w-4" />
            )}
            {!collapsed && <span className="ml-3 text-xs">Collapse</span>}
          </Button>

          <div className="h-px bg-white/5 mx-2 my-2" />

          <Button
            variant="ghost"
            size="sm"
            onClick={handleSignOut}
            className="w-full justify-start px-3 text-red-400 hover:text-red-300 hover:bg-red-500/10"
          >
            <LogOut className="h-4 w-4 shrink-0" />
            {!collapsed && <span className="ml-3">Sign Out</span>}
          </Button>
        </div>
      </div>
    </aside>
  );
}
