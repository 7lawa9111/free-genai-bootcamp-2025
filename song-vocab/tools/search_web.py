from duckduckgo_search import DDGS
from typing import List, Dict
import time

def search_web(query: str, max_results: int = 5, retries: int = 3) -> List[Dict]:
    """
    Search the web for lyrics using DuckDuckGo
    
    Args:
        query (str): Search query (e.g., "lyrics Shape of You Ed Sheeran")
        max_results (int): Maximum number of results to return
        retries (int): Number of retries if no results are found
        
    Returns:
        List[Dict]: List of search results with title, link, and snippet
    """
    # Return empty list for empty or whitespace-only queries
    if not query or not query.strip():
        return []

    results = []
    attempt = 0
    
    while attempt < retries and not results:
        try:
            with DDGS() as ddgs:
                # Add "lyrics" to the query if not present
                if "lyrics" not in query.lower():
                    query = f"lyrics {query}"
                
                print(f"\nAttempt {attempt + 1}: Searching for '{query}'")
                
                # Use the text method with proper parameters
                search_results = ddgs.text(
                    query,
                    max_results=max_results * 2,  # Request more results than needed in case some are filtered
                    region='wt-wt',  # Worldwide results
                    safesearch='off'
                )
                
                # Convert generator to list and process results
                raw_results = list(search_results)
                print(f"Found {len(raw_results)} raw results")
                
                # Debug: print first raw result structure
                if raw_results:
                    print("\nFirst raw result structure:")
                    print(raw_results[0])
                
                for result in raw_results:
                    try:
                        # Get title and link (using href as link)
                        title = str(result.get('title', '')).strip()
                        link = str(result.get('href', '')).strip()  # Use href instead of link
                        snippet = str(result.get('body', '')).strip()
                        
                        if title and link:  # Only include results with both title and link
                            results.append({
                                'title': title,
                                'link': link,
                                'snippet': snippet
                            })
                            print(f"\nValid result found:")
                            print(f"Title: {title}")
                            print(f"Link: {link}")
                    except Exception as e:
                        print(f"Error processing result: {str(e)}")
                        print(f"Raw result: {result}")
                        continue
                
                print(f"Processed {len(results)} valid results")
                        
                if results:
                    break
                    
                # Add an increasing delay between retries
                time.sleep(1 * (attempt + 1))
                
        except Exception as e:
            print(f"Error on attempt {attempt + 1}: {str(e)}")
            time.sleep(1 * (attempt + 1))
            
        attempt += 1
    
    return results[:max_results]  # Ensure we don't exceed max_results 