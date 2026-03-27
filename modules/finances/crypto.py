import requests

def get_data(config):
    tokens = config.get('tokens', ['bitcoin', 'ethereum', 'solana'])
    master_curr = (config.get('master_currency') or 'USD').lower()
    results = []
    
    try:
        ids = ",".join(tokens)
        url = f"https://api.coingecko.com/api/v3/simple/price?ids={ids}&vs_currencies={master_curr}&include_24hr_change=true"
        resp = requests.get(url, timeout=10).json()
        
        for token in tokens:
            data = resp.get(token)
            if data:
                price = data.get(master_curr)
                change_pct = data.get(f"{master_curr}_24h_change", 0)
                
                sign = "+" if change_pct >= 0 else ""
                if price is not None:
                    results.append(f"{token.title()}: {price:,.2f} {master_curr.upper()} ({sign}{change_pct:.2f}%)")
                else:
                    results.append(f"{token.title()}: Price unavailable")
            else:
                results.append(f"{token.title()}: Not found")
                
    except Exception as e:
        results.append(f"Crypto Error: {e}")
        
    return results
