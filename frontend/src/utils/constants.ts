import { Item, TFTSetData, Set } from "./types"; // Import Set type

// Remove the local Set interface definition

interface RawGameData {
  items: Item[]; // The raw array can contain both
  setData: Set[];
  // Add other top-level fields from the JSON if needed
}

/**
 * Loads and filters set data (champions, traits) from a JSON file based on a target mutator.
 * @param filePath - The path to the JSON file.
 * @param targetMutator - The mutator string to filter by.
 * @returns A promise that resolves with the filtered set data or rejects with an error.
 */
export async function loadSetDataFromFile(
  filePath: string,
  targetMutator: string,
): Promise<TFTSetData> {
  try {
    const response = await fetch(filePath);
    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }
    // Parse the entire structure first
    const fullData: RawGameData = await response.json();

    if (!fullData || !Array.isArray(fullData.setData)) {
      throw new Error(
        "Invalid set data format: 'setData' array not found or not an array.",
      );
    }

    const filteredSetData = fullData.setData.filter(
      (set) => set.mutator === targetMutator,
    );

    if (filteredSetData.length === 0) {
      const availableMutators = fullData.setData.map((set) => set.mutator);
      const mutatorsStr =
        availableMutators.length > 0
          ? ` Available mutators: ${availableMutators.join(", ")}`
          : "";
      throw new Error(
        `No data with mutator '${targetMutator}' found.${mutatorsStr}`,
      );
    }

    console.log(filteredSetData);

    // Return the data structure containing only the filtered sets
    // Ensure champions and traits are included
    return { setData: filteredSetData };
  } catch (error) {
    console.error("Error loading or parsing set data:", error);
    throw new Error(
      `Error loading set data from ${filePath}: ${error instanceof Error ? error.message : String(error)}`,
    );
  }
}

/**
 * Loads and separates item and augment data from a JSON file.
 * Assumes the JSON structure has a top-level "items" array containing both.
 * Differentiates based on apiName prefix.
 * @param filePath - The path to the JSON file.
 * @returns A promise that resolves with an object containing separate arrays for items and augments, or rejects with an error.
 */
export async function loadItemsAndAugmentsFromFile(
  filePath: string,
): Promise<Item[]> {
  try {
    const response = await fetch(filePath);
    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }
    // Assuming the JSON structure is { "items": [...] }
    const data: { items: Item[] } = await response.json();

    if (!data || !Array.isArray(data.items)) {
      throw new Error(
        "Invalid item/augment data format: 'items' array not found or not an array.",
      );
    }

    const allEntries = data.items as Item[];
    return allEntries;
  } catch (error) {
    console.error("Error loading or parsing item/augment data:", error);
    throw new Error(
      `Error loading item/augment data from ${filePath}: ${error instanceof Error ? error.message : String(error)}`,
    );
  }
}

// Remove the old loadAugmentDataFromFile function
/*
export async function loadAugmentDataFromFile(filePath: string): Promise<Augment[]> {
    // ... (old code removed)
}
*/

// // Mock champions data
// export const MOCK_CHAMPIONS: Champion[] = [
//   {
//     apiName: 'zac',
//     name: 'Zac',
//     cost: 4,
//     traits: ['Bruiser', 'Ooze'],
//     icon: '/placeholder.svg',
//     stars: 1
//   },
//   {
//     apiName: 'poppy',
//     name: 'Poppy',
//     cost: 1,
//     traits: ['Yordle', 'Knight'],
//     icon: '/placeholder.svg',
//     stars: 1
//   },
//   {
//     apiName: 'missfortune',
//     name: 'Miss Fortune',
//     cost: 5,
//     traits: ['Gunslinger', 'Ace'],
//     icon: '/placeholder.svg',
//     stars: 2,
//     items: [
//       { apiName: 'ie', name: 'Infinity Edge', desc: 'Critical strikes deal more damage', icon: '/placeholder.svg', type: 'completed' },
//       { apiName: 'gs', name: 'Giant Slayer', desc: 'Deal bonus damage to high health targets', icon: '/placeholder.svg', type: 'completed' }
//     ]
//   },
//   {
//     apiName: 'akali',
//     name: 'Akali',
//     cost: 4,
//     traits: ['Assassin', 'Ninja'],
//     icon: '/placeholder.svg',
//     stars: 1
//   },
//   {
//     apiName: 'ahri',
//     name: 'Ahri',
//     cost: 5,
//     traits: ['Spirit', 'Mage'],
//     icon: '/placeholder.svg',
//     stars: 1
//   },
//   {
//     apiName: 'lux',
//     name: 'Lux',
//     cost: 3,
//     traits: ['Sorcerer', 'Light'],
//     icon: '/placeholder.svg',
//     stars: 1
//   },
//   {
//     apiName: 'yasuo',
//     name: 'Yasuo',
//     cost: 4,
//     traits: ['Exile', 'Blademaster'],
//     icon: '/placeholder.svg',
//     stars: 1
//   },
//   {
//     apiName: 'jinx',
//     name: 'Jinx',
//     cost: 4,
//     traits: ['Scrap', 'Gunslinger'],
//     icon: '/placeholder.svg',
//     stars: 1
//   }
// ];

