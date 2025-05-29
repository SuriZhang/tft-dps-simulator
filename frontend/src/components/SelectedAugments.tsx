import { Plus } from "lucide-react";
import { useSimulator } from "../context/SimulatorContext";
import { cn } from "../lib/utils";

const SelectedAugments = () => {
  const { state} = useSimulator();
  const { selectedAugments } = state;

  return (
    <div className="ml-6 p-2 border rounded-lg bg-card flex flex-col gap-2">
      <span>Selected Augments</span>
       <div className="flex flex-row items-center gap-2">
            {/* Selected augments */}
            {[0, 1, 2, 3, 4].map((slot) => {
              const augment = selectedAugments[slot];

              return (
                <div
                  key={`slot-${slot}`}
                  className={cn(
                    "p-2 rounded-lg border border-dashed cursor-pointer w-16 h-16 flex items-center justify-center",
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
                      <Plus className="w-6 h-6" />
                    </div>
                  )}
                </div>
              );
            })}
          </div>
    </div>
  );
};

export default SelectedAugments;
