import os
from google_play_scraper import app as gplay_app

def get_data(config):
    package_name = config.get('package_name') or os.getenv('PACKAGE_NAME')
    results = []
    
    if not package_name:
        results.append("Google Play: Package name not configured.")
        return results

    try:
        play_data = gplay_app(package_name, lang='en', country='in')
        downloads = play_data.get('installs', 'Unknown')
        name = play_data.get('title', 'App')
        results.append(f"{name} Downloads: {downloads}")
    except:
        results.append(f"Google Play ({package_name}): Data unavailable.")
        
    return results
