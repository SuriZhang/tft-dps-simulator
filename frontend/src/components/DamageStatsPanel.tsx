import React from "react";
import { useSimulator } from "../context/SimulatorContext";
import { cn } from "../lib/utils";
import {
  Bar,
  BarChart,
  ResponsiveContainer,
  XAxis,
  YAxis,
  Tooltip as RechartsTooltip,
  Legend,
  CartesianGrid,
} from "recharts";
import { Card, CardContent, CardHeader, CardTitle } from "./ui/card";
import {
  Tooltip,
  TooltipContent,
  TooltipProvider,
  TooltipTrigger,
} from "./ui/tooltip";
import { Info } from "lucide-react";
import { ScrollArea } from "./ui/scroll-area";

const damageStatsPanel = () => {
  const { state } = useSimulator();
  const { simulationResults, boardChampions } = state;

  console.log("Simulation Results:", simulationResults);

  // Format data for the stacked bar chart
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
          // Access properties with correct lowercase naming
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

  console.log("Chart Data:", chartData);

  // Custom tooltip for the stacked bar chart
  const CustomTooltip = ({ active, payload, label }: any) => {
    if (active && payload && payload.length) {
      const data = payload[0].payload;
      return (
        <div className="bg-background border border-border p-3 rounded-md shadow-md">
          <p className="font-bold text-lg mb-1">{label}</p>
          <p className="text-sm mb-2">DPS: {data.dps.toFixed(1)}</p>

          <div className="space-y-1">
            <div className="flex justify-between items-center">
              <span className="flex items-center">
                <div className="w-3 h-3 bg-red-500 mr-2"></div>
                Physical:
              </span>
              <span className="font-medium">{data.AD}</span>
            </div>

            <div className="flex justify-between items-center">
              <span className="flex items-center">
                <div className="w-3 h-3 bg-blue-500 mr-2"></div>
                Magic:
              </span>
              <span className="font-medium">{data.AP}</span>
            </div>

            <div className="flex justify-between items-center">
              <span className="flex items-center">
                <div className="w-3 h-3 bg-purple-500 mr-2"></div>
                True:
              </span>
              <span className="font-medium">{data.True}</span>
            </div>

            <div className="flex justify-between items-center border-t border-border pt-1 mt-1">
              <span>Total:</span>
              <span className="font-bold">{data.total}</span>
            </div>
          </div>

          <div className="mt-2 text-xs text-muted-foreground">
            <div className="flex justify-between">
              <span>Auto Attacks: {data.autoAttacks}</span>
              <span>Spell Casts: {data.spellCasts}</span>
            </div>
          </div>
        </div>
      );
    }
    return null;
  };

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

      <CardContent>
        <ScrollArea className="h-[400px] pr-4">
          {!simulationResults || simulationResults.length === 0 ? (
            <div className="text-center text-muted-foreground italic py-8">
              {boardChampions.length === 0
                ? "No champions on the board to display statistics."
                : "Run a simulation to see damage statistics."}
            </div>
          ) : (
            <>
              <div className="h-[300px] mb-8">
                <ResponsiveContainer width="100%" height="100%">
                  <BarChart
                    data={chartData}
                    margin={{ top: 10, right: 30, left: 0, bottom: 5 }}
                  >
                    <CartesianGrid strokeDasharray="3 3" opacity={0.1} />
                    <XAxis dataKey="name" />
                    <YAxis />
                    <RechartsTooltip content={<CustomTooltip />} />
                    <Legend />
                    <Bar
                      dataKey="AD"
                      stackId="a"
                      fill="#ef4444"
                      name="Physical"
                    />
                    <Bar dataKey="AP" stackId="a" fill="#3b82f6" name="Magic" />
                    <Bar
                      dataKey="True"
                      stackId="a"
                      fill="#8b5cf6"
                      name="True"
                    />
                  </BarChart>
                </ResponsiveContainer>
              </div>

              <div className="space-y-4 mt-8">
                <h3 className="font-semibold text-lg">Detailed Stats</h3>
                {chartData.map((champion, i) => {
                  // Find champion details to get cost and stars
                  const boardChampion = boardChampions.find(
                    (c) => c.apiName === champion.apiName,
                  );
                  const cost = boardChampion?.cost || 1;
                  const stars = boardChampion?.stars || 1;

                  return (
                    <div
                      key={i}
                      className="border border-border rounded-lg p-3"
                    >
                      <div className="flex items-center justify-between mb-2">
                        <div className="flex items-center">
                          <div
                            className={cn(
                              "w-6 h-6 rounded-full flex items-center justify-center mr-2",
                              getCostColor(cost),
                            )}
                          >
                            <span className="text-xs font-bold">{cost}</span>
                          </div>
                          <span className="font-medium">{champion.name}</span>
                          <div className="flex ml-2">
                            {Array.from({ length: stars }).map((_, i) => (
                              <div
                                key={i}
                                className="w-3 h-3 bg-warning ml-0.5"
                                style={{
                                  clipPath:
                                    "polygon(50% 0%, 61% 35%, 98% 35%, 68% 57%, 79% 91%, 50% 70%, 21% 91%, 32% 57%, 2% 35%, 39% 35%)",
                                }}
                              />
                            ))}
                          </div>
                        </div>
                        <div className="text-right">
                          <div className="font-semibold">
                            {champion.total} damage
                          </div>
                          <div className="text-sm text-muted-foreground">
                            {champion.dps} DPS
                          </div>
                        </div>
                      </div>

                      <div className="flex h-6 rounded-md overflow-hidden bg-secondary mb-1">
                        {champion.total > 0 && (
                          <>
                            <div
                              className="bg-red-500 h-full"
                              style={{
                                width: `${(champion.AD / champion.total) * 100}%`,
                              }}
                              title={`Physical: ${champion.AD}`}
                            />
                            <div
                              className="bg-blue-500 h-full"
                              style={{
                                width: `${(champion.AP / champion.total) * 100}%`,
                              }}
                              title={`Magic: ${champion.AP}`}
                            />
                            <div
                              className="bg-purple-500 h-full"
                              style={{
                                width: `${(champion.True / champion.total) * 100}%`,
                              }}
                              title={`True: ${champion.True}`}
                            />
                          </>
                        )}
                      </div>

                      <div className="flex justify-between text-xs text-muted-foreground mt-1">
                        <div className="flex space-x-3">
                          <span className="flex items-center">
                            <div className="w-2 h-2 bg-red-500 mr-1"></div>
                            {champion.AD}
                          </span>
                          <span className="flex items-center">
                            <div className="w-2 h-2 bg-blue-500 mr-1"></div>
                            {champion.AP}
                          </span>
                          <span className="flex items-center">
                            <div className="w-2 h-2 bg-purple-500 mr-1"></div>
                            {champion.True}
                          </span>
                        </div>
                        <div>
                          Auto: {champion.autoAttacks} / Spell:{" "}
                          {champion.spellCasts}
                        </div>
                      </div>
                    </div>
                  );
                })}
              </div>
            </>
          )}
        </ScrollArea>
      </CardContent>
    </Card>
  );
};

// Function to get champion cost color
const getCostColor = (cost: number) => {
  switch (cost) {
    case 1:
      return "bg-gray-500 text-white";
    case 2:
      return "bg-green-500 text-white";
    case 3:
      return "bg-blue-500 text-white";
    case 4:
      return "bg-purple-500 text-white";
    case 5:
      return "bg-amber-500 text-white";
    default:
      return "bg-gray-500 text-white";
  }
};

export default damageStatsPanel;
