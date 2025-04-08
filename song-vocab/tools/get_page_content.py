import requests
from bs4 import BeautifulSoup
from typing import Optional
import re

def get_page_content(url: str) -> Optional[str]:
    """
    Fetch and parse the content of a webpage, attempting to extract lyrics
    
    Args:
        url (str): URL of the webpage to fetch
        
    Returns:
        Optional[str]: Extracted lyrics text or None if failed
    """
    try:
        headers = {
            'User-Agent': 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36'
        }
        response = requests.get(url, headers=headers, timeout=10)
        response.raise_for_status()
        
        soup = BeautifulSoup(response.text, 'html.parser')
        
        # Remove script and style elements
        for script in soup(["script", "style", "header", "footer", "nav"]):
            script.decompose()
        
        # Try to find lyrics content
        # First, look for common lyrics container classes/IDs
        lyrics_containers = soup.find_all(class_=re.compile(r'lyrics?|songtekst|letra|paroles', re.I))
        if not lyrics_containers:
            lyrics_containers = soup.find_all(id=re.compile(r'lyrics?|songtekst|letra|paroles', re.I))
        
        if lyrics_containers:
            # Use the largest lyrics container
            lyrics_text = max(lyrics_containers, key=lambda x: len(str(x))).get_text(separator='\n', strip=True)
        else:
            # Fallback: get all text content and try to identify lyrics section
            text = soup.get_text(separator='\n', strip=True)
            lines = text.split('\n')
            
            # Look for a section that looks like lyrics (multiple short lines)
            lyrics_lines = []
            lyrics_section = False
            line_count = 0
            short_line_count = 0
            
            for line in lines:
                line = line.strip()
                if not line:
                    continue
                    
                # Count consecutive short lines (typical for lyrics)
                if len(line) < 50:
                    short_line_count += 1
                else:
                    short_line_count = 0
                
                # If we find 4+ consecutive short lines, it's probably lyrics
                if short_line_count >= 4:
                    lyrics_section = True
                    
                if lyrics_section:
                    lyrics_lines.append(line)
                    line_count += 1
                    
                    # End lyrics section if we hit a long paragraph
                    if len(line) > 100:
                        break
                        
                    # End lyrics section if we have enough lines
                    if line_count > 50:
                        break
            
            lyrics_text = '\n'.join(lyrics_lines) if lyrics_lines else text
        
        # Clean up the lyrics
        # Remove common ads/annotations
        lyrics_text = re.sub(r'\[.*?\]', '', lyrics_text)  # Remove [Verse], [Chorus], etc.
        lyrics_text = re.sub(r'\(.*?\)', '', lyrics_text)  # Remove (x2), etc.
        lyrics_text = re.sub(r'^\d+×$', '', lyrics_text, flags=re.MULTILINE)  # Remove line numbers
        lyrics_text = re.sub(r'(\n\s*)+', '\n', lyrics_text)  # Remove multiple newlines
        
        return lyrics_text.strip()
        
    except Exception as e:
        print(f"Error fetching page content: {str(e)}")
        return None 