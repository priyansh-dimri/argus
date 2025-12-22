"use client";

import {
  Accordion,
  AccordionContent,
  AccordionItem,
  AccordionTrigger,
} from "@/components/ui/accordion";
import { Card } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { adrs } from "./adr-data";

export function ADRAccordion() {
  return (
    <Accordion type="single" collapsible className="space-y-4">
      {adrs.map((adr, idx) => (
        <AccordionItem
          key={adr.id}
          value={adr.id}
          className="glass border border-white/10 rounded-xl overflow-hidden"
        >
          <AccordionTrigger className="px-6 py-4 hover:bg-white/5 transition-colors">
            <div className="flex items-center gap-3 text-left">
              <Badge variant="outline" className="shrink-0">
                ADR-{idx + 1}
              </Badge>
              <span className="font-semibold">{adr.title}</span>
            </div>
          </AccordionTrigger>
          <AccordionContent className="px-6 py-4 text-sm">
            <div className="space-y-4">
              <div>
                <h4 className="font-semibold text-foreground mb-2">Problem</h4>
                <p className="text-muted-foreground">{adr.problem}</p>
              </div>

              <div>
                <h4 className="font-semibold text-foreground mb-2">
                  Alternatives Considered
                </h4>
                <ul className="space-y-2">
                  {adr.alternatives.map((alt, i) => (
                    <li key={i} className="flex gap-3">
                      <span className="text-muted-foreground shrink-0">•</span>
                      <div>
                        <span className="font-medium text-foreground">
                          {alt.name}:
                        </span>{" "}
                        <span className="text-muted-foreground">
                          {alt.tradeoff}
                        </span>
                      </div>
                    </li>
                  ))}
                </ul>
              </div>

              <div>
                <h4 className="font-semibold text-foreground mb-2">Decision</h4>
                <p className="text-muted-foreground">{adr.decision}</p>
              </div>

              <div>
                <h4 className="font-semibold text-foreground mb-2">
                  Consequences
                </h4>
                <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                  <Card className="bg-green-500/10 border-green-500/20 p-3">
                    <p className="text-xs font-semibold text-green-400 mb-2">
                      Positive
                    </p>
                    <ul className="space-y-1">
                      {adr.consequences.positive.map((pos, i) => (
                        <li
                          key={i}
                          className="text-xs text-muted-foreground flex gap-2"
                        >
                          <span className="text-green-400">✓</span>
                          {pos}
                        </li>
                      ))}
                    </ul>
                  </Card>
                  <Card className="bg-orange-500/10 border-orange-500/20 p-3">
                    <p className="text-xs font-semibold text-orange-400 mb-2">
                      Trade-offs
                    </p>
                    <ul className="space-y-1">
                      {adr.consequences.negative.map((neg, i) => (
                        <li
                          key={i}
                          className="text-xs text-muted-foreground flex gap-2"
                        >
                          <span className="text-orange-400">⚠</span>
                          {neg}
                        </li>
                      ))}
                    </ul>
                  </Card>
                </div>
              </div>
            </div>
          </AccordionContent>
        </AccordionItem>
      ))}
    </Accordion>
  );
}
