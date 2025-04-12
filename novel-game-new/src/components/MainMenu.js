import React from 'react';
import { useNavigate } from 'react-router-dom';
import {
  Box,
  Button,
  Container,
  Typography,
  Paper,
  Stack,
} from '@mui/material';
import { motion } from 'framer-motion';

const MainMenu = () => {
  const navigate = useNavigate();

  const containerVariants = {
    hidden: { opacity: 0 },
    visible: {
      opacity: 1,
      transition: {
        duration: 0.5,
        staggerChildren: 0.1
      }
    }
  };

  const itemVariants = {
    hidden: { opacity: 0, y: 20 },
    visible: {
      opacity: 1,
      y: 0,
      transition: { duration: 0.5 }
    }
  };

  return (
    <Box
      sx={{
        minHeight: '100vh',
        background: 'linear-gradient(to bottom, #1a1a1a, #000000)',
        display: 'flex',
        alignItems: 'center',
        justifyContent: 'center',
        p: 3,
      }}
    >
      <Container maxWidth="sm">
        <motion.div
          variants={containerVariants}
          initial="hidden"
          animate="visible"
        >
          <Paper
            elevation={3}
            sx={{
              p: 4,
              textAlign: 'center',
              backgroundColor: 'rgba(30, 30, 30, 0.9)',
              backdropFilter: 'blur(10px)',
              border: '1px solid rgba(255, 255, 255, 0.1)',
            }}
          >
            <motion.div variants={itemVariants}>
              <Typography
                variant="h1"
                sx={{
                  mb: 2,
                  background: 'linear-gradient(45deg, #90caf9 30%, #f48fb1 90%)',
                  backgroundClip: 'text',
                  textFillColor: 'transparent',
                  WebkitBackgroundClip: 'text',
                  WebkitTextFillColor: 'transparent',
                }}
              >
                日本語の冒険
              </Typography>
              <Typography
                variant="h2"
                sx={{ mb: 4, color: 'rgba(255, 255, 255, 0.7)' }}
              >
                Japanese Adventure
              </Typography>
            </motion.div>

            <Stack spacing={2}>
              <motion.div variants={itemVariants}>
                <Button
                  variant="contained"
                  color="primary"
                  fullWidth
                  size="large"
                  onClick={() => navigate('/game')}
                  sx={{ mb: 2 }}
                >
                  Start Game
                </Button>
              </motion.div>

              <motion.div variants={itemVariants}>
                <Button
                  variant="outlined"
                  color="secondary"
                  fullWidth
                  size="large"
                  onClick={() => {/* TODO: Implement load game */}}
                >
                  Load Game
                </Button>
              </motion.div>

              <motion.div variants={itemVariants}>
                <Button
                  variant="outlined"
                  color="primary"
                  fullWidth
                  size="large"
                  onClick={() => {/* TODO: Implement settings */}}
                >
                  Settings
                </Button>
              </motion.div>
            </Stack>
          </Paper>
        </motion.div>
      </Container>
    </Box>
  );
};

export default MainMenu; 