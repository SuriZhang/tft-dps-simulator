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
  | "other";

const getItemCategory = (item: Item): ItemCategory => {
  // Use the existing 'type' field first
  if (item.tags?.includes("component")) return "component";
  if (item.tags?.includes("{7ea41d13}")) return "craftable"; // composition

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
    // Start with 'all'
    const allCats: ItemCategory[] = [
      "all",
      "craftable",
      "radiant",
      "ornn",
      "support",
      "component",
      "other",
    ];
    const presentCats = new Set(items.map(getItemCategory));
    // Filter other categories based on presence, but always keep 'all'
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
        <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2 px-4 pt-4">
          <CardTitle className="text-base font-semibold">Items</CardTitle>

          <TabsList className="bg-muted">
            {availableCategories.map((category) => (
              <TabsTrigger
                key={category}
                value={category}
                className="capitalize text-xs px-3 py-1 h-auto data-[state=active]:bg-background data-[state=active]:text-foreground"
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
                          : category === "component"
                            ? "Components"
                            : "Other"}
              </TabsTrigger>
            ))}
          </TabsList>
        </CardHeader>
        <CardContent className="flex-1 overflow-hidden flex flex-col p-4 pt-0">
          <Input
            type="text"
            placeholder="Search items..."
            value={searchTerm}
            onChange={(e) => setSearchTerm(e.target.value)}
            className="mb-2 h-8 bg-muted border-none"
          />

          {/* Fixed height container with explicit overflow handling */}
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
          {/* </div> */}
        </CardContent>
      </Tabs>
    </Card>
  );
};

export default ItemTray;
