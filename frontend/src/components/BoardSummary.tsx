import React from "react";
import { useSimulator } from "../context/SimulatorContext";

const BoardSummary: React.FC = () => {
  const { state } = useSimulator();
  const { boardChampions, gold } = state;

  const championCount = boardChampions.length;

  return (
    <div className="rounded-lg p-4 h-12 mt-20">
      <div className="flex item-stretch gap-4">
        <div className="items-center flex flex-row gap-4">
          <img
            src="./board_unit_count.png"
            alt="Champion Icon"
            className="w-4 h-4 mx-auto"
          />
          <span className="font-semibold text-white">{championCount}</span>
        </div>
        <div className="items-center flex flex-row gap-4">
          <img src="./gold.png" alt="Gold Icon" className="w-4 h-4 mx-auto" />
          <span className="font-semibold text-yellow-400">{gold}</span>
        </div>
      </div>
    </div>
  );
};

export default BoardSummary;
