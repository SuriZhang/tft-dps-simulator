import React from "react";
import HexCell from "./HexCell";
import { BOARD_ROWS, BOARD_COLS } from "../utils/constants";
import { useSimulator } from "../context/SimulatorContext";
import { cn } from "../lib/utils";
import { Button } from "./ui/button";
import { Play, Plus } from "lucide-react";

const HexBoard = () => {
  const { state } = useSimulator();
  const { boardChampions } = state;

  const hexWidth = 80; // Match the width in HexCell.tsx
  const hexHeight = hexWidth * (Math.sqrt(3) / 2); // Calculate height for point-topped hex
  const verticalSpacing = hexHeight * 1;

  return (
    <div className="relative w-full h-full p-4 bg-indigo-950/20 shadow-inner">
      <div
        className="relative mx-auto"
        style={{
          width: `${BOARD_COLS * hexWidth + hexWidth / 2}px`,
          height: `${BOARD_ROWS * verticalSpacing + hexHeight / 4}px`,
        }}
      >
        {Array.from({ length: BOARD_ROWS }).map((_, row) =>
          Array.from({ length: BOARD_COLS }).map((_, col) => {
            const champion = boardChampions.find(
              (c) => c.position.row === row && c.position.col === col,
            );

            // Calculate position for honeycomb layout
            const xOffset = (row % 2) * (hexWidth / 2);
            const left = col * hexWidth + xOffset;
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
      <div className="absolute bottom-4 right-4">
        <Button
          variant="outline"
          className=" bg-primary/20 text-primary hover:bg-primary/30"
        >
          <Play className="h-4 w-4" /> Start Combat
        </Button>
      </div>
    </div>
  );
};

export default HexBoard;
