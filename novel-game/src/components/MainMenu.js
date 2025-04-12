import React from 'react';
import { useNavigate } from 'react-router-dom';
import { Box, Button, Typography, Container } from '@mui/material';
import { motion } from 'framer-motion';

const MainMenu = () => {
  const navigate = useNavigate();

  return (
    <Container maxWidth="sm">
      <Box
        sx={{
          minHeight: '100vh',
          display: 'flex',
          flexDirection: 'column',
          justifyContent: 'center',
          alignItems: 'center',
          textAlign: 'center',
          gap: 4,
        }}
      >
        <motion.div
          initial={{ opacity: 0, y: -20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.8 }}
        >
          <Typography variant="h2" component="h1" gutterBottom>
            Japanese Language Adventure
          </Typography>
          <Typography variant="h5" color="text.secondary" gutterBottom>
            Learn Japanese through an immersive visual novel experience
          </Typography>
        </motion.div>

        <motion.div
          initial={{ opacity: 0 }}
          animate={{ opacity: 1 }}
          transition={{ delay: 0.5 }}
        >
          <Button
            variant="contained"
            size="large"
            onClick={() => navigate('/game')}
            sx={{
              px: 4,
              py: 2,
              fontSize: '1.2rem',
              borderRadius: 2,
            }}
          >
            Start Game
          </Button>
        </motion.div>
      </Box>
    </Container>
  );
};

export default MainMenu; 