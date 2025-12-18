"use client";

import { useProjects } from "@/hooks/use-projects";
import { ProjectConfig } from "@/components/settings/project-config";
import { DangerZone } from "@/components/settings/danger-zone";
import { Loader2, LayoutGrid } from "lucide-react";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { CreateProjectDialog } from "@/components/dashboard/create-project-dialog";

export default function SettingsPage() {
  const {
    projects,
    selectedProject,
    setSelectedProject,
    createProject,
    updateProjectName,
    rotateApiKey,
    deleteProject,
    loading,
  } = useProjects();

  if (loading) {
    return (
      <div className="flex h-[50vh] w-full items-center justify-center">
        <Loader2 className="h-8 w-8 animate-spin text-neon-blue" />
      </div>
    );
  }

  if (projects.length === 0) {
    return (
      <div className="flex flex-col h-[70vh] items-center justify-center text-center space-y-6">
        <div className="h-20 w-20 rounded-full bg-white/5 flex items-center justify-center border border-white/10">
          <LayoutGrid className="h-10 w-10 text-muted-foreground" />
        </div>
        <div className="max-w-md space-y-2">
          <h2 className="text-2xl font-bold text-white">No Projects Found</h2>
          <p className="text-muted-foreground">
            Create your first project to configure settings.
          </p>
        </div>
        <CreateProjectDialog onCreate={(name) => createProject(name, false)} />
      </div>
    );
  }

  if (!selectedProject) {
    return (
      <div className="flex h-[50vh] w-full items-center justify-center text-muted-foreground">
        Select a project to manage settings.
      </div>
    );
  }

  return (
    <div className="max-w-4xl mx-auto space-y-8 animate-in fade-in duration-500">
      <div className="flex flex-col sm:flex-row items-start sm:items-center justify-between gap-4 border-b border-white/5 pb-6">
        <div>
          <h1 className="text-3xl font-bold tracking-tight text-white mb-2">
            Settings
          </h1>
          <p className="text-muted-foreground">
            Manage configuration for{" "}
            <span className="text-white font-medium">
              {selectedProject.name}
            </span>
            .
          </p>
        </div>

        <div className="flex items-center gap-3">
          <Select
            value={selectedProject.id}
            onValueChange={(val) =>
              setSelectedProject(projects.find((p) => p.id === val) || null)
            }
          >
            <SelectTrigger className="w-[200px] bg-white/5 border-white/10 text-white">
              <SelectValue placeholder="Select Project" />
            </SelectTrigger>
            <SelectContent className="bg-black border-white/10 text-white">
              {projects.map((p) => (
                <SelectItem key={p.id} value={p.id}>
                  {p.name}
                </SelectItem>
              ))}
            </SelectContent>
          </Select>

          <CreateProjectDialog
            onCreate={(name) => createProject(name, false)}
          />
        </div>
      </div>

      <ProjectConfig
        key={selectedProject.id}
        project={selectedProject}
        onUpdate={updateProjectName}
        onRotate={rotateApiKey}
      />

      <DangerZone project={selectedProject} onDelete={deleteProject} />
    </div>
  );
}
