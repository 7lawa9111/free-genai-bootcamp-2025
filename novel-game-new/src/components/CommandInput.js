import React, { useState, useRef, useEffect } from 'react';
import {
  Box,
  TextField,
  IconButton,
  Paper,
  List,
  ListItem,
  ListItemText,
  Typography,
} from '@mui/material';
import { Send as SendIcon, Clear as ClearIcon } from '@mui/icons-material';
import { motion, AnimatePresence } from 'framer-motion';
import useGameStore from '../stores/gameStore';

const CommandInput = ({ onSubmit }) => {
  const [input, setInput] = useState('');
  const inputRef = useRef(null);
  const { commandHistory, clearHistory } = useGameStore();

  // Focus input on mount
  useEffect(() => {
    inputRef.current?.focus();
  }, []);

  const handleSubmit = (e) => {
    e.preventDefault();
    if (input.trim()) {
      onSubmit(input.trim());
      setInput('');
    }
  };

  const handleClear = () => {
    clearHistory();
  };

  return (
    <Box sx={{ width: '100%' }}>
      {/* Command History */}
      <Paper
        elevation={3}
        sx={{
          backgroundColor: 'rgba(0, 0, 0, 0.8)',
          color: 'white',
          mb: 2,
          maxHeight: '200px',
          overflowY: 'auto',
          display: commandHistory.length ? 'block' : 'none',
        }}
      >
        <List>
          <AnimatePresence>
            {commandHistory.map((item, index) => (
              <motion.div
                key={index}
                initial={{ opacity: 0, x: -20 }}
                animate={{ opacity: 1, x: 0 }}
                exit={{ opacity: 0, x: 20 }}
                transition={{ duration: 0.2 }}
              >
                <ListItem
                  sx={{
                    borderBottom: '1px solid rgba(255, 255, 255, 0.1)',
                    py: 1,
                  }}
                >
                  <ListItemText
                    primary={
                      <Typography
                        component="span"
                        sx={{ color: 'primary.main', fontWeight: 'bold' }}
                      >
                        {'> ' + item.command}
                      </Typography>
                    }
                    secondary={
                      <Typography
                        component="span"
                        sx={{ color: 'rgba(255, 255, 255, 0.7)' }}
                      >
                        {item.response}
                      </Typography>
                    }
                  />
                </ListItem>
              </motion.div>
            ))}
          </AnimatePresence>
        </List>
      </Paper>

      {/* Command Input */}
      <Paper
        component="form"
        onSubmit={handleSubmit}
        sx={{
          p: '2px 4px',
          display: 'flex',
          alignItems: 'center',
          backgroundColor: 'rgba(0, 0, 0, 0.8)',
          border: '1px solid rgba(255, 255, 255, 0.1)',
        }}
      >
        <TextField
          fullWidth
          value={input}
          onChange={(e) => setInput(e.target.value)}
          placeholder="Type your command..."
          inputRef={inputRef}
          variant="standard"
          sx={{
            mx: 1,
            '& .MuiInputBase-input': {
              color: 'white',
            },
            '& .MuiInput-underline:before': {
              borderBottomColor: 'rgba(255, 255, 255, 0.1)',
            },
            '& .MuiInput-underline:hover:not(.Mui-disabled):before': {
              borderBottomColor: 'rgba(255, 255, 255, 0.2)',
            },
          }}
        />
        <IconButton
          type="submit"
          color="primary"
          sx={{ p: '10px' }}
          disabled={!input.trim()}
        >
          <SendIcon />
        </IconButton>
        {commandHistory.length > 0 && (
          <IconButton
            color="secondary"
            sx={{ p: '10px' }}
            onClick={handleClear}
          >
            <ClearIcon />
          </IconButton>
        )}
      </Paper>
    </Box>
  );
};

export default CommandInput; 