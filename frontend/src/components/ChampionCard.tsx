import React from 'react';
import { Champion } from '../utils/types';
import { cn } from '../lib/utils';
import { useSimulator } from '../context/SimulatorContext';
import { Card, CardContent } from './ui/card'; // Import Card components
import { Badge } from './ui/badge'; // Import Badge
import { Avatar, AvatarFallback, AvatarImage } from './ui/avatar'; // Import Avatar components

interface ChampionCardProps {
  champion: Champion;
  size?: 'sm' | 'md' | 'lg';
}

const ChampionCard: React.FC<ChampionCardProps> = ({ champion, size = 'md' }) => {
  const { dispatch, state } = useSimulator();
  const isSelected = state.selectedChampion?.id === champion.id;

  const costColors = {
    1: 'border-gray-500 bg-gray-500/20',
    2: 'border-green-500 bg-green-500/20',
    3: 'border-blue-500 bg-blue-500/20',
    4: 'border-purple-500 bg-purple-500/20',
    5: 'border-amber-500 bg-amber-500/20',
  };

  const sizeClasses = {
    sm: 'w-10 h-10 text-xs',
    md: 'w-14 h-14 text-sm',
    lg: 'w-20 h-20 text-base',
  };

  const handleClick = () => {
    dispatch({
      type: 'SELECT_CHAMPION',
      champion: isSelected ? undefined : champion,
    });
  };

  const handleDragStart = (e: React.DragEvent) => {
    e.dataTransfer.setData('application/json', JSON.stringify({
      type: 'champion',
      champion
    }));
  };

  return (
    // Use Card component
    <Card
      className={cn(
        'cursor-pointer transition-all duration-200 card-hover',
        costColors[champion.cost] || 'border-gray-500',
        'border-2 bg-card/80'
      )}
      onClick={handleClick}
      draggable
      onDragStart={handleDragStart}
    >
      {/* Use CardContent */}
      <CardContent className="p-2 flex items-center space-x-2">
        {/* Use Avatar component */}
        <Avatar className="h-10 w-10 rounded-sm">
          <AvatarImage src={champion.image || '/placeholder.svg'} alt={champion.name} />
          <AvatarFallback>{champion.name.substring(0, 1)}</AvatarFallback>
        </Avatar>
        <div className="flex-1 min-w-0">
          <p className="text-sm font-medium text-foreground truncate">{champion.name}</p>
          <div className="flex space-x-1 mt-1">
            {/* Use Badge component for traits */}
            {champion.traits.map((trait) => (
              <Badge key={trait} variant="secondary" className="text-xs px-1.5 py-0.5">
                {trait}
              </Badge>
            ))}
          </div>
        </div>
        {/* Use Badge for cost */}
        <Badge variant="outline" className={cn("text-xs font-bold", costColors[champion.cost])}>
          ${champion.cost}
        </Badge>
      </CardContent>
    </Card>
  );
};

export default ChampionCard;
