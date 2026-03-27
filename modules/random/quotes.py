import os
import json
import random

def get_data(config):
    quotes_file = config.get('quotes_file', 'quotes.json')
    try:
        script_dir = os.path.dirname(os.path.dirname(os.path.dirname(os.path.abspath(__file__))))
        json_path = os.path.join(script_dir, quotes_file)
        
        if os.path.exists(json_path):
            with open(json_path, 'r') as f:
                data = json.load(f)
                
            if isinstance(data, list):
                quote = random.choice(data)
                text = quote.get('text', '')
                attr = quote.get('attribution') or quote.get('author', 'Unknown')
                if text:
                    return [f"{text} ~ {attr}"]

            elif isinstance(data, dict):
                for key in ['quotes', 'lyrics', 'lines']:
                    if key in data and isinstance(data[key], list):
                        quote = random.choice(data[key])
                        if isinstance(quote, dict):
                            text = quote.get('text', '')
                            attr = quote.get('attribution') or quote.get('author', 'Unknown')
                            return [f"{text} ~ {attr}"]
                        elif isinstance(quote, str):
                            return [f"{quote}"]

    except Exception as e:
        return [f"Quote error: {e}"]
    return []
