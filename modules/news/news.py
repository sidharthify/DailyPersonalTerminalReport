import requests
import html
import re
from xml.etree import ElementTree

USER_AGENT = 'DPPR-Report-Bot/1.0'

PREDEFINED_FEEDS = {
    # Global / Major
    'bbc': 'http://feeds.bbci.co.uk/news/world/rss.xml',
    'aljazeera': 'https://www.aljazeera.com/xml/rss/all.xml',
    'reuters': 'https://www.reutersagency.com/en/reuters-best/rss-feeds/',
    'ap': 'https://apnews.com/feed',
    'npr': 'https://www.npr.org/rss/rss.php?id=1001',
    'nyt': 'https://rss.nytimes.com/services/xml/rss/nyt/HomePage.xml',

    # Europe
    'dw': 'https://rss.dw.com/rdf/rss-en-all',
    'france24': 'https://www.france24.com/en/rss',
    'euronews': 'https://www.euronews.com/rss?format=google-news&level=theme&name=news',
    'guardian': 'https://www.theguardian.com/world/rss',

    # Asia
    'cna': 'https://www.channelnewsasia.com/rss/news/asia/rss.xml',
    'bangkokpost': 'https://www.bangkokpost.com/rss/data/topstories.xml',
    'hindu': 'https://www.thehindu.com/news/national/feeder/default.rss',
    'toi': 'https://timesofindia.indiatimes.com/rssfeeds/-2128936835.cms',
    'ndtv': 'https://feeds.feedburner.com/ndtvnews-top-stories',
    'yourstory': 'https://yourstory.com/feed',
    'scmp': 'https://www.scmp.com/rss/2/feed.xml',
    'technode': 'https://technode.com/feed/',
    'caixin': 'https://www.caixinglobal.com/rss/all.xml',

    # Oceania
    'abc_au': 'https://www.abc.net.au/news/feed/2942460/rss.xml',
    'rnz': 'https://www.rnz.co.nz/rss/news.xml',

    # LATAM
    'mercopress': 'https://en.mercopress.com/rss/',
    'batimes': 'https://www.batimes.com.ar/rss',

    # Tech & Self-Hosted
    'techcrunch': 'https://techcrunch.com/feed/',
    'verge': 'https://www.theverge.com/rss/index.xml',
    'theverge': 'https://www.theverge.com/rss/index.xml',
    'wired': 'https://www.wired.com/feed/rss',
    'arstechnica': 'https://feeds.arstechnica.com/arstechnica/index',
    'theregister': 'https://www.theregister.com/headlines.rss',
    'hackernews': 'https://news.ycombinator.com/rss',
    'noted': 'https://noted.lol/rss/',
    'selfh_st': 'https://selfh.st/rss/',
}

def get_data(config):
    feeds = config.get('feeds', [])
    results = []
    
    for feed in feeds:
        feed_key = feed.get('feed')
        subreddit = feed.get('subreddit')
        url = feed.get('url')
        
        display_title = ""
        if subreddit:
            url = f"https://www.reddit.com/r/{subreddit}/top.rss?t=day"
            display_title = f"r/{subreddit}"
        elif feed_key and feed_key.lower() in PREDEFINED_FEEDS:
            url = PREDEFINED_FEEDS[feed_key.lower()]
            display_title = f"{feed_key.upper()} News"
        elif url:
            display_title = feed.get('title', 'Custom News')
        else:
            results.append(f"Feed error: No URL, subreddit, or known key for '{feed_key or 'unknown'}'")
            continue

        count = feed.get('count', 1)
        
        try:
            headers = {'User-Agent': USER_AGENT}
            resp = requests.get(url, timeout=10, headers=headers)
            if not resp.content:
                results.append(f"Feed error ({display_title}): Empty response")
                continue
                
            root = ElementTree.fromstring(resp.content)
            
            # Handle both RSS (item) and Atom (entry)
            items = root.findall('.//item')
            if not items:
                items = root.findall('.//{http://www.w3.org/2005/Atom}entry')
                
            items = items[:count]
            for item in items:
                title_elem = item.find('title')
                if title_elem is None:
                    title_elem = item.find('{http://www.w3.org/2005/Atom}title')
                    
                if title_elem is not None:
                    t_text = html.unescape(title_elem.text)
                    t_text = " ".join(t_text.split())
                    if len(t_text) > 110:
                        t_text = t_text[:107] + "..."
                    results.append(f">> {t_text}")
                    
                    desc_elem = item.find('description')
                    if desc_elem is None:
                        desc_elem = item.find('{http://www.w3.org/2005/Atom}summary')
                    if desc_elem is None:
                        desc_elem = item.find('{http://www.w3.org/2005/Atom}content')
                    
                    if desc_elem is not None and desc_elem.text:
                        desc = html.unescape(desc_elem.text or "")
                        # Remove HTML tags
                        desc = re.sub('<[^<]+?>', '', desc)
                        desc = " ".join(desc.split())
                        if len(desc) > 120:
                            desc = desc[:117] + "..."
                        results.append(f"   {desc}")
        except Exception as e:
            results.append(f"Feed error ({display_title}): {e}")
            
    return results
