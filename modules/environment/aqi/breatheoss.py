import requests

def get_data(config):
    # Docs: https://api.breatheoss.app/docs
    city = config.get('city', 'srinagar').lower().replace(' ', '_')
    results = []
    
    try:
        url = f"https://api.breatheoss.app/aqi/{city}"
        resp = requests.get(url, timeout=10)
        if resp.status_code != 200:
            results.append(f"BreatheOSS Error: City '{city}' not found or API error.")
            return results
            
        data = resp.json()
        us_aqi = data.get('us_aqi', '??')
        main_pollutant = data.get('us_main_pollutant', '??')
        
        breakdown = data.get('aqi_breakdown', {})
        pm25 = breakdown.get('pm2_5', '??')
        pm10 = breakdown.get('pm10', '??')
        
        status = "Unknown"
        if isinstance(us_aqi, (int, float)):
            if us_aqi <= 50: status = "Good"
            elif us_aqi <= 100: status = "Moderate"
            elif us_aqi <= 150: status = "Unhealthy (Sensitive Groups)"
            elif us_aqi <= 200: status = "Unhealthy"
            elif us_aqi <= 300: status = "Very Unhealthy"
            else: status = "Hazardous"
            
        results.append(f"US AQI ({city.title()}): {us_aqi} ({status})")
        pm25_s = str(pm25).strip()
        pm10_s = str(pm10).strip()
        results.append(f"PM2.5: {pm25_s} | PM10: {pm10_s} | Main: {main_pollutant.upper()}")
        
        source = data.get('source', 'BreatheOSS')
        results.append(f"  └─> Source: {source}")
        
    except Exception as e:
        results.append(f"BreatheOSS Error: {e}")
        
    return results
