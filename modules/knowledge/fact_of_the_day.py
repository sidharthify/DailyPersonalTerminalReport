import requests

def get_data(config):
    
    results = []
    try:
        url = "https://uselessfacts.jsph.pl/api/v2/facts/today"
        headers = {
            'Accept': 'application/json'
        }
        resp = requests.get(url, headers=headers, timeout=5).json()
        
        fact = resp.get('text')
        if fact:
            results.append(fact)
        else:
            results.append("Could not fetch a fact for today.")
            
    except Exception as e:
        results.append(f"Fact Error: {e}")
        
    return results
