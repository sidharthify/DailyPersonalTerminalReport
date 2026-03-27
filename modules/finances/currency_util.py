import requests

def get_exchange_rate(from_curr, to_curr):
    if not from_curr or not to_curr:
        return 1.0
        
    from_curr = from_curr.upper()
    to_curr = to_curr.upper()
    
    if from_curr == to_curr:
        return 1.0
        
    try:
        symbol = f"{from_curr}{to_curr}=X"
        url = f"https://query1.finance.yahoo.com/v8/finance/chart/{symbol}?interval=1d&range=1d"
        headers = {
            'User-Agent': 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36'
        }
        resp = requests.get(url, headers=headers, timeout=5).json()
        
        result = resp.get('chart', {}).get('result', [{}])[0]
        price = result.get('meta', {}).get('regularMarketPrice')
        
        return price if price is not None else 1.0
    except Exception:
        return 1.0

def convert_price(price, from_curr, to_curr):
    if price is None or not isinstance(price, (int, float)):
        return price
        
    rate = get_exchange_rate(from_curr, to_curr)
    return price * rate
def get_currency_prefix(currency_code):
    if not currency_code:
        return ""
        
    code = currency_code.upper()
    mapping = {
        "INR": "Rs.",
        "USD": "Usd.",
        "EUR": "Eur.",
        "GBP": "Gbp.",
        "JPY": "Jpy.",
    }
    

    return mapping.get(code, f"{code.capitalize()}.")
