import React, { useState, useRef, useEffect } from 'react';
import { useGameStateStore } from '../store/gameStateStore';
import { handleCommand } from '../utils/commandHandler';
import '../styles/CommandInput.css';

const CommandInput = () => {
  const [input, setInput] = useState('');
  const [commandHistory, setCommandHistory] = useState([]);
  const inputRef = useRef(null);
  const { currentScene } = useGameStateStore();

  useEffect(() => {
    // Focus the input when the component mounts
    inputRef.current?.focus();
  }, []);

  const handleSubmit = (e) => {
    e.preventDefault();
    if (!input.trim()) return;

    // Process the command using our command handler
    const response = handleCommand(input.toLowerCase(), currentScene);

    // Add the command and its response to history
    setCommandHistory(prev => [...prev, {
      command: input,
      response: response.text,
      translation: response.translation
    }]);

    // Clear the input
    setInput('');
  };

  return (
    <div className="command-input-container">
      <div className="command-history">
        {commandHistory.map((entry, index) => (
          <div key={index} className="command-entry">
            <div className="command">{'>'} {entry.command}</div>
            <div className="response">{entry.response}</div>
            <div className="translation">{entry.translation}</div>
          </div>
        ))}
      </div>
      <form onSubmit={handleSubmit} className="command-form">
        <input
          ref={inputRef}
          type="text"
          value={input}
          onChange={(e) => setInput(e.target.value)}
          placeholder="Type a command (help for available commands)"
          className="command-input"
        />
        <button type="submit" className="command-submit">Send</button>
      </form>
    </div>
  );
};

export default CommandInput; 