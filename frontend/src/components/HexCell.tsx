import React from "react";
import { BoardPosition, Champion, Item } from "../utils/types";
import { useSimulator } from "../context/SimulatorContext";
import { cn } from "../lib/utils";
import {
	ContextMenu,
	ContextMenuContent,
	ContextMenuItem,
	ContextMenuSeparator,
	ContextMenuTrigger,
} from "./ui/context-menu";
import { Star, Trash2 } from "lucide-react";

interface HexCellProps {
	row: number;
	col: number;
	champion?: Champion | null;
}

const HexCell : React.FC<HexCellProps> =({ row, col, champion } ) => {

	const position: BoardPosition = { row, col };
	const { state, dispatch } = useSimulator();
	const { selectedChampion, selectedItem } = state;

	const getHexBackground = () => {
		let base = "bg-gray-900/40 border border-gray-700/30";
		if (champion) {
			switch (champion.cost) {
				case 1:
					base = "bg-gray-800/60 border border-gray-600/80";
					break;
				case 2:
					base = "bg-green-900/60 border border-green-600/50";
					break;
				case 3:
					base = "bg-blue-900/60 border border-blue-500/50";
					break;
				case 4:
					base = "bg-purple-900/60 border border-purple-500/50";
					break;
				case 5:
					base = "bg-amber-900/60 border border-amber-500/50";
					break;
			}
		}
		return base;
	};

	// Click / drag handlers
	const handleCellClick = () => {
		if (selectedChampion && !champion) {
			dispatch({
				type: "ADD_CHAMPION_TO_BOARD",
				champion: selectedChampion,
				position,
      });
      console.log("Adding champion to board", selectedChampion, position);
		} else if (champion && selectedItem) {
			dispatch({
				type: "ADD_ITEM_TO_CHAMPION",
				item: selectedItem,
				position,
      });
      console.log("Adding item to champion", selectedItem, position);
      console.log("Champion", champion);
		}
	};
  const handleStarUp = () => 
    champion && dispatch({ type: "STAR_UP_CHAMPION", position });
  
  
	const handleRemove = () =>
    champion && dispatch({ type: "REMOVE_CHAMPION_FROM_BOARD", position });
  
  const handleDragOver = (e: React.DragEvent) => e.preventDefault();
  
	const handleDragStart = (e: React.DragEvent) => {
		if (champion) {
			e.dataTransfer.setData(
				"application/json",
				JSON.stringify({ type: "boardChampion", position })
			);
		}
	};
	const handleDrop = (e: React.DragEvent) => {
		e.preventDefault();
		try {
			const data = JSON.parse(e.dataTransfer.getData("application/json"));
			if (data.type === "champion" && !champion) {
				dispatch({
					type: "ADD_CHAMPION_TO_BOARD",
					champion: data.champion,
					position,
        });
        console.log("Adding champion to board", data.champion, position);
			} else if (data.type === "boardChampion" && data.position) {
				dispatch({
					type: "MOVE_CHAMPION",
					from: data.position,
					to: position,
        });
        console.log("Moving champion", data.position, position);
			} else if (data.type === "item" && champion) {
				dispatch({
					type: "ADD_ITEM_TO_CHAMPION",
					item: data.item,
					position,
        });
        console.log("Adding item to champion", data.item, position);
        console.log("Champion", champion);
			}
		} catch {
			/* ignore parsing errors */
		}
	};

	return (
		<div className={cn("relative aspect-[1/1]", `col-start-${col} rows-start-${row}`)}>
			<div
				className={cn(
					"w-[80px] h-[80px] inset-0 clip-hexagon border shadow-md transition-all cursor-pointer",
					getHexBackground(),
					!champion && selectedChampion
						? "border-primary border-2 hover:border-opacity-100"
						: "",
					champion && selectedItem
						? "border-accent border-2 hover:border-opacity-100"
						: ""
				)}
				title={`Row ${row}, Col ${col}`}
				onClick={handleCellClick}
				onDragOver={handleDragOver}
				onDrop={handleDrop}
        data-position={`${position.row}-${position.col}`}>
        
				{champion && (
					<ContextMenu>
						<ContextMenuTrigger asChild>
							<div
								className="absolute w-full h-full flex flex-col items-center justify-center rotate-[-90deg]"
								draggable
								onDragStart={handleDragStart}>
								{/* cost‐tint background */}
								{/* <div
									className={cn(
										"absolute inset-0 opacity-50",
										champion.cost === 1 && "bg-gray-500",
										champion.cost === 2 && "bg-green-500",
										champion.cost === 3 && "bg-blue-500",
										champion.cost === 4 && "bg-purple-500",
										champion.cost === 5 && "bg-amber-500"
									)}
								/> */}

								{/* name */}
								<p className="z-10 text-xs font-bold text-white">
									{champion.name}
								</p>

								{/* stars */}
								<div className="absolute bottom-0 left-0 w-full flex justify-center gap-0.5 z-20">
									{Array.from({
										length: champion.stars || 1,
									}).map((_, i) => (
										<div
											key={i}
											className="star bg-warning w-2.5 h-2.5"
											title={`${champion.stars || 1}★`}
										/>
									))}
								</div>

								{/* items */}
								{champion && champion.items && champion.items.length > 0 && (
									<div className="top-1 left-0 w-full flex justify-center gap-0.5 z-20">
										{champion.items.map((item: Item, i: number) => (
											<div
												key={i}
												className="w-3 h-3 rounded-sm bg-yellow-500"
												title={item.name}
											/>
										))}
									</div>
								)}
							</div>
						</ContextMenuTrigger>

						<ContextMenuContent>
							<ContextMenuItem onClick={handleStarUp}>
								<Star className="mr-2 h-4 w-4" />
								Star Up ({champion.stars || 1}★)
							</ContextMenuItem>
							<ContextMenuSeparator />
							<ContextMenuItem
								onClick={handleRemove}
								className="text-destructive">
								<Trash2 className="mr-2 h-4 w-4" />
								Remove
							</ContextMenuItem>
						</ContextMenuContent>
					</ContextMenu>
				)}
			</div>
		</div>
	);
};

export default HexCell;