// // Mock items data
// export const MOCK_ITEMS: Item[] = [
//   {
//     apiName: 'bf',
//     name: 'B.F. Sword',
//     desc: '+10 Attack Damage',
//     icon: '/placeholder.svg',
//     type: 'component'
//   },
//   {
//     apiName: 'bow',
//     name: 'Recurve Bow',
//     desc: '+10% Attack Speed',
//     icon: '/placeholder.svg',
//     type: 'component'
//   },
//   {
//     apiName: 'rod',
//     name: 'Needlessly Large Rod',
//     desc: '+10 Ability Power',
//     icon: '/placeholder.svg',
//     type: 'component'
//   },
//   {
//     apiName: 'ie',
//     name: 'Infinity Edge',
//     desc: 'Critical strikes deal more damage',
//     icon: '/placeholder.svg',
//     type: 'completed'
//   },
//   {
//     apiName: 'gs',
//     name: 'Giant Slayer',
//     desc: 'Deal bonus damage to high health targets',
//     icon: '/placeholder.svg',
//     type: 'completed'
//   },
//   {
//     apiName: 'hoj',
//     name: 'Hand of Justice',
//     desc: 'Grants random bonuses each round',
//     icon: '/placeholder.svg',
//     type: 'completed'
//   }
// ];

// // Mock traits data
// export const MOCK_TRAITS: Trait[] = [
//   {
//     apiName: 'virus',
//     name: 'Virus',
//     desc: 'Infect enemies with a virus that deals damage over time',
//     icon: '/placeholder.svg',
//     bonuses: [
//       { count: 1, effect: 'Viral infections deal 20% damage' },
//       { count: 3, effect: 'Viral infections deal 40% damage' },
//       { count: 5, effect: 'Viral infections deal 80% damage' }
//     ],
//     active: 1,
//     style: 'border-orange-500'
//   },
//   {
//     apiName: 'bastion',
//     name: 'Bastion',
//     desc: 'Gain shields at the start of combat',
//     icon: '/placeholder.svg',
//     bonuses: [
//       { count: 2, effect: '+15% shield strength' },
//       { count: 4, effect: '+40% shield strength' },
//       { count: 6, effect: '+70% shield strength' }
//     ],
//     active: 2,
//     style: 'border-blue-500'
//   },
//   {
//     apiName: 'cyberboss',
//     name: 'Cyberboss',
//     desc: 'Deal bonus damage based on opponent health',
//     icon: '/placeholder.svg',
//     bonuses: [
//       { count: 2, effect: '+10% damage vs high health enemies' },
//       { count: 4, effect: '+25% damage vs high health enemies' },
//       { count: 6, effect: '+50% damage vs high health enemies' }
//     ],
//     active: 2,
//     style: 'border-cyan-500'
//   },
//   {
//     apiName: 'dynamo',
//     name: 'Dynamo',
//     desc: 'Generate energy over time, powering up abilities',
//     icon: '/placeholder.svg',
//     bonuses: [
//       { count: 2, effect: '+10 energy per second' },
//       { count: 4, effect: '+25 energy per second' },
//       { count: 6, effect: '+50 energy per second' }
//     ],
//     active: 2,
//     style: 'border-yellow-500'
//   },
//   {
//     apiName: 'syndicate',
//     name: 'Syndicate',
//     desc: 'Gain lifesteal and spell vamp',
//     icon: '/placeholder.svg',
//     bonuses: [
//       { count: 3, effect: '+10% lifesteal and spell vamp' },
//       { count: 5, effect: '+25% lifesteal and spell vamp' },
//       { count: 7, effect: '+45% lifesteal and spell vamp' }
//     ],
//     active: 3,
//     style: 'border-purple-500'
//   }
// ];

// // Mock augments data
// export const MOCK_AUGMENTS: Augment[] = [
//   {
//     apiName: 'rich-get-richer',
//     name: 'Rich Get Richer',
//     desc: 'Gain +10 gold now, but interest rate reduces to 0%',
//     icon: '/placeholder.svg',
//     tier: 'silver'
//   },
//   {
//     apiName: 'cybernetic-implants',
//     name: 'Cybernetic Implants',
//     desc: 'Champions with at least 1 item gain +150 Health and +15 Attack Damage',
//     icon: '/placeholder.svg',
//     tier: 'gold'
//   },
//   {
//     apiName: 'celestial-blessing',
//     name: 'Celestial Blessing',
//     desc: 'All units gain 15% Omnivamp (healing for damage dealt)',
//     icon: '/placeholder.svg',
//     tier: 'prismatic'
//   },
//   {
//     apiName: 'built-different',
//     name: 'Built Different',
//     desc: 'Champions without any traits active gain bonus health and attack speed',
//     icon: '/placeholder.svg',
//     tier: 'gold'
//   }
// ];

// Board size constants
export const BOARD_ROWS = 4;
export const BOARD_COLS = 7;
