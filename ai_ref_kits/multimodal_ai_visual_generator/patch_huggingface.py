"""
This script patches the huggingface_hub module to add the cached_download function.
Run this script before running the main application.
"""

import huggingface_hub

# Check if cached_download is already defined
if not hasattr(huggingface_hub, 'cached_download'):
    # Add cached_download function to huggingface_hub
    def cached_download(*args, **kwargs):
        return huggingface_hub.hf_hub_download(*args, **kwargs)
    huggingface_hub.cached_download = cached_download
    print("Successfully patched huggingface_hub with cached_download function")
else:
    print("huggingface_hub already has cached_download function") 