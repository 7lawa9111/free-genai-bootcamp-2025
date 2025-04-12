// Command handler utility for the visual novel game

import { useGameStateStore } from '../store/gameStateStore';

const COMMANDS = {
  LOOK: ['look', 'l', 'examine', 'x'],
  TAKE: ['take', 'get', 'pick up'],
  DROP: ['drop', 'put down'],
  INVENTORY: ['inventory', 'i', 'inv'],
  TALK: ['talk', 'speak', 't'],
  OPEN: ['open', 'unlock'],
  CLOSE: ['close', 'lock'],
  HELP: ['help', '?', 'h'],
  QUIT: ['quit', 'q', 'exit'],
};

/**
 * Process a command and return the appropriate response and game state updates
 * @param {string} command - The command to process
 * @param {Object} currentScene - The current scene object
 * @param {Object} gameState - The current game state
 * @returns {Object} - Response and game state updates
 */
const handleCommand = (command, currentScene) => {
  const gameState = useGameStateStore.getState();
  const { location, inventory, flags, items } = gameState;

  // Split command into verb and object
  const [verb, ...objectParts] = command.split(' ');
  const object = objectParts.join(' ');

  // Help command
  if (COMMANDS.HELP.includes(verb)) {
    return {
      text: 'Available commands:\n' +
        'look/l - Look around or examine something\n' +
        'take/get - Pick up an item\n' +
        'drop - Drop an item\n' +
        'inventory/i - Check your inventory\n' +
        'talk/t - Talk to someone\n' +
        'open - Open something\n' +
        'close - Close something\n' +
        'help/? - Show this help message\n' +
        'quit/q - Exit the game',
      translation: '利用可能なコマンド:\n' +
        'look/l - 周りを見る、または何かを調べる\n' +
        'take/get - アイテムを拾う\n' +
        'drop - アイテムを置く\n' +
        'inventory/i - インベントリを確認する\n' +
        'talk/t - 誰かと話す\n' +
        'open - 何かを開ける\n' +
        'close - 何かを閉める\n' +
        'help/? - このヘルプメッセージを表示\n' +
        'quit/q - ゲームを終了する'
    };
  }

  // Look command
  if (COMMANDS.LOOK.includes(verb)) {
    if (!object) {
      // Look around
      const sceneDescription = currentScene?.description || "You are in a post office. The clerk is waiting to help you.";
      const itemsInLocation = gameState.getItemsInLocation(location) || [];
      const npcsInLocation = currentScene?.npcs || [];
      
      let description = sceneDescription;
      if (itemsInLocation.length > 0) {
        description += '\nYou see: ' + itemsInLocation.map(item => item.name).join(', ');
      }
      if (npcsInLocation.length > 0) {
        description += '\nPeople here: ' + npcsInLocation.map(npc => npc.name).join(', ');
      }

      return {
        text: description,
        translation: currentScene?.descriptionJP || "郵便局にいます。係員があなたを待っています。"
      };
    } else {
      // Examine specific object
      const item = items.find(i => 
        i.name.toLowerCase() === object.toLowerCase() &&
        (i.location === location || inventory.includes(i.id))
      );

      if (item) {
        return {
          text: item.description || `It's a ${item.name}.`,
          translation: item.descriptionJP || `${item.nameJP}です。`
        };
      } else {
        return {
          text: "You don't see that here.",
          translation: 'それはここには見当たりません。'
        };
      }
    }
  }

  // Take command
  if (COMMANDS.TAKE.includes(verb)) {
    if (!object) {
      return {
        text: "What do you want to take?",
        translation: '何を取りたいですか？'
      };
    }

    const item = items.find(i => 
      i.name.toLowerCase() === object.toLowerCase() &&
      i.location === location &&
      !i.isFixed
    );

    if (item) {
      gameState.addToInventory(item.id);
      gameState.setItemState(item.id, 'location', 'inventory');
      return {
        text: `You take the ${item.name}.`,
        translation: `${item.nameJP}を取りました。`
      };
    } else {
      return {
        text: "You can't take that.",
        translation: 'それは取れません。'
      };
    }
  }

  // Drop command
  if (COMMANDS.DROP.includes(verb)) {
    if (!object) {
      return {
        text: "What do you want to drop?",
        translation: '何を置きたいですか？'
      };
    }

    const item = items.find(i => 
      i.name.toLowerCase() === object.toLowerCase() &&
      inventory.includes(i.id)
    );

    if (item) {
      gameState.removeFromInventory(item.id);
      gameState.setItemState(item.id, 'location', location);
      return {
        text: `You drop the ${item.name}.`,
        translation: `${item.nameJP}を置きました。`
      };
    } else {
      return {
        text: "You don't have that.",
        translation: 'それは持っていません。'
      };
    }
  }

  // Inventory command
  if (COMMANDS.INVENTORY.includes(verb)) {
    if (inventory.length === 0) {
      return {
        text: "Your inventory is empty.",
        translation: 'インベントリは空です。'
      };
    }

    const inventoryItems = inventory
      .map(id => items.find(i => i.id === id))
      .map(item => item.name);

    return {
      text: "Inventory:\n" + inventoryItems.join('\n'),
      translation: 'インベントリ:\n' + inventoryItems.join('\n')
    };
  }

  // Talk command
  if (COMMANDS.TALK.includes(verb)) {
    if (!object) {
      return {
        text: "Who do you want to talk to?",
        translation: '誰と話したいですか？'
      };
    }

    const npc = currentScene.npcs?.find(n => 
      n.name.toLowerCase() === object.toLowerCase()
    );

    if (npc) {
      return {
        text: npc.dialogue,
        translation: npc.dialogueJP
      };
    } else {
      return {
        text: "There's no one here by that name.",
        translation: 'その名前の人はここにいません。'
      };
    }
  }

  // Open command
  if (COMMANDS.OPEN.includes(verb)) {
    if (!object) {
      return {
        text: "What do you want to open?",
        translation: '何を開けたいですか？'
      };
    }

    const item = items.find(i => 
      i.name.toLowerCase() === object.toLowerCase() &&
      i.location === location &&
      i.isOpenable
    );

    if (item) {
      if (item.isLocked) {
        return {
          text: "It's locked.",
          translation: '鍵がかかっています。'
        };
      }
      gameState.setItemState(item.id, 'isOpen', true);
      return {
        text: `You open the ${item.name}.`,
        translation: `${item.nameJP}を開けました。`
      };
    } else {
      return {
        text: "You can't open that.",
        translation: 'それは開けられません。'
      };
    }
  }

  // Close command
  if (COMMANDS.CLOSE.includes(verb)) {
    if (!object) {
      return {
        text: "What do you want to close?",
        translation: '何を閉めたいですか？'
      };
    }

    const item = items.find(i => 
      i.name.toLowerCase() === object.toLowerCase() &&
      i.location === location &&
      i.isOpenable
    );

    if (item) {
      gameState.setItemState(item.id, 'isOpen', false);
      return {
        text: `You close the ${item.name}.`,
        translation: `${item.nameJP}を閉めました。`
      };
    } else {
      return {
        text: "You can't close that.",
        translation: 'それは閉められません。'
      };
    }
  }

  // Quit command
  if (COMMANDS.QUIT.includes(verb)) {
    return {
      text: "Thanks for playing!",
      translation: 'プレイしていただき、ありがとうございました！'
    };
  }

  // Unknown command
  return {
    text: "I don't understand that command.",
    translation: 'そのコマンドは理解できません。'
  };
};

export { handleCommand, COMMANDS };

// Helper functions
const getItemsInLocation = (currentScene, gameState) => {
  if (!currentScene.items) return [];
  
  return currentScene.items.filter(item => {
    // Check if the item is in the location (not in inventory)
    return !gameState.items || gameState.items[item.id] !== false;
  });
};

const getNPCsInLocation = (currentScene) => {
  return currentScene.npcs || [];
};

const getOpenableObjects = (currentScene, gameState) => {
  if (!currentScene.openableObjects) return [];
  
  return currentScene.openableObjects.filter(obj => {
    // Check if the object is in the location
    return true;
  });
};

const getCloseableObjects = (currentScene, gameState) => {
  if (!currentScene.openableObjects) return [];
  
  return currentScene.openableObjects.filter(obj => {
    // Check if the object is in the location and is open
    return gameState.flags && gameState.flags[obj.openFlag];
  });
};

const getDirectionTranslation = (direction) => {
  switch (direction) {
    case 'north': return '北';
    case 'up': return '上';
    case 'left': return '左';
    case 'right': return '右';
    case 'forward': return '前';
    default: return direction;
  }
}; 