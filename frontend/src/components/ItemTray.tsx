import React, { useState, useMemo } from "react";
import { useSimulator } from "../context/SimulatorContext";
import ItemIcon from "./ItemIcon";
import { Input } from "./ui/input";
import { ScrollArea } from "./ui/scroll-area";
import { Card, CardContent, CardHeader, CardTitle } from "./ui/card"; // Removed CardHeader, CardTitle
import { Tabs, TabsContent, TabsList, TabsTrigger } from "./ui/tabs";
import { Item } from "../utils/types";

// Define item categories based on available type and potential naming conventions
type ItemCategory =
  | "all"
  | "component"
  | "craftable"
  | "radiant"
  | "ornn"
  | "support"
  | "emblem"
  | "other";

const getItemCategory = (item: Item): ItemCategory => {
  // Use the existing 'type' field first
  
  // Note: this might change if the item tags changes in json file
  if (item.name?.includes("Emblem")) return "emblem"; // emblem
  if (item.tags?.includes("component")) return "component";
  if (item.tags?.includes("{7ea41d13}")) return "craftable"; // composition
  if (item.tags?.includes("{27557a09}")) return "support"; // support
  if (item.tags?.includes("{44ace175}")) return "ornn"; // artifact
  if (item.tags?.includes("{6ef5c598}") || item.tags?.includes("{ec243f6b}"))
    return "radiant"; // radiant

  // Fallback category
  return "other";
};

const ItemTray: React.FC = () => {
  const { state } = useSimulator();
  const { items } = state;
  const [searchTerm, setSearchTerm] = useState("");
  const [activeTab, setActiveTab] = useState<ItemCategory>("all");

  const filteredItems = useMemo(() => {
    return items.filter((item) => {
      // Exclude items categorized as "other"
      if (getItemCategory(item) === "other") {
        return false;
      }
      const nameMatch = item?.name
        ?.toLowerCase()
        .includes(searchTerm.toLowerCase());
      // Only apply category filter if the active tab is NOT 'all'
      const categoryMatch =
        activeTab === "all" ? true : getItemCategory(item) === activeTab;
      return nameMatch && categoryMatch;
    });
  }, [items, searchTerm, activeTab]);

  // Define tab order and filter out categories with no items
  const availableCategories = useMemo(() => {
    // Start with 'all' and define the order of other categories, excluding "other"
    const allCats: ItemCategory[] = [
      "all",
      "craftable",
      "radiant",
      "ornn",
      "support",
      // "component",
      "emblem"
      // "other", // "other" category is removed from here
    ];
    const presentCats = new Set(items.map(getItemCategory));
    // Filter other categories based on presence, but always keep 'all'
    // Categories not in allCats (like "other") will be ignored here.
    return allCats.filter((cat) => cat === "all" || presentCats.has(cat));
  }, [items]);

  // Adjust active tab if the current one becomes unavailable
  React.useEffect(() => {
    if (!availableCategories.includes(activeTab)) {
      setActiveTab(availableCategories[0] || "craftable");
    }
  }, [availableCategories, activeTab]);

  return (
    <Card className="h-full flex flex-col border-none bg-transparent shadow-none">
      <Tabs
        value={activeTab}
        onValueChange={(value) => setActiveTab(value as ItemCategory)}
        className="flex flex-col flex-1"
      >
        <CardHeader className="flex flex-col justify-between space-y-0 pb-2 px-0 pt-4">
          <CardTitle className="text-base font-semibold items-start ">Items</CardTitle>
          <div>
          <TabsList className="bg-muted items-stretch w-full">
            {availableCategories.map((category) => (
              <TabsTrigger
                key={category}
                value={category}
                className="capitalize text-xs px-1 py-1 h-auto data-[state=active]:bg-background data-[state=active]:text-foreground"
              >
                {category === "all"
                  ? "All"
                  : category === "craftable"
                    ? "Craftable"
                    : category === "ornn"
                      ? "Ornn"
                      : category === "support"
                        ? "Support"
                        : category === "radiant"
                          ? "Radiant"
                          : category === "emblem"
                              ? "Emblems"
                              : category}
              </TabsTrigger>
            ))}
            </TabsList>
            </div>
        </CardHeader>
        <CardContent className="flex-1 overflow-hidden flex flex-col p-2 pt-0">
          <ScrollArea className="flex-1 h-40 mt-0 focus-visible:ring-0 focus-visible:ring-offset-0 p-0">
            {availableCategories.map((category) => (
              <TabsContent
                key={category}
                value={category}
                className="mt-0 pt-0 h-40 data-[state=inactive]:hidden"
              >
                <div className="grid grid-cols-5 sm:grid-cols-6 md:grid-cols-7 lg:grid-cols-8 gap-2">
                  {filteredItems.length === 0 ? (
                    <p className="text-center text-muted-foreground text-sm mt-4 col-span-full">
                      No items found.
                    </p>
                  ) : (
                    filteredItems.map((item) => (
                      <ItemIcon key={item.apiName} item={item} size="sm" />
                    ))
                  )}
                </div>
              </TabsContent>
            ))}
          </ScrollArea>
        </CardContent>
      </Tabs>
    </Card>
  );
};

export default ItemTray;
