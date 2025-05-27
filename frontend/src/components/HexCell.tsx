// import React from "react";
import { BoardPosition, Champion, Item } from "../utils/types";
import { useSimulator } from "../context/SimulatorContext";
import { cn } from "../lib/utils";
import {
  ContextMenu,
  ContextMenuContent,
  ContextMenuItem,
  ContextMenuSeparator,
  ContextMenuTrigger,
} from "./ui/context-menu";
import { Star, Trash2 } from "lucide-react";
import {
  Tooltip,
  TooltipContent,
  TooltipProvider,
  TooltipTrigger,
} from "./ui/tooltip";

interface HexCellProps {
  row: number;
  col: number;
  champion?: Champion | null;
}

const HexCell: React.FC<HexCellProps> = ({ row, col, champion }) => {
  const position: BoardPosition = { row, col };
  const { state, dispatch } = useSimulator();
  const { selectedChampion, selectedItem, hoveredTrait } = state;

  // Check if the champion has the hovered trait - add debugging to see values
  const hasHoveredTrait =
    champion &&
    hoveredTrait &&
    champion.traits &&
    champion.traits.some((trait) => {
      // Case-insensitive comparison or normalized comparison
      return trait.toLowerCase() === hoveredTrait.toLowerCase();
    });

  const getHexBackground = () => {
    let base = "bg-gray-900/40 border border-gray-700/30";
    if (champion) {
      switch (champion.cost) {
        case 1:
          base = "bg-gray-800/60 border border-gray-600/80";
          break;
        case 2:
          base = "bg-green-900/60 border border-green-600/50";
          break;
        case 3:
          base = "bg-blue-900/60 border border-blue-500/50";
          break;
        case 4:
          base = "bg-purple-900/60 border border-purple-500/50";
          break;
        case 5:
          base = "bg-amber-900/60 border border-amber-500/50";
          break;
      }
    }
    return base;
  };

  // Click / drag handlers
  const handleCellClick = () => {
    if (selectedChampion && !champion) {
      dispatch({
        type: "ADD_CHAMPION_TO_BOARD",
        champion: selectedChampion,
        position,
      });
      console.log("Adding champion to board", selectedChampion, position);
      // Deselect the champion after placing it
      dispatch({ type: "SELECT_CHAMPION", champion: undefined });
      console.log("deselecting:", selectedChampion, position);
    } else if (champion && selectedItem) {
      dispatch({
        type: "ADD_ITEM_TO_CHAMPION",
        item: selectedItem,
        position,
      });
      console.log("Adding item to champion", selectedItem, position);
      console.log("Champion", champion);
    }
  };

  const handleSetStarLevel = (level: number) => () => {
    champion &&
      dispatch({
        type: "SET_CHAMPION_STAR_LEVEL",
        position,
        level,
      });
  };

  const handleRemove = () =>
    champion && dispatch({ type: "REMOVE_CHAMPION_FROM_BOARD", position });

  const handleDragOver = (e: React.DragEvent) => e.preventDefault();

  const handleDragStart = (e: React.DragEvent) => {
    if (champion) {
      e.dataTransfer.setData(
        "application/json",
        JSON.stringify({ type: "boardChampion", position }),
      );
    }
  };
  const handleDrop = (e: React.DragEvent) => {
    e.preventDefault();
    try {
      const data = JSON.parse(e.dataTransfer.getData("application/json"));
      if (data.type === "champion" && !champion) {
        dispatch({
          type: "ADD_CHAMPION_TO_BOARD",
          champion: data.champion,
          position,
        });
        console.log("Adding champion to board", data.champion, position);
      } else if (data.type === "boardChampion" && data.position) {
        dispatch({
          type: "MOVE_CHAMPION",
          from: data.position,
          to: position,
        });
        console.log("Moving champion", data.position, position);
      } else if (data.type === "item" && champion) {
        dispatch({
          type: "ADD_ITEM_TO_CHAMPION",
          item: data.item,
          position,
        });
        console.log("Adding item to champion", data.item, position);
        console.log("Champion", champion);
      }
    } catch {
      /* ignore parsing errors */
    }
  };

  return (
    <div
      className={cn(
        "relative aspect-[1/1] cursor-pointer",
        `col-start-${col} rows-start-${row}`,
        // Move the ring styling here to the OUTER container instead
        hasHoveredTrait && "z-10",
      )}
    >
      {/* Stars positioned outside the hexagon but visually on top */}
      {champion && champion.stars && (
        <div className="absolute top-1 left-0 w-full flex justify-center gap-0.5 z-30">
          {Array.from({ length: champion.stars }).map((_, i) => (
            <Star
              key={i}
              fill="yellow"
              className="h-4 w-4 text-yellow-400 drop-shadow-[0_0_2px_rgba(0,0,0,0.8)]"
            />
          ))}
        </div>
      )}

      {/* Items positioned outside the hexagon at the bottom */}
      {champion && champion.items && champion.items.length > 0 && (
        <div className="absolute bottom-1 left-0 w-full flex justify-center gap-0.5 z-30">
          {champion.items.map((item: Item, i: number) => (
            <img
              key={i}
              src={`/tft-item/${item.icon}`}
              alt={item.name}
              className="w-5 h-5 object-cover rounded-sm border border-gray-800 drop-shadow-[0_0_2px_rgba(0,0,0,0.8)]"
              // title={item.name} // Item tooltips can be handled separately if needed, or removed if champion tooltip covers it
            />
          ))}
        </div>
      )}

      {/* Add a visible ring effect on the outer div when trait is matched */}
      {hasHoveredTrait && (
        <div
          className="absolute inset-0 rounded-md ring-3 ring-primary animate-pulse"
          style={{
            width: "90px",
            height: "90px",
            left: "-5px",
            top: "-5px",
            zIndex: 5,
          }}
        ></div>
      )}

      <TooltipProvider delayDuration={100}>
        <Tooltip>
          <TooltipTrigger asChild>
<div className={cn(
  "hexagon-border", 
  champion ? `cost-${champion.cost}-border` : "empty-hex-border"
)}>
              <div
                className={cn(
                  "w-[80px] h-[80px] inset-0 clip-hexagon shadow-md transition-all",
                  getHexBackground(),
                  !champion && selectedChampion
                    ? "border-primary border-3 hover:border-opacity-100"
                    : "",
                  champion && selectedItem
                    ? "border-accent border-3 hover:border-opacity-100"
                    : "",
                  hoveredTrait && champion && !hasHoveredTrait && "opacity-40",
                )}
                onClick={handleCellClick}
                onDragOver={handleDragOver}
                onDrop={handleDrop}
                data-position={`${position.row}-${position.col}`}
                style={{ aspectRatio: "1" }}
              >
                {/* If champion has the hovered trait, add a glow effect */}
                {hasHoveredTrait && (
                  <div className="absolute inset-0 bg-primary/20 clip-hexagon animate-pulse z-10"></div>
                )}

                {champion && (
                  <ContextMenu>
                    <ContextMenuTrigger asChild>
                      <div
                        className="absolute w-full h-full"
                        draggable
                        onDragStart={handleDragStart}
                      >
                        {/* Background champion image/name */}
                        {champion.icon && (
                          <div className="absolute inset-0 flex items-center justify-center overflow-hidden">
                            <img
                              src={`/tft-champion-icons/${champion.icon.toLowerCase()}`}
                              alt={champion.name}
                              className="w-full h-full object-cover rotate-[-90deg]"
                            />
                          </div>
                        )}
                      </div>
                    </ContextMenuTrigger>

                    <ContextMenuContent>
                      <ContextMenuItem onClick={handleSetStarLevel(1)}>
                        <Star className="mr-2 h-4 w-4" />
                        1-Star
                      </ContextMenuItem>
                      <ContextMenuItem onClick={handleSetStarLevel(2)}>
                        <Star className="mr-2 h-4 w-4" />
                        2-Star
                      </ContextMenuItem>
                      <ContextMenuItem onClick={handleSetStarLevel(3)}>
                        <Star className="mr-2 h-4 w-4" />
                        3-Star
                      </ContextMenuItem>
                      <ContextMenuSeparator />
                      <ContextMenuItem
                        onClick={handleRemove}
                        className="text-destructive"
                      >
                        <Trash2 className="mr-2 h-4 w-4" />
                        Remove
                      </ContextMenuItem>
                    </ContextMenuContent>
                  </ContextMenu>
                )}
              </div>
            </div>
          </TooltipTrigger>
          <TooltipContent>
            {champion ? (
              <>
                <p className="font-bold">{champion.name}</p>
                {champion.traits && champion.traits.length > 0 && (
                  <p className="text-xs text-muted-foreground">
                    {champion.traits.join(", ")}
                  </p>
                )}
                <p className="text-xs text-muted-foreground italic">
                  Click to view champion details
                </p>
              </>
            ) : (
              <p>{`Row ${row}, Col ${col}`}</p>
            )}
          </TooltipContent>
        </Tooltip>
      </TooltipProvider>
    </div>
  );
};

export default HexCell;
