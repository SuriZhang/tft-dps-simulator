import React from "react";
import { useSimulator } from "../context/SimulatorContext";
import { Card, CardContent, CardHeader, CardTitle } from "./ui/card";
import { Separator } from "./ui/separator"; // Add Separator import
import {
  TooltipProvider,
  Tooltip,
  TooltipTrigger,
  TooltipContent,
} from "./ui/tooltip";
import { Info } from "lucide-react";
import { ScrollArea } from "./ui/scroll-area";

const DamageStatsPanel = () => {
  const { state } = useSimulator();
  const { simulationResults, boardChampions } = state;

  // Format data for the damage display
  const chartData = React.useMemo(() => {
    if (!simulationResults) return [];

    return simulationResults
      .map((result) => {
        // Find champion name from apiName
        const champion = boardChampions.find(
          (c) => c.apiName === result.championApiName,
        );
        // Use champion name or extract from API name if not found
        const name =
          champion?.name || result.championApiName.replace("TFT14_", "");

        return {
          name,
          apiName: result.championApiName,
          icon: champion?.icon,
          cost: champion?.cost || 1,
          stars: champion?.stars || 1,
          AD: result.damageStats.totalADDamage,
          AP: result.damageStats.totalAPDamage,
          True: result.damageStats.totalTrueDamage,
          total: result.damageStats.totalDamage,
          dps: result.damageStats.dps,
          autoAttacks: result.damageStats.totalAutoAttackCounts,
          spellCasts: result.damageStats.totalSpellCastCounts,
        };
      })
      .sort((a, b) => b.total - a.total);
  }, [simulationResults, boardChampions]);

  // Calculate the maximum damage for proper scaling
  const maxDamage = React.useMemo(() => {
    if (!chartData.length) return 0;
    return Math.max(...chartData.map((champ) => champ.total));
  }, [chartData]);

  return (
    <Card className="mb-4">
      <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
        <CardTitle className="text-xl font-bold">Damage Statistics</CardTitle>
        <TooltipProvider>
          <Tooltip>
            <TooltipTrigger asChild>
              <Info className="h-4 w-4 text-muted-foreground cursor-help" />
            </TooltipTrigger>
            <TooltipContent>
              <p>Damage breakdown from simulation results</p>
            </TooltipContent>
          </Tooltip>
        </TooltipProvider>
      </CardHeader>

      {/* Add separator between header and content */}
      <Separator className="bg-gray-700 border" />

      <CardContent className="p-3">
        <ScrollArea className="h-[400px]">
          {!simulationResults || simulationResults.length === 0 ? (
            <div className="text-center text-muted-foreground italic py-8">
              {boardChampions.length === 0
                ? "No champions on the board to display statistics."
                : "Run a simulation to see damage statistics."}
            </div>
          ) : (
            <div className="space-y-3">
              {chartData.map((champion, i) => (
                <TooltipProvider key={i}>
                  <Tooltip>
                    <TooltipTrigger asChild>
                      <div className="flex items-center gap-3 relative cursor-pointer">
                        {/* Champion portrait with stars - updated to match image 2 */}
                        <div className="flex-shrink-0 relative">
                          <div
                            className={`w-12 h-12 border ${getBorderColor(champion.cost)} rounded-sm p-0.5 bg-gray-900 relative overflow-hidden`}
                          >
                            {champion.icon ? (
                              <img
                                src={`/tft-champion-icons/${champion.icon.toLowerCase()}`}
                                alt={champion.name}
                                className="w-full h-full object-cover"
                              />
                            ) : (
                              <div
                                className={`w-full h-full ${getBgColor(champion.cost)} flex items-center justify-center`}
                              >
                                {champion.name.charAt(0)}
                              </div>
                            )}

                            {/* Stars below the champion portrait - styled like image 2 */}
                            <div className="absolute bottom-0 left-0 right-0 flex justify-center bg-black bg-opacity-50 py-0.5">
                              {Array.from({ length: champion.stars }).map(
                                (_, i) => (
                                  <div
                                    key={i}
                                    className="w-2 h-2 bg-yellow-400 mx-0.5"
                                    style={{
                                      clipPath:
                                        "polygon(50% 0%, 61% 35%, 98% 35%, 68% 57%, 79% 91%, 50% 70%, 21% 91%, 32% 57%, 2% 35%, 39% 35%)",
                                    }}
                                  />
                                ),
                              )}
                            </div>
                          </div>
                        </div>

                        {/* Damage bar and info - updated to match image 1 */}
                        <div className="flex-grow">
                          <div className="flex justify-between mb-1">
                            <div className="text-sm font-medium">
                              {champion.name}
                            </div>
                            <div className="text-sm font-bold">
                              {Math.round(champion.total)}
                            </div>
                          </div>

                          <div className="h-4 bg-gray-800 rounded-sm w-full overflow-hidden relative">
                            {/* Bar segments with dynamic width based on damage percentage */}
                            <div
                              className="absolute h-full bg-orange-500"
                              style={{
                                width: `${(champion.total / maxDamage) * 100}%`,
                              }}
                            >
                              {/* Inside the orange bar, we'll position the blue and white segments */}
                              <div
                                className="absolute h-full bg-blue-500 right-0"
                                style={{
                                  width: `${(champion.AP / champion.total) * 100}%`,
                                }}
                              />
                              <div
                                className="absolute h-full bg-white right-0"
                                style={{
                                  width: `${(champion.True / champion.total) * 100}%`,
                                }}
                              />
                            </div>
                          </div>
                        </div>
                      </div>
                    </TooltipTrigger>
                    <TooltipContent
                      side="right"
                      className="p-3 w-[220px] bg-gray-900 border-gray-700"
                    >
                      <div className="font-bold text-lg mb-2">
                        {champion.name}
                      </div>

                      <div className="grid grid-cols-2 gap-2 mb-3">
                        <div className="flex flex-col">
                          <span className="text-xs text-gray-400">DPS</span>
                          <span className="font-semibold text-white">
                            {champion.dps.toFixed(1)}
                          </span>
                        </div>
                        <div className="flex flex-col">
                          <span className="text-xs text-gray-400">
                            Total Damage
                          </span>
                          <span className="font-semibold text-white">
                            {champion.total.toFixed(0)}
                          </span>
                        </div>
                      </div>

                      <div className="space-y-1.5 mb-3">
                        <div className="flex justify-between">
                          <span className="flex items-center text-gray-200">
                            <div className="w-3 h-3 bg-orange-500 mr-2"></div>
                            Physical
                          </span>
                          <span className="font-medium text-white">
                            {champion.AD.toFixed(0)}
                          </span>
                        </div>

                        <div className="flex justify-between">
                          <span className="flex items-center text-gray-200">
                            <div className="w-3 h-3 bg-blue-500 mr-2"></div>
                            Magic
                          </span>
                          <span className="font-medium text-white">
                            {champion.AP.toFixed(0)}
                          </span>
                        </div>

                        <div className="flex justify-between">
                          <span className="flex items-center text-gray-200">
                            <div className="w-3 h-3 bg-white mr-2"></div>
                            True
                          </span>
                          <span className="font-medium text-white">
                            {champion.True.toFixed(0)}
                          </span>
                        </div>
                      </div>

                      <div className="border-t border-gray-700 pt-2 text-sm">
                        <div className="flex justify-between">
                          <span className="text-gray-300">Auto Attacks:</span>
                          <span className="font-medium text-white">
                            {champion.autoAttacks}
                          </span>
                        </div>

                        <div className="flex justify-between">
                          <span className="text-gray-300">Spell Casts:</span>
                          <span className="font-medium text-white">
                            {champion.spellCasts}
                          </span>
                        </div>
                      </div>
                    </TooltipContent>
                  </Tooltip>
                </TooltipProvider>
              ))}
            </div>
          )}
        </ScrollArea>
      </CardContent>
    </Card>
  );
};

// Function to get champion border color based on cost
const getBorderColor = (cost: number) => {
  switch (cost) {
    case 1:
      return "border-gray-500";
    case 2:
      return "border-green-500";
    case 3:
      return "border-blue-500";
    case 4:
      return "border-purple-500";
    case 5:
      return "border-amber-500";
    default:
      return "border-gray-500";
  }
};

// Function to get champion background color based on cost
const getBgColor = (cost: number) => {
  switch (cost) {
    case 1:
      return "bg-gray-700";
    case 2:
      return "bg-green-900";
    case 3:
      return "bg-blue-900";
    case 4:
      return "bg-purple-900";
    case 5:
      return "bg-amber-900";
    default:
      return "bg-gray-700";
  }
};

export default DamageStatsPanel;
