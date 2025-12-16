"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";
import { cn } from "@/lib/utils";
import { Button } from "@/components/ui/button";
import { Menu } from "lucide-react";
import { FaGithub } from "react-icons/fa";
import { useState } from "react";
import {
  Sheet,
  SheetContent,
  SheetTrigger,
  SheetTitle,
  SheetHeader,
} from "@/components/ui/sheet";
import { ArgusLogo } from "../shared/argus-logo";

const navItems = [
  { name: "Docs", href: "/docs" },
  { name: "Architecture", href: "#architecture" },
  { name: "Benchmarks", href: "#benchmarks" },
];

export function Navbar() {
  const pathname = usePathname();
  const [isOpen, setIsOpen] = useState(false);

  return (
    <header className="fixed top-4 w-full z-50 px-4 md:px-0">
      <div className="mx-auto max-w-5xl">
        <nav className="glass rounded-full px-6 h-14 flex items-center justify-between transition-all duration-300">
          <Link href="/" className="flex items-center gap-2 group">
            <ArgusLogo />
          </Link>

          <div className="hidden md:flex items-center gap-8">
            {navItems.map((item) => (
              <Link
                key={item.name}
                href={item.href}
                className={cn(
                  "text-sm font-medium text-muted-foreground hover:text-foreground transition-colors relative",
                  pathname === item.href && "text-foreground"
                )}
              >
                {item.name}
              </Link>
            ))}
          </div>

          <div className="hidden md:flex items-center gap-3">
            <Link
              href="https://github.com/priyansh-dimri/argus"
              target="_blank"
              className="text-muted-foreground hover:text-foreground transition-colors"
            >
              <FaGithub className="h-5 w-5" />
            </Link>

            <div className="h-4 w-[1px] bg-border mx-1" />

            <Link href="/login">
              <Button
                variant="ghost"
                size="sm"
                className="hover:bg-transparent hover:text-neon-blue rounded-full px-5"
              >
                Log in
              </Button>
            </Link>

            <Link href="/signup">
              <Button
                size="sm"
                className="rounded-full px-5 bg-white text-black hover:bg-white/90 font-semibold shadow-[0_0_15px_rgba(255,255,255,0.3)] transition-shadow"
              >
                Sign Up
              </Button>
            </Link>
          </div>

          <div className="md:hidden">
            <Sheet open={isOpen} onOpenChange={setIsOpen}>
              <SheetTrigger asChild>
                <Button variant="ghost" size="icon" className="md:hidden">
                  <Menu className="h-5 w-5" />
                </Button>
              </SheetTrigger>
              <SheetContent
                side="top"
                className="w-full pt-20 border-b border-border bg-background/95 backdrop-blur-xl"
              >
                <SheetHeader>
                  <SheetTitle className="sr-only">Navigation Menu</SheetTitle>
                </SheetHeader>
                <div className="flex flex-col items-center gap-6">
                  {navItems.map((item) => (
                    <Link
                      key={item.name}
                      href={item.href}
                      onClick={() => setIsOpen(false)}
                      className="text-lg font-medium text-foreground/80 hover:text-foreground"
                    >
                      {item.name}
                    </Link>
                  ))}
                  <div className="flex flex-col gap-4 w-full max-w-xs mt-4">
                    <Link href="/login" onClick={() => setIsOpen(false)}>
                      <Button variant="outline" className="w-full">
                        Log in
                      </Button>
                    </Link>
                    <Link href="/signup" onClick={() => setIsOpen(false)}>
                      <Button className="w-full bg-white text-black hover:bg-white/90">
                        Deploy Proxy
                      </Button>
                    </Link>
                  </div>
                </div>
              </SheetContent>
            </Sheet>
          </div>
        </nav>
      </div>
    </header>
  );
}
