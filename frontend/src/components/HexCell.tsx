import React from "react";
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

interface HexCellProps {
  row: number;
  col: number;
  champion?: Champion | null;
}

const HexCell: React.FC<HexCellProps> = ({ row, col, champion }) => {
  const position: BoardPosition = { row, col };
  const { state, dispatch } = useSimulator();
  const { selectedChampion, selectedItem } = state;

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

  if (champion) {
    console.log(champion);
  }

  return (
    <div
      className={cn(
        "relative aspect-[1/1]",
        `col-start-${col} rows-start-${row}`,
      )}
    >
      <div
        className={cn(
          "w-[80px] h-[80px] inset-0 clip-hexagon shadow-md transition-all cursor-pointer",
          getHexBackground(),
          !champion && selectedChampion
            ? "border-primary border-3 hover:border-opacity-100"
            : "",
          champion && selectedItem
            ? "border-accent border-3 hover:border-opacity-100"
            : "",
        )}
        title={`Row ${row}, Col ${col}`}
        onClick={handleCellClick}
        onDragOver={handleDragOver}
        onDrop={handleDrop}
        data-position={`${position.row}-${position.col}`}
      >
        {champion && (
          <ContextMenu>
            <ContextMenuTrigger asChild>
              <div
                className="absolute w-full h-full"
                draggable
                onDragStart={handleDragStart}
              >
                {/* Background champion image/name */}
                {champion.icon ? (
                  <div className="absolute inset-0 flex items-center justify-center overflow-hidden">
                    <img
                      src={`/tft-champion/${champion.icon.toLowerCase()}`}
                      alt={champion.name}
                      className="w-full h-full object-cover rotate-[-90deg]"
                      style={{ objectPosition: "60% 45%" }}
                      title={champion.name}
                    />
                  </div>
                ) : (
                  <div className="absolute inset-0 flex items-center justify-center">
                    <p className="text-xs font-bold text-white">
                      {champion.name}
                    </p>
                  </div>
                )}

                {/* Overlay elements (stars and items) */}
                <div className="absolute w-full h-full flex flex-col items-center justify-center rotate-[-90deg]">
                  <div className="absolute top-1 w-full flex justify-center gap-0 z-20">
                    {champion.stars && Array.from({ length : champion.stars }).map((_, i) => (
                      <Star
                      key={i}
                      fill="yellow"
                      className="h-4 w-4 text-yellow-400"
                      />
                    ))}
                  </div>

                  {/* items at the bottom */}
                  {champion.items && champion.items.length > 0 && (
                    <div className="absolute bottom-2 w-full flex justify-center gap-0.5 z-20">
                      {champion.items.map((item: Item, i: number) => (
                        <img
                          key={i}
                          src={`/tft-item/${item.icon}`}
                          alt={item.name}
                          className="w-4 h-4 rounded-sm object-cover border-2"
                          title={item.name}
                        />
                      ))}
                    </div>
                  )}
                </div>
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
  );
};

export default HexCell;
