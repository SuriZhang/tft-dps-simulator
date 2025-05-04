// Define interfaces for the structure of the set data JSON
export interface Set {
  mutator: string;
  name: string;
  number: number;
  champions: any[];
  items: string[]; // Replace 'any' with a specific Item type if available
  augments: string[]; // Replace 'any' with a specific Augment type if available
  traits: any[]; // Replace 'any' with a specific Trait type if available
  // Add other fields present in the 'Set' object if needed
}

export interface TFTSetData {
  setData: Set[];
  // Add other top-level fields from the JSON if needed
}

export interface Champion {
  apiName: string;
  name: string;
  cost: number;
  traits: string[];
  icon: string;
  squareIcon: string;
  stars?: 1 | 2 | 3;
  items?: Item[];
}

export interface Item {
  apiName: string;
  associatedTraits?: string[];
  composition?: string[]; // Array of item API names
  desc: string;
  effects?: Record<string, any>; // Using 'any' for flexibility as effect values vary
  from?: string[]; // Array of item API names
  icon: string;
  id?: number; // mostly null
  incompatibleTraits?: string[];
  name: string;
  tags?: string[];
  unique: boolean;
}

export interface Trait {
  apiName: string;
  name: string;
  desc: string;
  icon: string;
  effects: TraitEffect[];
  active: number;
  style?: string;
}

export interface TraitEffect {
  maxUnits: number;
  minUnits: number;
  style: number;
  variables: Record<string, any>;
}

export interface BoardPosition {
  row: number;
  col: number;
}

export interface BoardChampion extends Champion {
  position: BoardPosition;
}

export interface SimulatorState {
  champions: Champion[];
  boardChampions: BoardChampion[];
  items: Item[];
  traits: Trait[];
  augments: Item[];
  selectedAugments: Item[];
  gold: number;
  level: number;
  selectedItem?: Item;
  selectedChampion?: Champion;
  loading: boolean; // Add loading state
  error?: string; // Add error state
}

export type SimulatorAction =
  | {
      type: "ADD_CHAMPION_TO_BOARD";
      champion: Champion;
      position: BoardPosition;
    }
  | { type: "REMOVE_CHAMPION_FROM_BOARD"; position: BoardPosition }
  | { type: "MOVE_CHAMPION"; from: BoardPosition; to: BoardPosition }
  | { type: "ADD_ITEM_TO_CHAMPION"; item: Item; position: BoardPosition }
  | {
      type: "REMOVE_ITEM_FROM_CHAMPION";
      itemApiName: string;
      position: BoardPosition;
    }
  | { type: "SELECT_ITEM"; item: Item | undefined }
  | { type: "SELECT_CHAMPION"; champion: Champion | undefined }
  | { type: "STAR_UP_CHAMPION"; position: BoardPosition }
  | { type: "SELECT_AUGMENT"; augment: Item; index: number }
  | { type: "CLEAR_BOARD" }
  | {
      type: "SET_LOADED_DATA";
      payload: {
        champions: Champion[];
        traits: Trait[];
        items: Item[];
        augments: Item[];
      };
  }
  | { type: "SET_CHAMPION_STAR_LEVEL"; position: BoardPosition; level: number }
  | { type: "SET_LOADING_ERROR"; error: string }; // Add action for loading error
