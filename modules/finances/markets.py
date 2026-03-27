import requests
try:
    from .currency_util import get_exchange_rate
except ImportError:
    def get_exchange_rate(f, t): return 1.0

def get_data(config):
    symbols = config.get('symbols', ['AAPL', 'MSFT'])
    master_curr = config.get('master_currency')
    results = []
    
    headers = {
        'User-Agent': 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36'
    }
    
    for symbol in symbols:
        try:
            url = f"https://query1.finance.yahoo.com/v8/finance/chart/{symbol}?interval=1m&range=1d"
            resp = requests.get(url, headers=headers, timeout=5).json()
            
            result = resp.get('chart', {}).get('result', [{}])[0]
            meta = result.get('meta', {})
            
            price = meta.get('regularMarketPrice')
            prev_close = (
                meta.get('regularMarketPreviousClose') or 
                meta.get('previousClose') or 
                meta.get('chartPreviousClose')
            )
            native_curr = meta.get('currency', 'USD')
            
            # Currency Conversion
            display_curr = native_curr
            if master_curr and master_curr.upper() != native_curr.upper():
                rate = get_exchange_rate(native_curr, master_curr)
                if price is not None: price *= rate
                if prev_close is not None: prev_close *= rate
                display_curr = master_curr.upper()

            if price is not None and prev_close is not None:
                change = price - prev_close
                pct_change = (change / prev_close) * 100
                sign = "+" if change >= 0 else ""
                results.append(f"{symbol}: {price:,.2f} {display_curr} ({sign}{change:,.2f} | {sign}{pct_change:.2f}%)")
            elif price is not None:
                results.append(f"{symbol}: {price:,.2f} {display_curr}")
            else:
                results.append(f"{symbol}: Data unavailable")
                
        except Exception as e:
            results.append(f"{symbol}: Error ({e})")
            
    return results
