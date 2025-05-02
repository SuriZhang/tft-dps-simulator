import React from 'react';
import { useSimulator } from '../context/SimulatorContext';
import { cn } from '../lib/utils';
import { Button } from './ui/button';
import { Card, CardContent, CardHeader, CardTitle } from './ui/card'; // Import Card components
import { Separator } from './ui/separator'; // Import Separator
import { Plus } from 'lucide-react'; // Import Plus icon
import { ScrollArea } from './ui/scroll-area';

const AugmentPanel: React.FC = () => {
  const { state, dispatch } = useSimulator();
  const { augments, selectedAugments } = state;

  const tierColors = {
    silver: 'border-gray-400 bg-gray-800/70 text-gray-300',
    gold: 'border-amber-400 bg-amber-900/30 text-amber-300',
    prismatic: 'border-cyan-400 bg-cyan-900/30 text-cyan-300',
  };

  const handleAugmentClick = (augment: typeof augments[0], slot: number) => {
    dispatch({
      type: 'SELECT_AUGMENT',
      augment,
      index: slot
    });
  };

  return (
    // Use Card component for the main container
    <Card className="mb-4">
      {/* Use CardHeader for the title and button */}
      <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
        {/* Use CardTitle for the heading */}
        <CardTitle className="text-xl font-bold">Augments</CardTitle>
        {/* Use lucide-react icon */}
        <Button variant="outline" size="icon" className="w-8 h-8 rounded-full bg-primary/20 text-primary hover:bg-primary/30">
          <Plus className="h-4 w-4" />
          <span className="sr-only">Add Augment</span>
        </Button>
      </CardHeader>
      {/* Use CardContent for the body */}
      <CardContent>
        <ScrollArea className="h-48 mb-4">
        <div className="space-y-2">
          {/* Selected augments */}
          {[0, 1, 2].map((slot) => {
            const augment = selectedAugments[slot];

            return (
              <div
                key={`slot-${slot}`}
                className={cn(
                  "p-2 rounded-lg border cursor-pointer",
                  augment
                    ? tierColors[augment.tier]
                    : "border-dashed border-gray-600 bg-gray-800/20"
                )}
                onClick={() => {
                  // Simple implementation - just cycles through augments
                  const nextIndex = augments.findIndex(a => a.id === augment?.id) + 1;
                  const nextAugment = augments[nextIndex % augments.length];
                  if (nextAugment) handleAugmentClick(nextAugment, slot);
                }}
              >
                {augment ? (
                  <div className="flex flex-col">
                    <div className="text-sm font-bold">{augment.name}</div>
                    <div className="text-xs">{augment.description}</div>
                  </div>
                ) : (
                  <div className="h-12 flex items-center justify-center text-sm text-gray-500">
                    Select an augment...
                  </div>
                )}
              </div>
            );
          })}
        </div>

        {/* Use Separator component */}
        <Separator className="my-4 bg-gray-700" />

        {/* Available Augments Section */}
        <div>
          <h3 className="text-sm font-semibold text-gray-400 mb-2">Available Augments</h3>
          <div className="grid grid-cols-3 gap-2">
            {augments.map((augment) => (
              <div
                key={augment.id}
                className={cn(
                  "p-2 rounded-md text-xs cursor-pointer",
                  tierColors[augment.tier]
                )}
                onClick={() => {
                  // Find first empty slot or overwrite last one
                  const emptySlot = [0, 1, 2].find(i => !selectedAugments[i]);
                  const slotToUse = emptySlot !== undefined ? emptySlot : 2;
                  handleAugmentClick(augment, slotToUse);
                }}
              >
                {augment.name}
              </div>
            ))}
          </div>
          </div>
          </ScrollArea>
      </CardContent>
    </Card>
  );
};

export default AugmentPanel;
