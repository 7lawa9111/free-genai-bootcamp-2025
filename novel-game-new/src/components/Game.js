import React, { useEffect } from 'react';
import { Box, Container, Paper, Typography } from '@mui/material';
import { motion, AnimatePresence } from 'framer-motion';
import CommandInput from './CommandInput';
import useGameStore from '../stores/gameStore';
import useLessonStore from '../stores/lessonStore';

const Game = () => {
  const {
    currentText,
    currentBackground,
    characters,
    commandHistory,
    addCommand,
    updateScene,
  } = useGameStore();

  const {
    lessons,
    currentLessonId,
    currentSceneIndex,
    startLesson,
    nextScene,
    previousScene,
  } = useLessonStore();

  // Get current lesson and scene
  const currentLesson = currentLessonId ? lessons[currentLessonId] : null;
  const currentScene = currentLesson?.scenes[currentSceneIndex];
  const isLastScene = currentLesson ? currentSceneIndex >= currentLesson.scenes.length - 1 : true;

  // Initialize the game state
  useEffect(() => {
    if (!currentLessonId) {
      startLesson('lesson1');
    }
  }, [currentLessonId, startLesson]);

  // Update scene state when scene changes
  useEffect(() => {
    if (!currentScene) return;

    updateScene({
      text: currentScene.text,
      background: currentLesson?.background,
      characterId: currentScene.character,
      position: currentScene.position,
      vocabulary: currentScene.vocabulary,
    });
  }, [currentScene, currentLesson, updateScene]);

  const handleCommand = (command) => {
    let response = '';
    
    switch (command.toLowerCase()) {
      case 'help':
        response = `Available commands:
- next: Go to next scene
- back: Go to previous scene
- look: Describe the current scene
- vocabulary: Show learned vocabulary
- grammar: Show learned grammar points
- clear: Clear command history`;
        break;
        
      case 'next':
        if (isLastScene) {
          response = "You've reached the end of this lesson!";
        } else {
          nextScene();
          response = "Moving to next scene...";
        }
        break;
        
      case 'back':
        previousScene();
        response = "Going back to previous scene...";
        break;
        
      case 'look':
        response = currentLesson 
          ? `You are in ${currentLesson.title} (${currentLesson.titleJp}).\n${currentText}`
          : "No lesson is currently active.";
        break;
        
      case 'vocabulary':
        const vocab = currentLesson?.vocabulary || [];
        response = vocab.length ? 
          "Vocabulary in this lesson:\n" + vocab.map(v => 
            `${v.word} (${v.romaji}) - ${v.meaning}`
          ).join('\n') :
          "No vocabulary words yet!";
        break;
        
      case 'grammar':
        const grammar = currentLesson?.grammar || [];
        response = grammar.length ?
          "Grammar points in this lesson:\n" + grammar.map(g =>
            `${g.point}:\n${g.explanation}\n${g.examples.map(ex =>
              `${ex.jp} (${ex.romaji}) - ${ex.en}\n${ex.note}`
            ).join('\n')}`
          ).join('\n\n') :
          "No grammar points yet!";
        break;

      case 'clear':
        response = 'Command history cleared.';
        break;
        
      default:
        // Check if this is an expected input for the current scene
        if (currentScene?.expectedInput && 
            command.toLowerCase() === currentScene.expectedInput.toLowerCase()) {
          response = "Correct! That's perfect!";
          nextScene();
        } else {
          response = `Command not recognized. Type "help" for available commands.`;
        }
    }
    
    addCommand(command, response);
  };

  return (
    <Box
      sx={{
        height: '100vh',
        width: '100vw',
        position: 'relative',
        overflow: 'hidden',
        background: '#000',
      }}
    >
      {/* Background Layer */}
      <AnimatePresence mode="wait">
        <motion.div
          key={currentBackground}
          initial={{ opacity: 0 }}
          animate={{ opacity: 1 }}
          exit={{ opacity: 0 }}
          transition={{ duration: 0.5 }}
          style={{
            position: 'absolute',
            width: '100%',
            height: '100%',
          }}
        >
          <Box
            component="img"
            src={currentBackground}
            alt="background"
            sx={{
              width: '100%',
              height: '100%',
              objectFit: 'cover',
              opacity: 0.8,
            }}
          />
        </motion.div>
      </AnimatePresence>

      {/* Character Layer */}
      <Box
        sx={{
          position: 'absolute',
          width: '100%',
          height: '100%',
          display: 'flex',
          justifyContent: 'center',
          alignItems: 'flex-end',
          pb: '20%',
        }}
      >
        {Object.entries(characters).map(([id, character]) => (
          character.visible && (
            <motion.div
              key={id}
              initial={{ opacity: 0, y: 20 }}
              animate={{ opacity: 1, y: 0 }}
              exit={{ opacity: 0, y: 20 }}
              transition={{ duration: 0.3 }}
              style={{
                position: 'absolute',
                left: character.position === 'left' ? '20%' : 
                      character.position === 'right' ? '60%' : '40%',
                transform: 'translateX(-50%)',
              }}
            >
              <Box
                component="img"
                src={character.image}
                alt={character.nameRomaji}
                sx={{
                  height: '70vh',
                  maxWidth: '100%',
                  objectFit: 'contain',
                }}
              />
            </motion.div>
          )
        ))}
      </Box>

      {/* Text and Command Interface */}
      <Container
        maxWidth="lg"
        sx={{
          position: 'absolute',
          bottom: 0,
          left: 0,
          right: 0,
          pb: 2,
        }}
      >
        <Paper
          elevation={3}
          sx={{
            backgroundColor: 'rgba(0, 0, 0, 0.8)',
            color: 'white',
            p: 2,
            mb: 2,
            borderRadius: 2,
            minHeight: '100px',
          }}
        >
          <AnimatePresence mode="wait">
            <motion.div
              key={currentText}
              initial={{ opacity: 0, y: 20 }}
              animate={{ opacity: 1, y: 0 }}
              exit={{ opacity: 0, y: -20 }}
              transition={{ duration: 0.3 }}
            >
              <Typography variant="body1">
                {currentText}
              </Typography>
            </motion.div>
          </AnimatePresence>
        </Paper>
        
        <CommandInput onSubmit={handleCommand} history={commandHistory} />
      </Container>
    </Box>
  );
};

export default Game; 