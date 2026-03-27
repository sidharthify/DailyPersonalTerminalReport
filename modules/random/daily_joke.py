import requests

def get_data(config):
    results = []
    try:
        headers = {"Accept": "application/json"}
        resp = requests.get("https://icanhazdadjoke.com/", headers=headers, timeout=5).json()
        joke = resp.get('joke', 'No joke found today.')
        
        # Simple text wrapping for PDF layout
        max_len = 50
        words = joke.split()
        current_line = ""
        for word in words:
            if len(current_line) + len(word) + 1 > max_len:
                results.append(current_line)
                current_line = word
            else:
                current_line = (current_line + " " + word).strip()
        if current_line:
            results.append(current_line)
    except Exception:
        results.append("Daily Joke: Service unavailable")
        
    return results
