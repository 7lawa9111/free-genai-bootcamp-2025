import { create } from 'zustand';
import { persist } from 'zustand/middleware';

const useVocabularyStore = create(
  persist(
    (set, get) => ({
      vocabulary: [],
      
      // Add a new word to the vocabulary
      addWord: (word) => {
        const { vocabulary } = get();
        
        // Check if word already exists
        if (!vocabulary.some(item => item.japanese === word.japanese)) {
          set({
            vocabulary: [...vocabulary, { ...word, learned: false, dateAdded: new Date().toISOString() }]
          });
        }
      },
      
      // Mark a word as learned
      markAsLearned: (japanese) => {
        const { vocabulary } = get();
        set({
          vocabulary: vocabulary.map(word => 
            word.japanese === japanese ? { ...word, learned: true } : word
          )
        });
      },
      
      // Mark a word as not learned
      markAsNotLearned: (japanese) => {
        const { vocabulary } = get();
        set({
          vocabulary: vocabulary.map(word => 
            word.japanese === japanese ? { ...word, learned: false } : word
          )
        });
      },
      
      // Remove a word from vocabulary
      removeWord: (japanese) => {
        const { vocabulary } = get();
        set({
          vocabulary: vocabulary.filter(word => word.japanese !== japanese)
        });
      },
      
      // Get all vocabulary words
      getAllWords: () => get().vocabulary,
      
      // Get learned words
      getLearnedWords: () => get().vocabulary.filter(word => word.learned),
      
      // Get unlearned words
      getUnlearnedWords: () => get().vocabulary.filter(word => !word.learned),
      
      // Clear all vocabulary
      clearVocabulary: () => set({ vocabulary: [] }),
    }),
    {
      name: 'vocabulary-storage',
    }
  )
);

export default useVocabularyStore; 