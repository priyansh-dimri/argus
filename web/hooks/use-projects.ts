"use client";

import { useEffect, useState, useCallback } from "react";
import { fetchAPI } from "@/lib/api";

export interface Project {
  id: string;
  user_id: string;
  name: string;
  api_key: string;
  created_at: string;
}

export function useProjects() {
  const [projects, setProjects] = useState<Project[]>([]);
  const [selectedProject, setSelectedProject] = useState<Project | null>(null);
  const [loading, setLoading] = useState(true);

  const refreshProjects = useCallback(async () => {
    try {
      const data = await fetchAPI("/projects");
      setProjects(data || []);

      if (data?.length > 0) {
        setSelectedProject((current) => current || data[0]);
      }
    } catch (err) {
      console.error("Failed to fetch projects", err);
    } finally {
      setLoading(false);
    }
  }, []);

  const createProject = async (name: string, skipRefresh = false) => {
    try {
      const res = await fetchAPI("/projects", {
        method: "POST",
        body: JSON.stringify({ name }),
      });
      if (!skipRefresh) {
        await refreshProjects();
      }
      return res.project;
    } catch (err) {
      throw err;
    }
  };

  useEffect(() => {
    refreshProjects();
  }, [refreshProjects]);

  return {
    projects,
    selectedProject,
    setSelectedProject,
    createProject,
    loading,
    refreshProjects,
  };
}
