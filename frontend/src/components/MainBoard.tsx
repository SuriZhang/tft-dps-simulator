import TraitTracker from "./TraitTracker";
import HexBoard from "./HexBoard";
import ChampionPool from "./ChampionPool";
import ItemTray from "./ItemTray";
import AugmentPanel from "./AugmentPanel";
import ControlBar from "./ControlBar";
import DamageStatsPanel from "./DamageStatsPanel";
// Removed Resizable imports
import { ScrollArea } from "./ui/scroll-area";

const MainBoard = () => {
	return (
		<div className="flex h-[90%]">
			<div className="flex flex-col h-screen w-[75%] p-4 gap-4 bg-card rounded-l-lg">
				<div className="flex h-[60%]">
					<div className="w-[25%] flex h-full flex-col p-2 gap-2 bg-card rounded-l-lg">
						<ScrollArea>
							<TraitTracker />
						</ScrollArea>
					</div>

					<div className="w-[75%] h-full flex flex-col items-center px-4">
						<ControlBar />
						<HexBoard />
					</div>
				</div>

				<div className="h-[40%] flex flex-row gap-4 bg-panel-bg p-2">
					{/* ChampionPool takes available space */}
					<div className="h-full w-[60%]">
						<ScrollArea className="h-full">
							<ChampionPool />
						</ScrollArea>
					</div>
					{/* ItemTray takes available space */}
					<div className="h-full w-[40%]">
						<ScrollArea className="h-full">
							<ItemTray />
						</ScrollArea>
					</div>
				</div>
			</div>

			<div className="w-[25%] flex flex-col p-4 gap-4 bg-card rounded-r-lg">
				<DamageStatsPanel />
				<AugmentPanel />
			</div>
		</div>
	);
};

export default MainBoard;
