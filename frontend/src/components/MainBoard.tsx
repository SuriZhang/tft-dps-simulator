import TraitTracker from "./TraitTracker";
import HexBoard from "./HexBoard";
import ChampionPool from "./ChampionPool";
import ItemTray from "./ItemTray";
import AugmentPanel from "./AugmentPanel";
import ControlBar from "./ControlBar";
import DamageStatsPanel from "./DamageStatsPanel";
import AugmentTray from "./AugmentTray";

const MainBoard = () => {
  return (
    <>
      <div className="flex flex-col h-[90%]">
        <div className="flex flex-row h-[55%]">
          <div className="flex flex-col w-[55%] p-2 gap-4 bg-card rounded-l-lg">
            <ControlBar />
            <div className="w-full h-full flex flex-row items-start">
                <div className="flex w-[20%] h-full">
                  <TraitTracker />
                </div> 
              <div className="flex-1 w-[80%]">
                <AugmentPanel />
                  <HexBoard />
                </div>
            </div>
          </div>
          <div className="flex-1 w-[45%] p-2 bg-card rounded-r-lg">
            <DamageStatsPanel />
            </div>

          </div>

        <div className="h-[45%] w-full flex flex-row gap-4 bg-panel-bg p-2">
          <div className="w-[45%] h-full">
            <ChampionPool />
          </div>
          <div className="w-[30%] h-full">
            <ItemTray />
          </div>
          <div className="w-[25%] f-ull">
            <AugmentTray />
            </div>
          </div>
        </div>
      
      {/* <div className="bg-background/50 p-6 rounded-lg shadow">
        <SimulationTimelineChart />
      </div> */}
    </>
  );
};

export default MainBoard;
