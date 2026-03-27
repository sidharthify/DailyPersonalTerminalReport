import os
import requests

def get_data(config):
    lat = config.get('lat') or os.getenv('LAT')
    lon = config.get('lon') or os.getenv('LON')
    results = []

    if not lat or not lon:
        results.append("AQI: Missing coordinates (LAT/LON)")
        return results

    try:
        aqi_url = f"https://air-quality-api.open-meteo.com/v1/air-quality?latitude={lat}&longitude={lon}&current=us_aqi,european_aqi,pm10,pm2_5"
        aqi_res = requests.get(aqi_url, timeout=5).json()
        current_aqi = aqi_res.get('current', {})
        
        us_aqi = current_aqi.get('us_aqi', '??')
        pm25 = current_aqi.get('pm2_5', '??')
        pm10 = current_aqi.get('pm10', '??')
        
        status = "Unknown"
        if isinstance(us_aqi, (int, float)):
            if us_aqi <= 50: status = "Good"
            elif us_aqi <= 100: status = "Moderate"
            elif us_aqi <= 150: status = "Unhealthy (Sensitive Groups)"
            elif us_aqi <= 200: status = "Unhealthy"
            elif us_aqi <= 300: status = "Very Unhealthy"
            else: status = "Hazardous"
            
        results.append(f"US AQI: {us_aqi} ({status}) | PM2.5: {pm25} | PM10: {pm10}")
    except Exception as e:
        results.append(f"AQI Error: {e}")

    return results
