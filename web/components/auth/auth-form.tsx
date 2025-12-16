"use client";

import * as React from "react";
import { cn } from "@/lib/utils";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Loader2, CheckCircle2 } from "lucide-react";
import { createClient } from "@/lib/supabase/client";

interface AuthFormProps extends React.HTMLAttributes<HTMLDivElement> {
  type: "login" | "signup";
}

export function AuthForm({ className, type, ...props }: AuthFormProps) {
  const [isLoading, setIsLoading] = React.useState<boolean>(false);
  const [isSuccess, setIsSuccess] = React.useState<boolean>(false);
  const [error, setError] = React.useState<string | null>(null);

  async function onSubmit(event: React.SyntheticEvent) {
    event.preventDefault();
    setIsLoading(true);
    setError(null);

    const target = event.target as typeof event.target & {
      email: { value: string };
    };
    const email = target.email.value;
    const supabase = createClient();

    const { error } = await supabase.auth.signInWithOtp({
      email,
      options: {
        emailRedirectTo: `${window.location.origin}/auth/callback`,
      },
    });

    setIsLoading(false);
    if (error) {
      setError(error.message);
      return;
    }

    setIsSuccess(true);
  }

  if (isSuccess) {
    return (
      <div className={cn("grid gap-6", className)}>
        <div className="flex flex-col items-center justify-center text-center space-y-4 p-8 border border-neon-green/20 bg-neon-green/5 rounded-xl">
          <div className="h-12 w-12 rounded-full bg-neon-green/20 flex items-center justify-center text-neon-green">
            <CheckCircle2 className="h-6 w-6" />
          </div>
          <h3 className="text-xl font-semibold text-white">Check your email</h3>
          <p className="text-sm text-gray-400">
            We sent a magic link to your inbox. Click it to log in instantly.
          </p>
          <Button
            variant="outline"
            className="w-full mt-4 border-white/10 hover:bg-white/5"
            onClick={() => setIsSuccess(false)}
          >
            Try another email
          </Button>
        </div>
      </div>
    );
  }

  return (
    <div className={cn("grid gap-6", className)} {...props}>
      <form onSubmit={onSubmit}>
        <div className="grid gap-4">
          <div className="grid gap-2">
            <Label className="text-gray-300" htmlFor="email">
              Email
            </Label>
            <Input
              id="email"
              placeholder="name@example.com"
              type="email"
              autoCapitalize="none"
              autoComplete="email"
              autoCorrect="off"
              disabled={isLoading}
              className="bg-white/5 border-white/10 text-white placeholder:text-gray-500 focus:border-neon-blue/50 transition-colors h-11"
              required
            />
          </div>

          {error && (
            <div className="text-sm text-red-500 bg-red-500/10 p-3 rounded border border-red-500/20">
              {error}
            </div>
          )}

          <Button
            disabled={isLoading}
            className="h-11 bg-white text-black hover:bg-gray-200 font-medium"
          >
            {isLoading && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
            {type === "login" ? "Send Magic Link" : "Create Account"}
          </Button>
        </div>
      </form>

      <div className="text-center text-xs text-gray-500 mt-4">
        Passwordless authentication via Supabase.
      </div>
    </div>
  );
}
