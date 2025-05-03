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
  LabelList,
} from "recharts";
import { Card, CardContent, CardHeader, CardTitle } from "./ui/card"; // Import Card components
import { Progress } from "./ui/progress"; // Import Progress
import {
  Tooltip,
  TooltipContent,
  TooltipProvider,
  TooltipTrigger,
} from "./ui/tooltip"; // Import Tooltip components
import { Info, Scroll } from "lucide-react"; // Import icon
import { ScrollArea } from "./ui/scroll-area";

const DamageStatsPanel = () => {
  const { state } = useSimulator();
  const { boardChampions } = state;

  // Generate mock damage stats for the champions
  const damageStats = boardChampions
    .map((champion) => {
      const stars = champion.stars || 1;
      const cost = champion.cost;
      // Calculate damage based on champion attributes
      const damage = Math.floor(
        cost * 100 * stars * (0.8 + Math.random() * 0.4),
      );

      return {
        name: champion.name,
        damage,
        cost: champion.cost,
        stars: champion.stars || 1,
      };
    })
    .sort((a, b) => b.damage - a.damage);

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

  // Custom tooltip for the chart
  const CustomTooltip = ({ active, payload }: any) => {
    if (active && payload && payload.length) {
      const data = payload[0].payload;
      return (
        <div className="bg-background border border-border p-2 rounded-md shadow-md">
          <p className="font-bold">{data.name}</p>
          <p className="text-sm">Damage: {data.damage}</p>
          <div className="flex items-center mt-1">
            <span className="mr-2">Cost: {data.cost}</span>
            <div className="flex">
              {Array.from({ length: data.stars }).map((_, i) => (
                <div key={i} className="w-3 h-3 bg-warning mr-0.5" />
              ))}
            </div>
          </div>
        </div>
      );
    }
    return null;
  };

  return (
    // Use Card component
    <Card className="mb-4">
      {/* Use CardHeader */}
      <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
        {/* Use CardTitle */}
        <CardTitle className="text-xl font-bold">Damage Statistics</CardTitle>
        {/* Use Tooltip for info icon */}
        <TooltipProvider>
          <Tooltip>
            <TooltipTrigger asChild>
              <Info className="h-4 w-4 text-muted-foreground cursor-help" />
            </TooltipTrigger>
            <TooltipContent>
              <p>Estimated damage based on current setup. (Simulation TBD)</p>
            </TooltipContent>
          </Tooltip>
        </TooltipProvider>
      </CardHeader>
      {/* Use CardContent */}
      <CardContent>
        <ScrollArea className="h-200">
          {damageStats.length === 0 ? (
            <div className="text-center text-muted-foreground italic py-8">
              No champions on the board to display statistics.
            </div>
          ) : (
            <>
              <div className="h-64 mb-4">
                <ResponsiveContainer width="100%" height="100%">
                  <BarChart
                    data={damageStats}
                    layout="vertical"
                    margin={{ top: 5, right: 30, left: 20, bottom: 5 }}
                  >
                    <XAxis type="number" />
                    <YAxis
                      type="category"
                      dataKey="name"
                      width={60}
                      tick={{ fontSize: 12 }}
                    />
                    <RechartsTooltip content={<CustomTooltip />} />
                    <Bar
                      dataKey="damage"
                      fill="#8884d8"
                      background={{ fill: "rgba(255, 255, 255, 0.1)" }}
                    >
                      <LabelList
                        dataKey="damage"
                        position="right"
                        fill="#fff"
                      />
                    </Bar>
                  </BarChart>
                </ResponsiveContainer>
              </div>

              <div className="space-y-2 mt-6">
                {damageStats.map((stat, i) => (
                  <div key={i} className="flex items-center justify-between">
                    <div className="flex items-center">
                      <div
                        className={cn(
                          "w-6 h-6 rounded-full flex items-center justify-center mr-2",
                          getCostColor(stat.cost),
                        )}
                      >
                        <span className="text-xs font-bold">{stat.cost}</span>
                      </div>
                      <span className="font-medium">{stat.name}</span>
                      <div className="flex ml-2">
                        {Array.from({ length: stat.stars }).map((_, i) => (
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
                    <span className="font-mono">{stat.damage}</span>
                  </div>
                ))}
              </div>
            </>
          )}
        </ScrollArea>
      </CardContent>
    </Card>
  );
};

export default DamageStatsPanel;
