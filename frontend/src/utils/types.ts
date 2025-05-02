
export interface Champion {
  apiName: string;
  name: string;
  cost: number;
  traits: string[];
  image: string;
  stars?: 1 | 2 | 3;
  items?: Item[];
}

export interface Item {
  apiName: string;
  name: string;
  description: string;
  image: string;
  type: 'component' | 'completed' | 'special';
}

export interface Trait {
  apiName: string;
  name: string;
  description: string;
  image: string;
  bonuses: TraitBonus[];
  active: number;
  style?: string;
}

export interface TraitBonus {
  count: number;
  effect: string;
}

export interface Augment {
  apiName: string;
  name: string;
  description: string;
  image: string;
  tier: 'silver' | 'gold' | 'prismatic';
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
  augments: Augment[];
  selectedAugments: Augment[];
  gold: number;
  level: number;
  selectedItem?: Item;
  selectedChampion?: Champion;
}

export type SimulatorAction = 
  | { type: 'ADD_CHAMPION_TO_BOARD'; champion: Champion; position: BoardPosition }
  | { type: 'REMOVE_CHAMPION_FROM_BOARD'; position: BoardPosition }
  | { type: 'MOVE_CHAMPION'; from: BoardPosition; to: BoardPosition }
  | { type: 'ADD_ITEM_TO_CHAMPION'; item: Item; position: BoardPosition }
  | { type: 'REMOVE_ITEM_FROM_CHAMPION'; itemApiName: string; position: BoardPosition }
  | { type: 'SELECT_ITEM'; item: Item | undefined }
  | { type: 'SELECT_CHAMPION'; champion: Champion | undefined }
  | { type: 'STAR_UP_CHAMPION'; position: BoardPosition }
  | { type: 'SELECT_AUGMENT'; augment: Augment; index: number }
  | { type: 'CLEAR_BOARD' };
