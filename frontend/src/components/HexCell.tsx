import React, { useState } from "react";
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
  const [isDragOver, setIsDragOver] = useState(false);
  const [isChampionBeingDragged, setIsChampionBeingDragged] = useState(false);

  // Listen for champion drag events from other cells
  React.useEffect(() => {
    const handleChampionDragStart = () => setIsChampionBeingDragged(true);
    const handleChampionDragEnd = () => setIsChampionBeingDragged(false);

    document.addEventListener('championDragStart', handleChampionDragStart);
    document.addEventListener('championDragEnd', handleChampionDragEnd);

    return () => {
      document.removeEventListener('championDragStart', handleChampionDragStart);
      document.removeEventListener('championDragEnd', handleChampionDragEnd);
    };
  }, []);

  // Check if the champion has the hovered trait
  const hasHoveredTrait =
    champion &&
    hoveredTrait &&
    champion.traits &&
    champion.traits.some((trait) => {
      return trait.toLowerCase() === hoveredTrait.toLowerCase();
    });

  // Click / drag handlers
  const handleCellClick = () => {
    if (selectedChampion && !champion) {
      dispatch({
        type: "ADD_CHAMPION_TO_BOARD",
        champion: selectedChampion,
        position,
      });
      console.log("Adding champion to board", selectedChampion, position);
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

  const handleDragOver = (e: React.DragEvent) => {
    e.preventDefault();
    setIsDragOver(true);
  };

  const handleDragEnter = (e: React.DragEvent) => {
    e.preventDefault();
    setIsDragOver(true);
  };

  const handleDragLeave = (e: React.DragEvent) => {
    e.preventDefault();
    if (!e.currentTarget.contains(e.relatedTarget as Node)) {
      setIsDragOver(false);
    }
  };

  const handleDragStart = (e: React.DragEvent) => {
    if (champion) {
      e.dataTransfer.setData(
        "application/json",
        JSON.stringify({ type: "boardChampion", position }),
      );
      e.dataTransfer.effectAllowed = "move";
      
      document.dispatchEvent(new CustomEvent('championDragStart'));
    }
  };

  const handleDragEnd = (e: React.DragEvent) => {
    document.dispatchEvent(new CustomEvent('championDragEnd'));
    
    setIsDragOver(false);
    
    const boardElement = document.querySelector('[data-board-area="true"]');
    if (boardElement && champion) {
      const rect = boardElement.getBoundingClientRect();
      const { clientX, clientY } = e;

      if (
        clientX < rect.left ||
        clientX > rect.right ||
        clientY < rect.top ||
        clientY > rect.bottom
      ) {
        dispatch({ type: "REMOVE_CHAMPION_FROM_BOARD", position });
        console.log("Removing champion dragged outside board", champion, position);
      }
    }
  };

  const handleDrop = (e: React.DragEvent) => {
    e.preventDefault();
    setIsDragOver(false);
    
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

  const handleItemClick = (e: React.MouseEvent, item: Item) => {
    e.stopPropagation(); // Prevent cell click from triggering
    dispatch({
      type: "REMOVE_ITEM_FROM_CHAMPION",
      position,
      itemApiName: item.apiName,
    });
    console.log("Removing item by click", item, position);
  };

  // Show highlight when a champion is being dragged over an empty cell
  const shouldShowDragHighlight = isChampionBeingDragged && !champion && isDragOver;

  return (
    <div
      className={cn(
        "relative aspect-[1/1]",
        `col-start-${col} rows-start-${row}`,
        hasHoveredTrait && "z-10",
        champion ? "cursor-pointer" : "cursor-default"
      )}
    >
      {/* Stars positioned outside the hexagon but visually on top */}
      {champion && champion.stars && (
        <div className="absolute top-1 left-0 w-full flex justify-center gap-0.5 z-30">
          {Array.from({ length: champion.stars }).map((_, i) => (
            <Star
              key={i}
              fill="yellow"
              className="h-3 w-3 text-yellow-400 drop-shadow-[0_0_2px_rgba(0,0,0,0.8)]"
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
              className="w-5 h-5 object-cover rounded-sm border border-gray-800 drop-shadow-[0_0_2px_rgba(0,0,0,0.8)] cursor-pointer hover:brightness-110 transition-all"
              onClick={(e) => handleItemClick(e, item)}
              title={`Click to remove ${item.name}`}
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

      {champion ? (
        <TooltipProvider delayDuration={100}>
          <Tooltip>
            <TooltipTrigger asChild>
              <div className={cn(
                "hexagon-border",
                champion ? `cost-${champion.cost}-border` : "empty-hex-border"
              )}>
                <div
                  className={cn(
                    "w-[60px] h-[69.3px] inset-0 clip-hexagon shadow-md transition-all bg-secondary",
                    !champion && selectedChampion
                      ? "border-primary border-3 hover:border-opacity-100"
                      : "",
                    champion && selectedItem
                      ? "border-accent border-3 hover:border-opacity-100"
                      : "",
                    // Add bright background when dragging champion over empty cell
                    shouldShowDragHighlight && "!bg-gray-400/60"
                  )}
                  onClick={handleCellClick}
                  onDragOver={handleDragOver}
                  onDragEnter={handleDragEnter}
                  onDragLeave={handleDragLeave}
                  onDrop={handleDrop}
                  data-position={`${position.row}-${position.col}`}
                  style={{
                    opacity: hoveredTrait && champion && !hasHoveredTrait ? 0.6 : 1
                  }}
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
                          onDragEnd={handleDragEnd}
                        >
                          {/* Background champion image/name */}
                          {champion.icon && (
                            <div className="absolute inset-0 flex items-center justify-center overflow-hidden">
                              <img
                                src={`/tft-champion-icons/${champion.icon.toLowerCase()}`}
                                alt={champion.name}
                                className="w-full h-full object-cover"
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
              <p className="font-bold">{champion.name}</p>
              {champion.traits && champion.traits.length > 0 && (
                <p className="text-xs text-muted-foreground">
                  {champion.traits.join(", ")}
                </p>
              )}
              <p className="text-xs text-muted-foreground italic">
                Click to view champion details
              </p>
            </TooltipContent>
          </Tooltip>
        </TooltipProvider>
      ) : (
        <div className={cn(
          "hexagon-border empty-hex-border"
        )}>
          <div
            className={cn(
              "w-[60px] h-[69.3px] inset-0 clip-hexagon shadow-md transition-all bg-secondary",
              !champion && selectedChampion
                ? "border-primary border-3 hover:border-opacity-100"
                : "",
              champion && selectedItem
                ? "border-accent border-3 hover:border-opacity-100"
                : "",
              // Add bright background when dragging champion over empty cell
              shouldShowDragHighlight && "!bg-gray-400/60"
            )}
            onClick={handleCellClick}
            onDragOver={handleDragOver}
            onDragEnter={handleDragEnter}
            onDragLeave={handleDragLeave}
            onDrop={handleDrop}
            data-position={`${position.row}-${position.col}`}
            style={{
              opacity: hoveredTrait && champion && !hasHoveredTrait ? 0.6 : 1
            }}
          >
          </div>
        </div>
      )}
    </div>
  );
};

export default HexCell;
