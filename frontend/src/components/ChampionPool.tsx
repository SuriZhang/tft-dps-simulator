import React, { useState, useMemo } from "react";
import { useSimulator } from "../context/SimulatorContext";
import { Input } from "./ui/input";
import { ScrollArea } from "./ui/scroll-area";
import { Card, CardContent, CardHeader, CardTitle } from "./ui/card";
import { Button } from "./ui/button";
import { SlidersHorizontal } from "lucide-react";
import ChampionIcon from "./ChampionIcon";

const ChampionPool: React.FC = () => {
  const { state } = useSimulator();
  const { champions } = state;
  const [searchTerm, setSearchTerm] = useState("");
  const [filterCost, setFilterCost] = useState<number | null>(null);
  const [filterTrait, _setFilterTrait] = useState<string | null>(null);

  const filteredChampions = useMemo(() => {
    return champions
      .filter((champion) => {
        // Check if search term matches champion name or any trait
        const searchTermLower = searchTerm.toLowerCase();
        const nameMatch = champion?.name
          ?.toLowerCase()
          .includes(searchTermLower);
        const traitMatch = champion.traits.some((trait) =>
          trait.toLowerCase().includes(searchTermLower),
        );

        // Apply cost filter
        const costMatch = filterCost === null || champion.cost === filterCost;

        // Apply trait filter dropdown (separate from search)
        const traitFilterMatch =
          filterTrait === null || champion.traits.includes(filterTrait);
        

        // Show champions that match by name OR trait AND match the filters
        return (nameMatch || traitMatch) && costMatch && traitFilterMatch &&
          champion.cost <= 6 && !champion.apiName.includes("Summon"); // exclude trait spawns
      })
      .sort((a, b) => a.cost - b.cost); // Sort by cost (ascending)
  }, [champions, searchTerm, filterCost, filterTrait]);

  return (
    <Card className="h-full flex flex-col border-none bg-transparent shadow-none">
      <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2 px-4 pt-4">
        <CardTitle className="text-base font-semibold">Champions</CardTitle>
        <div className="flex flex-row">
        <Input
          type="text"
          placeholder="Search by name/trait..."
          value={searchTerm}
          onChange={(e) => setSearchTerm(e.target.value)}
          className="mb-2 h-8 w-300 bg-muted border-none"
        />
        <div className="flex items-center space-x-1">
          {[1, 2, 3, 4, 5].map((cost) => (
            <Button
              key={cost}
              variant={filterCost === cost ? "secondary" : "ghost"}
              size="sm"
              className="w-6 h-6 p-0 text-xs"
              onClick={() =>
                setFilterCost((prev) => (prev === cost ? null : cost))
              }
            >
              {cost}
            </Button>
          ))}
          <Button
            variant="ghost"
            size="icon"
            className="w-6 h-6 ml-2"
            // onClick={() => {/* TODO: Open trait filter dropdown */}}
          >
            <SlidersHorizontal className="h-3 w-3" />
            </Button>
            </div>
        </div>
      </CardHeader>
      <CardContent className="flex-1 overflow-hidden flex flex-col p-4 pt-0">
        <ScrollArea className="flex-1 h-40 mt-0 focus-visible:ring-0 focus-visible:ring-offset-0 p-0">
          <div className="grid sm:grid-cols-6 md:grid-cols-10 lg:grid-cols-10 gap-2 overflow-hidden ">
            {filteredChampions.length === 0 ? (
              <p className="text-center text-muted-foreground text-sm mt-4">
                No champions found.
              </p>
            ) : (
              filteredChampions.map((champion) => (
                <ChampionIcon key={champion.apiName} champion={champion} />
              ))
            )}
          </div>
        </ScrollArea>
      </CardContent>
    </Card>
  );
};

export default ChampionPool;
