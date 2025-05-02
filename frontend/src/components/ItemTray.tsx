import React, { useState, useMemo } from 'react';
import { useSimulator } from '../context/SimulatorContext';
import ItemIcon from './ItemIcon';
import { Input } from './ui/input';
import { ScrollArea } from './ui/scroll-area';
import { Card, CardContent } from './ui/card'; // Removed CardHeader, CardTitle
import { Tabs, TabsContent, TabsList, TabsTrigger } from "./ui/tabs";
import { Item } from '../utils/types';

// Define item categories based on available type and potential naming conventions
type ItemCategory = 'component' | 'craftable' | 'radiant' | 'ornn' | 'support' | 'other';

const getItemCategory = (item: Item): ItemCategory => {
  // Use the existing 'type' field first
  if (item.type === 'component') return 'component';
  if (item.type === 'completed') {
    // Further differentiate completed items if possible (example using name patterns)
    // This is placeholder logic - adjust based on actual data/naming conventions
    const lowerName = item.name.toLowerCase();
    if (lowerName.includes('(radiant)') || item.id.startsWith('TFT_Item_Radiant')) return 'radiant';
    if (lowerName.includes('(ornn)') || item.id.startsWith('TFT_Item_Ornn')) return 'ornn';
    if (lowerName.includes('(support)') || item.id.startsWith('TFT_Item_Support')) return 'support';
    return 'craftable'; // Default completed items
  }
  if (item.type === 'special') {
    // Categorize special items (e.g., Emblems, Tactician's Crown)
    // Placeholder logic - adjust as needed
     const lowerName = item.name.toLowerCase();
    if (lowerName.includes('emblem')) return 'other'; // Or a specific 'emblem' category?
    if (lowerName.includes('crown')) return 'other';
    return 'other';
  }

  // Fallback category
  return 'other';
};


const ItemTray: React.FC = () => {
  const { state } = useSimulator();
  const { items } = state;
  const [searchTerm, setSearchTerm] = useState('');
  const [activeTab, setActiveTab] = useState<ItemCategory>('craftable');

  const filteredItems = useMemo(() => {
    return items.filter(item => {
      const nameMatch = item.name.toLowerCase().includes(searchTerm.toLowerCase());
      const categoryMatch = getItemCategory(item) === activeTab;
      return nameMatch && categoryMatch;
    });
  }, [items, searchTerm, activeTab]);

  // Define tab order and filter out categories with no items
  const availableCategories = useMemo(() => {
    const allCats: ItemCategory[] = ['craftable', 'radiant', 'ornn', 'support', 'component', 'other'];
    const presentCats = new Set(items.map(getItemCategory));
    return allCats.filter(cat => presentCats.has(cat));
  }, [items]);

  // Adjust active tab if the current one becomes unavailable
  React.useEffect(() => {
    if (!availableCategories.includes(activeTab)) {
      setActiveTab(availableCategories[0] || 'craftable');
    }
  }, [availableCategories, activeTab]);

  return (
    <Card className="h-full flex flex-col border-none bg-transparent shadow-none">
      <Tabs value={activeTab} onValueChange={(value) => setActiveTab(value as ItemCategory)} className="flex flex-col flex-1">
        <div className="flex justify-between items-center px-4 pt-4 pb-2"> {/* Container for TabsList and Search */}
          <TabsList className="bg-muted">
            {availableCategories.map(category => (
              <TabsTrigger key={category} value={category} className="capitalize text-xs px-3 py-1 h-auto data-[state=active]:bg-background data-[state=active]:text-foreground">
                {/* Improved Naming */}
                {category === 'craftable' ? 'Craftable' :
                 category === 'ornn' ? 'Ornn' :
                 category === 'support' ? 'Support' :
                 category === 'radiant' ? 'Radiant' :
                 category === 'component' ? 'Components' :
                 'Other'}
              </TabsTrigger>
            ))}
          </TabsList>
          <Input
            type="text"
            placeholder="Search items..."
            value={searchTerm}
            onChange={(e) => setSearchTerm(e.target.value)}
            className="h-8 w-40 ml-4 bg-muted border-none placeholder:text-muted-foreground/80 text-sm"
          />
        </div>

        {availableCategories.map(category => (
          <TabsContent key={category} value={category} className="flex-1 mt-0 focus-visible:ring-0 focus-visible:ring-offset-0">
            <ScrollArea className="h-full px-4 pb-4">
              <div className="grid grid-cols-6 md:grid-cols-8 lg:grid-cols-10 gap-2">
                {filteredItems.map((item) => (
                  <ItemIcon key={item.id} item={item} size="sm" />
                ))}
              </div>
              {filteredItems.length === 0 && (
                 <p className="text-center text-muted-foreground text-sm mt-4">No items found.</p>
              )}
            </ScrollArea>
          </TabsContent>
        ))}
      </Tabs>
    </Card>
  );
};

export default ItemTray;
