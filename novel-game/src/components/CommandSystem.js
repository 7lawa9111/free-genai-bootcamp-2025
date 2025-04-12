import React, { useState, useEffect, useRef } from 'react';
import {
  Box,
  TextField,
  IconButton,
  Paper,
  Typography,
  List,
  ListItem,
  ListItemText,
  Divider,
  Tooltip,
  Collapse,
  Button,
} from '@mui/material';
import {
  Send as SendIcon,
  Help as HelpIcon,
  KeyboardArrowDown as KeyboardArrowDownIcon,
  KeyboardArrowUp as KeyboardArrowUpIcon,
} from '@mui/icons-material';

const CommandSystem = ({ onCommand, currentScene }) => {
  const [command, setCommand] = useState('');
  const [commandHistory, setCommandHistory] = useState([]);
  const [showHelp, setShowHelp] = useState(false);
  const [suggestions, setSuggestions] = useState([]);
  const inputRef = useRef(null);

  // Available commands
  const availableCommands = [
    { command: 'look', description: 'Examine surroundings or objects' },
    { command: 'move', description: 'Go in a direction' },
    { command: 'north', description: 'Move north' },
    { command: 'up', description: 'Move up' },
    { command: 'left', description: 'Move left' },
    { command: 'right', description: 'Move right' },
    { command: 'take', description: 'Pick up items' },
    { command: 'drop', description: 'Discard items' },
    { command: 'talk', description: 'Communicate with NPCs' },
    { command: 'say', description: 'Say something to NPCs' },
    { command: 'use', description: 'Interact with items' },
    { command: 'give', description: 'Transfer items to others' },
    { command: 'open', description: 'Open doors, chests, etc.' },
    { command: 'close', description: 'Close doors, containers' },
    { command: 'eat', description: 'Consume food items' },
    { command: 'inventory', description: 'Check carried items' },
    { command: 'drink', description: 'Consume liquids' },
    { command: 'help', description: 'View commands/instructions' },
  ];

  // Focus input on mount
  useEffect(() => {
    if (inputRef.current) {
      inputRef.current.focus();
    }
  }, []);

  // Update suggestions based on current input
  useEffect(() => {
    if (command.trim() === '') {
      setSuggestions([]);
      return;
    }

    const filtered = availableCommands.filter(cmd => 
      cmd.command.startsWith(command.toLowerCase())
    );
    setSuggestions(filtered);
  }, [command]);

  const handleCommandChange = (e) => {
    setCommand(e.target.value);
  };

  const handleCommandSubmit = (e) => {
    e.preventDefault();
    
    if (command.trim() === '') return;
    
    // Add to history
    setCommandHistory([...commandHistory, { 
      command: command, 
      timestamp: new Date().toISOString() 
    }]);
    
    // Process command
    processCommand(command);
    
    // Clear input
    setCommand('');
    setSuggestions([]);
  };

  const processCommand = (cmd) => {
    const normalizedCmd = cmd.toLowerCase().trim();
    
    // Check if it's a help command
    if (normalizedCmd === 'help') {
      setShowHelp(!showHelp);
      return;
    }
    
    // Check if it's a valid command
    const validCommand = availableCommands.find(c => c.command === normalizedCmd);
    
    if (validCommand) {
      // Pass the command to the parent component
      onCommand(normalizedCmd, currentScene);
    } else {
      // Check for partial matches
      const partialMatch = availableCommands.find(c => 
        normalizedCmd.startsWith(c.command)
      );
      
      if (partialMatch) {
        onCommand(partialMatch.command, currentScene);
      } else {
        // Command not recognized
        onCommand('unknown', currentScene);
      }
    }
  };

  const handleSuggestionClick = (suggestion) => {
    setCommand(suggestion.command);
    inputRef.current.focus();
  };

  const toggleHelp = () => {
    setShowHelp(!showHelp);
  };

  return (
    <Paper 
      elevation={3} 
      sx={{ 
        p: 2, 
        position: 'relative',
        backgroundColor: 'rgba(255, 255, 255, 0.9)',
        backdropFilter: 'blur(5px)',
      }}
    >
      <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 1 }}>
        <Typography variant="subtitle1" sx={{ fontWeight: 'bold' }}>
          Command Console
        </Typography>
        <Tooltip title="Toggle help">
          <IconButton size="small" onClick={toggleHelp}>
            {showHelp ? <KeyboardArrowUpIcon /> : <KeyboardArrowDownIcon />}
          </IconButton>
        </Tooltip>
      </Box>
      
      <Collapse in={showHelp}>
        <Paper variant="outlined" sx={{ p: 2, mb: 2, maxHeight: '200px', overflow: 'auto' }}>
          <Typography variant="subtitle2" gutterBottom>
            Available Commands:
          </Typography>
          <List dense>
            {availableCommands.map((cmd, index) => (
              <ListItem key={index} sx={{ py: 0.5 }}>
                <ListItemText 
                  primary={
                    <Typography variant="body2">
                      <strong>{cmd.command}</strong> - {cmd.description}
                    </Typography>
                  }
                />
              </ListItem>
            ))}
          </List>
        </Paper>
      </Collapse>
      
      <form onSubmit={handleCommandSubmit}>
        <Box sx={{ display: 'flex', alignItems: 'center' }}>
          <TextField
            inputRef={inputRef}
            fullWidth
            variant="outlined"
            placeholder="Enter command..."
            value={command}
            onChange={handleCommandChange}
            size="small"
            sx={{ mr: 1 }}
          />
          <Tooltip title="Send command">
            <IconButton 
              type="submit" 
              color="primary"
              disabled={command.trim() === ''}
            >
              <SendIcon />
            </IconButton>
          </Tooltip>
        </Box>
      </form>
      
      {suggestions.length > 0 && (
        <Paper 
          variant="outlined" 
          sx={{ 
            mt: 1, 
            maxHeight: '150px', 
            overflow: 'auto',
            position: 'absolute',
            width: '100%',
            zIndex: 10,
            left: 0,
          }}
        >
          <List dense>
            {suggestions.map((suggestion, index) => (
              <ListItem 
                key={index} 
                button 
                onClick={() => handleSuggestionClick(suggestion)}
                sx={{ py: 0.5 }}
              >
                <ListItemText 
                  primary={
                    <Typography variant="body2">
                      <strong>{suggestion.command}</strong> - {suggestion.description}
                    </Typography>
                  }
                />
              </ListItem>
            ))}
          </List>
        </Paper>
      )}
      
      {commandHistory.length > 0 && (
        <Box sx={{ mt: 2, maxHeight: '150px', overflow: 'auto' }}>
          <Typography variant="caption" color="text.secondary">
            Command History:
          </Typography>
          <List dense>
            {commandHistory.slice(-5).map((item, index) => (
              <ListItem key={index} sx={{ py: 0.5 }}>
                <ListItemText 
                  primary={
                    <Typography variant="body2">
                      <strong>{item.command}</strong>
                      <Typography 
                        component="span" 
                        variant="caption" 
                        color="text.secondary"
                        sx={{ ml: 1 }}
                      >
                        {new Date(item.timestamp).toLocaleTimeString()}
                      </Typography>
                    </Typography>
                  }
                />
              </ListItem>
            ))}
          </List>
        </Box>
      )}
    </Paper>
  );
};

export default CommandSystem; 