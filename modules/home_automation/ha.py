import os
import requests

def get_data(config):
    token = os.getenv('HA_ACCESS_TOKEN')
    ha_url = config.get('url') or os.getenv('HA_URL')
    
    if not token or not ha_url:
        return ["Home Assistant: Credentials or URL missing."]
    
    ha_url = ha_url.rstrip('/')
    if not ha_url.endswith('/api'):
        ha_url = f"{ha_url}/api"
        
    headers = {
        "Authorization": f"Bearer {token}",
        "Content-Type": "application/json",
    }
    
    entities = config.get('entities', [])
    if not entities:
        return ["Home Assistant: No entities configured in layout."]

    output = []
    try:
        resp = requests.get(f"{ha_url}/states", headers=headers, timeout=10)
        
        if resp.status_code != 200:
            return [f"HA Error: Status {resp.status_code}"]
            
        states = resp.json()
        states_dict = {s['entity_id']: s for s in states}
        
        for ent in entities:
            eid = ent.get('entity_id') or ent.get('id') # Support both
            label = ent.get('label', eid)
            
            state = states_dict.get(eid)
            if state:
                val = state.get('state', 'Unknown')
                unit = state.get('attributes', {}).get('unit_of_measurement', '')
                output.append(f"{label}: {val} {unit}".strip())
            else:
                output.append(f"{label}: Not found.")
                
    except Exception as e:
        output.append(f"HA Error: {e}")
        
    return output
