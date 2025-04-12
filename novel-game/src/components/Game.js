import React, { useState, useEffect } from 'react';
import { Box, Typography, Button, Container, Fade, IconButton, Drawer } from '@mui/material';
import { motion, AnimatePresence } from 'framer-motion';
import { getScene } from '../data/scenes';
import useGameStore from '../store/gameStore';
import VocabularyList from './VocabularyList';
import VocabularyExtractor from './VocabularyExtractor';
import MenuBookIcon from '@mui/icons-material/MenuBook';
import CommandInput from './CommandInput';

// This will be moved to a separate file later
const initialScene = {
  id: 'start',
  location: 'Post Office',
  character: {
    name: 'Tanaka Hiroshi',
    image: 'https://via.placeholder.com/800x600/4a90e2/ffffff?text=Post+Office+Clerk',
  },
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
};

const Game = () => {
  const [currentScene, setCurrentScene] = useState(getScene('start'));
  const [showTranslation, setShowTranslation] = useState(false);
  const [isTransitioning, setIsTransitioning] = useState(false);
  const [vocabularyDrawerOpen, setVocabularyDrawerOpen] = useState(false);
  const { addVisitedScene, addChoice } = useGameStore();

  useEffect(() => {
    // Initialize the game with the first scene
    addVisitedScene('start');
  }, []);

  const handleChoice = (nextScene) => {
    // Start transition
    setIsTransitioning(true);
    
    // Add the current scene to visited scenes
    addVisitedScene(currentScene.id);
    
    // Add the choice to the game progress
    const choice = currentScene.choices.find(c => c.nextScene === nextScene);
    if (choice) {
      addChoice({
        sceneId: currentScene.id,
        choice: choice.text,
        nextScene: nextScene
      });
    }
    
    // Wait for transition animation to complete
    setTimeout(() => {
      // Get the next scene from our scenes data
      const nextSceneData = getScene(nextScene);
      setCurrentScene(nextSceneData);
      setIsTransitioning(false);
    }, 500);
  };

  const toggleVocabularyDrawer = () => {
    setVocabularyDrawerOpen(!vocabularyDrawerOpen);
  };

  return (
    <Container maxWidth="lg" className="game-container">
      <Box
        sx={{
          position: 'relative',
          minHeight: '80vh',
          display: 'flex',
          flexDirection: 'column',
          overflow: 'hidden',
          borderRadius: 2,
          boxShadow: '0 4px 20px rgba(0,0,0,0.2)',
          marginBottom: '100px', // Add margin to make room for CommandInput
        }}
      >
        {/* Vocabulary Drawer Toggle Button */}
        <IconButton
          onClick={toggleVocabularyDrawer}
          sx={{
            position: 'absolute',
            top: 20,
            right: 20,
            zIndex: 10,
            backgroundColor: 'rgba(255, 255, 255, 0.8)',
            '&:hover': {
              backgroundColor: 'rgba(255, 255, 255, 0.9)',
            },
          }}
        >
          <MenuBookIcon />
        </IconButton>

        {/* Vocabulary Drawer */}
        <Drawer
          anchor="right"
          open={vocabularyDrawerOpen}
          onClose={toggleVocabularyDrawer}
          PaperProps={{
            sx: {
              width: '400px',
              maxWidth: '100%',
            },
          }}
        >
          <VocabularyList />
        </Drawer>

        {/* Background Image */}
        <AnimatePresence mode="wait">
          <motion.div
            key={currentScene.id}
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            exit={{ opacity: 0 }}
            transition={{ duration: 0.5 }}
            style={{
              position: 'absolute',
              top: 0,
              left: 0,
              right: 0,
              bottom: 0,
              zIndex: -1,
            }}
          >
            <Box
              sx={{
                position: 'absolute',
                top: 0,
                left: 0,
                right: 0,
                bottom: 0,
                backgroundImage: `url(${currentScene.location.image})`,
                backgroundSize: 'cover',
                backgroundPosition: 'center',
                filter: 'brightness(0.7)',
              }}
            />
          </motion.div>
        </AnimatePresence>

        {/* Location Header */}
        <Typography
          variant="h6"
          sx={{
            position: 'absolute',
            top: 20,
            left: 20,
            color: 'white',
            textShadow: '2px 2px 4px rgba(0,0,0,0.5)',
            zIndex: 1,
            backgroundColor: 'rgba(0,0,0,0.3)',
            padding: '8px 16px',
            borderRadius: 1,
          }}
        >
          {currentScene.location.name}
        </Typography>

        {/* Character Image */}
        <Box
          sx={{
            flex: 1,
            display: 'flex',
            justifyContent: 'center',
            alignItems: 'center',
            position: 'relative',
          }}
        >
          <AnimatePresence mode="wait">
            <motion.div
              key={`${currentScene.id}-${currentScene.character.name}`}
              initial={{ opacity: 0, scale: 0.9, y: 20 }}
              animate={{ opacity: 1, scale: 1, y: 0 }}
              exit={{ opacity: 0, scale: 0.9, y: -20 }}
              transition={{ duration: 0.5 }}
              style={{
                maxHeight: '70vh',
                maxWidth: '100%',
                display: 'flex',
                justifyContent: 'center',
                alignItems: 'center',
              }}
            >
              <img
                src={currentScene.character.image}
                alt={currentScene.character.name}
                style={{
                  maxHeight: '100%',
                  maxWidth: '100%',
                  objectFit: 'contain',
                  filter: isTransitioning ? 'brightness(0.7)' : 'none',
                }}
              />
            </motion.div>
          </AnimatePresence>
        </Box>

        {/* Dialogue Box */}
        <motion.div
          initial={{ y: 100 }}
          animate={{ y: 0 }}
          className="dialog-box"
          style={{
            position: 'relative',
            zIndex: 2,
            opacity: isTransitioning ? 0.7 : 1,
            transition: 'opacity 0.3s',
          }}
        >
          <Typography variant="h6" className="character-name">
            {currentScene.character.name}
          </Typography>
          
          <Box sx={{ position: 'relative' }}>
            <Typography
              variant="body1"
              className="dialog-text"
              onClick={() => setShowTranslation(!showTranslation)}
              sx={{ cursor: 'pointer' }}
            >
              {currentScene.dialogue[0].text}
              <Fade in={showTranslation}>
                <motion.div
                  initial={{ opacity: 0 }}
                  animate={{ opacity: 1 }}
                  className="translation"
                >
                  <Typography
                    variant="body2"
                    color="text.secondary"
                    sx={{ mt: 1 }}
                  >
                    {currentScene.dialogue[0].translation}
                  </Typography>
                </motion.div>
              </Fade>
            </Typography>
            
            {/* Vocabulary Extractor */}
            <VocabularyExtractor 
              text={currentScene.dialogue[0].text} 
              translation={currentScene.dialogue[0].translation} 
            />
          </Box>

          {/* Choices */}
          <Box sx={{ mt: 2 }}>
            {currentScene.choices.map((choice, index) => (
              <motion.div
                key={index}
                initial={{ opacity: 0, x: -20 }}
                animate={{ opacity: 1, x: 0 }}
                transition={{ delay: index * 0.2 }}
              >
                <Button
                  variant="outlined"
                  fullWidth
                  onClick={() => handleChoice(choice.nextScene)}
                  sx={{ 
                    mb: 1,
                    transition: 'all 0.2s',
                    '&:hover': {
                      backgroundColor: 'rgba(74, 144, 226, 0.1)',
                      transform: 'translateY(-2px)',
                    },
                  }}
                  disabled={isTransitioning}
                >
                  <Box sx={{ textAlign: 'left', width: '100%' }}>
                    <Typography variant="body1">{choice.text}</Typography>
                    <Typography variant="caption" color="text.secondary">
                      {choice.translation}
                    </Typography>
                  </Box>
                </Button>
              </motion.div>
            ))}
          </Box>
        </motion.div>
      </Box>
      
      {/* Add CommandInput component */}
      <CommandInput />
    </Container>
  );
};

export default Game; 