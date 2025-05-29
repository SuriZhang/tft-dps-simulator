import { Card, CardContent, CardHeader, CardTitle } from "./ui/card";
import { Tabs, TabsList, TabsTrigger, TabsContent } from "./ui/tabs";
import {
  TooltipProvider,
  Tooltip,
  TooltipTrigger,
  TooltipContent,
} from "./ui/tooltip";
import { Info } from "lucide-react";
import { ScrollArea } from "./ui/scroll-area";
import SimulationTimelineChart from "./SimulationTimelineChart";
import DamageBarChart from "./DamageBarChart";

const DamageStatsPanel = () => {
  return (
    <Card className="mb-4">
      <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
        <CardTitle className="text-xl font-bold">Damage Analysis</CardTitle>
        <TooltipProvider>
          <Tooltip>
            <TooltipTrigger asChild>
              <Info className="h-4 w-4 text-muted-foreground cursor-help" />
            </TooltipTrigger>
            <TooltipContent>
              <p>
                View damage statistics and timeline visualization of simulation
                results
              </p>
            </TooltipContent>
          </Tooltip>
        </TooltipProvider>
      </CardHeader>

      <CardContent className="p-4">
        <Tabs defaultValue="damage" className="w-full">
          <TabsList className="grid w-full grid-cols-2">
            <TabsTrigger value="damage">Damage Statistics</TabsTrigger>
            <TabsTrigger value="timeline">Timeline Chart</TabsTrigger>
          </TabsList>

          <TabsContent value="damage" className="mt-4">
            <DamageBarChart />
          </TabsContent>

          <TabsContent value="timeline" className="mt-4">
            <ScrollArea className="h-[420px]">
              <SimulationTimelineChart />
            </ScrollArea>
          </TabsContent>
        </Tabs>
      </CardContent>
    </Card>
  );
};

export default DamageStatsPanel;
