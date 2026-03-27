import requests
import re
import html
from xml.etree import ElementTree

def get_data(config):
    results = _get_merriam_webster()
    
    if not results or "unavailable" in results[0]:
        results = _get_random_fallback()
        
    return results

def _get_merriam_webster():
    url = "https://www.merriam-webster.com/wotd/feed/rss2"
    headers = {'User-Agent': 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36'}
    
    try:
        resp = requests.get(url, headers=headers, timeout=10)
        if resp.status_code != 200:
            return ["Word of the Day: Feed unavailable (HTTP 403/404)"]
            
        root = ElementTree.fromstring(resp.content)
        item = root.find('.//item')
        if item is not None:
            raw_title = item.find('title').text or ""
            word = raw_title.split(':')[-1].strip()
            
            raw_desc = item.find('description').text or ""
            desc = html.unescape(raw_desc)
            clean_desc = re.sub('<[^<]+?>', '', desc)
            
            meaning = ""
            if "What It Means" in clean_desc:
                meaning = clean_desc.split("What It Means")[-1].split("Examples")[0].strip()
            elif "is a " in clean_desc:
                for line in clean_desc.split('.'):
                    if "is a" in line[:30] or "describes" in line[:30]:
                        meaning = line.strip() + "."
                        break
            
            if not meaning:
                text_lines = [l.strip() for l in clean_desc.split('\n') if l.strip()]
                for line in text_lines:
                    if "Word of the Day for" in line: continue
                    if "What It Means" in line: continue
                    if word.lower() in line.lower() and ("is a" in line or "describes" in line or "means" in line):
                        meaning = line
                        break
                
                if not meaning and text_lines:
                    for line in text_lines:
                        if len(line) > 20 and ":" not in line and "Merriam-Webster" not in line:
                            meaning = line
                            break
            
            meaning = " ".join(meaning.split())
            if not meaning: meaning = "Definition unavailable"
            
            return [f"Word: {word.upper()}", f"Def: {meaning}"]
    except Exception:
        pass
    return None

def _get_random_fallback():
    try:
        word_resp = requests.get("https://random-word-api.herokuapp.com/word", timeout=5).json()
        word = word_resp[0]
        def_url = f"https://api.dictionaryapi.dev/api/v2/entries/en/{word}"
        def_resp = requests.get(def_url, timeout=5).json()
        
        if isinstance(def_resp, list) and len(def_resp) > 0:
            entry = def_resp[0]
            meanings = entry.get('meanings', [])
            definition = meanings[0].get('definitions', [{}])[0].get('definition', "No definition.")
            return [f"Word: {word.upper()} (Fallback)", f"Def: {definition}"]
    except Exception:
        pass
    return ["Word of the Day: Service unavailable"]
