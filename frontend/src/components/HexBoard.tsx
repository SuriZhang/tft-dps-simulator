import React from "react";
import HexCell from "./HexCell";
import { BOARD_ROWS, BOARD_COLS } from "../utils/constants";
import { useSimulator } from "../context/SimulatorContext";
import { Button } from "./ui/button";

const HexBoard = () => {
  const { state } = useSimulator();
  const { boardChampions } = state;

  const hexWidth = 80; // Match the width in HexCell.tsx
  const hexHeight = hexWidth * (Math.sqrt(3) / 2); // Calculate height for point-topped hex

  // Add spacing between cells (adjust this value to control spacing)
  const spacing = 7; // Spacing in pixels

  // Adjust spacing calculations
  const horizontalSpacing = hexWidth + spacing;
  const verticalSpacing = hexHeight + 0.5*spacing;

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
      <div className="absolute bottom-1 right-4">
        <Button
          variant="outline"
          className=" bg-primary/20 text-primary hover:bg-primary/30"
        >
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
        </Button>
      </div>
    </div>
  );
};

export default HexBoard;
