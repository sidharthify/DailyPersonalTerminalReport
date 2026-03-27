import requests

def get_data(config):
    count = config.get('count', 5)
    results = []
    try:
        top_ids_url = "https://hacker-news.firebaseio.com/v0/topstories.json"
        top_ids = requests.get(top_ids_url, timeout=10).json()[:count]
        
        for item_id in top_ids:
            item_url = f"https://hacker-news.firebaseio.com/v0/item/{item_id}.json"
            item = requests.get(item_url, timeout=5).json()
            if item:
                title = item.get('title')
                score = item.get('score', 0)
                results.append(f"[{score}] {title}")
    except Exception as e:
        results.append(f"Hacker News Error: {e}")
    return results
