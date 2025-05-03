import React, { useState, useMemo } from 'react';
import { useSimulator } from '../context/SimulatorContext';
import { Input } from './ui/input';
import { ScrollArea } from './ui/scroll-area';
import { Card, CardContent, CardHeader, CardTitle } from './ui/card';
import { Button } from './ui/button';
import { SlidersHorizontal } from 'lucide-react';
import ChampionIcon from './ChampionIncon';

const ChampionPool: React.FC = () => {
  const { state } = useSimulator();
  const { champions } = state;
  const [searchTerm, setSearchTerm] = useState('');
  const [filterCost, setFilterCost] = useState<number | null>(null);
  const [filterTrait, setFilterTrait] = useState<string | null>(null);

  const filteredChampions = useMemo(() => {
    return champions.filter(champion => {
      const nameMatch = champion?.name?.toLowerCase().includes(searchTerm.toLowerCase());
      const costMatch = filterCost === null || champion.cost === filterCost;
      const traitMatch = filterTrait === null || champion.traits.includes(filterTrait);
      return nameMatch && costMatch && traitMatch;
    });
  }, [champions, searchTerm, filterCost, filterTrait]);

  // Extract unique traits for filtering
  const uniqueTraits = useMemo(() => {
    const traits = new Set<string>();
    champions.forEach(c => c.traits.forEach(t => traits.add(t)));
    return Array.from(traits).sort();
  }, [champions]);

  return (
    // Use Card component - Make card transparent
    <Card className="h-full flex flex-col border-none bg-transparent shadow-none">
      {/* Use CardHeader */}
      <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2 px-4 pt-4">
        {/* Use CardTitle */}
        <CardTitle className="text-base font-semibold">Champions</CardTitle> {/* Adjusted size */}
        {/* Filter buttons/tabs can go here - simplified for now */}
        <div className="flex items-center space-x-1">
          {[1, 2, 3, 4, 5].map(cost => (
            <Button
              key={cost}
              variant={filterCost === cost ? 'secondary' : 'ghost'}
              size="sm"
              className="w-6 h-6 p-0 text-xs"
              onClick={() => setFilterCost(prev => prev === cost ? null : cost)}
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
            <span className="sr-only">Filters</span>
          </Button>
        </div>
      </CardHeader>
      {/* Use CardContent */}
      <CardContent className="flex-1 flex flex-col p-4 pt-0">
        {/* Use Input component for search */}
        <Input
          type="text"
          placeholder="Search by name/trait..."
          value={searchTerm}
          onChange={(e) => setSearchTerm(e.target.value)}
          className="mb-2 h-8 bg-muted border-none" // Adjusted style
        />
        {/* Use ScrollArea component */}
        <ScrollArea className="flex-1 mt-0 focus-visible:ring-0 focus-visible:ring-offset-0">
          <div className="grid grid-cols-6 md:grid-cols-8 lg:grid-cols-10 gap-2 "> 
            {filteredChampions.map((champion) => (
              <ChampionIcon key={champion.apiName} champion={champion} />
            ))}
          </div>
           {filteredChampions.length === 0 && (
             <p className="text-center text-muted-foreground text-sm mt-4">No champions found.</p>
           )}
        </ScrollArea>
      </CardContent>
    </Card>
  );
};

export default ChampionPool;
