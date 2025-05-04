import React, {
  createContext,
  useReducer,
  useContext,
  ReactNode,
  useEffect,
  useState,
} from "react";
import {
  Champion,
  BoardPosition,
  Item,
  SimulatorState,
  SimulatorAction,
  BoardChampion,
  Trait,
} from "../utils/types";
import {
  loadSetDataFromFile,
  loadItemsAndAugmentsFromFile,
} from "../utils/constants";

const initialState: SimulatorState = {
  champions: [], // Start with empty data
  boardChampions: [],
  items: [],
  traits: [],
  augments: [],
  selectedAugments: [],
  gold: 0, // Start with 0 gold (no champions)
  level: 2, // Start at base level 2 (no champions)
  selectedItem: undefined,
  selectedChampion: undefined,
  loading: true,
  error: undefined,
  hoveredTrait: "",
};

// Define the context shape
interface SimulatorContextValue {
  state: SimulatorState;
  dispatch: React.Dispatch<SimulatorAction>;
  setHoveredTrait: (trait: string) => void;
  // ...other methods
}

// Create the context with a default value
const SimulatorContext = createContext<SimulatorContextValue | undefined>(
  undefined,
);

// Reducer function
function simulatorReducer(
  state: SimulatorState,
  action: SimulatorAction,
): SimulatorState {
  switch (action.type) {
    case "SET_LOADED_DATA": {
      // Handle loading success
      const { champions, traits, items, augments } = action.payload;
      return {
        ...state,
        champions,
        traits: traits.map((trait) => ({ ...trait, active: 0 })), // Initialize active count
        items,
        augments,
        loading: false,
        error: undefined,
        gold: 0, // Start with 0 gold (no champions)
        level: 2, // Start at base level 2 (no champions)
      };
    }
    case "SET_LOADING_ERROR": {
      // Handle loading error
      return {
        ...state,
        loading: false,
        error: action.error,
      };
    }
    case "ADD_CHAMPION_TO_BOARD": {
      const { champion, position } = action;

      // Check if position is already occupied
      const isOccupied = state.boardChampions.some(
        (c) =>
          c.position.row === position.row && c.position.col === position.col,
      );

      if (isOccupied) {
        return state;
      }

      const boardChampion: BoardChampion = {
        ...champion,
        position,
        stars: 1,
        items: [],
      };

      const updatedBoardChampions = [...state.boardChampions, boardChampion];

      // Update traits based on added champion
      const updatedTraits = updateTraits(updatedBoardChampions, state.traits);

      return {
        ...state,
        boardChampions: updatedBoardChampions,
        traits: updatedTraits,
        gold: calculateTotalGold(updatedBoardChampions),
        level: calculateLevel(updatedBoardChampions),
      };
    }

    case "REMOVE_CHAMPION_FROM_BOARD": {
      const { position } = action;
      const updatedBoardChampions = state.boardChampions.filter(
        (c) =>
          !(c.position.row === position.row && c.position.col === position.col),
      );

      // Update traits after removing the champion
      const updatedTraits = updateTraits(updatedBoardChampions, state.traits);

      return {
        ...state,
        boardChampions: updatedBoardChampions,
        traits: updatedTraits,
        gold: calculateTotalGold(updatedBoardChampions),
        level: calculateLevel(updatedBoardChampions),
      };
    }

    case "MOVE_CHAMPION": {
      const { from, to } = action;

      // Find the champion at the 'from' position
      const championToMove = state.boardChampions.find(
        (c) => c.position.row === from.row && c.position.col === from.col,
      );

      // Check if destination is already occupied
      const isDestinationOccupied = state.boardChampions.some(
        (c) => c.position.row === to.row && c.position.col === to.col,
      );

      if (!championToMove) {
        return state;
      }

      let updatedBoardChampions = [...state.boardChampions];

      // Remove the champion from old position
      updatedBoardChampions = updatedBoardChampions.filter(
        (c) => !(c.position.row === from.row && c.position.col === from.col),
      );

      // If destination is occupied, swap champions
      if (isDestinationOccupied) {
        const championAtDest = state.boardChampions.find(
          (c) => c.position.row === to.row && c.position.col === to.col,
        )!;

        // Create new champion object with updated position (from -> to)
        const updatedChampionToMove = {
          ...championToMove,
          position: to,
        };

        // Create new champion object with updated position (to -> from)
        const updatedChampionAtDest = {
          ...championAtDest,
          position: from,
        };

        // Add both champions with swapped positions
        updatedBoardChampions = updatedBoardChampions.filter(
          (c) => !(c.position.row === to.row && c.position.col === to.col),
        );

        updatedBoardChampions.push(
          updatedChampionToMove,
          updatedChampionAtDest,
        );
      } else {
        // Just move the champion to the new position
        updatedBoardChampions.push({
          ...championToMove,
          position: to,
        });
      }

      return {
        ...state,
        boardChampions: updatedBoardChampions,
        // Moving champions doesn't change gold or level
      };
    }

    case "ADD_ITEM_TO_CHAMPION": {
      const { item, position } = action;

      const updatedBoardChampions = state.boardChampions.map((c) => {
        if (
          c.position.row === position.row &&
          c.position.col === position.col
        ) {
          // Check if champion already has 3 items
          if (c.items && c.items.length >= 3) {
            return c;
          }

          return {
            ...c,
            items: c.items ? [...c.items, item] : [item],
          };
        }
        return c;
      });

      return {
        ...state,
        boardChampions: updatedBoardChampions,
        selectedItem: undefined,
      };
    }

    case "REMOVE_ITEM_FROM_CHAMPION": {
      const { itemApiName: itemId, position } = action;

      const updatedBoardChampions = state.boardChampions.map((c) => {
        if (
          c.position.row === position.row &&
          c.position.col === position.col &&
          c.items
        ) {
          return {
            ...c,
            items: c.items.filter((item) => item.apiName !== itemId),
          };
        }
        return c;
      });

      return {
        ...state,
        boardChampions: updatedBoardChampions,
      };
    }

    case "SELECT_ITEM": {
      return {
        ...state,
        selectedItem: action.item,
        selectedChampion: undefined,
      };
    }

    case "SELECT_CHAMPION": {
      return {
        ...state,
        selectedChampion: action.champion,
        selectedItem: undefined,
      };
    }

    case "STAR_UP_CHAMPION": {
      const { position } = action;

      const updatedBoardChampions = state.boardChampions.map((c) => {
        if (
          c.position.row === position.row &&
          c.position.col === position.col
        ) {
          // Only allow star up to 3 stars
          const currentStars = c.stars || 1;
          if (currentStars >= 3) {
            return c;
          }

          return {
            ...c,
            stars: (currentStars + 1) as 1 | 2 | 3,
          };
        }
        return c;
      });

      return {
        ...state,
        boardChampions: updatedBoardChampions,
        gold: calculateTotalGold(updatedBoardChampions),
        // Level doesn't change on star up, but included for consistency
        level: calculateLevel(updatedBoardChampions),
      };
    }

    case "SELECT_AUGMENT": {
      const { augment, index } = action;
      const updatedSelectedAugments = [...state.selectedAugments];

      // Replace augment at index or add it if it doesn't exist
      if (index >= updatedSelectedAugments.length) {
        updatedSelectedAugments.push(augment);
      } else {
        updatedSelectedAugments[index] = augment;
      }

      return {
        ...state,
        selectedAugments: updatedSelectedAugments,
      };
    }

    case "CLEAR_BOARD": {
      return {
        ...state,
        boardChampions: [],
        traits: state.traits.map((trait) => ({
          // Reset active counts on clear
          ...trait,
          active: 0,
        })),
        gold: 0, // Reset gold to 0 when board is cleared
        level: 2, // Reset level to base level 2 when board is cleared
      };
    }
    case "SET_CHAMPION_STAR_LEVEL": {
      // First, create the updated boardChampions array
      const updatedBoardChampions = state.boardChampions.map(
        (boardChampion) => {
          if (
            boardChampion &&
            boardChampion.position.row === action.position.row &&
            boardChampion.position.col === action.position.col
          ) {
            return {
              ...boardChampion,
              stars: action.level as 1 | 2 | 3, // Cast to the correct type (1, 2, or 3)
            };
          }
          return boardChampion;
        },
      );

      // Then use the updated array for gold and level calculations
      return {
        ...state,
        boardChampions: updatedBoardChampions,
        gold: calculateTotalGold(updatedBoardChampions),
        level: calculateLevel(updatedBoardChampions),
      };
    }

    default:
      return state;
  }
}

