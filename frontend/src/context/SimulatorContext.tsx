import React, { createContext, useReducer, useContext, ReactNode, useEffect } from 'react';
import { 
  Champion, 
  BoardPosition, 
  Item, 
  SimulatorState, 
  SimulatorAction, 
  BoardChampion, 
  Trait 
} from '../utils/types';
import { loadSetDataFromFile, loadItemsAndAugmentsFromFile } from '../utils/constants';

const initialState: SimulatorState = {
  champions: [], // Start with empty data
  boardChampions: [],
  items: [],
  traits: [],
  augments: [],
  selectedAugments: [],
  gold: 50,
  level: 7,
  selectedItem: undefined,
  selectedChampion: undefined,
  loading: true, // Set initial loading state to true
  error: undefined
};

// Create context
const SimulatorContext = createContext<{
  state: SimulatorState;
  dispatch: React.Dispatch<SimulatorAction>;
} | undefined>(undefined);

// Reducer function
function simulatorReducer(state: SimulatorState, action: SimulatorAction): SimulatorState {
  switch (action.type) {
    case 'SET_LOADED_DATA': { // Handle loading success
      const { champions, traits, items, augments } = action.payload;
      return {
        ...state,
        champions,
        traits: traits.map(trait => ({ ...trait, active: 0 })), // Initialize active count
        items,
        augments,
        loading: false,
        error: undefined
      };
    }
    case 'SET_LOADING_ERROR': { // Handle loading error
        return {
            ...state,
            loading: false,
            error: action.error
        };
    }
    case 'ADD_CHAMPION_TO_BOARD': {
      const { champion, position } = action;
      
      // Check if position is already occupied
      const isOccupied = state.boardChampions.some(
        c => c.position.row === position.row && c.position.col === position.col
      );
      
      if (isOccupied) {
        return state;
      }
      
      const boardChampion: BoardChampion = {
        ...champion,
        position,
        stars: 1,
        items: []
      };
      
      // Update traits based on added champion
      const updatedTraits = updateTraits([...state.boardChampions, boardChampion], state.traits);
      
      return {
        ...state,
        boardChampions: [...state.boardChampions, boardChampion],
        traits: updatedTraits
      };
    }
    
    case 'REMOVE_CHAMPION_FROM_BOARD': {
      const { position } = action;
      const updatedBoardChampions = state.boardChampions.filter(
        c => !(c.position.row === position.row && c.position.col === position.col)
      );
      
      // Update traits after removing the champion
      const updatedTraits = updateTraits(updatedBoardChampions, state.traits);
      
      return {
        ...state,
        boardChampions: updatedBoardChampions,
        traits: updatedTraits
      };
    }
    
    case 'MOVE_CHAMPION': {
      const { from, to } = action;
      
      // Find the champion at the 'from' position
      const championToMove = state.boardChampions.find(
        c => c.position.row === from.row && c.position.col === from.col
      );
      
      // Check if destination is already occupied
      const isDestinationOccupied = state.boardChampions.some(
        c => c.position.row === to.row && c.position.col === to.col
      );
      
      if (!championToMove) {
        return state;
      }
      
      let updatedBoardChampions = [...state.boardChampions];
      
      // Remove the champion from old position
      updatedBoardChampions = updatedBoardChampions.filter(
        c => !(c.position.row === from.row && c.position.col === from.col)
      );
      
      // If destination is occupied, swap champions
      if (isDestinationOccupied) {
        const championAtDest = state.boardChampions.find(
          c => c.position.row === to.row && c.position.col === to.col
        )!;
        
        // Create new champion object with updated position (from -> to)
        const updatedChampionToMove = {
          ...championToMove,
          position: to
        };
        
        // Create new champion object with updated position (to -> from)
        const updatedChampionAtDest = {
          ...championAtDest,
          position: from
        };
        
        // Add both champions with swapped positions
        updatedBoardChampions = updatedBoardChampions.filter(
          c => !(c.position.row === to.row && c.position.col === to.col)
        );
        
        updatedBoardChampions.push(updatedChampionToMove, updatedChampionAtDest);
      } else {
        // Just move the champion to the new position
        updatedBoardChampions.push({
          ...championToMove,
          position: to
        });
      }
      
      return {
        ...state,
        boardChampions: updatedBoardChampions
      };
    }
    
    case 'ADD_ITEM_TO_CHAMPION': {
      const { item, position } = action;
      
      const updatedBoardChampions = state.boardChampions.map(c => {
        if (c.position.row === position.row && c.position.col === position.col) {
          // Check if champion already has 3 items
          if (c.items && c.items.length >= 3) {
            return c;
          }
          
          return {
            ...c,
            items: c.items ? [...c.items, item] : [item]
          };
        }
        return c;
      });
      
      return {
        ...state,
        boardChampions: updatedBoardChampions,
        selectedItem: undefined
      };
    }
    
    case 'REMOVE_ITEM_FROM_CHAMPION': {
      const { itemApiName: itemId, position } = action;
      
      const updatedBoardChampions = state.boardChampions.map(c => {
        if (c.position.row === position.row && c.position.col === position.col && c.items) {
          return {
            ...c,
            items: c.items.filter(item => item.apiName !== itemId)
          };
        }
        return c;
      });
      
      return {
        ...state,
        boardChampions: updatedBoardChampions
      };
    }
    
    case 'SELECT_ITEM': {
      return {
        ...state,
        selectedItem: action.item,
        selectedChampion: undefined
      };
    }
    
    case 'SELECT_CHAMPION': {
      return {
        ...state,
        selectedChampion: action.champion,
        selectedItem: undefined
      };
    }
    
    case 'STAR_UP_CHAMPION': {
      const { position } = action;
      
      const updatedBoardChampions = state.boardChampions.map(c => {
        if (c.position.row === position.row && c.position.col === position.col) {
          // Only allow star up to 3 stars
          const currentStars = c.stars || 1;
          if (currentStars >= 3) {
            return c;
          }
          
          return {
            ...c,
            stars: (currentStars + 1) as 1 | 2 | 3
          };
        }
        return c;
      });
      
      return {
        ...state,
        boardChampions: updatedBoardChampions
      };
    }
    
    case 'SELECT_AUGMENT': {
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
        selectedAugments: updatedSelectedAugments
      };
    }
    
    case 'CLEAR_BOARD': {
      return {
        ...state,
        boardChampions: [],
        traits: state.traits.map(trait => ({ // Reset active counts on clear
          ...trait,
          active: 0
        }))
      };
    }
    
    default:
      return state;
  }
}

