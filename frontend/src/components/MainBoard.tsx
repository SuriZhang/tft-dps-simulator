import TraitTracker from "./TraitTracker";
import HexBoard from "./HexBoard";
import ChampionPool from "./ChampionPool";
import ItemTray from "./ItemTray";
import AugmentPanel from "./AugmentPanel";
import ControlBar from "./ControlBar";
import DamageStatsPanel from "./DamageStatsPanel";
import SimulationTimelineChart from "./SimulationTimelineChart";

const MainBoard = () => {
  return (
    <>
    <div className="flex h-[90%]">
      <div className="flex flex-col h-screen w-[80%] p-4 gap-4 bg-card rounded-l-lg">
        <div className="flex h-[60%]">
          <div className="w-[30%] flex h-full flex-col p-2 gap-2 bg-card rounded-l-lg">
            <TraitTracker />
          </div>

          <div className="w-[70%] h-full flex flex-col items-center px-4">
            <ControlBar />
            <HexBoard />
          </div>
        </div>

        <div className="h-[40%] flex flex-row gap-4 bg-panel-bg p-2">
          <div className="h-full w-[60%]">
            <ChampionPool />
          </div>
          <div className="h-full w-[40%]">
            <ItemTray />
          </div>
        </div>
      </div>

      <div className="w-[20%] flex flex-col p-4 gap-4 bg-card rounded-r-lg">
        <div className="flex-2">
          <DamageStatsPanel />
        </div>
        <div className="flex-1">
          <AugmentPanel />
        </div>
      </div>
      </div>
              {/* Add the timeline chart */}
              <div className="bg-background/50 p-6 rounded-lg shadow">
        <SimulationTimelineChart />
      </div>
    </>
  );
};

export default MainBoard;
