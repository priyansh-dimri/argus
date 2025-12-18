"use client";

import { useState } from "react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Copy, Check, RotateCw, Loader2 } from "lucide-react";
import { Project } from "@/hooks/use-projects";

interface ProjectConfigProps {
  project: Project;
  onUpdate: (id: string, name: string) => Promise<void>;
  onRotate: (id: string) => Promise<string>;
}

export function ProjectConfig({
  project,
  onUpdate,
  onRotate,
}: ProjectConfigProps) {
  const [name, setName] = useState(project.name);
  const [loading, setLoading] = useState(false);
  const [rotating, setRotating] = useState(false);
  const [copied, setCopied] = useState(false);

  const handleSave = async () => {
    setLoading(true);
    await onUpdate(project.id, name);
    setLoading(false);
  };

  const handleRotate = async () => {
    if (!confirm("This will invalidate the old key immediately. Continue?"))
      return;
    setRotating(true);
    await onRotate(project.id);
    setRotating(false);
  };

  const copyKey = () => {
    navigator.clipboard.writeText(project.api_key);
    setCopied(true);
    setTimeout(() => setCopied(false), 2000);
  };

  return (
    <div className="space-y-6">
      <Card className="bg-black/40 border-white/10 backdrop-blur-md">
        <CardHeader>
          <CardTitle className="text-white">Project Configuration</CardTitle>
          <CardDescription>
            Manage your project identifiers and credentials.
          </CardDescription>
        </CardHeader>
        <CardContent className="space-y-4">
          <div className="space-y-2">
            <Label className="text-zinc-400">Project Name</Label>
            <Input
              value={name}
              onChange={(e) => setName(e.target.value)}
              className="bg-zinc-900 border-zinc-800 text-white focus:border-neon-blue/50"
            />
          </div>

          <div className="space-y-2">
            <Label className="text-zinc-400">API Key</Label>
            <div className="flex gap-2">
              <div className="relative flex-1">
                <Input
                  value={project.api_key}
                  readOnly
                  className="bg-zinc-900/50 border-zinc-800 text-zinc-500 font-mono pr-10"
                />
                <Button
                  size="icon"
                  variant="ghost"
                  onClick={copyKey}
                  className="absolute right-1 top-1 h-7 w-7 text-zinc-400 hover:text-white"
                >
                  {copied ? (
                    <Check className="h-3 w-3" />
                  ) : (
                    <Copy className="h-3 w-3" />
                  )}
                </Button>
              </div>
              <Button
                variant="outline"
                onClick={handleRotate}
                disabled={rotating}
                className="border-zinc-700 text-zinc-300 hover:bg-zinc-800 hover:text-white"
              >
                {rotating ? (
                  <Loader2 className="h-4 w-4 animate-spin" />
                ) : (
                  <RotateCw className="h-4 w-4 mr-2" />
                )}
                Rotate
              </Button>
            </div>
            <p className="text-[10px] text-red-400/80">
              Rotating this key will immediately block all traffic using the old
              key.
            </p>
          </div>
        </CardContent>
        <CardFooter className="border-t border-white/5 bg-white/5 py-3 flex justify-end">
          <Button
            onClick={handleSave}
            disabled={loading || name === project.name}
            className="bg-white text-black hover:bg-white/90"
          >
            {loading && <Loader2 className="h-4 w-4 mr-2 animate-spin" />}
            Save Changes
          </Button>
        </CardFooter>
      </Card>
    </div>
  );
}
