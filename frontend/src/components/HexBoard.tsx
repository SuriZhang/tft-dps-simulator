import React, { useState } from "react";
import HexCell from "./HexCell";
import { BOARD_ROWS, BOARD_COLS } from "../utils/constants";
import { useSimulator } from "../context/SimulatorContext";
import { Button } from "./ui/button";
import { Loader2 } from "lucide-react"; // Import a loading icon

const HexBoard = () => {
  const { state, dispatch } = useSimulator();
  const { boardChampions } = state;
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);


  // const fetchData = () => {
  //   fetch(`http://localhost:${import.meta.env.VITE_PORT}/`)
  //     .then((response) => response.text())
  //     .then((data) => setMessage(data))
  //     .catch((error) => console.error("Error fetching data:", error));
  // };

  // Function to handle the combat simulation
  const handleStartCombat = async () => {
    // Reset state
    setIsLoading(true);
    setError(null);

    try {
      // Call the mock endpoint
      const response = await fetch(
        `http://localhost:${import.meta.env.VITE_PORT}/api/v1/simulation/mock-run`,
        {
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
        },
      );

      if (!response.ok) {
        throw new Error(`Server responded with status: ${response.status}`);
      }

      const data = await response.json();

      // Update context with simulation results
      dispatch({
        type: "SET_SIMULATION_RESULTS",
        payload: data.results,
      });

      // Optionally: Show a success message or navigate to results view
      console.log("Simulation completed successfully:", data);
    } catch (err) {
      console.error("Error running simulation:", err);
      setError(
        err instanceof Error ? err.message : "An unknown error occurred",
      );
    } finally {
      setIsLoading(false);
    }
  };

  const hexWidth = 80; // Match the width in HexCell.tsx
  const hexHeight = hexWidth * (Math.sqrt(3) / 2); // Calculate height for point-topped hex

  // Add spacing between cells (adjust this value to control spacing)
  const spacing = 7; // Spacing in pixels

  // Adjust spacing calculations
  const horizontalSpacing = hexWidth + spacing;
  const verticalSpacing = hexHeight + 0.5 * spacing;

  return (
    <div className="relative w-full h-full p-4 bg-indigo-950/20 shadow-inner">
      <div
        className="relative mx-auto"
        style={{
          width: `${BOARD_COLS * horizontalSpacing + hexWidth / 2}px`,
          height: `${BOARD_ROWS * verticalSpacing + hexHeight / 4}px`,
        }}
      >
        {Array.from({ length: BOARD_ROWS }).map((_, row) =>
          Array.from({ length: BOARD_COLS }).map((_, col) => {
            const champion = boardChampions.find(
              (c) => c.position.row === row && c.position.col === col,
            );

            // Calculate position for honeycomb layout with spacing
            const xOffset = (row % 2) * (horizontalSpacing / 2);
            const left = col * horizontalSpacing + xOffset;
            const top = row * verticalSpacing;

            return (
              <div
                key={`${row}-${col}`}
                className="absolute transition-all duration-200"
                style={{ left: `${left}px`, top: `${top}px` }}
              >
                <HexCell row={row} col={col} champion={champion} />
              </div>
            );
          }),
        )}
      </div>

      {/* Error message display */}
      {error && (
        <div className="absolute bottom-12 right-4 text-red-500 bg-red-100 p-2 rounded">
          {error}
        </div>
      )}

      {/* Combat button */}
      <div className="absolute bottom-1 right-4">
        <Button
          variant="outline"
          className="bg-primary/20 text-primary hover:bg-primary/30"
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
                className="h-6 w-6 mr-2 bg-primary"
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
              Start Combat
            </>
          )}
        </Button>
      </div>
    </div>
  );
};

export default HexBoard;
