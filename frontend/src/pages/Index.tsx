import { SimulatorProvider } from "../context/SimulatorContext";
import MainBoard from "../components/MainBoard";
import { ScrollArea } from "@radix-ui/react-scroll-area";
import { Button } from "../components/ui/button";

const Index = () => {
  // TODO: Implement search functionality
  // const [globalSearchTerm, setGlobalSearchTerm] = React.useState("");

  return (
    <SimulatorProvider>
      <div className="min-h-screen bg-dark-bg text-foreground flex flex-col">
        <ScrollArea>
          {/* Main container adjusted for flex column */}
          <div className="flex-1 max-w-[1400px] w-full mx-auto p-4 px-6 md:px-6 lg:px-8 xl:px-10 pt-2 flex flex-col">
            <div className="my-4 p-3 rounded-lg bg-muted text-center text-xl font-bold text-primary-foreground shrink-0">
              TFT Simulator
            </div>
            <MainBoard />

            <div className="mt-4 p-1 rounded-lg bg-muted text-center text-xs text-muted-foreground shrink-0">
              <Button variant="ghost" size="sm" asChild>
                <a
                  href="https://discord.gg/your-server"
                  target="_blank"
                  rel="noopener noreferrer"
                  className="flex items-center gap-2 hover:text-foreground transition-colors"
                >
                  <img src="/discord.png" alt="Discord" className="w-4 h-4" />
                  Discord
                </a>
              </Button>
              {" • "}
              <Button variant="ghost" size="sm" asChild>
                <a
                  href="https://github.com/SuriZhang/tft-dps-simulator/"
                  target="_blank"
                  rel="noopener noreferrer"
                  className="flex items-center gap-2 hover:text-foreground transition-colors"
                >
                  <img src="/github.png" alt="GitHub" className="w-4 h-4" />
                  GitHub
                </a>
              </Button>
            </div>
          </div>
        </ScrollArea>
      </div>
    </SimulatorProvider>
  );
};

export default Index;
