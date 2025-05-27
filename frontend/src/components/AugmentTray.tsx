import {useState} from "react";
import { useSimulator } from "../context/SimulatorContext";
import { cn } from "../lib/utils";
import { Button } from "./ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "./ui/card"; 
import {  SlidersHorizontal } from "lucide-react"; 
import { ScrollArea } from "./ui/scroll-area";
import { Item } from "../utils/types";
import { Input } from "./ui/input";

const AugmentTray = () => {
  const { state, dispatch } = useSimulator();
    const { augments, selectedAugments } = state;
    const [searchTerm, setSearchTerm] = useState("");

  const tierColors = {
    silver: "border-gray-400 bg-gray-800/70 text-gray-300",
    gold: "border-amber-400 bg-amber-900/30 text-amber-300",
    prismatic: "border-cyan-400 bg-cyan-900/30 text-cyan-300",
  };

  const handleAugmentClick = (augment: Item, slot: number) => {
    dispatch({
      type: "SELECT_AUGMENT",
      augment,
      index: slot,
    });
    };
    
// const filteredChampions = useMemo(() => {
//     return champions
//       .filter((champion) => {
//         // Check if search term matches champion name or any trait
//         const searchTermLower = searchTerm.toLowerCase();
//         const nameMatch = champion?.name
//           ?.toLowerCase()
//           .includes(searchTermLower);
//         const traitMatch = champion.traits.some((trait) =>
//           trait.toLowerCase().includes(searchTermLower),
//         );

//         // Apply cost filter
//         const costMatch = filterCost === null || champion.cost === filterCost;

//         // Apply trait filter dropdown (separate from search)
//         const traitFilterMatch =
//           filterTrait === null || champion.traits.includes(filterTrait);
        

//         // Show champions that match by name OR trait AND match the filters
//         return (nameMatch || traitMatch) && costMatch && traitFilterMatch &&
//           champion.cost <= 6 && !champion.apiName.includes("Summon"); // exclude trait spawns
//       })
//       .sort((a, b) => a.cost - b.cost); // Sort by cost (ascending)
//   }, [champions, searchTerm, filterCost, filterTrait]);

  return (
    // Use Card component for the main container
    <Card className="h-full flex flex-col border-none bg-transparent shadow-none">
      <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2 px-4 pt-4">
        <CardTitle className="text-base font-semibold">Augments</CardTitle>
        <div className="flex flex-row">
        <Input
          type="text"
          placeholder="Search by name"
          value={searchTerm}
          onChange={(e) => setSearchTerm(e.target.value)}
          className="mb-2 h-8 w-300 bg-muted border-none"
                  />
                <Button
            variant="ghost"
            size="icon"
            className="w-6 h-6 ml-2"
            // onClick={() => {/* TODO: Open trait filter dropdown */}}
          >
            <SlidersHorizontal className="h-3 w-3" />
            </Button>
              </div>
              
      </CardHeader>
      <CardContent>
        <ScrollArea className="flex-1 h-40 mt-0 focus-visible:ring-0 focus-visible:ring-offset-0 p-0">
          <div>
            <h3 className="text-sm font-semibold text-gray-400 mb-2">
              Available Augments
            </h3>
            <div className="grid grid-cols-3 gap-2">
              {augments.map((augment) => (
                <div
                  key={augment.apiName}
                  className={cn(
                    "p-2 rounded-md text-xs cursor-pointer flex items-center space-x-1",
                    tierColors.silver, // TODO: fix, Default to silver for all augments
                  )}
                  onClick={() => {
                    // Find first empty slot or overwrite last one
                    const emptySlot = [0, 1, 2].find(
                      (i) => !selectedAugments[i],
                    );
                    const slotToUse = emptySlot !== undefined ? emptySlot : 2;
                    handleAugmentClick(augment, slotToUse);
                  }}
                >
                  {augment.icon && (
                    <img
                      src={`/tft-augment/${augment.icon}`}
                      alt=""
                      className="w-4 h-4 object-cover"
                    />
                  )}
                  <span>{augment.name}</span>
                </div>
              ))}
            </div>
          </div>
        </ScrollArea>
      </CardContent>
    </Card>
  );
};

export default AugmentTray;
