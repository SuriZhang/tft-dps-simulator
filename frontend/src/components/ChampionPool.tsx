import React, { useState, useMemo } from "react";
import { useSimulator } from "../context/SimulatorContext";
import { Input } from "./ui/input";
import { ScrollArea } from "./ui/scroll-area";
import { Card, CardContent, CardHeader, CardTitle } from "./ui/card";
import { Button } from "./ui/button";
import { SlidersHorizontal, CircleX } from "lucide-react";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuTrigger,
  DropdownMenuCheckboxItem,
} from "./ui/dropdown-menu";
import ChampionIcon from "./ChampionIcon";
import { Tooltip, TooltipContent, TooltipTrigger } from "./ui/tooltip";
import { TooltipProvider } from "@radix-ui/react-tooltip";

const ChampionPool: React.FC = () => {
  const { state } = useSimulator();
  const { champions } = state;
  const [searchTerm, setSearchTerm] = useState("");
  const [filterCost, setFilterCost] = useState<number | null>(null);
  const [selectedTraits, setSelectedTraits] = useState<string[]>([]);

  // Get all unique traits from champions
  const allTraits = useMemo(() => {
    const traits = new Set<string>();
    champions.forEach((champion) => {
      champion.traits.forEach((trait) => traits.add(trait));
    });
    return Array.from(traits).sort();
  }, [champions]);

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
          selectedTraits.length === 0 ||
          champion.traits.some((trait) => selectedTraits.includes(trait));

        // Show champions that match by name OR trait AND match the filters
        return (
          (nameMatch || traitMatch) &&
          costMatch &&
          traitFilterMatch &&
          champion.cost <= 6 &&
          !champion.apiName.includes("Summon")
        ); // exclude trait spawns
      })
      .sort((a, b) => a.cost - b.cost); // Sort by cost (ascending)
  }, [champions, searchTerm, filterCost, selectedTraits]);

  const handleTraitToggle = (trait: string) => {
    setSelectedTraits((prev) =>
      prev.includes(trait) ? prev.filter((t) => t !== trait) : [...prev, trait],
    );
  };

  const clearTraitFilters = () => {
    setSelectedTraits([]);
  };

  return (
    <Card className="h-full flex flex-col border-none bg-transparent shadow-none">
      <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2 px-4 pt-4 h-18">
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
            <DropdownMenu>
              <DropdownMenuTrigger asChild>
                <Button
                  variant={selectedTraits.length > 0 ? "secondary" : "ghost"}
                  size="icon"
                  className="w-6 h-6 ml-2"
                >
                  <SlidersHorizontal className="h-3 w-3" />
                </Button>
              </DropdownMenuTrigger>
              <DropdownMenuContent className="w-56 max-h-96 overflow-y-auto">
                {allTraits.map((trait) => (
                  <DropdownMenuCheckboxItem
                    key={trait}
                    checked={selectedTraits.includes(trait)}
                    onCheckedChange={() => handleTraitToggle(trait)}
                  >
                    {trait}
                  </DropdownMenuCheckboxItem>
                ))}
              </DropdownMenuContent>
            </DropdownMenu>
            {selectedTraits.length > 0 && (
              <TooltipProvider>
                <Tooltip>
                  <TooltipTrigger asChild>
                    <Button
                      variant="ghost"
                      size="icon"
                      className="w-6 h-6 ml-1"
                      onClick={clearTraitFilters}
                    >
                      <CircleX className="h-3 w-3" />
                    </Button>
                  </TooltipTrigger>
                  <TooltipContent>
                    <p>Reset trait filters</p>
                  </TooltipContent>
                </Tooltip>
              </TooltipProvider>
            )}
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
