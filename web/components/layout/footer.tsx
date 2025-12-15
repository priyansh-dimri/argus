import Link from "next/link";
import { ShieldCheck, BookOpen } from "lucide-react";
import { FaGithub, FaLinkedin } from "react-icons/fa";
import { StatusIndicator } from "./status-indicator";

export function Footer() {
  return (
    <footer className="border-t border-white/10 bg-black pt-16 pb-8">
      <div className="container mx-auto px-4 max-w-6xl">
        <div className="grid grid-cols-1 md:grid-cols-4 gap-12 md:gap-8 mb-16">
          <div className="col-span-1 md:col-span-2">
            <Link href="/" className="flex items-center gap-2 mb-4 group w-fit">
              <div className="relative flex h-8 w-8 items-center justify-center rounded-lg bg-primary/10">
                <ShieldCheck className="h-5 w-5 text-neon-blue" />
              </div>
              <span className="font-bold text-lg">Argus</span>
            </Link>
            <p className="text-sm text-muted-foreground max-w-xs mb-6">
              Engineering secure systems from first principles. The hybrid WAF
              that balances latency and context for modern applications.
            </p>
            <StatusIndicator />
          </div>

          <div className="flex flex-col gap-4">
            <h4 className="font-semibold text-white">Product</h4>
            <Link
              href="#features"
              className="text-sm text-muted-foreground hover:text-white transition-colors"
            >
              Features
            </Link>
            <Link
              href="#integration"
              className="text-sm text-muted-foreground hover:text-white transition-colors"
            >
              Integration
            </Link>
            <Link
              href="/docs"
              className="text-sm text-muted-foreground hover:text-white transition-colors"
            >
              Documentation
            </Link>
          </div>

          <div className="flex flex-col gap-4">
            <h4 className="font-semibold text-white">Connect</h4>
            <Link
              href="https://github.com/priyansh-dimri"
              target="_blank"
              className="flex items-center gap-2 text-sm text-muted-foreground hover:text-white transition-colors"
            >
              <FaGithub className="w-4 h-4" /> GitHub
            </Link>
            <Link
              href="https://linkedin.com/in/priyanshdimri"
              target="_blank"
              className="flex items-center gap-2 text-sm text-muted-foreground hover:text-white transition-colors"
            >
              <FaLinkedin className="w-4 h-4" /> LinkedIn
            </Link>
            <Link
              href="https://hashnode.com/@priyanshdimri"
              target="_blank"
              className="flex items-center gap-2 text-sm text-muted-foreground hover:text-white transition-colors"
            >
              <BookOpen className="w-4 h-4" /> Blog
            </Link>
          </div>
        </div>

        <div className="pt-8 border-t border-white/5 flex flex-col md:flex-row justify-between items-center gap-4">
          <p className="text-xs text-muted-foreground text-center md:text-left">
            © 2025 Priyansh Dimri. Open Source under MIT License.
          </p>
          <div className="flex gap-6 text-xs text-muted-foreground">
            <Link href="#" className="hover:text-white transition-colors">
              Privacy Policy
            </Link>
            <Link href="#" className="hover:text-white transition-colors">
              Terms of Service
            </Link>
          </div>
        </div>
      </div>
    </footer>
  );
}
