"use client";

import { Button } from "@/components/ui/button";
import { Bell, Menu, Search } from "lucide-react";
import { Input } from "@/components/ui/input";
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

        <div className="hidden md:flex items-center relative">
          <Search className="absolute left-2.5 top-2.5 h-4 w-4 text-muted-foreground" />
          <Input
            type="search"
            placeholder="Search logs..."
            className="w-64 pl-9 bg-white/5 border-white/10 focus:border-neon-blue/50 transition-colors h-9 text-sm"
          />
        </div>
      </div>

      <div className="flex items-center gap-4">
        <Button
          variant="ghost"
          size="icon"
          className="text-muted-foreground hover:text-white relative"
        >
          <Bell className="h-5 w-5" />
          <span className="absolute top-2 right-2 h-2 w-2 rounded-full bg-neon-orange animate-pulse" />
        </Button>

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
