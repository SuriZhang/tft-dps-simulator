import React, { useState } from "react";
import { useSimulator } from "../context/SimulatorContext";
import { cn } from "../lib/utils";
import { Trait, TraitEffect } from "../utils/types";
import {
  Card,
  CardContent,
  CardFooter,
  CardHeader,
  CardTitle,
} from "./ui/card";
import { ScrollArea, ScrollBar } from "./ui/scroll-area";

// Define SVG backgrounds for trait tiers
const traitTierColors = {
  inactive: "./trait-inactive.svg", 
  bronze: "./trait-bronze.svg", 
  silver: "./trait-silver.svg", 
  gold: "./trait-gold.svg", 
  prismatic: "./trait-prismatic.svg", 
};

// Find the activated bonus for a trait
const getActiveEffect = (trait: Trait): TraitEffect | undefined => {
  if (trait.active === 0) return undefined;

  // Find the highest effect that's activated
  const sortedEffects = [...trait.effects].sort(
    (a, b) => b.minUnits - a.minUnits,
  ); // Sort descending by minUnits
  return sortedEffects.find((effect) => trait.active >= effect.minUnits);
};

// Function to determine the highest active tier based on TraitEffect.style
const getHighestActiveTier = (trait: Trait): keyof typeof traitTierColors => {
  if (trait.active === 0) return "inactive";

  const activeEffect = getActiveEffect(trait);
  if (!activeEffect) return "inactive";

  // Map the effect style to a tier
  switch (activeEffect.style) {
    case 1:
      return "bronze";
    case 3:
      return "silver";
    case 4:
      return "gold";
    case 5:
      return "gold";
    case 6:
      return "prismatic";
    default:
      return "inactive";
  }
};

const TraitTracker: React.FC = () => {
  const { state, setHoveredTrait } = useSimulator();
  const { traits, boardChampions } = state;
  const [localHoveredTrait, setLocalHoveredTrait] = useState<string>("");

  // Only show traits when there are champions on board
  const hasChampions = boardChampions.length > 0;

  // Handle trait hover
  const handleTraitMouseEnter = (traitName: string) => {
    setLocalHoveredTrait(traitName);
    if (setHoveredTrait) {
      setHoveredTrait(traitName);
    }
  };

  const handleTraitMouseLeave = () => {
    setLocalHoveredTrait("");
    if (setHoveredTrait) {
      setHoveredTrait("");
    }
  };

  // Sort traits: active first, then by name
  const sortedTraits = [...traits]
    .filter(() => hasChampions) // Show no traits when no champions
    .sort((a, b) => {
      if (a.active > 0 && b.active === 0) return -1;
      if (a.active === 0 && b.active > 0) return 1;
      // Sort active traits by tier first, then count
      const tierA = getHighestActiveTier(a);
      const tierB = getHighestActiveTier(b);
      const tierOrder: (keyof typeof traitTierColors)[] = [
        "prismatic",
        "gold",
        "silver",
        "bronze",
        "inactive",
      ];
      if (tierOrder.indexOf(tierA) < tierOrder.indexOf(tierB)) return -1;
      if (tierOrder.indexOf(tierA) > tierOrder.indexOf(tierB)) return 1;
      // If same tier, sort by active count descending
      if (a.active > b.active) return -1;
      if (a.active < b.active) return 1;
      // Finally by name
      return a.name.localeCompare(b.name);
    });

  // Find current active effect
  const getCurrentActiveEffect = (trait: Trait): TraitEffect | undefined => {
    if (trait.active === 0) return undefined;
    return trait.effects.find((effect) => trait.active >= effect.minUnits && trait.active <= effect.maxUnits);
  };

  return (
    <Card className="bg-card rounded-lg shadow-lg p-0 h-full w-full flex flex-col">
      <div className="flex-1 overflow-hidden flex flex-col p-2">
        <div className="flex-1 overflow-x-auto">
          {!hasChampions ? (
            <div className="text-gray-400 text-center py-4">
              Add champions to activate traits
            </div>
          ) : sortedTraits.some((trait) => trait.active > 0) ? (
            <ScrollArea className="w-full">
              <div className="flex flex-wrap gap-2">
                {sortedTraits.map((trait) => {
                  // const activeEffect = getActiveEffect(trait); // Not directly used in new render
                  const currentEffect = getCurrentActiveEffect(trait); // For title attribute
                  // const isActive = trait.active > 0; // Redundant if filtering by tier or active count
                  const currentTier = getHighestActiveTier(trait);
                  const tierSvg = traitTierColors[currentTier];

                  // Only show active traits
                  if (trait.active === 0) return null; 
                  return (
                    <div
                      key={trait.apiName}
                      className={cn(
                        "flex items-center gap-2 p-1.5 rounded-md bg-neutral-800/80 cursor-pointer w-40", 
                        "transition-all hover:bg-neutral-700/90",
                        localHoveredTrait === trait.name &&
                          "ring-1 ring-primary ring-opacity-70",
                      )}
                      onMouseEnter={() => handleTraitMouseEnter(trait.name)}
                      onMouseLeave={handleTraitMouseLeave}
                      title={`${trait.name} (${trait.active}${currentEffect ? `/${currentEffect.minUnits}` : ''})`}
                    >
                      {/* 1. Hexagon Icon */}
                      <div className="relative flex-shrink-0 w-7 h-7">
                        <img
                          src={tierSvg}
                          alt={`${currentTier} tier`}
                          className="absolute inset-0 w-full h-full object-cover" // SVG fills the container
                        />
                        {trait.icon && (
                          <img
                            src={`/tft-trait/${trait.icon}`}
                            alt={trait.name}
                            className={cn(
                              "absolute top-1/2 left-1/2 transform -translate-x-1/2 -translate-y-1/2",
                              currentTier === 'inactive' ? "brightness-0 invert" : "brightness-0" 
                            )}
                            style={{ width: '16px', height: '16px' }} // Trait icon size: 16x16px
                          />
                        )}
                      </div>

                      {/* 2. Active Count Box */}
                      <div 
                        className="flex items-center justify-center rounded px-1.5 py-0.5 text-xs" // Smaller text for count
                        style={{ backgroundColor: 'rgba(10,10,20,0.5)', minWidth: '24px' }} // Darker, slightly transparent
                      >
                        <span className="font-bold text-white">{trait.active}</span>
                      </div>

                      {/* 3. Text Info (Name + Thresholds) */}
                      <div className="flex flex-col items-start leading-tight whitespace-nowrap">
                        <span className="text-xs font-semibold text-white">{trait.name}</span>
                        <div className="text-[10px]"> {/* Smaller text for thresholds */}
                          {trait.effects
                            .sort((a, b) => a.minUnits - b.minUnits)
                            .map((effect, index, arr) => (
                              <span
                                key={effect.minUnits}
                                className={cn(
                                  "font-medium",
                                  trait.active >= effect.minUnits && trait.active <= effect.maxUnits ? "text-white" : "text-neutral-500",
                                )}
                              >
                                {effect.minUnits}
                                {index < arr.length - 1 ? <span className="text-neutral-600 mx-0.5">{">"}</span> : ""}
                              </span>
                            ))}
                        </div>
                      </div>
                    </div>
                  );
                })}
                </div>
                <ScrollBar orientation="horizontal" />
            </ScrollArea>
          ) : (
            <div className="text-gray-400 text-center py-4">
              No active traits
            </div>
          )}
        </div>
      </div>
    </Card>
  );
};

export default TraitTracker;