// Helper function to update traits based on board champions
function updateTraits(boardChampions: BoardChampion[], traits: Trait[]): Trait[] {
  // Count champions per trait
  const traitCounts: { [key: string]: number } = {};
  
  boardChampions.forEach(champion => {
    champion.traits.forEach(traitName => {
      traitCounts[traitName] = (traitCounts[traitName] || 0) + 1;
    });
  });
  
  // Update trait active values
  return traits.map(trait => ({
    ...trait,
    active: traitCounts[trait.name] || 0
  }));
}

// Provider component
export function SimulatorProvider({ children }: { children: ReactNode }) {
  const [state, dispatch] = useReducer(simulatorReducer, initialState);

  useEffect(() => {
    const fetchData = async () => {
      try {
        const dataFilePath = './en_us_14.1b.json'; // Assuming this file contains all needed data
        const targetMutator = 'TFTSet14'; // Example Mutator - Adjust if needed
        const targetPrefix = 'TFT14_'; // Example prefix for champions
        
        // Load all data concurrently
        const [setDataResult, itemsAndAugmenets] = await Promise.all([
          loadSetDataFromFile(dataFilePath, targetMutator),
          loadItemsAndAugmentsFromFile(dataFilePath), // Adjust path as needed
        ]);

        const allChampions = setDataResult.setData[0]?.champions as Champion[]|| [];
        const champions = allChampions.filter(champion => champion.apiName && champion.apiName.startsWith(targetPrefix));
        const traits = setDataResult.setData[0]?.traits || [];

        const setActiveItems = setDataResult.setData[0]?.items || [];
        const setActiveAugments = setDataResult.setData[0]?.augments || [];

        const items = itemsAndAugmenets.filter(item => setActiveItems.includes(item.apiName)) as Item[];
        const augments = itemsAndAugmenets.filter(item => setActiveAugments.includes(item.apiName)) as Item[];

        console.log("Loaded data:", { champions, traits, items, augments });
      
        dispatch({ 
          type: 'SET_LOADED_DATA', 
          payload: { champions, traits, items, augments } 
        });
      } catch (err) {
        console.error("Failed to load simulator data:", err);
        dispatch({ type: 'SET_LOADING_ERROR', error: err instanceof Error ? err.message : String(err) });
      }
    };

    fetchData();
  }, []);
  
  return (
    <SimulatorContext.Provider value={{ state, dispatch }}>
      {children}
    </SimulatorContext.Provider>
  );
}

// Custom hook for using the simulator context
export function useSimulator() {
  const context = useContext(SimulatorContext);
  
  if (context === undefined) {
    throw new Error("useSimulator must be used within a SimulatorProvider");
  }
  
  return context;
}