// Helper function to update traits based on board champions
function updateTraits(
  boardChampions: BoardChampion[],
  traits: Trait[],
): Trait[] {
  // Count unique champions per trait
  const traitCounts: { [key: string]: number } = {};
  const traitContributors: { [key: string]: Set<string> } = {}; // Track unique champions by apiName

  boardChampions.forEach((champion) => {
    champion.traits.forEach((traitName) => {
      // Initialize set if not exists
      if (!traitContributors[traitName]) {
        traitContributors[traitName] = new Set();
      }

      // Only count unique champions for each trait
      if (!traitContributors[traitName].has(champion.apiName)) {
        traitContributors[traitName].add(champion.apiName);
        traitCounts[traitName] = (traitCounts[traitName] || 0) + 1;
      }
    });
  });

  // Update trait active values only - we'll use traitTierColors in the TraitTracker
  return traits.map((trait) => {
    const active = traitCounts[trait.name] || 0;

    return {
      ...trait,
      active,
    };
  });
}

// Helper function to calculate total gold based on champions' cost and star level
function calculateTotalGold(boardChampions: BoardChampion[]): number {
  return boardChampions.reduce((total, champion) => {
    // Cost multiplier based on star level: 1★ = 1x, 2★ = 3x, 3★ = 9x
    const starMultiplier =
      champion.stars === 1 ? 1 : champion.stars === 2 ? 3 : 9;
    return total + champion.cost * starMultiplier;
  }, 0);
}

