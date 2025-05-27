import React from "react";
import { useSimulator } from "../context/SimulatorContext";
import { cn } from "../lib/utils";

const AugmentPanel: React.FC = () => {
  const { state} = useSimulator();
  const { selectedAugments } = state;

  return (
    <div className="ml-6 p-2">
      <span>Selected Augments</span>
       <div className="flex flex-row items-center gap-2">
            {/* Selected augments */}
            {[0, 1, 2].map((slot) => {
              const augment = selectedAugments[slot];

              return (
                <div
                  key={`slot-${slot}`}
                  className={cn(
                    "p-2 rounded-lg border cursor-pointer w-16 h-16 flex items-center justify-center",
                  )}
                >
                  {augment ? (
                    <div className="flex items-center justify-center">
                      {augment.icon && (
                        <img
                          src={`/tft-augment/${augment.icon}`}
                          alt={augment.name}
                          className="w-8 h-8 object-cover rounded"
                        />
                      )}
                    </div>
                  ) : (
                    <div className="text-xs text-gray-500 text-center">
                      Select augment
                    </div>
                  )}
                </div>
              );
            })}
          </div>
    </div>
  );
};

export default AugmentPanel;
