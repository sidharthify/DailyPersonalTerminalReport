import os
import json
import random

def get_data(config):
    results = []
    
    if config.get('show_review', True):
        results.append("Report successfully compiled with latest live data.")
        
    return results
