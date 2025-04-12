import React, { useState } from 'react';
import {
  Box,
  Typography,
  Button,
  TextField,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  IconButton,
  Tooltip,
  Chip,
} from '@mui/material';
import {
  Add as AddIcon,
  Close as CloseIcon,
  Bookmark as BookmarkIcon,
} from '@mui/icons-material';
import useVocabularyStore from '../store/vocabularyStore';

const VocabularyExtractor = ({ text, translation }) => {
  const [open, setOpen] = useState(false);
  const [selectedText, setSelectedText] = useState('');
  const [wordData, setWordData] = useState({
    japanese: '',
    english: '',
    kanji: '',
    example: '',
  });
  const { addWord } = useVocabularyStore();

  const handleOpen = () => {
    setOpen(true);
  };

  const handleClose = () => {
    setOpen(false);
    setWordData({
      japanese: '',
      english: '',
      kanji: '',
      example: '',
    });
  };

  const handleTextSelection = () => {
    const selection = window.getSelection();
    if (selection.toString().trim()) {
      setSelectedText(selection.toString().trim());
      setWordData({
        ...wordData,
        japanese: selection.toString().trim(),
        example: text,
      });
    }
  };

  const handleInputChange = (e) => {
    const { name, value } = e.target;
    setWordData({
      ...wordData,
      [name]: value,
    });
  };

  const handleAddWord = () => {
    if (wordData.japanese && wordData.english) {
      addWord(wordData);
      handleClose();
    }
  };

  return (
    <>
      <Tooltip title="Add to vocabulary">
        <IconButton
          size="small"
          onClick={handleOpen}
          sx={{ position: 'absolute', top: 8, right: 8 }}
        >
          <BookmarkIcon fontSize="small" />
        </IconButton>
      </Tooltip>

      <Dialog open={open} onClose={handleClose} maxWidth="sm" fullWidth>
        <DialogTitle>
          Add to Vocabulary
          <IconButton
            aria-label="close"
            onClick={handleClose}
            sx={{
              position: 'absolute',
              right: 8,
              top: 8,
            }}
          >
            <CloseIcon />
          </IconButton>
        </DialogTitle>
        <DialogContent dividers>
          <Box sx={{ mb: 2 }}>
            <Typography variant="subtitle2" gutterBottom>
              Select text from the dialogue to add to your vocabulary
            </Typography>
            <Box
              sx={{
                p: 2,
                border: '1px dashed #ccc',
                borderRadius: 1,
                backgroundColor: '#f9f9f9',
                cursor: 'text',
              }}
              onMouseUp={handleTextSelection}
            >
              <Typography variant="body1">{text}</Typography>
              <Typography variant="body2" color="text.secondary">
                {translation}
              </Typography>
            </Box>
          </Box>

          {selectedText && (
            <Chip
              label={selectedText}
              onDelete={() => setSelectedText('')}
              color="primary"
              sx={{ mb: 2 }}
            />
          )}

          <TextField
            autoFocus
            margin="dense"
            name="japanese"
            label="Japanese"
            type="text"
            fullWidth
            variant="outlined"
            value={wordData.japanese}
            onChange={handleInputChange}
            sx={{ mb: 2 }}
          />
          <TextField
            margin="dense"
            name="english"
            label="English Translation"
            type="text"
            fullWidth
            variant="outlined"
            value={wordData.english}
            onChange={handleInputChange}
            sx={{ mb: 2 }}
          />
          <TextField
            margin="dense"
            name="kanji"
            label="Kanji (optional)"
            type="text"
            fullWidth
            variant="outlined"
            value={wordData.kanji}
            onChange={handleInputChange}
            sx={{ mb: 2 }}
          />
          <TextField
            margin="dense"
            name="example"
            label="Example Sentence"
            type="text"
            fullWidth
            variant="outlined"
            value={wordData.example}
            onChange={handleInputChange}
            multiline
            rows={2}
          />
        </DialogContent>
        <DialogActions>
          <Button onClick={handleClose}>Cancel</Button>
          <Button
            onClick={handleAddWord}
            variant="contained"
            color="primary"
            startIcon={<AddIcon />}
            disabled={!wordData.japanese || !wordData.english}
          >
            Add to Vocabulary
          </Button>
        </DialogActions>
      </Dialog>
    </>
  );
};

export default VocabularyExtractor; 