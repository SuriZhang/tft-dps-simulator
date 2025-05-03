import React from "react";
import { useSimulator } from "../context/SimulatorContext";
import { cn } from "../lib/utils";
import { Trait, TraitEffect } from "../utils/types";
import { Card, CardContent, CardHeader, CardTitle } from "./ui/card";
import { ScrollArea } from "./ui/scroll-area";

// Define colors for trait tiers
const traitTierColors = {
  inactive: "border-gray-700 bg-gray-700/10",
  bronze: "border-amber-700 bg-amber-700/20",
  silver: "border-gray-400 bg-gray-400/20",
  gold: "border-yellow-500 bg-yellow-500/20",
  prismatic: "border-zinc-50 bg-zinc-50/20",
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
  if (trait.active === 0) return 'inactive';

  const activeEffect = getActiveEffect(trait);
  console.log(activeEffect);
  if (!activeEffect) return 'inactive';

  // Map the effect style to a tier
  switch (activeEffect.style) {
    case 1: return 'bronze';
    case 3: return 'silver';
    case 4: return 'gold';
    case 5: return 'gold';
    case 6: return 'prismatic';
    default: return 'inactive';
  }
};

const TraitTracker: React.FC = () => {
  const { state } = useSimulator();
  const { traits, boardChampions } = state;
  
  // Only show traits when there are champions on board
  const hasChampions = boardChampions.length > 0;

  // Sort traits: active first, then by name
  const sortedTraits = [...traits]
    .filter(trait => hasChampions ? true : false) // Show no traits when no champions
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
    return sortedEffects.find(effect => trait.active < effect.minUnits);
  };

  return (
    <Card className="bg-card rounded-lg shadow-lg p-4 h-full flex flex-col">
      <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
        <CardTitle className="text-xl font-bold mb-4 text-white">
          Traits 
        </CardTitle>
      </CardHeader>
      <CardContent className="flex-1 overflow-hidden flex flex-col p-4 pt-0">
        <div className="flex-1 min-h-0 overflow-hidden">
          {!hasChampions && "Add champions to activate traits"}
          <ScrollArea className="flex-1 h-full mt-0 focus-visible:ring-0 focus-visible:ring-offset-0 p-0">
            <div className="space-y-3 pr-2">
              {sortedTraits.map((trait) => {
                const activeEffect = getActiveEffect(trait);
                const nextEffect = getNextEffect(trait);
                const isActive = trait.active > 0;
                
                // Only show active traits when there are champions
                if (hasChampions && trait.active === 0) return null;

                return (
                  <div
                    key={trait.apiName}
                    className={cn(
                      "p-3 rounded-md w-full transition-all",
                      isActive
                        ? "bg-gray-800/60 border-2 " + traitTierColors[getHighestActiveTier(trait)]
                        : "bg-gray-800/30",
                    )}
                  >
                    <div className="flex items-center justify-between">
                      {/* Trait name and count */}
                      <div className="flex items-center space-x-2">
                        <div
                          className={cn(
                            "w-6 h-6 rounded-full flex items-center justify-center text-xs font-bold",
                            isActive
                              ? "bg-primary text-black"
                              : "bg-gray-700 text-gray-300",
                          )}
                        >
                          {trait.active}
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
                        <div className="text-xs text-gray-400">
                          {trait.active}/{nextEffect.minUnits}
                        </div>
                      )}
                    </div>

                    {/* Active bonus effect */}
                    {activeEffect && (
                      <div className="mt-1 text-xs text-primary">
                        {trait.desc}
                        {activeEffect.variables ? 
                          Object.entries(activeEffect.variables).map(([key, value]) => 
                            `${key.replace(/([A-Z])/g, ' $1').trim()}: ${value}`
                          ).join(', ') :
                          "Active effect"
                        }
                      </div>
                    )}

                    {/* Progress bar to next bonus */}
                    {nextEffect && (
                      <div className="mt-2 h-1 bg-gray-700 rounded-full overflow-hidden">
                        <div 
                          className="h-full bg-primary"
                          style={{ 
                            width: `${Math.min(100, (trait.active / nextEffect.minUnits) * 100)}%` 
                          }}
                        ></div>
                      </div>
                    )}
                  </div>
                );
              })}
            </div>

            {/* Gold counter */}
            <div className="mt-6 flex items-center justify-between">
              <div className="flex items-center space-x-2">
                <div className="w-6 h-6 rounded-full bg-amber-500 flex items-center justify-center">
                  <span className="text-black font-bold">{state.gold}</span>
                </div>
                <span className="text-amber-400 font-semibold">Gold</span>
              </div>

              <div className="flex items-center space-x-2">
                <div className="w-6 h-6 rounded-full bg-primary flex items-center justify-center">
                  <span className="text-black font-bold">{state.level}</span>
                </div>
                <span className="text-primary font-semibold">Level</span>
              </div>
            </div>
          </ScrollArea>
        </div>
      </CardContent>
    </Card>
  );
};

export default TraitTracker;
