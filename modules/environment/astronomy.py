import requests

def get_data(config):
    lat = config.get('lat')
    lon = config.get('lon')
    results = []
    try:
        url = f"https://api.open-meteo.com/v1/forecast?latitude={lat}&longitude={lon}&daily=sunrise,sunset&timezone=auto"
        resp = requests.get(url, timeout=10).json()
        
        daily = resp.get('daily', {})
        sunrise = daily.get('sunrise', ['N/A'])[0]
        sunset = daily.get('sunset', ['N/A'])[0]
        
        if 'T' in sunrise: sunrise = sunrise.split('T')[1]
        if 'T' in sunset: sunset = sunset.split('T')[1]
        
        results.append(f"Sunrise: {sunrise}")
        results.append(f"Sunset:  {sunset}")
    except Exception as e:
        results.append(f"Astronomy Error: {e}")
    return results
