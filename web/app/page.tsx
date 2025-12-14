import { Hero } from "@/components/marketing/hero";
import { Navbar } from "@/components/marketing/navbar";
import { PlaygroundSection } from "@/components/marketing/playground/playground-section";

export default function Home() {
  return (
    <>
      <Navbar />
      <Hero />
      <PlaygroundSection />
    </>
  );
}
