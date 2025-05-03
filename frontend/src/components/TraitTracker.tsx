import React from "react";
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
import { ScrollArea } from "./ui/scroll-area";
import {
  Accordion,
  AccordionContent,
  AccordionItem,
  AccordionTrigger,
} from "./ui/accordion";

// Define colors for trait tiers
const traitTierColors = {
  inactive: "border-gray-700 bg-gray-700/10",
  bronze: "border-yellow-700 bg-yellow-600/50",
  silver: "border-gray-400 bg-gray-400/20",
  gold: "border-yellow-500 bg-yellow-500/20",
  prismatic:
    "bg-clip-text bg-gradient-to-r from-pink-500 via-purple-500 to-indigo-500",
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
  const { state } = useSimulator();
  const { traits, boardChampions } = state;

  // Only show traits when there are champions on board
  const hasChampions = boardChampions.length > 0;

  // Sort traits: active first, then by name
  const sortedTraits = [...traits]
    .filter((trait) => (hasChampions ? true : false)) // Show no traits when no champions
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

  // Find the next trait tier to achieve
  const getNextEffect = (trait: Trait): TraitEffect | undefined => {
    if (trait.active === 0) {
      // If no units for this trait, next effect is the first one
      const sortedEffects = [...trait.effects].sort(
        (a, b) => a.minUnits - b.minUnits,
      );
      return sortedEffects[0];
    }

    // Find the next effect tier that hasn't been reached
    const sortedEffects = [...trait.effects].sort(
      (a, b) => a.minUnits - b.minUnits,
    );
    return sortedEffects.find((effect) => trait.active < effect.minUnits);
  };

  return (
    <Card className="bg-card rounded-lg shadow-lg p-0 h-full flex flex-col">
      <CardHeader className="flex flex-row items-center justify-start space-y-0 pb-0">
        <CardTitle className="text-xl font-bold mb-4 text-white">
          Traits
        </CardTitle>
      </CardHeader>
      <CardContent className="flex-1 overflow-hidden flex flex-col pr-0 p-4 pt-0">
        <div className="flex-1 min-h-0 overflow-x-auto">
          {!hasChampions ? (
            <div className="text-gray-400 text-center py-4">
              Add champions to activate traits
            </div>
          ) : sortedTraits.some((trait) => trait.active > 0) ? (
            <ScrollArea className="flex-1 h-full mt-0 mb-0 focus-visible:ring-0 focus-visible:ring-offset-0 mr-0">
              <div className="space-y-3 pr-2">
                {sortedTraits.map((trait) => {
                  const activeEffect = getActiveEffect(trait);
                  const nextEffect = getNextEffect(trait);
                  const isActive = trait.active > 0;

                  // Only show active traits
                  if (trait.active === 0) return null;

                  return (
                    <Accordion
                      key={trait.apiName}
                      type="single"
                      collapsible
                      disabled={!isActive}
                      className="w-full mr-2"
                    >
                      <AccordionItem
                        value={trait.apiName}
                        className={cn(
                          "border-0 mb-0",
                          isActive
                            ? "bg-gray-800/60 border-2 rounded-md " +
                                traitTierColors[getHighestActiveTier(trait)]
                            : "bg-gray-800/30 rounded-md",
                        )}
                      >
                        <AccordionTrigger
                          className={cn(
                            "px-3 py-2 hover:no-underline",
                            !isActive && "cursor-default",
                          )}
                        >
                          <div className="flex flex-col w-full">
                            <div className="flex items-center justify-between w-full">
                              {/* Trait name and count */}
                              <div className="flex items-center space-x-2">
                                <div className="w-6 h-6 rounded-full flex items-center justify-center overflow-hidden">
                                  {trait.icon ? (
                                    <img
                                      src={`/tft-trait/${trait.icon}`}
                                      alt={trait.name}
                                      className={cn(
                                        "w-full h-full object-cover",
                                        !isActive && "opacity-40",
                                      )}
                                    />
                                  ) : (
                                    // Fallback if icon is missing
                                    <div
                                      className={cn(
                                        "w-full h-full rounded-full flex items-center justify-center text-xs font-bold",
                                        isActive
                                          ? "bg-primary text-black"
                                          : "bg-gray-700 text-gray-300",
                                      )}
                                    >
                                      {trait.active}
                                    </div>
                                  )}
                                </div>
                                <span
                                  className={cn(
                                    "font-semibold",
                                    isActive ? "text-white" : "text-gray-400",
                                  )}
                                >
                                  {trait.name}
                                </span>
                              </div>

                              {/* Progress to next bonus */}
                              {nextEffect && (
                                <div className="text-sm text-gray-400">
                                  {trait.active}/{nextEffect.minUnits}
                                </div>
                              )}
                            </div>

                            {/* Progress bar to next bonus */}
                            {nextEffect && (
                              <div className="mt-2 h-1 bg-gray-700 rounded-full overflow-hidden w-full">
                                <div
                                  className="h-full bg-primary"
                                  style={{
                                    width: `${Math.min(
                                      100,
                                      (trait.active / nextEffect.minUnits) *
                                        100,
                                    )}%`,
                                  }}
                                ></div>
                              </div>
                            )}
                          </div>
                        </AccordionTrigger>

                        <AccordionContent className="px-3 pb-3">
                          {/* Description and active bonus effect */}
                          <div className="text-xs text-primary overflow-hidden">
                            {trait.desc && (
                              <div className="mb-2 text-gray-300 break-all whitespace-pre-wrap overflow-hidden max-w-full">
                                {trait.desc
                                  .replace(/&nbsp;/g, " ")
                                  .replace(/<br\s*\/?>/g, "\n")
                                  .replace(/@/g, "@ ")}
                              </div>
                            )}

                            {activeEffect && (
                              <div className="text-primary">
                                <div className="font-semibold mb-1">
                                  Active Effect:
                                </div>
                                {activeEffect.variables ? (
                                  <div className="flex flex-wrap gap-1">
                                    {Object.entries(activeEffect.variables).map(
                                      ([key, value], idx) => (
                                        <div
                                          key={idx}
                                          className="text-primary break-all"
                                        >
                                          {`${key.replace(/([A-Z])/g, " $1").trim()}: ${value}`}
                                          {idx <
                                          Object.entries(activeEffect.variables)
                                            .length -
                                            1
                                            ? ", "
                                            : ""}
                                        </div>
                                      ),
                                    )}
                                  </div>
                                ) : (
                                  "Active effect"
                                )}
                              </div>
                            )}
                          </div>
                        </AccordionContent>
                      </AccordionItem>
                    </Accordion>
                  );
                })}
              </div>
            </ScrollArea>
          ) : (
            <div className="text-gray-400 text-center py-4">
              No active traits
            </div>
          )}
        </div>
      </CardContent>
      <CardFooter>
        {/* Gold counter */}
        <div className="mt-2 flex items-center justify-between">
          <div className="flex items-center space-x-2">
            <div className="w-6 h-6 rounded-full bg-amber-500 flex items-center justify-center">
              <span className="text-black font-bold">{state.gold}</span>
            </div>
            <span className="text-amber-400 font-semibold">Gold</span>
          </div>

          <div className="flex items-center space-x-2">
            <div className="w-6 h-6 rounded-full bg-primary flex items-center justify-center">
              <span className="text-black font-bold">
                {state.boardChampions.length}
              </span>
            </div>
            <span className="text-primary font-semibold">Champions</span>
          </div>
        </div>
      </CardFooter>
    </Card>
  );
};

export default TraitTracker;
