import React from "react";
import { SimulatorProvider } from "../context/SimulatorContext";
import MainBoard from "../components/MainBoard";
import { Scroll } from "lucide-react";
import { ScrollArea } from "@radix-ui/react-scroll-area";

const Index = () => {
	// TODO: Implement search functionality
	const [globalSearchTerm, setGlobalSearchTerm] = React.useState("");

	return (
		<SimulatorProvider>
      <div className="min-h-screen bg-dark-bg text-foreground flex flex-col">
        <ScrollArea >
				{/* Main container adjusted for flex column */}
				<div className="flex-1 max-w-[1800px] w-full mx-auto p-4 pt-2 flex flex-col">
					<div className="my-4 p-3 rounded-lg bg-muted text-center text-xl font-bold text-muted-foreground shrink-0">
						TFT Simulator
					</div>
					<MainBoard />

					<div className="mt-4 p-1 rounded-lg bg-muted text-center text-xs text-muted-foreground shrink-0">
						<span className="mr-3">
							Press{" "}
							<kbd className="px-1.5 py-0.5 text-[10px] font-semibold text-foreground bg-background rounded border">
								Shift
							</kbd>{" "}
							+ Click to place 2-star units
						</span>
						<span className="mr-3">
							Press{" "}
							<kbd className="px-1.5 py-0.5 text-[10px] font-semibold text-foreground bg-background rounded border">
								Ctrl
							</kbd>{" "}
							+ Click to place 3-star units
						</span>
						<span>
							Press{" "}
							<kbd className="px-1.5 py-0.5 text-[10px] font-semibold text-foreground bg-background rounded border">
								Alt
							</kbd>{" "}
							+ Click to swap units
						</span>
					</div>
          </div>
          </ScrollArea>
			</div>
		</SimulatorProvider>
	);
};

export default Index;