// Helper function to calculate level based on number of champions
function calculateLevel(boardChampions: BoardChampion[]): number {
  // Base level is 2, each champion adds +1, max level is 9
  return Math.min(9, 2 + boardChampions.length);
}

// Provider component
export const SimulatorProvider: React.FC<{ children: React.ReactNode }> = ({
  children,
}) => {
  const [state, dispatch] = useReducer(simulatorReducer, initialState);
  const [hoveredTrait, setHoveredTraitValue] = useState<string>("");

  // Combine the reducer state with the hoveredTrait state
  const combinedState = {
    ...state,
    hoveredTrait,
  };

  // Create the context value
  const value: SimulatorContextValue = {
    state: combinedState,
    dispatch,
    setHoveredTrait: (trait: string) => setHoveredTraitValue(trait),
    // ...other methods
  };

  useEffect(() => {
    const fetchData = async () => {
      try {
        const dataFilePath = "./en_us_14.1b.json"; // Assuming this file contains all needed data
        const targetMutator = "TFTSet14"; // Example Mutator - Adjust if needed
        const targetPrefix = "TFT14_"; // Example prefix for champions

        // Load all data concurrently
        const [setDataResult, itemsAndAugmenets] = await Promise.all([
          loadSetDataFromFile(dataFilePath, targetMutator),
          loadItemsAndAugmentsFromFile(dataFilePath), // Adjust path as needed
        ]);

        const allChampions =
          (setDataResult.setData[0]?.champions as Champion[]) || [];
        const champions = allChampions
          .filter(
            (champion) =>
              champion.apiName &&
              champion.apiName.startsWith(targetPrefix) &&
              !champion.apiName.includes("NPC"),
          )
          .map((champion) => {
            if (champion.squareIcon && champion.icon) {
              const squareIconFileName = champion.squareIcon
                .split("/")
                .pop()
                ?.replace(".tex", ".png");
              const iconFileName = champion.icon
                .split("/")
                .pop()
                ?.replace(".tex", ".png");
              return {
                ...champion,
                squareIcon: squareIconFileName,
                icon: iconFileName,
              };
            }
            return champion;
          }) as Champion[];

        const traits =
          setDataResult.setData[0]?.traits.map((trait) => {
            if (trait.icon) {
              const iconFileName = trait.icon
                .split("/")
                .pop()
                ?.replace(".tex", ".png");
              return {
                ...trait,
                icon: iconFileName,
              };
            }
            return trait;
          }) || [];

        const setActiveItems = setDataResult.setData[0]?.items || [];
        const setActiveAugments = setDataResult.setData[0]?.augments || [];

        const items = itemsAndAugmenets
          .filter(
            (item) =>
              setActiveItems.includes(item.apiName) &&
              item.apiName.startsWith("TFT_Item_") &&
              !item?.apiName?.startsWith("TFT_Item_Grant") &&
              !item?.apiName?.includes("Anvil"),
          )
          .map((item) => {
            if (item.icon) {
              const iconFileName = item.icon
                .split("/")
                .pop()
                ?.replace(/(?:\.TFT_Set\d+)?\.tex$/, ".png");

              return {
                ...item,
                icon: iconFileName,
              };
            }
            return item;
          }) as Item[];

        const augments = itemsAndAugmenets
          .filter((item) => setActiveAugments.includes(item.apiName))
          .map((augment) => {
            if (augment.icon) {
              const iconFileName = augment.icon
                .split("/")
                .pop()
                ?.replace(".tex", ".png");
              return {
                ...augment,
                icon: iconFileName,
              };
            }
            return augment;
          }) as Item[];

        console.log("Loaded data:", { champions, traits, items, augments });

        dispatch({
          type: "SET_LOADED_DATA",
          payload: { champions, traits, items, augments },
        });
      } catch (err) {
        console.error("Failed to load simulator data:", err);
        dispatch({
          type: "SET_LOADING_ERROR",
          error: err instanceof Error ? err.message : String(err),
        });
      }
    };

    fetchData();
  }, []);

  return (
    <SimulatorContext.Provider value={value}>
      {children}
    </SimulatorContext.Provider>
  );
};

// Custom hook for using the simulator context
export function useSimulator() {
  const context = useContext(SimulatorContext);

  if (context === undefined) {
    throw new Error("useSimulator must be used within a SimulatorProvider");
  }

  return context;
}
