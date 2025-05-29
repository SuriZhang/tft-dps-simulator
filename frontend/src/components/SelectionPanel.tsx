import { Tabs, TabsContent, TabsList, TabsTrigger } from "./ui/tabs";
import { Card } from "./ui/card";
import ChampionPool from "./ChampionPool";
import ItemTray from "./ItemTray";
import AugmentTray from "./AugmentTray";

const SelectionPanel = () => {
  return (
    <Card className="h-full w-full border-none bg-panel-bg shadow-none">
      <Tabs defaultValue="champions-items" className="h-full flex flex-col">
        <TabsList className="bg-muted w-fit mx-4 mt-4">
          <TabsTrigger 
            value="champions-items" 
            className="data-[state=active]:bg-background data-[state=active]:text-foreground"
          >
            Champions + Items
          </TabsTrigger>
          <TabsTrigger 
            value="augments"
            className="data-[state=active]:bg-background data-[state=active]:text-foreground"
          >
            Augments
          </TabsTrigger>
        </TabsList>
        
        <TabsContent value="champions-items" className="flex-1 mt-0 p-4 pt-2">
          <div className="h-full flex flex-row gap-4 bg-card rounded-lg p-2">
            <div className="w-[60%] h-full">
              <ChampionPool />
            </div>
            <div className="w-[40%] h-full">
              <ItemTray />
            </div>
          </div>
        </TabsContent>
        
        <TabsContent value="augments" className="h-full flex-1 gap-4 bg-card rounded-lg p-2">
          <div className="h-full">
            <AugmentTray />
          </div>
        </TabsContent>
      </Tabs>
    </Card>
  );
};

export default SelectionPanel;