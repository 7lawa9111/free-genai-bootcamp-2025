import React, { useState } from 'react';
import {
  Box,
  Typography,
  List,
  ListItem,
  ListItemText,
  ListItemSecondaryAction,
  IconButton,
  Chip,
  Tabs,
  Tab,
  TextField,
  InputAdornment,
  Paper,
  Divider,
  Tooltip,
  Button,
} from '@mui/material';
import {
  CheckCircle as CheckCircleIcon,
  Cancel as CancelIcon,
  Delete as DeleteIcon,
  Search as SearchIcon,
  Sort as SortIcon,
  Quiz as QuizIcon,
} from '@mui/icons-material';
import useVocabularyStore from '../store/vocabularyStore';
import VocabularyQuiz from './VocabularyQuiz';

const VocabularyList = () => {
  const [tabValue, setTabValue] = useState(0);
  const [searchTerm, setSearchTerm] = useState('');
  const [sortBy, setSortBy] = useState('dateAdded'); // 'dateAdded', 'japanese', 'english'
  const [showQuiz, setShowQuiz] = useState(false);
  
  const {
    vocabulary,
    markAsLearned,
    markAsNotLearned,
    removeWord,
    getLearnedWords,
    getUnlearnedWords,
  } = useVocabularyStore();

  const handleTabChange = (event, newValue) => {
    setTabValue(newValue);
  };

  const handleSearchChange = (event) => {
    setSearchTerm(event.target.value);
  };

  const handleSortChange = () => {
    if (sortBy === 'dateAdded') setSortBy('japanese');
    else if (sortBy === 'japanese') setSortBy('english');
    else setSortBy('dateAdded');
  };

  const getSortIcon = () => {
    switch (sortBy) {
      case 'dateAdded': return '📅';
      case 'japanese': return 'あ';
      case 'english': return 'A';
      default: return '📅';
    }
  };

  const filterAndSortWords = (words) => {
    // Filter by search term
    let filtered = words.filter(word => 
      word.japanese.toLowerCase().includes(searchTerm.toLowerCase()) ||
      word.english.toLowerCase().includes(searchTerm.toLowerCase())
    );
    
    // Sort words
    return filtered.sort((a, b) => {
      if (sortBy === 'dateAdded') {
        return new Date(b.dateAdded) - new Date(a.dateAdded);
      } else if (sortBy === 'japanese') {
        return a.japanese.localeCompare(b.japanese);
      } else if (sortBy === 'english') {
        return a.english.localeCompare(b.english);
      }
      return 0;
    });
  };

  const displayWords = tabValue === 0 
    ? filterAndSortWords(getUnlearnedWords())
    : filterAndSortWords(getLearnedWords());

  const toggleQuiz = () => {
    setShowQuiz(!showQuiz);
  };

  return (
    <Paper elevation={3} sx={{ p: 2, maxHeight: '80vh', overflow: 'auto' }}>
      <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 2 }}>
        <Typography variant="h5">
          Vocabulary List
        </Typography>
        <Button
          variant="contained"
          color="primary"
          startIcon={<QuizIcon />}
          onClick={toggleQuiz}
          size="small"
        >
          Quiz
        </Button>
      </Box>
      
      {showQuiz ? (
        <VocabularyQuiz />
      ) : (
        <>
          <Box sx={{ display: 'flex', justifyContent: 'space-between', mb: 2 }}>
            <Tabs value={tabValue} onChange={handleTabChange} sx={{ mb: 2 }}>
              <Tab label={`To Learn (${getUnlearnedWords().length})`} />
              <Tab label={`Learned (${getLearnedWords().length})`} />
            </Tabs>
            
            <Tooltip title={`Sort by ${sortBy}`}>
              <IconButton onClick={handleSortChange}>
                <SortIcon />
                <Typography variant="caption" sx={{ ml: 0.5 }}>
                  {getSortIcon()}
                </Typography>
              </IconButton>
            </Tooltip>
          </Box>
          
          <TextField
            fullWidth
            variant="outlined"
            placeholder="Search vocabulary..."
            value={searchTerm}
            onChange={handleSearchChange}
            InputProps={{
              startAdornment: (
                <InputAdornment position="start">
                  <SearchIcon />
                </InputAdornment>
              ),
            }}
            sx={{ mb: 2 }}
          />
          
          <Divider sx={{ mb: 2 }} />
          
          {displayWords.length === 0 ? (
            <Typography variant="body1" color="text.secondary" align="center">
              No vocabulary words found.
            </Typography>
          ) : (
            <List>
              {displayWords.map((word) => (
                <ListItem
                  key={word.japanese}
                  sx={{
                    borderLeft: word.learned ? '4px solid #4caf50' : '4px solid #ff9800',
                    mb: 1,
                    borderRadius: 1,
                    backgroundColor: 'rgba(0, 0, 0, 0.02)',
                  }}
                >
                  <ListItemText
                    primary={
                      <Box sx={{ display: 'flex', alignItems: 'center' }}>
                        <Typography variant="body1" sx={{ fontWeight: 'bold', mr: 1 }}>
                          {word.japanese}
                        </Typography>
                        {word.kanji && (
                          <Chip 
                            label={word.kanji} 
                            size="small" 
                            color="primary" 
                            variant="outlined"
                            sx={{ mr: 1 }}
                          />
                        )}
                      </Box>
                    }
                    secondary={
                      <Box>
                        <Typography variant="body2" color="text.secondary">
                          {word.english}
                        </Typography>
                        {word.example && (
                          <Typography variant="caption" color="text.secondary" sx={{ display: 'block', mt: 0.5 }}>
                            Example: {word.example}
                          </Typography>
                        )}
                        <Typography variant="caption" color="text.secondary" sx={{ display: 'block', mt: 0.5 }}>
                          Added: {new Date(word.dateAdded).toLocaleDateString()}
                        </Typography>
                      </Box>
                    }
                  />
                  <ListItemSecondaryAction>
                    {word.learned ? (
                      <Tooltip title="Mark as not learned">
                        <IconButton 
                          edge="end" 
                          onClick={() => markAsNotLearned(word.japanese)}
                          color="warning"
                        >
                          <CancelIcon />
                        </IconButton>
                      </Tooltip>
                    ) : (
                      <Tooltip title="Mark as learned">
                        <IconButton 
                          edge="end" 
                          onClick={() => markAsLearned(word.japanese)}
                          color="success"
                        >
                          <CheckCircleIcon />
                        </IconButton>
                      </Tooltip>
                    )}
                    <Tooltip title="Remove word">
                      <IconButton 
                        edge="end" 
                        onClick={() => removeWord(word.japanese)}
                        color="error"
                      >
                        <DeleteIcon />
                      </IconButton>
                    </Tooltip>
                  </ListItemSecondaryAction>
                </ListItem>
              ))}
            </List>
          )}
        </>
      )}
    </Paper>
  );
};

export default VocabularyList; 