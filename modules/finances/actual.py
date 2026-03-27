import os
from actual import Actual
from actual.queries import get_accounts
try:
    from .currency_util import get_currency_prefix
except ImportError:
    def get_currency_prefix(code): return ""

def get_data(config):
    results = []
    
    ACTUAL_SERVER_URL = os.getenv('ACTUAL_SERVER_URL')
    ACTUAL_PASSWORD = os.getenv('ACTUAL_PASSWORD')
    ACTUAL_SYNC_ID = os.getenv('ACTUAL_SYNC_ID')
    ACTUAL_ENCRYPTION_PASSWORD = os.getenv('ACTUAL_ENCRYPTION_PASSWORD')

    if not ACTUAL_SERVER_URL or not ACTUAL_PASSWORD:
        results.append("Actual Budget: Missing credentials")
        return results

    main_match = config.get('main_account_match')

    try:
        with Actual(
            base_url=ACTUAL_SERVER_URL, 
            password=ACTUAL_PASSWORD, 
            encryption_password=ACTUAL_ENCRYPTION_PASSWORD if ACTUAL_ENCRYPTION_PASSWORD and ACTUAL_ENCRYPTION_PASSWORD != 'your_encryption_password' else None,
            file=ACTUAL_SYNC_ID
        ) as actual:
            accounts = get_accounts(actual.session)
            active_accounts = [acc for acc in accounts if not acc.closed and not acc.offbudget]
            
            main_bank = None
            others = []
            
            for acc in active_accounts:
                if main_match and main_match in acc.name:
                    main_bank = acc
                else:
                    others.append(acc)
            
            master_curr = config.get('master_currency', 'USD')
            prefix = get_currency_prefix(master_curr)
            
            if main_bank:
                results.append(f"{main_bank.name}: {prefix} {main_bank.balance:,.2f}")
            
            if not main_bank and not others:
                results.append("No active accounts found.")
            else:
                for acc in others:
                    results.append(f"{acc.name}: {prefix} {acc.balance:,.2f}")
                    
    except Exception as e:
        results.append(f"Actual Budget Error: {e}")
        
    return results
