"use client";

import { useProjects } from "@/hooks/use-projects";
import { useThreats } from "@/hooks/use-threats";
import { StatsCards } from "@/components/dashboard/stats-cards";
import { ThreatTable } from "@/components/dashboard/threat-table";
import { ThreatChart } from "@/components/dashboard/threat-chart";
import { CreateProjectDialog } from "@/components/dashboard/create-project-dialog";
import { Loader2, LayoutGrid } from "lucide-react";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";

export default function DashboardPage() {
  const {
    projects,
    selectedProject,
    setSelectedProject,
    createProject,
    loading: projectsLoading,
  } = useProjects();
  const { threats } = useThreats(selectedProject?.id);

  if (projectsLoading) {
    return (
      <div className="flex h-[80vh] w-full items-center justify-center">
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
            Create your first project to generate an API key and start
            monitoring threats.
          </p>
        </div>
        <CreateProjectDialog onCreate={createProject} />
      </div>
    );
  }

  return (
    <div className="space-y-6 animate-in fade-in duration-500">
      <div className="flex flex-col sm:flex-row items-start sm:items-center justify-between gap-4">
        <div className="flex items-center gap-4">
          <h1 className="text-2xl font-bold tracking-tight text-white hidden md:block">
            Dashboard
          </h1>
          <div className="h-6 w-px bg-white/10 hidden md:block" />

          <Select
            value={selectedProject?.id}
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
        </div>

        <CreateProjectDialog onCreate={createProject} />
      </div>

      <StatsCards threats={threats} />

      <div className="grid grid-cols-1 lg:grid-cols-7 gap-6">
        <div className="lg:col-span-4">
          <ThreatChart threats={threats} />
        </div>
        <div className="lg:col-span-3">
          <div className="h-full min-h-[300px] rounded-xl border border-white/10 bg-black/40 backdrop-blur-md p-6 flex flex-col items-center justify-center text-center">
            <span className="text-neon-blue font-mono text-sm mb-2">
              ● LIVE FEED
            </span>
            {threats.length === 0 ? (
              <p className="text-muted-foreground text-sm">
                Waiting for incoming traffic for{" "}
                <strong>{selectedProject?.name}</strong>...
              </p>
            ) : (
              <div className="text-left w-full space-y-2">
                {threats.slice(0, 5).map((t) => (
                  <div
                    key={t.id}
                    className="text-xs font-mono text-gray-400 border-b border-white/5 pb-2"
                  >
                    <span
                      className={
                        t.is_threat ? "text-red-500" : "text-green-500"
                      }
                    >
                      {t.is_threat ? "BLOCK" : "ALLOW"}
                    </span>{" "}
                    {t.ip} - {t.route}
                  </div>
                ))}
              </div>
            )}
          </div>
        </div>
      </div>

      <ThreatTable threats={threats} />
    </div>
  );
}
