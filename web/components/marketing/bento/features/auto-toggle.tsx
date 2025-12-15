"use client";

import { motion } from "framer-motion";
import { useEffect, useState } from "react";

interface AutoToggleProps {
  interval?: number;
  className?: string;
}

export function AutoToggle({ interval = 3000, className = "" }: AutoToggleProps) {
  const [isChecked, setIsChecked] = useState(false);

  useEffect(() => {
    const timer = setInterval(() => {
      setIsChecked((prev) => {
        const newState = !prev;
        return newState;
      });
    }, interval);

    return () => clearInterval(timer);
  }, [interval]);

  const activeColor = "#f97316";
  const inactiveColor = "#20f06cff";

  return (
    <div
      className={`relative aspect-[292/142] h-[30px] ${className}`}
      style={
        {
          "--active-color": activeColor,
          "--inactive-color": inactiveColor,
        } as React.CSSProperties
      }
    >
      <svg
        className="h-full w-full overflow-visible"
        viewBox="0 0 292 142"
        xmlns="http://www.w3.org/2000/svg"
      >
        <motion.path
          d="M71 142C31.7878 142 0 110.212 0 71C0 31.7878 31.7878 0 71 0C110.212 0 119 30 146 30C173 30 182 0 221 0C260 0 292 31.7878 292 71C292 110.212 260.212 142 221 142C181.788 142 173 112 146 112C119 112 110.212 142 71 142Z"
          animate={{
            fill: isChecked ? activeColor : inactiveColor,
          }}
          transition={{ duration: 0.4 }}
        />

        <motion.rect
          x="64"
          y="39"
          width="12"
          height="64"
          rx="6"
          animate={{
            fill: isChecked ? "hsl(var(--muted-foreground))" : inactiveColor,
          }}
          transition={{ duration: 0.4 }}
        />

        <motion.path
          fillRule="evenodd"
          d="M221 91C232.046 91 241 82.0457 241 71C241 59.9543 232.046 51 221 51C209.954 51 201 59.9543 201 71C201 82.0457 209.954 91 221 91ZM221 103C238.673 103 253 88.6731 253 71C253 53.3269 238.673 39 221 39C203.327 39 189 53.3269 189 71C189 88.6731 203.327 103 221 103Z"
          animate={{
            fill: isChecked ? activeColor : "hsl(var(--muted-foreground))",
          }}
          transition={{ duration: 0.4 }}
        />

        <g filter="url('#goo')">
          <motion.rect
            className="toggle-circle-center"
            y="42"
            width="116"
            height="58"
            rx="29"
            fill="hsl(var(--muted-foreground))"
            animate={{
              x: isChecked ? 163 : 13,
            }}
            transition={{ duration: 0.6, ease: "easeInOut" }}
          />
          <motion.rect
            className="toggle-circle-left"
            x="14"
            y="14"
            width="114"
            height="114"
            rx="58"
            fill="hsl(var(--muted-foreground))"
            style={{ transformOrigin: "center" }}
            animate={{
              scale: isChecked ? 0 : 1,
            }}
            transition={{ duration: 0.45, ease: "easeInOut" }}
          />
          <motion.rect
            className="toggle-circle-right"
            x="164"
            y="14"
            width="114"
            height="114"
            rx="58"
            fill="hsl(var(--muted-foreground))"
            style={{ transformOrigin: "center" }}
            animate={{
              scale: isChecked ? 1 : 0,
            }}
            transition={{ duration: 0.45, ease: "easeInOut" }}
          />
        </g>

        <defs>
          <filter id="goo">
            <feGaussianBlur in="SourceGraphic" result="blur" stdDeviation="10" />
            <feColorMatrix
              in="blur"
              mode="matrix"
              values="1 0 0 0 0  0 1 0 0 0  0 0 1 0 0  0 0 0 18 -7"
              result="goo"
            />
          </filter>
        </defs>
      </svg>
    </div>
  );
}

export default AutoToggle;