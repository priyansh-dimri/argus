"use client";

import { useState } from "react";
import { Button } from "@/components/ui/button";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Loader2, UserX, AlertTriangle } from "lucide-react";
import { createClient } from "@/lib/supabase/client";
import { fetchAPI } from "@/lib/api";
import { useRouter } from "next/navigation";

interface DeleteAccountDialogProps {
  trigger?: React.ReactNode;
}

export function DeleteAccountDialog({ trigger }: DeleteAccountDialogProps) {
  const [open, setOpen] = useState(false);
  const [confirmation, setConfirmation] = useState("");
  const [loading, setLoading] = useState(false);
  const router = useRouter();
  const supabase = createClient();

  const handleDeleteAccount = async (e: React.FormEvent) => {
    e.preventDefault();
    if (confirmation !== "DELETE MY ACCOUNT") return;

    setLoading(true);
    try {
      await fetchAPI("/account", { method: "DELETE" });
      await supabase.auth.signOut();
      router.push("/login");
    } catch (err) {
      console.error(err);
      alert("Failed to delete account");
      setLoading(false);
    }
  };

  return (
    <Dialog open={open} onOpenChange={setOpen}>
      <DialogTrigger asChild>
        {trigger || (
          <Button
            variant="ghost"
            className="text-red-500 hover:text-red-400 hover:bg-red-500/10"
          >
            <UserX className="h-4 w-4 mr-2" /> Delete Account
          </Button>
        )}
      </DialogTrigger>
      <DialogContent className="sm:max-w-[425px] bg-red-950/30 border-red-500/20 text-white backdrop-blur-xl">
        <DialogHeader>
          <DialogTitle className="text-red-500 flex items-center gap-2">
            <AlertTriangle className="h-5 w-5" /> Delete Account
          </DialogTitle>
          <DialogDescription className="text-red-200/70">
            This action is irreversible. It will permanently delete your account
            and all associated projects.
          </DialogDescription>
        </DialogHeader>

        <form onSubmit={handleDeleteAccount} className="space-y-4 py-4">
          <div className="space-y-2">
            <Label htmlFor="confirmation" className="text-red-200">
              Type{" "}
              <span className="font-mono font-bold select-all">
                DELETE MY ACCOUNT
              </span>{" "}
              to confirm
            </Label>
            <Input
              id="confirmation"
              value={confirmation}
              onChange={(e) => setConfirmation(e.target.value)}
              className="bg-black/50 border-red-500/30 text-red-500 placeholder:text-red-500/30 focus:border-red-500"
              placeholder="DELETE MY ACCOUNT"
              autoComplete="off"
            />
          </div>

          <DialogFooter>
            <Button
              type="button"
              variant="ghost"
              onClick={() => setOpen(false)}
              className="hover:bg-white/5 text-zinc-400"
            >
              Cancel
            </Button>
            <Button
              type="submit"
              disabled={loading || confirmation !== "DELETE MY ACCOUNT"}
              className="bg-red-600 hover:bg-red-700 text-white border border-red-500"
            >
              {loading && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
              Permanently Delete
            </Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  );
}
