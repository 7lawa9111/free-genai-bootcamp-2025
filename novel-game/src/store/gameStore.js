import create from 'zustand';

const useGameStore = create((set) => ({
  currentScene: null,
  gameProgress: {
    visitedScenes: [],
    choices: [],
    vocabulary: new Set(),
  },
  setCurrentScene: (scene) => set({ currentScene: scene }),
  addVisitedScene: (sceneId) =>
    set((state) => ({
      gameProgress: {
        ...state.gameProgress,
        visitedScenes: [...state.gameProgress.visitedScenes, sceneId],
      },
    })),
  addChoice: (choice) =>
    set((state) => ({
      gameProgress: {
        ...state.gameProgress,
        choices: [...state.gameProgress.choices, choice],
      },
    })),
  addVocabulary: (word) =>
    set((state) => ({
      gameProgress: {
        ...state.gameProgress,
        vocabulary: new Set([...state.gameProgress.vocabulary, word]),
      },
    })),
  resetGame: () =>
    set({
      currentScene: null,
      gameProgress: {
        visitedScenes: [],
        choices: [],
        vocabulary: new Set(),
      },
    }),
}));

export default useGameStore; 