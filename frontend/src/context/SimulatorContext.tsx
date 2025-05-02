
import React, { createContext, useReducer, useContext, ReactNode } from 'react';
import { 
  Champion, 
  BoardPosition, 
  Item, 
  SimulatorState, 
  SimulatorAction, 
  Augment, 
  BoardChampion, 
  Trait 
} from '../utils/types';
import { MOCK_CHAMPIONS, MOCK_ITEMS, MOCK_TRAITS, MOCK_AUGMENTS } from '../utils/constants';

const initialState: SimulatorState = {
  champions: MOCK_CHAMPIONS,
  boardChampions: [],
  items: MOCK_ITEMS,
  traits: MOCK_TRAITS,
  augments: MOCK_AUGMENTS,
  selectedAugments: [],
  gold: 50,
  level: 7,
  selectedItem: undefined,
  selectedChampion: undefined
};

// Create context
const SimulatorContext = createContext<{
  state: SimulatorState;
  dispatch: React.Dispatch<SimulatorAction>;
} | undefined>(undefined);

// Reducer function
function simulatorReducer(state: SimulatorState, action: SimulatorAction): SimulatorState {
  switch (action.type) {
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
      const { itemId, position } = action;
      
      const updatedBoardChampions = state.boardChampions.map(c => {
        if (c.position.row === position.row && c.position.col === position.col && c.items) {
          return {
            ...c,
            items: c.items.filter(item => item.id !== itemId)
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
        traits: state.traits.map(trait => ({
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
