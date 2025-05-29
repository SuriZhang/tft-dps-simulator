import { useState } from "react";
import { useSimulator } from "../context/SimulatorContext";
import { Button } from "./ui/button";
import { Trash2, Share2, Copy, UploadCloud, Loader2 } from "lucide-react";
import {
  Tooltip,
  TooltipContent,
  TooltipProvider,
  TooltipTrigger,
} from "./ui/tooltip";

const ControlBar = () => {
  const { state, dispatch } = useSimulator();

  const { boardChampions } = state;
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  // Handle clear board action
  const handleClearBoard = () => {
    dispatch({ type: "CLEAR_BOARD" });
  };

  // Function to handle the combat simulation
  const handleStartCombat = async () => {
    // Reset state
    setIsLoading(true);
    setError(null);

    try {
      // Call the mock endpoint
      const response = await fetch(`/api/v1/simulation/run`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          boardChampions: boardChampions.map((champion) => ({
            apiName: champion.apiName,
            stars: champion.stars,
            position: champion.position,
            items: champion.items || [],
          })),
        }),
      });

      // console request body
      console.log("Request body:", {
        boardChampions: boardChampions.map((champion) => ({
          apiName: champion.apiName,
          stars: champion.stars,
          position: champion.position,
          items: champion.items || [],
        })),
      });

      if (!response.ok) {
        throw new Error(`Server responded with status: ${response.status}`);
      }

      const data = await response.json();

      // Update context with simulation results and events
      dispatch({
        type: "SET_SIMULATION_DATA",
        payload: {
          results: data.results,
          events: data.archieveEvents || [], // Note: keeping the same spelling as backend
        },
      });

      // Optionally: Show a success message or navigate to results view
      console.log("Simulation completed successfully:", {
        results: data.results,
        events: data.archieveEvents,
      });
    } catch (err) {
      console.error("Error running simulation:", err);
      setError(
        err instanceof Error ? err.message : "An unknown error occurred",
      );
    } finally {
      setIsLoading(false);
    }
  };

  // Placeholder handlers
  const handleShare = () => console.log("Share TBD");
  const handleCopyCode = () => console.log("Copy Code TBD");
  const handleImportCode = () => console.log("Import Code TBD");

  return (
    <div className="bg-card shadow-none rounded-lg  p-4 mb-4">
      <div className="flex flex-wrap justify-between items-center gap-4">
        {/* Set info */}
        <div className="flex items-center space-x-4">
          {/* Consider using Badge component here if appropriate */}
          <div className="bg-accent rounded-md px-3 py-1">
            <span className="text-white font-medium">SET 14</span>
          </div>
          <div className="text-gray-400 text-sm">
            SET 13 {/* Potentially a Button or Link */}
          </div>
        </div>
        <div className="space-x-2">
          <TooltipProvider delayDuration={100}>
            <Tooltip>
              <TooltipTrigger asChild>
                <Button variant="outline" size="sm" onClick={handleCopyCode}>
                  <Copy className="h-4 w-4 mr-1" />
                </Button>
              </TooltipTrigger>
              <TooltipContent>Copy Code</TooltipContent>
            </Tooltip>
          </TooltipProvider>

          <TooltipProvider delayDuration={100}>
            <Tooltip>
              <TooltipTrigger asChild>
                <Button variant="outline" size="sm" onClick={handleImportCode}>
                  <UploadCloud className="h-4 w-4 mr-1" />
                </Button>
              </TooltipTrigger>
              <TooltipContent>Import Code</TooltipContent>
            </Tooltip>
          </TooltipProvider>

          <TooltipProvider delayDuration={100}>
            <Tooltip>
              <TooltipTrigger asChild>
                <Button variant="outline" size="sm" onClick={handleShare}>
                  <Share2 className="h-4 w-4 mr-1" />
                </Button>
              </TooltipTrigger>
              <TooltipContent>Share Code</TooltipContent>
            </Tooltip>
          </TooltipProvider>

          <Button variant="destructive" size="sm" onClick={handleClearBoard}>
            <Trash2 className="h-4 w-4 mr-1" /> CLEAR BOARD
          </Button>

          {error && (
            <div className="absolute bottom-12 right-4 text-red-500 bg-red-100 p-2 rounded">
              {error}
            </div>
          )}
          <Button
            variant="outline"
            className="bg-emerald-700 text-white hover:bg-emerald-600 p-2"
            onClick={handleStartCombat}
            disabled={isLoading || boardChampions.length === 0}
          >
            {isLoading ? (
              <>
                <Loader2 className="h-5 w-5 mr-2 animate-spin" />
                Simulating...
              </>
            ) : (
              <>
                <div
                  className="h-4 w-4 bg-white"
                  style={{
                    WebkitMaskImage: "url(./TFTM_ModeIcon_Normal.png)",
                    maskImage: "url(./TFTM_ModeIcon_Normal.png)",
                    WebkitMaskSize: "contain",
                    maskSize: "contain",
                    WebkitMaskRepeat: "no-repeat",
                    maskRepeat: "no-repeat",
                    WebkitMaskPosition: "center",
                    maskPosition: "center",
                  }}
                ></div>
                START COMBAT
              </>
            )}
          </Button>
        </div>
      </div>
    </div>
  );
};

export default ControlBar;
