"use client";

import { useState } from "react";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Loader2, Copy, Check, Plus } from "lucide-react";
import type { Project } from "@/hooks/use-projects";

interface CreateProjectDialogProps {
  onCreate: (name: string) => Promise<Project>;
  onClose?: () => void;
}

export function CreateProjectDialog({
  onCreate,
  onClose,
}: CreateProjectDialogProps) {
  const [open, setOpen] = useState(false);
  const [name, setName] = useState("");
  const [loading, setLoading] = useState(false);
  const [newProject, setNewProject] = useState<{ api_key: string } | null>(
    null
  );
  const [copied, setCopied] = useState(false);

  const handleOpenChange = (isOpen: boolean) => {
    setOpen(isOpen);
    if (!isOpen) {
      setNewProject(null);
      setName("");
      if (onClose) onClose();
    }
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);
    try {
      const project = await onCreate(name);
      setNewProject(project);
      setName("");
    } catch (error) {
      console.error(error);
    } finally {
      setLoading(false);
    }
  };

  const copyKey = () => {
    navigator.clipboard.writeText(newProject?.api_key || "");
    setCopied(true);
    setTimeout(() => setCopied(false), 2000);
  };

  const handleDone = () => {
    handleOpenChange(false);
  };

  return (
    <Dialog open={open} onOpenChange={handleOpenChange}>
      <DialogTrigger asChild>
        <Button
          size="sm"
          className="gap-2 bg-white text-black hover:bg-white/90"
        >
          <Plus className="h-4 w-4" /> New Project
        </Button>
      </DialogTrigger>
      <DialogContent className="sm:max-w-[425px] bg-black/95 border-white/10 text-white backdrop-blur-xl">
        {!newProject ? (
          <form onSubmit={handleSubmit}>
            <DialogHeader>
              <DialogTitle>Create Project</DialogTitle>
              <DialogDescription className="text-gray-400">
                Create a new project to get an API Key for the SDK.
              </DialogDescription>
            </DialogHeader>
            <div className="grid gap-4 py-6">
              <div className="grid gap-2">
                <Label htmlFor="name" className="text-gray-300">
                  Project Name
                </Label>
                <Input
                  id="name"
                  value={name}
                  onChange={(e) => setName(e.target.value)}
                  className="bg-white/5 border-white/10 text-white focus:border-neon-blue/50"
                  placeholder="e.g. Production API"
                  autoFocus
                />
              </div>
            </div>
            <DialogFooter>
              <Button
                type="submit"
                disabled={loading || !name}
                className="bg-neon-blue text-white hover:bg-neon-blue/80"
              >
                {loading && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
                Create Project
              </Button>
            </DialogFooter>
          </form>
        ) : (
          <div className="space-y-6">
            <DialogHeader>
              <DialogTitle className="text-neon-green">
                Project Created!
              </DialogTitle>
              <DialogDescription className="text-gray-400">
                Here is your API Key. Copy it now, you won&#39;t see it again.
              </DialogDescription>
            </DialogHeader>

            <div className="p-4 rounded-lg bg-white/5 border border-neon-green/20 relative group">
              <code className="text-sm font-mono text-neon-green break-all">
                {newProject.api_key}
              </code>
              <Button
                size="icon"
                variant="ghost"
                className="absolute top-2 right-2 h-6 w-6 text-gray-400 hover:text-white"
                onClick={copyKey}
              >
                {copied ? (
                  <Check className="h-3 w-3" />
                ) : (
                  <Copy className="h-3 w-3" />
                )}
              </Button>
            </div>

            <DialogFooter>
              <Button
                onClick={handleDone}
                className="w-full bg-white/10 hover:bg-white/20 text-white"
              >
                Done
              </Button>
            </DialogFooter>
          </div>
        )}
      </DialogContent>
    </Dialog>
  );
}
