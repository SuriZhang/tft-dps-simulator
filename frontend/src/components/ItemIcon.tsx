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
import { formatDescription } from "../utils/helpers";

// // Helper function to format item description
// const formatEffectValueHelper = (value: any): string => {
//   if (typeof value === "number") {
//     // Display whole numbers without decimal, others as is (or to a fixed precision if needed)
//     return value % 1 === 0 ? value.toFixed(0) : value.toString();
//   }
//   return String(value);
// };

// const formatDescription = (
//   desc: string | undefined,
//   effects: Record<string, any> | undefined,
// ): string => {
//   if (!desc) return "";
//   let formattedDesc = desc;

//   // 1. Replace placeholders like @EffectName@ or @EffectName*100@
//   if (effects) {
//     const placeholderRegex = /@([^@]+)@/g;
//     let resultString = "";
//     let lastIndex = 0;
//     let match;

//     while ((match = placeholderRegex.exec(formattedDesc)) !== null) {
//       resultString += formattedDesc.substring(lastIndex, match.index); // Append text before placeholder

//       const placeholderContent = match[1]; // e.g., "CritDamageToGive*100" or "HexRange"
//       let replacementValue = match[0]; // Default to original placeholder if not processed

//       if (placeholderContent.endsWith("*100")) {
//         const key = placeholderContent.slice(0, -4); // Remove '*100'
//         if (effects[key] !== undefined && typeof effects[key] === "number") {
//           replacementValue = (effects[key] * 100).toFixed(0);
//         }
//       } else {
//         const key = placeholderContent;
//         if (effects[key] !== undefined) {
//           replacementValue = formatEffectValueHelper(effects[key]);
//         }
//       }
//       resultString += replacementValue;
//       lastIndex = placeholderRegex.lastIndex;
//     }
//     resultString += formattedDesc.substring(lastIndex); // Append remaining text
//     formattedDesc = resultString;
//   }

//   // 2. Handle specific TFT tags and other HTML-like tags

//   // Styling for specific keywords/concepts
//   formattedDesc = formattedDesc.replace(
//     /<TFTKeyword>(.*?)<\/TFTKeyword>/gi,
//     '<span class="text-yellow-500 font-medium">$1</span>',
//   );
//   formattedDesc = formattedDesc.replace(
//     /<tftbold>(.*?)<\/tftbold>/gi,
//     '<em class="font-semibold">$1</em>', // Italic + Semi-bold
//   );
//   formattedDesc = formattedDesc.replace(
//     /<magicDamage>(.*?)<\/magicDamage>/gi,
//     '<span class="text-purple-400">$1</span>',
//   );
//   formattedDesc = formattedDesc.replace(
//     /<physicalDamage>(.*?)<\/physicalDamage>/gi,
//     '<span class="text-red-500">$1</span>',
//   );
//   formattedDesc = formattedDesc.replace(
//     /<trueDamage>(.*?)<\/trueDamage>/gi,
//     '<span class="text-white font-semibold">$1</span>',
//   );
//   formattedDesc = formattedDesc.replace(
//     /<TFTBonus>(.*?)<\/TFTBonus>/gi,
//     '<span class="text-green-400">$1</span>',
//   );
//   formattedDesc = formattedDesc.replace(
//     /<TFTShadowItemPenalty>(.*?)<\/TFTShadowItemPenalty>/gi,
//     '<span class="text-red-400 italic">$1</span>', // Example: Red and italic for penalty
//   );
//   formattedDesc = formattedDesc.replace(
//     /<healing>(.*?)<\/healing>/gi,
//     '<span class="text-green-300">$1</span>',
//   );
//   formattedDesc = formattedDesc.replace(
//     /<TFTHighlight>(.*?)<\/TFTHighlight>/gi,
//     '<span class="text-blue-300 font-medium">$1</span>',
//   );

//   // Basic handling for <li> items (convert to bullet points on new lines)
//   // Ensuring <ul> or <ol> wrappers would require more complex parsing.
//   formattedDesc = formattedDesc.replace(/<li>(.*?)<\/li>/gi, "<br />&bull; $1");

