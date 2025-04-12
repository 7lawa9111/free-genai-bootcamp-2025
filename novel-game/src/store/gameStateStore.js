import { create } from 'zustand';
import { persist } from 'zustand/middleware';

// Initial game state
const initialState = {
  location: 'Post Office',
  inventory: [],
  flags: {},
  items: [
    {
      id: 'stamp',
      name: 'Stamp',
      nameJP: '切手',
      description: 'A standard postal stamp.',
      descriptionJP: '普通の郵便切手です。',
      location: 'Post Office',
      isFixed: false,
      isOpenable: false
    },
    {
      id: 'package',
      name: 'Package',
      nameJP: '小包',
      description: 'A small package ready to be sent.',
      descriptionJP: '送る準備ができた小包です。',
      location: 'Post Office',
      isFixed: false,
      isOpenable: true,
      isLocked: false
    }
  ],
  npcs: [
    {
      id: 'clerk',
      name: 'Post Office Clerk',
      nameJP: '郵便局員',
      location: 'Post Office',
      dialogue: 'How may I help you today?',
      dialogueJP: '本日はどのようなご用件でしょうか？'
    }
  ],
  currentScene: {
    id: 'post-office',
    description: 'You are in a post office. The clerk is waiting to help you.',
    descriptionJP: '郵便局にいます。係員があなたを待っています。',
    location: 'Post Office'
  },
  commandHistory: []
};

const useGameStateStore = create(
  persist(
    (set, get) => ({
      // Initialize with initial state
      ...initialState,

      // State management functions
      setLocation: (newLocation) => set({ location: newLocation }),
      
      addToInventory: (itemId) => set((state) => ({
        inventory: [...state.inventory, itemId]
      })),
      
      removeFromInventory: (itemId) => set((state) => ({
        inventory: state.inventory.filter(id => id !== itemId)
      })),
      
      setFlag: (flag, value) => set((state) => ({
        flags: { ...state.flags, [flag]: value }
      })),
      
      setItemState: (itemId, property, value) => set((state) => ({
        items: state.items.map(item => 
          item.id === itemId ? { ...item, [property]: value } : item
        )
      })),

      // Helper functions for getting game state info
      getItemsInLocation: (location) => {
        const state = get();
        return state.items.filter(item => item.location === location);
      },

      getNPCsInLocation: (location) => {
        const state = get();
        return state.npcs.filter(npc => npc.location === location);
      },

      // Scene management
      setCurrentScene: (scene) => set({ currentScene: scene }),
      
      // Command history management
      addCommandToHistory: (command, response) => set((state) => ({
        commandHistory: [...state.commandHistory, { command, response, timestamp: Date.now() }]
      })),
      
      clearCommandHistory: () => set({ commandHistory: [] }),
      
      // Reset game state
      resetGame: () => set(initialState),
      
      // Helper functions
      hasItem: (itemId) => {
        const state = get();
        return state.inventory.includes(itemId);
      },
      
      hasFlag: (flag) => {
        const state = get();
        return state.flags[flag] === true;
      },
      
      isItemInLocation: (itemId) => {
        const state = get();
        return state.items.some(item => item.id === itemId && item.location === state.location);
      },
      
      getOpenableObjects: (scene) => {
        const state = get();
        return state.items.filter(item => item.isOpenable && item.location === state.location);
      },
      
      getCloseableObjects: (scene) => {
        const state = get();
        return state.items.filter(item => 
          item.isOpenable && 
          item.location === state.location && 
          state.flags[`${item.id}_open`]
        );
      }
    }),
    {
      name: 'game-state-storage',
      partialize: (state) => ({
        inventory: state.inventory,
        location: state.location,
        flags: state.flags,
        items: state.items,
        npcs: state.npcs,
        currentScene: state.currentScene
      })
    }
  )
);

export { useGameStateStore }; 