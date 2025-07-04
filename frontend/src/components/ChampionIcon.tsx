import React from "react";
import { Champion } from "../utils/types";
import { cn } from "../lib/utils";
import { useSimulator } from "../context/SimulatorContext";
import { Avatar, AvatarFallback, AvatarImage } from "./ui/avatar"; // Import Avatar components
import {
  Tooltip,
  TooltipContent,
  TooltipProvider,
  TooltipTrigger,
} from "./ui/tooltip"; // Import Tooltip components

interface ChampionIconProps {
  champion: Champion;
  size?: "sm" | "md" | "lg";
  draggable?: boolean;
}

const ChampionIcon: React.FC<ChampionIconProps> = ({
  champion,
  size = "md",
  draggable = true,
}) => {
  const { dispatch, state } = useSimulator();
  const isSelected = state.selectedChampion?.apiName === champion.apiName;

  const costColors = {
    1: "border-gray-500 bg-gray-500/20",
    2: "border-green-500 bg-green-500/20",
    3: "border-blue-500 bg-blue-500/20",
    4: "border-purple-500 bg-purple-500/20",
    5: "border-amber-500 bg-amber-500/20",
  };

  const sizeClasses = {
    sm: "w-8 h-8",
    md: "w-10 h-10",
    lg: "w-14 h-14",
  };

  const handleClick = () => {
    dispatch({
      type: "SELECT_CHAMPION",
      champion: isSelected ? undefined : champion,
    });
  };

  const handleDragStart = (e: React.DragEvent) => {
    e.dataTransfer.setData(
      "application/json",
      JSON.stringify({
        type: "champion",
        champion,
      }),
    );
  };

  return (
    // Use TooltipProvider and Tooltip
    <TooltipProvider delayDuration={200}>
      <Tooltip>
        <TooltipTrigger asChild>
          <Avatar
            className={cn(
              "relative rounded-md border cursor-pointer transition-all duration-200 flex items-center justify-center",
              sizeClasses[size],
              costColors[champion.cost as keyof typeof costColors],
              isSelected
                ? "ring-2 ring-accent shadow-lg shadow-accent/20"
                : "hover:scale-105",
            )}
            draggable={draggable}
            onDragStart={handleDragStart}
            onClick={handleClick}
          >
            <AvatarImage
              src={`/tft-champion-icons/${champion.icon.toLowerCase()}`}
              alt={champion.name}
            />
            <AvatarFallback>{champion.name.substring(0, 4)}</AvatarFallback>
          </Avatar>
        </TooltipTrigger>

        <TooltipContent>
          <p className="font-bold">{champion.name}</p>
          <p className="text-xs text-muted-foreground">
            {champion.traits.join(", ")}
          </p>
        </TooltipContent>
      </Tooltip>
    </TooltipProvider>
  );
};

export default ChampionIcon;
