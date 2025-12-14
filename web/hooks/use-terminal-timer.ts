"use client";

import { useState, useEffect, RefObject } from "react";
import { useInView } from "framer-motion";

export function useTerminalTimer(ref: RefObject<HTMLElement | null>) {
  const isInView = useInView(ref, { once: true, margin: "-100px" });
  const [elapsedTime, setElapsedTime] = useState(0);
  const [restartKey, setRestartKey] = useState(0);

  useEffect(() => {
    if (!isInView) return;

    const startTime = Date.now();
    let animationFrameId: number;

    const update = () => {
      setElapsedTime(Date.now() - startTime);
      animationFrameId = requestAnimationFrame(update);
    };

    update();
    return () => cancelAnimationFrame(animationFrameId);
  }, [isInView, restartKey]);

  const restart = () => {
    setRestartKey((prev) => prev + 1);
    setElapsedTime(0);
  };

  return { elapsedTime, restartKey, restart, isInView };
}
