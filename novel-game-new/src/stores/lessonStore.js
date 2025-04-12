import { create } from 'zustand';

const initialLessons = {
  lesson1: {
    id: 'lesson1',
    title: 'Basic Greetings',
    titleJp: 'あいさつ',
    background: '/images/backgrounds/classroom.jpg',
    scenes: [
      {
        text: "Welcome to your first Japanese lesson! I'm your sensei, and I'll be teaching you basic Japanese greetings.",
        textJp: "日本語の授業へようこそ！私があなたの先生です。基本的なあいさつを教えましょう。",
        character: 'sensei',
        position: 'center',
      },
      {
        text: "Let's start with 'Hello' - In Japanese, we say 'Konnichiwa'",
        textJp: "「こんにちは」から始めましょう。",
        character: 'sensei',
        position: 'center',
        vocabulary: {
          word: 'こんにちは',
          romaji: 'konnichiwa',
          meaning: 'hello (during the day)',
        },
      },
      {
        text: "Now, let's practice! Try saying 'Konnichiwa'",
        textJp: "では、練習しましょう！「こんにちは」と言ってみてください。",
        character: 'sensei',
        position: 'center',
        expectedInput: 'konnichiwa',
      },
    ],
    vocabulary: [
      {
        word: 'こんにちは',
        romaji: 'konnichiwa',
        meaning: 'hello (during the day)',
      },
      {
        word: 'おはようございます',
        romaji: 'ohayou gozaimasu',
        meaning: 'good morning (polite)',
      },
      {
        word: 'こんばんは',
        romaji: 'konbanwa',
        meaning: 'good evening',
      },
      {
        word: 'さようなら',
        romaji: 'sayounara',
        meaning: 'goodbye',
      },
    ],
    grammar: [
      {
        point: 'Greetings Time of Day',
        explanation: 'Japanese greetings change based on the time of day:',
        examples: [
          {
            jp: 'おはようございます',
            romaji: 'ohayou gozaimasu',
            en: 'good morning',
            note: 'Used in the morning',
          },
          {
            jp: 'こんにちは',
            romaji: 'konnichiwa',
            en: 'hello/good afternoon',
            note: 'Used during the day',
          },
          {
            jp: 'こんばんは',
            romaji: 'konbanwa',
            en: 'good evening',
            note: 'Used in the evening',
          },
        ],
      },
    ],
  },
  // Add more lessons here
};

const useLessonStore = create((set, get) => ({
  // State
  lessons: initialLessons,
  currentLessonId: null,
  currentSceneIndex: 0,
  
  // Actions
  startLesson: (lessonId) => set(() => ({ 
    currentLessonId: lessonId,
    currentSceneIndex: 0,
  })),

  nextScene: () => set((state) => {
    const currentLesson = state.lessons[state.currentLessonId];
    if (!currentLesson) return state;

    const nextIndex = state.currentSceneIndex + 1;
    if (nextIndex >= currentLesson.scenes.length) {
      return state;
    }

    return { currentSceneIndex: nextIndex };
  }),

  previousScene: () => set((state) => ({
    currentSceneIndex: Math.max(0, state.currentSceneIndex - 1),
  })),

  resetLesson: () => set(() => ({
    currentSceneIndex: 0,
  })),
}));

export default useLessonStore; 