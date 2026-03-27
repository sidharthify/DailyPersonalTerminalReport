import requests
from datetime import datetime, timedelta

def get_data(config):
    count = config.get('count', 5)
    lang = config.get('language', '')
    results = []
    try:
        last_week = (datetime.now() - timedelta(days=7)).strftime('%Y-%m-%d')
        query = f"created:>{last_week}"
        if lang:
            query += f" language:{lang}"
            
        url = f"https://api.github.com/search/repositories?q={query}&sort=stars&order=desc"
        headers = {'Accept': 'application/vnd.github.v3+json'}
        resp = requests.get(url, headers=headers, timeout=10).json()
        
        items = resp.get('items', [])[:count]
        for item in items:
            name = item.get('full_name')
            stars = item.get('stargazers_count')
            desc = item.get('description', '')
            if desc and len(desc) > 60:
                desc = desc[:57] + "..."
            results.append(f"★{stars} {name}")
            if desc:
                results.append(f"  {desc}")
    except Exception as e:
        results.append(f"GitHub Error: {e}")
    return results
