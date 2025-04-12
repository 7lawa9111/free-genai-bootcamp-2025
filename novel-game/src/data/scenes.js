import { CHARACTERS, LOCATIONS } from './constants';

export const scenes = {
  start: {
    id: 'start',
    location: LOCATIONS.POST_OFFICE,
    character: CHARACTERS.TANAKA_HIROSHI,
    dialogue: [
      {
        text: 'いらっしゃいませ。何かお手伝いできますか？',
        translation: 'Welcome. How may I help you?',
      },
    ],
    choices: [
      {
        text: '切手を買いたいです。',
        translation: 'I want to buy stamps.',
        nextScene: 'buy-stamps',
      },
      {
        text: '小包を送りたいです。',
        translation: 'I want to send a package.',
        nextScene: 'send-package',
      },
    ],
  },
  'buy-stamps': {
    id: 'buy-stamps',
    location: LOCATIONS.POST_OFFICE,
    character: CHARACTERS.TANAKA_HIROSHI,
    dialogue: [
      {
        text: '切手は何枚必要ですか？',
        translation: 'How many stamps do you need?',
      },
    ],
    choices: [
      {
        text: '5枚ください。',
        translation: 'Please give me 5 stamps.',
        nextScene: 'stamps-5',
      },
      {
        text: '10枚ください。',
        translation: 'Please give me 10 stamps.',
        nextScene: 'stamps-10',
      },
    ],
  },
  'send-package': {
    id: 'send-package',
    location: LOCATIONS.POST_OFFICE,
    character: CHARACTERS.TANAKA_HIROSHI,
    dialogue: [
      {
        text: '小包の送り先はどこですか？',
        translation: 'Where is the package going?',
      },
    ],
    choices: [
      {
        text: 'アメリカです。',
        translation: 'To America.',
        nextScene: 'package-america',
      },
      {
        text: 'イギリスです。',
        translation: 'To England.',
        nextScene: 'package-uk',
      },
    ],
  },
  'stamps-5': {
    id: 'stamps-5',
    location: LOCATIONS.POST_OFFICE,
    character: CHARACTERS.TANAKA_HIROSHI,
    dialogue: [
      {
        text: '5枚の切手で250円です。',
        translation: '5 stamps will be 250 yen.',
      },
    ],
    choices: [
      {
        text: 'はい、お願いします。',
        translation: 'Yes, please.',
        nextScene: 'stamps-payment',
      },
      {
        text: '他の切手はありますか？',
        translation: 'Do you have other stamps?',
        nextScene: 'other-stamps',
      },
    ],
  },
  'stamps-10': {
    id: 'stamps-10',
    location: LOCATIONS.POST_OFFICE,
    character: CHARACTERS.TANAKA_HIROSHI,
    dialogue: [
      {
        text: '10枚の切手で500円です。',
        translation: '10 stamps will be 500 yen.',
      },
    ],
    choices: [
      {
        text: 'はい、お願いします。',
        translation: 'Yes, please.',
        nextScene: 'stamps-payment',
      },
      {
        text: '他の切手はありますか？',
        translation: 'Do you have other stamps?',
        nextScene: 'other-stamps',
      },
    ],
  },
  'package-america': {
    id: 'package-america',
    location: LOCATIONS.POST_OFFICE,
    character: CHARACTERS.TANAKA_HIROSHI,
    dialogue: [
      {
        text: 'アメリカへの小包は航空便で3,500円、船便で1,200円です。どちらにしますか？',
        translation: 'Shipping to America is 3,500 yen by air mail and 1,200 yen by sea mail. Which would you prefer?',
      },
    ],
    choices: [
      {
        text: '航空便でお願いします。',
        translation: 'Air mail, please.',
        nextScene: 'package-air',
      },
      {
        text: '船便でお願いします。',
        translation: 'Sea mail, please.',
        nextScene: 'package-sea',
      },
    ],
  },
  'package-uk': {
    id: 'package-uk',
    location: LOCATIONS.POST_OFFICE,
    character: CHARACTERS.TANAKA_HIROSHI,
    dialogue: [
      {
        text: 'イギリスへの小包は航空便で3,800円、船便で1,300円です。どちらにしますか？',
        translation: 'Shipping to England is 3,800 yen by air mail and 1,300 yen by sea mail. Which would you prefer?',
      },
    ],
    choices: [
      {
        text: '航空便でお願いします。',
        translation: 'Air mail, please.',
        nextScene: 'package-air',
      },
      {
        text: '船便でお願いします。',
        translation: 'Sea mail, please.',
        nextScene: 'package-sea',
      },
    ],
  },
  // Add more scenes as needed
};

export const getScene = (sceneId) => {
  return scenes[sceneId] || scenes.start;
}; 