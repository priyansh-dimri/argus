"use client";

import { Button } from "@/components/ui/button";
import { Menu } from "lucide-react";
import { Avatar, AvatarFallback } from "@/components/ui/avatar";
import { Sheet, SheetContent, SheetTrigger } from "@/components/ui/sheet";
import { Sidebar } from "./sidebar";
import { useState } from "react";

export function Header() {
  const [mobileOpen, setMobileOpen] = useState(false);

  return (
    <header className="sticky top-0 z-30 flex h-16 w-full items-center justify-between border-b border-white/5 bg-background/50 px-6 backdrop-blur-xl">
      <div className="flex items-center gap-4">
        <Sheet open={mobileOpen} onOpenChange={setMobileOpen}>
          <SheetTrigger asChild>
            <Button
              variant="ghost"
              size="icon"
              className="lg:hidden text-muted-foreground"
            >
              <Menu className="h-5 w-5" />
            </Button>
          </SheetTrigger>
          <SheetContent
            side="left"
            className="p-0 bg-black border-r border-white/10 w-72"
          >
            <Sidebar collapsed={false} setCollapsed={() => {}} />
          </SheetContent>
        </Sheet>
      </div>

      <div className="flex items-center gap-4">
        <div className="h-8 w-[1px] bg-white/10" />

        <div className="flex items-center gap-3">
          <div className="text-right hidden sm:block">
            <p className="text-sm font-medium leading-none text-white">
              Admin User
            </p>
            <p className="text-xs text-muted-foreground mt-1">
              admin@argus.com
            </p>
          </div>
          <Avatar className="h-9 w-9 border border-white/10">
            <AvatarFallback className="bg-neon-blue/10 text-neon-blue text-xs font-bold">
              AD
            </AvatarFallback>
          </Avatar>
        </div>
      </div>
    </header>
  );
}
