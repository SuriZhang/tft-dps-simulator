
import { Champion, Item, Trait, Augment } from './types';

// Mock champions data
export const MOCK_CHAMPIONS: Champion[] = [
  {
    id: 'zac',
    name: 'Zac',
    cost: 4,
    traits: ['Bruiser', 'Ooze'],
    image: '/placeholder.svg',
    stars: 1
  },
  {
    id: 'poppy',
    name: 'Poppy',
    cost: 1,
    traits: ['Yordle', 'Knight'],
    image: '/placeholder.svg',
    stars: 1
  },
  {
    id: 'missfortune',
    name: 'Miss Fortune',
    cost: 5,
    traits: ['Gunslinger', 'Ace'],
    image: '/placeholder.svg',
    stars: 2,
    items: [
      { id: 'ie', name: 'Infinity Edge', description: 'Critical strikes deal more damage', image: '/placeholder.svg', type: 'completed' },
      { id: 'gs', name: 'Giant Slayer', description: 'Deal bonus damage to high health targets', image: '/placeholder.svg', type: 'completed' }
    ]
  },
  {
    id: 'akali',
    name: 'Akali',
    cost: 4,
    traits: ['Assassin', 'Ninja'],
    image: '/placeholder.svg',
    stars: 1
  },
  {
    id: 'ahri',
    name: 'Ahri',
    cost: 5,
    traits: ['Spirit', 'Mage'],
    image: '/placeholder.svg',
    stars: 1
  },
  {
    id: 'lux',
    name: 'Lux',
    cost: 3,
    traits: ['Sorcerer', 'Light'],
    image: '/placeholder.svg',
    stars: 1
  },
  {
    id: 'yasuo',
    name: 'Yasuo',
    cost: 4,
    traits: ['Exile', 'Blademaster'],
    image: '/placeholder.svg',
    stars: 1
  },
  {
    id: 'jinx',
    name: 'Jinx',
    cost: 4,
    traits: ['Scrap', 'Gunslinger'],
    image: '/placeholder.svg',
    stars: 1
  }
];

// Mock items data
export const MOCK_ITEMS: Item[] = [
  {
    id: 'bf',
    name: 'B.F. Sword',
    description: '+10 Attack Damage',
    image: '/placeholder.svg',
    type: 'component'
  },
  {
    id: 'bow',
    name: 'Recurve Bow',
    description: '+10% Attack Speed',
    image: '/placeholder.svg',
    type: 'component'
  },
  {
    id: 'rod',
    name: 'Needlessly Large Rod',
    description: '+10 Ability Power',
    image: '/placeholder.svg',
    type: 'component'
  },
  {
    id: 'ie',
    name: 'Infinity Edge',
    description: 'Critical strikes deal more damage',
    image: '/placeholder.svg',
    type: 'completed'
  },
  {
    id: 'gs',
    name: 'Giant Slayer',
    description: 'Deal bonus damage to high health targets',
    image: '/placeholder.svg',
    type: 'completed'
  },
  {
    id: 'hoj',
    name: 'Hand of Justice',
    description: 'Grants random bonuses each round',
    image: '/placeholder.svg',
    type: 'completed'
  }
];

// Mock traits data
export const MOCK_TRAITS: Trait[] = [
  {
    id: 'virus',
    name: 'Virus',
    description: 'Infect enemies with a virus that deals damage over time',
    image: '/placeholder.svg',
    bonuses: [
      { count: 1, effect: 'Viral infections deal 20% damage' },
      { count: 3, effect: 'Viral infections deal 40% damage' },
      { count: 5, effect: 'Viral infections deal 80% damage' }
    ],
    active: 1,
    style: 'border-orange-500'
  },
  {
    id: 'bastion',
    name: 'Bastion',
    description: 'Gain shields at the start of combat',
    image: '/placeholder.svg',
    bonuses: [
      { count: 2, effect: '+15% shield strength' },
      { count: 4, effect: '+40% shield strength' },
      { count: 6, effect: '+70% shield strength' }
    ],
    active: 2,
    style: 'border-blue-500'
  },
  {
    id: 'cyberboss',
    name: 'Cyberboss',
    description: 'Deal bonus damage based on opponent health',
    image: '/placeholder.svg',
    bonuses: [
      { count: 2, effect: '+10% damage vs high health enemies' },
      { count: 4, effect: '+25% damage vs high health enemies' },
      { count: 6, effect: '+50% damage vs high health enemies' }
    ],
    active: 2,
    style: 'border-cyan-500'
  },
  {
    id: 'dynamo',
    name: 'Dynamo',
    description: 'Generate energy over time, powering up abilities',
    image: '/placeholder.svg',
    bonuses: [
      { count: 2, effect: '+10 energy per second' },
      { count: 4, effect: '+25 energy per second' },
      { count: 6, effect: '+50 energy per second' }
    ],
    active: 2,
    style: 'border-yellow-500'
  },
  {
    id: 'syndicate',
    name: 'Syndicate',
    description: 'Gain lifesteal and spell vamp',
    image: '/placeholder.svg',
    bonuses: [
      { count: 3, effect: '+10% lifesteal and spell vamp' },
      { count: 5, effect: '+25% lifesteal and spell vamp' },
      { count: 7, effect: '+45% lifesteal and spell vamp' }
    ],
    active: 3,
    style: 'border-purple-500'
  }
];

// Mock augments data
export const MOCK_AUGMENTS: Augment[] = [
  {
    id: 'rich-get-richer',
    name: 'Rich Get Richer',
    description: 'Gain +10 gold now, but interest rate reduces to 0%',
    image: '/placeholder.svg',
    tier: 'silver'
  },
  {
    id: 'cybernetic-implants',
    name: 'Cybernetic Implants',
    description: 'Champions with at least 1 item gain +150 Health and +15 Attack Damage',
    image: '/placeholder.svg',
    tier: 'gold'
  },
  {
    id: 'celestial-blessing',
    name: 'Celestial Blessing',
    description: 'All units gain 15% Omnivamp (healing for damage dealt)',
    image: '/placeholder.svg',
    tier: 'prismatic'
  },
  {
    id: 'built-different',
    name: 'Built Different',
    description: 'Champions without any traits active gain bonus health and attack speed',
    image: '/placeholder.svg',
    tier: 'gold'
  }
];

// Board size constants
export const BOARD_ROWS = 4;
export const BOARD_COLS = 7;
