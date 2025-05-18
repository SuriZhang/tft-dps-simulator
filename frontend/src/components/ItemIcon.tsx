import React from "react";
import { Item } from "../utils/types";
import { cn } from "../lib/utils";
import { useSimulator } from "../context/SimulatorContext";
import { Avatar, AvatarFallback, AvatarImage } from "./ui/avatar";
import {
  Tooltip,
  TooltipContent,
  TooltipProvider,
  TooltipTrigger,
} from "./ui/tooltip";

// Helper function to format item description
const formatEffectValueHelper = (value: any): string => {
  if (typeof value === "number") {
    // Display whole numbers without decimal, others as is (or to a fixed precision if needed)
    return value % 1 === 0 ? value.toFixed(0) : value.toString();
  }
  return String(value);
};

const formatDescription = (
  desc: string | undefined,
  effects: Record<string, any> | undefined,
): string => {
  if (!desc) return "";
  let formattedDesc = desc;

  // 1. Replace placeholders like @EffectName@
  if (effects) {
    for (const [key, value] of Object.entries(effects)) {
      const placeholder = new RegExp(`@${key}(?![a-zA-Z0-9_])@`, "g"); // Ensure full placeholder match
      const displayValue =
        value !== null && value !== undefined
          ? formatEffectValueHelper(value)
          : "";
      formattedDesc = formattedDesc.replace(placeholder, displayValue);
    }
  }

  // 2. Handle specific TFT tags
  // <TFTKeyword>Sunder</TFTKeyword> -> <span class="text-yellow-500 font-medium">Sunder</span>
  formattedDesc = formattedDesc.replace(
    /<TFTKeyword>(.*?)<\/TFTKeyword>/g,
    '<span class="text-yellow-500 font-medium">$1</span>',
  );

  // <tftbold>Sunder</tftbold> -> <strong>Sunder</strong>
  formattedDesc = formattedDesc.replace(
    /<tftbold>(.*?)<\/tftbold>/g,
    "<strong>$1</strong>",
  );

  // <br> tags
  formattedDesc = formattedDesc.replace(/<br\s*\/?>/gi, "<br />");

  // Strip <tftitemrules> tags but keep content (content should have been processed for <tftbold>)
  formattedDesc = formattedDesc.replace(/<\/?tftitemrules>/g, "");

  // Handle general stat tags like <scaleAP> or <scaleAD> by making them slightly distinct if needed
  // For now, this example focuses on the provided tags.
  // You might want to add more rules for other custom tags if they exist.

  return formattedDesc;
};

const primaryStatsConfig: Record<
  string,
  { name: string; colorClass?: string; isPercentage?: boolean; icon?: string }
> = {
  Health: { name: "Health", colorClass: "text-green-400", icon:"/icon_health_max.png" },
  Mana: { name: "Mana", colorClass: "text-sky-400", icon:"/icon_mana.png" },
  Armor: { name: "Armor", colorClass: "text-yellow-400", icon:"/icon_armor.png" },
  MagicResist: { name: "Magic Resist", colorClass: "text-cyan-400", icon:"/icon_mr.png" },
  AttackDamage: { name: "Attack Damage", colorClass: "text-red-400", icon:"/icon_damage.png" },
  AbilityPower: { name: "Ability Power", colorClass: "text-pink-400", icon:"/icon_ap.png" },
  AttackSpeed: {
    name: "Attack Speed",
    colorClass: "text-teal-400",
    isPercentage: true,
    icon:"/icon_as.png",
  },
  CritChance: {
    name: "Crit Chance",
    colorClass: "text-orange-400",
    isPercentage: true,
    icon:"/icon_crit.png",
  },
  CritDamage: {
    name: "Crit Damage",
    colorClass: "text-orange-500",
    isPercentage: true,
    icon:"/icon_critmult.png",
  },
  // Add more stats as needed, e.g., Omnivamp, etc.
  // Example with icon: Health: { name: "Health", colorClass: "text-green-400", icon: "/icons/health.svg" },
};

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
    <TooltipProvider delayDuration={200}>
      <Tooltip>
        <TooltipTrigger asChild>
          <Avatar
            className={cn(
              "relative rounded-md border cursor-pointer transition-all duration-200 flex items-center justify-center",
              sizeClasses[size],
              isSelected
                ? "ring-2 ring-accent shadow-lg shadow-accent/20"
                : "hover:scale-105",
            )}
            draggable={draggable}
            onDragStart={handleDragStart}
            onClick={handleClick}
          >
            <AvatarImage src={`/tft-item/${item.icon}`} alt={item.name} />
            <AvatarFallback className="text-xs">
              {item.name.substring(0, 5)}
            </AvatarFallback>
          </Avatar>
        </TooltipTrigger>
        <TooltipContent className="w-72 md:w-80 p-3">
          {" "}
          {/* Adjusted width and padding */}
          <p className="font-bold text-base mb-2">{item.name}</p>
          {/* Display primary stats */}
          <div className="space-y-0.5 mb-3 text-sm">
            {
              item.effects &&
                Object.entries(item.effects)
                  .map(([key, value]) => {
                    const statConfig = primaryStatsConfig[key];
                    if (
                      statConfig &&
                      typeof value === "number" &&
                      value !== 0
                    ) {
                      let displayValue =
                        value > 0 ? `+${value.toFixed(0)}` : value.toFixed(0);
                      if (statConfig.isPercentage) {
                        // Assuming 'value' for percentage stats is the direct number, e.g., 25 for 25%
                        // If it's a decimal like 0.25, you'd multiply by 100.
                        displayValue = `${value > 0 ? "+" : ""}${value.toFixed(0)}%`;
                      }

                      return (
                        <div key={key} className="flex items-center">
                          {/* Placeholder for icon, if you have them mapped in primaryStatsConfig */}
                          {statConfig.icon && <img src={statConfig.icon}  className="w-4 h-4 mr-1.5" />}
                          <span
                            className={`font-medium ${statConfig.colorClass || "text-foreground"} mr-1.5`}
                          >
                            {displayValue}
                          </span>
                          <span className="text-muted-foreground">
                            {statConfig.name}
                          </span>
                        </div>
                      );
                    }
                    return null;
                  })
                  .filter(Boolean) // Remove null entries
            }
          </div>
          {/* Formatted Description */}
          {item.desc && (
            <div
              className="text-sm leading-normal" // Uses default popover foreground color
              dangerouslySetInnerHTML={{
                __html: formatDescription(item.desc, item.effects),
              }}
            />
          )}
        </TooltipContent>
      </Tooltip>
    </TooltipProvider>
  );
};

export default ItemIcon;
