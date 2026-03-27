import requests

def get_data(config):
    
    count = config.get('count', 3)
    results = []
    
    try:
        url = f"https://www.reddit.com/r/showerthoughts/top.json?limit={count}&t=day"
        headers = {
            'User-Agent': 'DPPR/1.0 (Personal Daily Report)'
        }
        resp = requests.get(url, headers=headers, timeout=10).json()
        
        posts = resp.get('data', {}).get('children', [])
        for post in posts:
            title = post.get('data', {}).get('title')
            if title:
                results.append(f"• {title}")
                
        if not results:
            results.append("No shower thoughts found today.")
            
    except Exception as e:
        results.append(f"Reddit Error: {e}")
        
    return results
