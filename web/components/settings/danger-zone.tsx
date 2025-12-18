"use client";

import { useState } from "react";
import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { AlertTriangle, Loader2 } from "lucide-react";
import { Project } from "@/hooks/use-projects";

interface DangerZoneProps {
  project: Project;
  onDelete: (id: string) => Promise<void>;
}

export function DangerZone({ project, onDelete }: DangerZoneProps) {
  const [deleting, setDeleting] = useState(false);

  const handleDelete = async () => {
    const confirmed = prompt(
      `To confirm deletion, type "${project.name}" below:`
    );
    if (confirmed !== project.name) return;

    setDeleting(true);
    await onDelete(project.id);
  };

  return (
    <Card className="bg-red-500/5 border-red-500/20 backdrop-blur-md">
      <CardHeader>
        <CardTitle className="text-red-400 flex items-center gap-2">
          <AlertTriangle className="h-5 w-5" /> Danger Zone
        </CardTitle>
        <CardDescription className="text-red-400/60">
          Irreversible actions for your project.
        </CardDescription>
      </CardHeader>
      <CardContent className="flex items-center justify-between">
        <div className="space-y-1">
          <h4 className="text-sm font-medium text-zinc-200">Delete Project</h4>
          <p className="text-xs text-zinc-500">
            Permanently remove this project and all its threat logs.
          </p>
        </div>
        <Button
          variant="outline"
          onClick={handleDelete}
          disabled={deleting}
          className="bg-red-500/10 text-red-500 hover:bg-red-500/20 border border-red-500/50"
        >
          {deleting ? (
            <Loader2 className="h-4 w-4 animate-spin" />
          ) : (
            "Delete Project"
          )}
        </Button>
      </CardContent>
    </Card>
  );
}
