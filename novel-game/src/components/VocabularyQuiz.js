import React, { useState, useEffect } from 'react';
import {
  Box,
  Typography,
  Button,
  Radio,
  RadioGroup,
  FormControlLabel,
  FormControl,
  Paper,
  LinearProgress,
  IconButton,
  Tooltip,
  Divider,
} from '@mui/material';
import {
  Refresh as RefreshIcon,
  Check as CheckIcon,
  Close as CloseIcon,
} from '@mui/icons-material';
import useVocabularyStore from '../store/vocabularyStore';

const VocabularyQuiz = () => {
  const [quizWords, setQuizWords] = useState([]);
  const [currentQuestionIndex, setCurrentQuestionIndex] = useState(0);
  const [selectedAnswer, setSelectedAnswer] = useState('');
  const [score, setScore] = useState(0);
  const [showResults, setShowResults] = useState(false);
  const [quizType, setQuizType] = useState('japanese-to-english'); // 'japanese-to-english' or 'english-to-japanese'
  const [quizMode, setQuizMode] = useState('all'); // 'all', 'learned', 'unlearned'
  
  const { getAllWords, getLearnedWords, getUnlearnedWords, markAsLearned } = useVocabularyStore();

  // Initialize quiz
  const initializeQuiz = () => {
    let words = [];
    
    // Get words based on quiz mode
    if (quizMode === 'all') {
      words = getAllWords();
    } else if (quizMode === 'learned') {
      words = getLearnedWords();
    } else if (quizMode === 'unlearned') {
      words = getUnlearnedWords();
    }
    
    // Shuffle words
    const shuffled = [...words].sort(() => Math.random() - 0.5);
    
    // Take first 10 words or all if less than 10
    const selectedWords = shuffled.slice(0, Math.min(10, shuffled.length));
    
    setQuizWords(selectedWords);
    setCurrentQuestionIndex(0);
    setSelectedAnswer('');
    setScore(0);
    setShowResults(false);
  };

  useEffect(() => {
    initializeQuiz();
  }, [quizMode]);

  const handleAnswerChange = (event) => {
    setSelectedAnswer(event.target.value);
  };

  const handleSubmit = () => {
    const currentWord = quizWords[currentQuestionIndex];
    const isCorrect = selectedAnswer === currentWord.english;
    
    if (isCorrect) {
      setScore(score + 1);
      
      // If the word wasn't marked as learned, mark it now
      if (!currentWord.learned) {
        markAsLearned(currentWord.japanese);
      }
    }
    
    // Move to next question or show results
    if (currentQuestionIndex < quizWords.length - 1) {
      setCurrentQuestionIndex(currentQuestionIndex + 1);
      setSelectedAnswer('');
    } else {
      setShowResults(true);
    }
  };

  const handleRestart = () => {
    initializeQuiz();
  };

  const handleQuizTypeChange = () => {
    setQuizType(quizType === 'japanese-to-english' ? 'english-to-japanese' : 'japanese-to-english');
    setSelectedAnswer('');
  };

  const handleQuizModeChange = (mode) => {
    setQuizMode(mode);
  };

  if (quizWords.length === 0) {
    return (
      <Paper elevation={3} sx={{ p: 3, textAlign: 'center' }}>
        <Typography variant="h6" gutterBottom>
          No vocabulary words available for quiz
        </Typography>
        <Typography variant="body2" color="text.secondary">
          Add some words to your vocabulary first!
        </Typography>
      </Paper>
    );
  }

  if (showResults) {
    return (
      <Paper elevation={3} sx={{ p: 3 }}>
        <Typography variant="h5" gutterBottom align="center">
          Quiz Results
        </Typography>
        <Box sx={{ display: 'flex', justifyContent: 'center', mb: 3 }}>
          <Typography variant="h3" color="primary">
            {score}/{quizWords.length}
          </Typography>
        </Box>
        <Box sx={{ mb: 3 }}>
          <LinearProgress 
            variant="determinate" 
            value={(score / quizWords.length) * 100} 
            sx={{ height: 10, borderRadius: 5 }}
          />
        </Box>
        <Typography variant="body1" align="center" gutterBottom>
          {score === quizWords.length 
            ? "Perfect score! Great job!" 
            : score > quizWords.length / 2 
              ? "Good job! Keep practicing!" 
              : "Keep studying and try again!"}
        </Typography>
        <Box sx={{ display: 'flex', justifyContent: 'center', mt: 3 }}>
          <Button 
            variant="contained" 
            color="primary" 
            onClick={handleRestart}
            startIcon={<RefreshIcon />}
          >
            Try Again
          </Button>
        </Box>
      </Paper>
    );
  }

  const currentWord = quizWords[currentQuestionIndex];
  const progress = ((currentQuestionIndex + 1) / quizWords.length) * 100;

  return (
    <Paper elevation={3} sx={{ p: 3 }}>
      <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 2 }}>
        <Typography variant="h6">
          Vocabulary Quiz
        </Typography>
        <Box>
          <Tooltip title="Change quiz type">
            <IconButton onClick={handleQuizTypeChange} size="small">
              <RefreshIcon />
            </IconButton>
          </Tooltip>
        </Box>
      </Box>
      
      <Box sx={{ mb: 2 }}>
        <LinearProgress 
          variant="determinate" 
          value={progress} 
          sx={{ height: 8, borderRadius: 4 }}
        />
        <Typography variant="caption" color="text.secondary" sx={{ mt: 0.5, display: 'block', textAlign: 'right' }}>
          Question {currentQuestionIndex + 1} of {quizWords.length}
        </Typography>
      </Box>
      
      <Box sx={{ mb: 3 }}>
        <Typography variant="subtitle1" gutterBottom>
          Quiz Mode:
        </Typography>
        <Box sx={{ display: 'flex', gap: 1 }}>
          <Button 
            variant={quizMode === 'all' ? 'contained' : 'outlined'} 
            size="small"
            onClick={() => handleQuizModeChange('all')}
          >
            All Words
          </Button>
          <Button 
            variant={quizMode === 'learned' ? 'contained' : 'outlined'} 
            size="small"
            onClick={() => handleQuizModeChange('learned')}
          >
            Learned
          </Button>
          <Button 
            variant={quizMode === 'unlearned' ? 'contained' : 'outlined'} 
            size="small"
            onClick={() => handleQuizModeChange('unlearned')}
          >
            To Learn
          </Button>
        </Box>
      </Box>
      
      <Divider sx={{ mb: 3 }} />
      
      <Box sx={{ mb: 3 }}>
        <Typography variant="h5" gutterBottom align="center">
          {quizType === 'japanese-to-english' 
            ? currentWord.japanese 
            : currentWord.english}
        </Typography>
        
        {currentWord.kanji && (
          <Typography variant="subtitle1" align="center" color="text.secondary" gutterBottom>
            Kanji: {currentWord.kanji}
          </Typography>
        )}
      </Box>
      
      <FormControl component="fieldset" sx={{ width: '100%' }}>
        <Typography variant="subtitle1" gutterBottom>
          {quizType === 'japanese-to-english' 
            ? "What is the English translation?" 
            : "What is the Japanese word?"}
        </Typography>
        
        <RadioGroup value={selectedAnswer} onChange={handleAnswerChange}>
          {quizWords.map((word, index) => (
            <FormControlLabel
              key={index}
              value={quizType === 'japanese-to-english' ? word.english : word.japanese}
              control={<Radio />}
              label={quizType === 'japanese-to-english' ? word.english : word.japanese}
              sx={{ 
                mb: 1,
                p: 1,
                borderRadius: 1,
                '&:hover': {
                  backgroundColor: 'rgba(0, 0, 0, 0.04)',
                },
              }}
            />
          ))}
        </RadioGroup>
      </FormControl>
      
      <Box sx={{ display: 'flex', justifyContent: 'center', mt: 3 }}>
        <Button
          variant="contained"
          color="primary"
          onClick={handleSubmit}
          disabled={!selectedAnswer}
          startIcon={<CheckIcon />}
        >
          Submit Answer
        </Button>
      </Box>
    </Paper>
  );
};

export default VocabularyQuiz; 