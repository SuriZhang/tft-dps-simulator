import HexCell from "./HexCell";
import { BOARD_ROWS, BOARD_COLS } from "../utils/constants";
import { useSimulator } from "../context/SimulatorContext";

const HexBoard = () => {
  const { state } = useSimulator();
  const { boardChampions } = state;

  const hexWidth = 60;
  const hexHeight = hexWidth * 1.155; // Make it square (60x60)
  const spacing = 10;

  // Adjust spacing calculations for perfect honeycomb
  const horizontalSpacing = hexWidth + spacing;
  const verticalSpacing = (hexWidth + spacing) * 1.732 / 2;

  return (
    <div className="relative w-full h-full p-4 bg-indigo-950/20 shadow-inner mt-4">
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
                className="absolute transition-all duration-200 board-container"
                style={{ left: `${left}px`, top: `${top}px` }}
                data-board-area="true"
              >
                <HexCell row={row} col={col} champion={champion} />
              </div>
            );
          }),
        )}
      </div>

      

      {/* Combat button
      <div className="absolute bottom-0 right-0">
        <Button
          variant="outline"
          className="bg-primary/20 text-primary hover:bg-primary/30 p-2"
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
                className="h-4 w-4 bg-primary"
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
      </div> */}
    </div>
  );
};

export default HexBoard;
