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
import { DeleteAccountDialog } from "./delete-account-dialog";

interface DangerZoneProps {
  project: Project;
  onDelete: (id: string) => Promise<void>;
}

export function DangerZone({ project, onDelete }: DangerZoneProps) {
  const [deletingProject, setDeletingProject] = useState(false);

  const handleDeleteProject = async () => {
    const confirmed = prompt(
      `To confirm PROJECT deletion, type "${project.name}" below:`
    );
    if (confirmed !== project.name) return;

    setDeletingProject(true);
    await onDelete(project.id);
    setDeletingProject(false);
  };

  return (
    <Card className="bg-red-500/5 border-red-500/20 backdrop-blur-md">
      <CardHeader>
        <CardTitle className="text-red-400 flex items-center gap-2">
          <AlertTriangle className="h-5 w-5" /> Danger Zone
        </CardTitle>
        <CardDescription className="text-red-400/60">
          Irreversible actions for your project and account.
        </CardDescription>
      </CardHeader>

      <CardContent className="space-y-6">
        <div className="flex items-center justify-between">
          <div className="space-y-1">
            <h4 className="text-sm font-medium text-zinc-200">
              Delete Project
            </h4>
            <p className="text-xs text-zinc-500">
              Permanently remove{" "}
              <span className="font-mono text-zinc-400">{project.name}</span>{" "}
              and its threat logs.
            </p>
          </div>
          <Button
            variant="outline"
            onClick={handleDeleteProject}
            disabled={deletingProject}
            className="text-red-500 hover:text-red-400 hover:bg-red-500/10"
          >
            {deletingProject ? (
              <Loader2 className="h-4 w-4 animate-spin" />
            ) : (
              "Delete Project"
            )}
          </Button>
        </div>

        <div className="h-px bg-red-500/10" />

        <div className="flex items-center justify-between">
          <div className="space-y-1">
            <h4 className="text-sm font-medium text-zinc-200">
              Delete Account
            </h4>
            <p className="text-xs text-zinc-500">
              Wipe your user account and ALL projects. This cannot be undone.
            </p>
          </div>

          <DeleteAccountDialog
            trigger={
              <Button
                variant="ghost"
                className="text-red-500 hover:text-red-400 hover:bg-red-500/10"
              >
                Delete Account
              </Button>
            }
          />
        </div>
      </CardContent>
    </Card>
  );
}
