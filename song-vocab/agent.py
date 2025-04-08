import ollama
import re
import os
from typing import Dict, List, Optional, Tuple
from tools.search_web import search_web
from tools.get_page_content import get_page_content
from tools.extract_vocabulary import extract_vocabulary
from tools.save_result import save_result

class LyricsAgent:
    def __init__(self, output_dir: str = "outputs"):
        self.client = ollama.Client()
        self.system_prompt = self._load_prompt("lyrics_agent.md")
        self.output_dir = output_dir
        os.makedirs(output_dir, exist_ok=True)

    def _load_prompt(self, prompt_name: str) -> str:
        """Load a prompt from the prompts directory"""
        prompt_path = os.path.join(os.path.dirname(__file__), 'prompts', prompt_name)
        try:
            with open(prompt_path, 'r') as f:
                return f.read().strip()
        except Exception as e:
            raise Exception(f"Failed to load prompt from {prompt_path}: {str(e)}")

    def run(self, query: str) -> str:
        """
        Run the agent to find lyrics and vocabulary.
        Returns a unique identifier for the saved files.
        """
        messages = [
            {
                "role": "system",
                "content": self.system_prompt
            },
            {
                "role": "user",
                "content": f"Find lyrics and vocabulary for: {query}"
            }
        ]
        
        max_steps = 10
        step = 0
        
        while step < max_steps:
            response = self.client.chat(
                model="deepseek-r1:latest",
                messages=messages
            )
            
            current_message = response['message']['content']
            messages.append({"role": "assistant", "content": current_message})
            
            if "Final Answer:" in current_message:
                # Use the save_result tool to save the final answer
                result_id = save_result(query, current_message, self.output_dir)
                if result_id:
                    return result_id
                else:
                    # If parsing failed, ask the model to retry
                    messages.append({
                        "role": "user",
                        "content": "The answer format was incorrect. Please provide the answer with clear LYRICS and VOCABULARY sections marked with START/END tags."
                    })
                    continue
            
            # Parse the next action
            action, action_input = self._parse_action(current_message)
            if action:
                try:
                    # Execute the tool
                    if action == "search_web":
                        result = search_web(action_input)
                    elif action == "get_page_content":
                        result = get_page_content(action_input)
                        if not result:
                            result = "Error: Could not extract lyrics from the page"
                        elif len(result.split('\n')) < 3:
                            result = "Error: Extracted content is too short to be lyrics"
                    elif action == "extract_vocabulary":
                        result = extract_vocabulary(action_input)
                        if not result:
                            result = "Error: Could not extract vocabulary items"
                    else:
                        result = "Unknown action"
                except Exception as e:
                    result = f"Error executing {action}: {str(e)}"
                
                # Add the observation to messages
                messages.append({
                    "role": "user",
                    "content": f"Observation: {str(result)}"
                })
            
            step += 1
        
        raise Exception("Max steps reached without finding lyrics")

    def _parse_action(self, message: str) -> Tuple[Optional[str], Optional[str]]:
        """Parse the action and action input from the message"""
        if "Action:" not in message or "Action Input:" not in message:
            return None, None
        
        action_start = message.find("Action:") + 7
        action_end = message.find("Action Input:")
        action = message[action_start:action_end].strip()
        
        input_start = message.find("Action Input:") + 13
        input_end = message.find("Observation:") if "Observation:" in message else len(message)
        action_input = message[input_start:input_end].strip()
        
        return action, action_input 