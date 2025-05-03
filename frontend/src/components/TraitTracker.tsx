import React from "react";
import { useSimulator } from "../context/SimulatorContext";
import { cn } from "../lib/utils";
import { Trait, TraitEffect } from "../utils/types";
import { Card, CardContent, CardHeader, CardTitle } from "./ui/card";
import { ScrollArea } from "./ui/scroll-area";

// Define colors for trait tiers
const traitTierColors = {
  inactive: "border-gray-700 bg-gray-700/10",
  bronze: "border-orange-700 bg-orange-700/20",
  silver: "border-gray-400 bg-gray-400/20",
  gold: "border-yellow-500 bg-yellow-500/20",
  prismatic: "border-purple-500 bg-purple-500/20", // Example, adjust as needed
};

// Function to determine the highest active tier
const getHighestActiveTier = (trait: Trait): keyof typeof traitTierColors => {
  // if (trait.active === 0) return 'inactive';

  // const sortedBonuses = [...trait.effects].sort((a, b) => b.count - a.count);
  // const activeBonus = sortedBonuses.find(bonus => bonus.count <= trait.active);

  // if (!activeBonus) return 'inactive'; // Should not happen if active > 0

  // // Determine tier based on bonus count thresholds (example, adjust based on actual game tiers)
  // // This requires knowing the max possible count for each tier
  // // For simplicity, let's map based on the number of bonuses achieved
  // const bonusIndex = sortedBonuses.findIndex(b => b.count === activeBonus.count);
  // const totalBonuses = sortedBonuses.length;

  // // Simple tier assignment based on index (adjust logic as needed)
  // if (totalBonuses <= 1) return 'bronze'; // Or maybe gold if it's the only tier
  // if (bonusIndex === 0) return totalBonuses > 2 ? 'prismatic' : 'gold'; // Highest tier
  // if (bonusIndex === 1) return totalBonuses > 2 ? 'gold' : 'silver';
  // if (bonusIndex === 2) return 'silver';
  return "bronze"; // Lowest active tier
};

const TraitTracker: React.FC = () => {
  const { state } = useSimulator();
  const { traits } = state;

  // Sort traits: active first, then by name
  const sortedTraits = [...traits].sort((a, b) => {
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

  // Find the activated bonus for a trait
  const getActiveEffect = (trait: Trait): TraitEffect | undefined => {
    if (trait.active === 0) return undefined;

    // Find the highest bonus that's activated
    // TODO: Implement logic to find the highest activated bonus based on trait.active and trait.effects
    const sortedEffects = [...trait.effects].sort(
      (a, b) => b.minUnits - a.minUnits,
    ); // Sort descending by count needed
    return sortedEffects.find((effect) => trait.active >= effect.minUnits);
  };

  const getTraitTierStyle = (
    count: number,
    bonuses: { count: number }[],
    activeCount: number,
  ): string => {
    let style = "text-gray-500"; // Default inactive/unreached
    let highestAchieved = -1;
    bonuses.forEach((b) => {
      if (activeCount >= b.count) {
        highestAchieved = Math.max(highestAchieved, b.count);
      }
    });

    if (activeCount >= count) {
      // This specific bonus count is active
      if (count === highestAchieved) {
        // Determine color based on tier (using simplified logic here)
        const tierIndex = bonuses.findIndex((b) => b.count === count);
        const totalTiers = bonuses.length;
        if (totalTiers <= 1)
          style = "text-yellow-400 font-bold"; // Single tier
        else if (tierIndex === bonuses.length - 1)
          style = "text-purple-400 font-bold"; // Highest (Prismatic?)
        else if (tierIndex === bonuses.length - 2)
          style = "text-yellow-400 font-bold"; // Gold
        else if (tierIndex === bonuses.length - 3)
          style = "text-gray-300 font-bold"; // Silver
        else style = "text-orange-500 font-bold"; // Bronze
      } else {
        style = "text-gray-400"; // Achieved, but not the highest
      }
    } else {
      // Not yet achieved
      // Find the next tier to achieve
      const nextTier = bonuses.find((b) => b.count > activeCount);
      if (nextTier && count === nextTier.count) {
        style = "text-gray-600"; // Next tier hint
      }
    }
    return style;
  };

  return (
    <Card className="bg-card rounded-lg shadow-lg p-4 h-full">
      <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
        <CardTitle className="text-xl font-bold mb-4 text-white">
          Traits
        </CardTitle>
      </CardHeader>
      {/* <ScrollArea className='flex-1 mr-0'> */}
      <CardContent>
        <div className="space-y-3">
          {sortedTraits.map((trait) => {
            const activeEffect = getActiveEffect(trait);
            const isActive = trait.active > 0;
            // const nextBonus = trait.bonuses.find(bonus => bonus.count > trait.active);

            return (
              <div
                key={trait.apiName}
                className={cn(
                  "p-3 rounded-md transition-all",
                  isActive
                    ? "bg-gray-800/60 border " + trait.style
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
                  {/* {nextBonus && (
                  <div className="text-xs text-gray-400">
                    {trait.active}/{nextBonus.count}
                  </div>
                )} */}
                </div>

                {/* Active bonus effect */}
                {activeEffect && (
                  <div className="mt-1 text-xs text-primary">
                    "get extra bonus"
                  </div>
                )}

                {/* Progress bar to next bonus */}
                {/* {nextBonus && (
                <div className="mt-2 progress-bar">
                  <div 
                    className="progress-fill bg-primary"
                    style={{ width: `${Math.min(100, (trait.active / nextBonus.count) * 100)}%` }}
                  ></div>
                </div>
              )} */}
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
      </CardContent>
      {/* </ScrollArea> */}
    </Card>
  );
};

export default TraitTracker;
