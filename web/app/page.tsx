import { Footer } from "@/components/layout/footer";
import { BentoGrid } from "@/components/marketing/bento/bento-grid";
import { Hero } from "@/components/marketing/hero";
import { IntegrationSection } from "@/components/marketing/integration/integration-section";
import { Navbar } from "@/components/marketing/navbar";
import { PlaygroundSection } from "@/components/marketing/playground/playground-section";

export default function Home() {
  return (
    <>
      <Navbar />
      <Hero />
      <PlaygroundSection />
      <BentoGrid />
      <IntegrationSection />
      <Footer />
    </>
  );
}
