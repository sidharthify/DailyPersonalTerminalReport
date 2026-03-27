import os
import requests

def get_data(config):
    lat = config.get('lat') or os.getenv('LAT')
    lon = config.get('lon') or os.getenv('LON')
    results = []

    if not lat or not lon:
        results.append("Weather: Missing coordinates (LAT/LON)")
        return results

    try:
        weather_url = f"https://api.open-meteo.com/v1/forecast?latitude={lat}&longitude={lon}&current=temperature_2m,relative_humidity_2m,weather_code&daily=temperature_2m_max,temperature_2m_min&timezone=auto"
        w_res = requests.get(weather_url, timeout=5).json()
        current = w_res.get('current', {})
        daily = w_res.get('daily', {})
        
        temp = current.get('temperature_2m', '??')
        hum = current.get('relative_humidity_2m', '??')
        max_t = daily.get('temperature_2m_max', ['??'])[0]
        min_t = daily.get('temperature_2m_min', ['??'])[0]
        
        results.append(f"Temperature: {temp}°C (Min: {min_t}°, Max: {max_t}°)")
        results.append(f"Humidity: {hum}%")
    except Exception as e:
        results.append(f"Weather Error: {e}")

    return results
