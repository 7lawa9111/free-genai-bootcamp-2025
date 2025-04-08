import unittest
import sys
import os

# Add the parent directory to the Python path
sys.path.insert(0, os.path.abspath(os.path.join(os.path.dirname(__file__), '..')))

# Import only search_web directly
from tools.search_web import search_web

class TestSearchWeb(unittest.TestCase):
    def setUp(self):
        """Setup test cases"""
        self.test_query = "YOASOBI 群青 lyrics"
    
    def test_basic_search(self):
        """Test basic search functionality"""
        results = search_web(self.test_query)
        
        # Check if we got results
        self.assertIsInstance(results, list)
        self.assertGreater(len(results), 0)
        
        # Check result structure and content
        first_result = results[0]
        self.assertIn('title', first_result)
        self.assertIn('link', first_result)
        self.assertIn('snippet', first_result)
        
        # Verify non-empty values
        self.assertTrue(first_result['title'].strip())
        self.assertTrue(first_result['link'].strip())
        
        # Print first result for manual verification
        print("\nFirst search result:")
        print(f"Title: {first_result['title']}")
        print(f"Link: {first_result['link']}")
        print(f"Snippet: {first_result['snippet'][:100]}...")
        
    def test_empty_query(self):
        """Test handling of empty query"""
        # Test completely empty query
        self.assertEqual(search_web(""), [])
        # Test whitespace-only query
        self.assertEqual(search_web("   "), [])
        # Test None query (should handle gracefully)
        self.assertEqual(search_web(None), [])
        
    def test_max_results(self):
        """Test max_results parameter"""
        max_results = 3
        results = search_web(self.test_query, max_results=max_results)
        self.assertLessEqual(len(results), max_results)
        # Verify each result has required fields
        for result in results:
            self.assertTrue(result['title'].strip())
            self.assertTrue(result['link'].strip())

if __name__ == '__main__':
    unittest.main(verbose=True) 