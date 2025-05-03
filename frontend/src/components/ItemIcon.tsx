import React from "react";
import { Item } from "../utils/types";
import { cn } from "../lib/utils";
import { useSimulator } from "../context/SimulatorContext";
import { Avatar, AvatarFallback, AvatarImage } from "./ui/avatar"; // Import Avatar components
import {
  Tooltip,
  TooltipContent,
  TooltipProvider,
  TooltipTrigger,
} from "./ui/tooltip"; // Import Tooltip components

interface ItemIconProps {
  item: Item;
  size?: "sm" | "md" | "lg";
  draggable?: boolean;
}

const ItemIcon: React.FC<ItemIconProps> = ({
  item,
  size = "md",
  draggable = true,
}) => {
  const { dispatch, state } = useSimulator();
  const isSelected = state.selectedItem?.apiName === item.apiName;

  const typeColors = {
    component: "border-blue-400 bg-blue-800/30",
    completed: "border-purple-400 bg-purple-800/30",
    special: "border-amber-400 bg-amber-800/30",
  };

  const sizeClasses = {
    sm: "w-8 h-8",
    md: "w-10 h-10",
    lg: "w-12 h-12",
  };

  const handleClick = () => {
    dispatch({
      type: "SELECT_ITEM",
      item: isSelected ? undefined : item,
    });
  };

  const handleDragStart = (e: React.DragEvent) => {
    e.dataTransfer.setData(
      "application/json",
      JSON.stringify({
        type: "item",
        item,
      }),
    );
  };

  return (
    // Use TooltipProvider and Tooltip
    <TooltipProvider delayDuration={200}>
      <Tooltip>
        <TooltipTrigger asChild>
          {/* Use Avatar component */}
          <Avatar
            className={cn(
              "relative rounded-md border cursor-pointer transition-all duration-200 flex items-center justify-center",
              sizeClasses[size],
              // typeColors[item.type],
              isSelected
                ? "ring-2 ring-accent shadow-lg shadow-accent/20"
                : "hover:scale-105",
            )}
            draggable={draggable}
            onDragStart={handleDragStart}
            onClick={handleClick}
          >
            <AvatarImage
              src={item.icon || "/placeholder.svg"}
              alt={item.name}
            />
            <AvatarFallback className="text-xs">
              {item.name.substring(0, 5)}
            </AvatarFallback>
          </Avatar>
        </TooltipTrigger>
        {/* Use TooltipContent */}
        <TooltipContent>
          <p className="font-bold">{item.name}</p>
          <p className="text-xs text-muted-foreground">{item.desc}</p>
          {/* {item.stats && (
            <div className="mt-1 text-xs">
              {Object.entries(item.stats).map(([stat, value]) => (
                <p key={stat}>{stat}: {value}</p>
              ))}
            </div>
          )} */}
        </TooltipContent>
      </Tooltip>
    </TooltipProvider>
  );
};

export default ItemIcon;