//   // Tags to strip while keeping their content.
//   // This will also strip attributes within these tags.
//   // Order might matter if tags are nested.
//   const tagsToStripKeepContent = [
//     "tftitemrules",
//     "rules",
//     "TFTRadiantItemBonus",
//     "TFTShadowItemBonus",
//     "TFTPassive",
//     "scaleLevel", // Handles <scaleLevel> and <scaleLevel enabled=...> by stripping tag, keeping content
//     "scaleShimmer",
//     "TFTTrackerLabel",
//     "spellPassive",
//     "spellActive", // Handles <spellActive> and <spellActive enabled=...> by stripping tag, keeping content
//     "scaleHealth",
//     "keyword", // If different from styled TFTKeyword
//     "tftrules",
//     "active",
//     // Complex conditional tags - content is kept as a fallback
//     "ShowIfNot\\.TFT14_Mob_IsActive_T3", // Escaped dot for regex
//     "ShowIf\\.TFT14_Mob_IsActive_T3", // Escaped dot for regex
//     "ShowIfNot\\.TFTUnitProperty\\.Item:TFT10_BlingActive",
//     "ShowIf\\.TFTUnitProperty\\.Item:TFT10_BlingActive",
//     "ShowIfCustom\\.Set=TFTSetEventCT",
//   ];

//   tagsToStripKeepContent.forEach((tag) => {
//     // Regex to match <tag ...attributes...>content</tag> and replace with content
//     const contentKeepingRegex = new RegExp(`<${tag}[^>]*>(.*?)</${tag}>`, "gi");
//     formattedDesc = formattedDesc.replace(contentKeepingRegex, "$1");
//     // Regex to match opening/closing tags like <tag ...attributes...> or </tag> and remove them
//     // This helps clean up tags that might not have been caught by the above or are self-closing/empty
//     const stripTagOnlyRegex = new RegExp(`</?${tag}[^>]*>`, "gi");
//     formattedDesc = formattedDesc.replace(stripTagOnlyRegex, "");
//   });

//   // Consolidate multiple <br> tags (and variants) into one
//   formattedDesc = formattedDesc.replace(/(<br\s*\/?>\s*)+/gi, "<br />");

//   return formattedDesc;
// };

const primaryStatsConfig: Record<
  string,
  { name: string; colorClass?: string; isPercentage?: boolean; icon?: string }
> = {
  Health: {
    name: "Health",
    colorClass: "text-green-400",
    icon: "/icon_health_max.png",
  },
  Mana: { name: "Mana", colorClass: "text-sky-400", icon: "/icon_mana.png" },
  Armor: {
    name: "Armor",
    colorClass: "text-yellow-400",
    icon: "/icon_armor.png",
  },
  MagicResist: {
    name: "Magic Resist",
    colorClass: "text-cyan-400",
    icon: "/icon_mr.png",
  },
  AD: {
    // Changed from AttackDamage to AD to match item.effects
    name: "Attack Damage",
    colorClass: "text-orange-400", // Orange for AD value as per Infinity Edge image
    isPercentage: true, // AD for Infinity Edge is a percentage
    icon: "/icon_damage.png",
  },
  AP: {
    name: "Ability Power",
    colorClass: "text-purple-400",
    icon: "/icon_ap.png",
  },
  AS: {
    name: "Attack Speed",
    colorClass: "text-amber-400",
    isPercentage: true,
    icon: "/icon_as.png",
  },
  CritChance: {
    name: "Critical Strike Chance", // Updated name
    colorClass: "text-red-400", // Red for Crit Chance value as per Infinity Edge image
    isPercentage: true,
    icon: "/icon_crit.png",
  },
  CritDamage: {
    name: "Crit Damage",
    colorClass: "text-orange-500", 
    isPercentage: true,
    icon: "/icon_critmult.png",
  },
  BonusDamage: {
    name: "Damage Amp",
    colorClass: "text-gray-300",
    isPercentage: true,
    icon: "/icon_damageamp.png",
  },
  StatOmnivamp: {
    name: "Omnivamp",
    colorClass: "text-red-400",
    isPercentage: true,
    icon: "/icon_omnivamp.png",
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
                      let displayValue: string;
                      if (statConfig.isPercentage) {
                        // If value is a decimal like 0.15, multiply by 100.
                        // Otherwise, assume it's already a whole percentage number (e.g., 25 for 25%).
                        const numericValue = Number(value); // Ensure it's a number
                        const percentage =
                          (numericValue < 1 &&
                            numericValue > -1 &&
                            numericValue !== 0) ||
                          (numericValue > 1 &&
                            numericValue % 1 !== 0 &&
                            String(numericValue).includes(".")) // handles 0.15 or 1.15 (115%)
                            ? numericValue * 100
                            : numericValue;
                        displayValue = `${percentage > 0 ? "+" : ""}${percentage.toFixed(0)}%`;
                      } else {
                        displayValue = `${value > 0 ? "+" : ""}${value.toFixed(0)}`;
                      }

                      return (
                        <div key={key} className="flex items-center">
                          {/* Placeholder for icon, if you have them mapped in primaryStatsConfig */}
                          {statConfig.icon && (
                            <img
                              src={statConfig.icon}
                              className="w-4 h-4 mr-1.5"
                            />
                          )}
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
