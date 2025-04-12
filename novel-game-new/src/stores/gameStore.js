import { create } from 'zustand';

const useGameStore = create((set) => ({
  // Text and Scene Management
  currentText: '',
  currentBackground: '/images/backgrounds/classroom.jpg',
  
  // Character States
  characters: {
    sensei: {
      nameRomaji: 'Sensei',
      nameJapanese: '先生',
      image: '/images/characters/sensei.jpg',
      visible: false,
      position: 'center', // 'left', 'center', 'right'
    },
    student: {
      nameRomaji: 'Student',
      nameJapanese: '学生',
      image: '/images/characters/student.jpg',
      visible: false,
      position: 'left',
    },
  },

  // Command History
  commandHistory: [],

  // Game Progress
  completedLessons: [],
  learnedVocabulary: [],
  learnedGrammar: [],

  // Actions
  updateScene: ({ text, background, characterId, position, vocabulary }) => set((state) => {
    const newState = {
      currentText: text || state.currentText,
      currentBackground: background || state.currentBackground,
      characters: { ...state.characters },
    };

    // Reset all characters visibility first
    Object.keys(newState.characters).forEach(id => {
      newState.characters[id] = {
        ...newState.characters[id],
        visible: false,
      };
    });

    // Update specific character if provided
    if (characterId) {
      newState.characters[characterId] = {
        ...newState.characters[characterId],
        visible: true,
        position: position || 'center',
      };
    }

    // Add vocabulary if present
    if (vocabulary) {
      newState.learnedVocabulary = [...state.learnedVocabulary, vocabulary];
    }

    return newState;
  }),

  addCommand: (command, response) => set((state) => ({
    commandHistory: [...state.commandHistory, { command, response }],
  })),

  clearHistory: () => set(() => ({ commandHistory: [] })),

  completeLesson: (lessonId) => set((state) => ({
    completedLessons: [...state.completedLessons, lessonId],
  })),

  addGrammar: (grammar) => set((state) => ({
    learnedGrammar: [...state.learnedGrammar, grammar],
  })),
}));

export default useGameStore; 