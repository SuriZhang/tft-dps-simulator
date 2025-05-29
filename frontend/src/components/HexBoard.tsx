import HexCell from "./HexCell";
import { BOARD_ROWS, BOARD_COLS } from "../utils/constants";
import { useSimulator } from "../context/SimulatorContext";
import { Trash2 } from "lucide-react";
import { useState, useEffect } from "react";

const HexBoard = () => {
  const { state } = useSimulator();
  const { boardChampions } = state;
  const [isDraggingChampion, setIsDraggingChampion] = useState(false);
  const [isMouseOutsideBoard, setIsMouseOutsideBoard] = useState(false);

  const hexWidth = 60;
  const hexHeight = hexWidth * 1.155; // Make it square (60x60)
  const spacing = 10;

  // Adjust spacing calculations for perfect honeycomb
  const horizontalSpacing = hexWidth + spacing;
  const verticalSpacing = (hexWidth + spacing) * 1.732 / 2;

  // Listen for custom events from HexCell
  useEffect(() => {
    const handleChampionDragStart = () => {
      setIsDraggingChampion(true);
      setIsMouseOutsideBoard(false);
    };
    const handleChampionDragEnd = () => {
      setIsDraggingChampion(false);
      setIsMouseOutsideBoard(false);
    };

    document.addEventListener('championDragStart', handleChampionDragStart);
    document.addEventListener('championDragEnd', handleChampionDragEnd);

    return () => {
      document.removeEventListener('championDragStart', handleChampionDragStart);
      document.removeEventListener('championDragEnd', handleChampionDragEnd);
    };
  }, []);

  const handleDragOver = (e: React.DragEvent) => {
    e.preventDefault();
    e.dataTransfer.dropEffect = "move";
  };

  const handleDrop = (e: React.DragEvent) => {
    e.preventDefault();
    setIsDraggingChampion(false);
    setIsMouseOutsideBoard(false);
    // The actual removal logic is handled in HexCell's handleDragEnd
  };

  const handleBoardDragEnter = (e: React.DragEvent) => {
    if (isDraggingChampion) {
      setIsMouseOutsideBoard(false);
    }
  };

  const handleBoardDragLeave = (e: React.DragEvent) => {
    if (isDraggingChampion) {
      // Check if we're actually leaving the board area (not just moving between child elements)
      const rect = e.currentTarget.getBoundingClientRect();
      const { clientX, clientY } = e;
      
      if (
        clientX < rect.left ||
        clientX > rect.right ||
        clientY < rect.top ||
        clientY > rect.bottom
      ) {
        setIsMouseOutsideBoard(true);
      }
    }
  };

  return (
    <div 
      className="relative w-full h-full p-4 bg-indigo-950/20 shadow-inner mt-4"
      onDragOver={handleDragOver}
      onDrop={handleDrop}
    >
      {/* Trash bin overlay when dragging champion outside board area */}
      {isDraggingChampion && isMouseOutsideBoard && (
        <div className="absolute inset-0 z-50 pointer-events-none">
          <div className="absolute inset-0 bg-red-500/20 animate-pulse rounded-lg">
            {/* Trash bin positioned in bottom right corner of outer div */}
            <div className="absolute bottom-2 right-2">
              <div className="bg-red-600/50 rounded-full p-2 animate-bounce">
                <Trash2 className="h-6 w-6 text-white" />
              </div>
            </div>
            {/* Instructions text positioned above trash bin */}
            <div className="absolute bottom-16 right-4">
              <div className="bg-black/80 text-white px-3 py-1 rounded-lg text-sm font-semibold whitespace-nowrap">
                Drop to remove
              </div>
            </div>
          </div>
        </div>
      )}

      <div
        className="relative mx-auto"
        style={{
          width: `${BOARD_COLS * horizontalSpacing + hexWidth / 2}px`,
          height: `${BOARD_ROWS * verticalSpacing + hexHeight / 4}px`,
        }}
        data-board-area="true"
        onDragEnter={handleBoardDragEnter}
        onDragLeave={handleBoardDragLeave}
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
                className="absolute transition-all duration-200 board-container"
                style={{ left: `${left}px`, top: `${top}px` }}
              >
                <HexCell row={row} col={col} champion={champion} />
              </div>
            );
          }),
        )}
      </div>
    </div>
  );
};

export default HexBoard;
