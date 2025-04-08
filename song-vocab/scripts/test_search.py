import sys
import os

# Add the parent directory to the Python path
sys.path.insert(0, os.path.abspath(os.path.join(os.path.dirname(__file__), '..')))

from tools.search_web import search_web

def test_search():
    # Test search with a Japanese song
    query = "YOASOBI 群青 lyrics"
    print(f"\nSearching for: {query}")
    results = search_web(query, max_results=3)
    
    # Print results in a readable format
    for i, result in enumerate(results, 1):
        print(f"\nResult {i}:")
        print(f"Title: {result['title']}")
        print(f"Link: {result['link']}")
        print(f"Snippet: {result['snippet'][:200]}...")

if __name__ == "__main__":
    test_search() 