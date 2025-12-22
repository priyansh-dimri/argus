import { Card } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";

interface TechStackCardProps {
  category: string;
  stack: string[];
}

export function TechStackCard({ category, stack }: TechStackCardProps) {
  return (
    <Card className="glass border border-white/10 p-6">
      <h3 className="text-lg font-semibold mb-4">{category}</h3>
      <div className="flex flex-wrap gap-2">
        {stack.map((tech) => (
          <Badge
            key={tech}
            variant="secondary"
            className="bg-white/5 hover:bg-white/10 border-white/10"
          >
            {tech}
          </Badge>
        ))}
      </div>
    </Card>
  );
}
