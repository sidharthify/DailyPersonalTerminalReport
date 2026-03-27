import requests

def get_data(config):
    base_curr = config.get('base_currency', 'USD').upper()
    target_curr = config.get('target_currency', 'INR').upper()
    
    try:
        url = config.get('exchange_rate_url')
        if not url:
            url = f"https://api.exchangerate-api.com/v4/latest/{base_curr}"
            
        res = requests.get(url, timeout=5).json()
        rates = res.get('rates', {})
        rate = rates.get(target_curr)
        
        if rate:
            return [f"{base_curr}/{target_curr}: {rate}"]
        else:
            return [f"{base_curr}/{target_curr}: N/A"]
    except Exception as e:
        return [f"{base_curr}/{target_curr}: Error ({e})"]
