# DPPR: Daily Personal Printed Report

A modular, configuration-driven Python system that generates a personalized daily briefing as a PDF and optionally prints it via CUPS. Every section of the report is a self-contained module that can be enabled, disabled, reordered, or swapped out entirely through a single YAML file.

## Features

- **Modular Architecture**: Add, remove, or reorder any section using `config.yaml`.
- **Extensible Module System**: Easily write and plug in your own Python modules for custom data.
- **Multi-OS Support**: Automatic font discovery for Linux, Windows, and macOS.
- **Automated Printing**: Direct integration with CUPS for daily physical reports.

## Why DPPR?

The inspiration for this project came from a Reddit post where someone was printing their daily reports using a thermal receipt printer. That sparked the idea for a similar, but more comprehensive system that works with any regular home printer.

The goal of DPPR is to provide a curated, physical report to start your morning. Instead of immediately diving into your phone and getting lost in notifications, you can sit down with a cup of coffee and read a one-page summary of everything that matters to you: your local weather, the latest news, your financial portfolios, and your own server's health. It's about taking back the first few minutes of your day with a focused, reading experience.

The system is built around three core ideas:

1. **Modules** are small, independent Python files that each fetch one type of data.
2. **Layout** is defined in `config.yaml`, which controls the order and configuration of modules.
3. **The Engine** takes module output and renders it into a formatted PDF.

---

## 📚 Documentation

For complete setup and configuration details, please see our guides:

*   [Getting Started](docs/GETTING_STARTED.md) - Installation and OS-specific requirements.
*   [Configuration Guide](docs/CONFIGURATION.md) - Deep dive into `config.yaml` and `.env`.
*   [Module Catalog](docs/MODULES.md) - Full list of all available modules and their options.
*   [Development Guide](docs/DEVELOPMENT.md) - How to write your own custom modules.

## Project Structure

```
DPPR/
├── main.py                  Entry point, loads config and runs modules
├── config.yaml              Your personal configuration (gitignored)
├── config.template.yaml     Reference config with all available options
├── .env                     API keys and credentials (gitignored)
├── requirements.txt         Python dependencies
├── holidays.json            Custom holiday dates (gitignored)
├── quotes.json              Custom quotes/lyrics for the footer (gitignored)
├── signature.png            Optional signature image (gitignored)
│
├── engine/
│   ├── pdf_generator.py     Renders report data into a PDF
│   └── printer.py           Sends the PDF to a CUPS printer
│
└── modules/
    ├── environment/
    │   ├── weather/weather.py       Current weather via Open-Meteo
    │   ├── astronomy.py             Sunrise and sunset times
    │   └── aqi/
    │       ├── openmeteo.py         US AQI via Open-Meteo (Global)
    │       └── breatheoss.py        US AQI via BreatheOSS (Jammu and Kashmir, India)
    │
    ├── news/
    │   ├── news.py                  RSS/Atom feed aggregator with Reddit support
    │   ├── hacker_news.py           Top stories from Hacker News
    │   └── github_trending.py       Trending GitHub repositories
    │
    ├── finances/
    │   ├── exchange_rate.py         Forex rates
    │   ├── markets.py               Stock prices via Yahoo Finance
    │   ├── crypto.py                Crypto prices via CoinGecko
    │   ├── commodities.py           Gold, Silver, Oil via Yahoo Finance
    │   ├── actual.py                Bank balances via Actual Budget
    │   └── currency_util.py         Shared currency conversion helpers
    │
    ├── planning/
    │   └── calendar.py              Upcoming holidays and next Sunday
    │
    ├── home_automation/
    │   └── ha.py                    Entity states from Home Assistant
    │
    ├── communication/
    │   └── email_inbox.py           Recent emails via IMAP
    │
    ├── android/
    │   └── gplay.py                 Google Play download stats
    │
    ├── knowledge/
    │   └── fact_of_the_day.py       Daily fact from Useless Facts API
    │
    ├── social/
    │   └── shower_thoughts.py       Top posts from r/ShowerThoughts
    │
    ├── random/
    │   ├── word_of_the_day.py       Word of the Day from Merriam-Webster
    │   ├── daily_joke.py            Dad joke from icanhazdadjoke
    │   └── quotes.py                Random quote/lyric for the footer
    │
    ├── system/
    │   └── system.py                Disk, RAM, and CPU stats
    │
    └── essentials/
        └── essentials.py            Static confirmation line
```

## How Modules Work

Every module is a Python file inside `modules/` that exposes a single function:

```python
def get_data(config):
    return ["Line 1", "Line 2", "Line 3"]
```

- `config` is a dictionary passed from the module's entry in `config.yaml`.
- The function returns a list of strings. Each string becomes one line in the report.
- The engine handles text wrapping, encoding, and page breaks automatically.

Modules are discovered by filename. If you create `modules/custom/my_thing.py`, you reference it in the layout as `module: "my_thing"`. The engine walks the `modules/` directory tree to find it.

### The `quotes` Module

The `quotes` module is special. Instead of appearing in the report body, its output is rendered as a footer quote at the bottom of the page. It does not need a `title` in the layout.

## Setup

### 1. Clone and Install

```bash
git clone <repo-url> && cd DPPR
python -m venv .venv && source .venv/bin/activate
pip install -r requirements.txt
```

### 2. Configure

```bash
cp config.template.yaml config.yaml
```

Edit `config.yaml` to:
- Set your `name` and `greeting` under `user`.
- Set your `lat` and `lon` under `user` (used by weather, AQI, and astronomy modules).
- Set `master_currency` under `settings` to auto-convert all financial data.
- Add, remove, or reorder modules in the `layout` list.

### 3. Environment Variables

Create a `.env` file for credentials that should not live in the config:

```env
EMAIL_USER=you@gmail.com
EMAIL_PASS=your_app_password
EMAIL_HOST=imap.gmail.com
HA_URL=http://homeassistant.local:8123
HA_ACCESS_TOKEN=your_long_lived_token
ACTUAL_SERVER_URL=https://actual.example.com
ACTUAL_PASSWORD=your_password
ACTUAL_SYNC_ID=your_sync_id
ACTUAL_ENCRYPTION_PASSWORD=optional
PACKAGE_NAME=com.example.app
```

Only the modules you enable need their corresponding variables. Unused modules will not throw errors for missing credentials.

### 4. Optional Files

| File | Purpose |
| :--- | :--- |
| `holidays.json` | A `{"YYYY-MM-DD": "Name"}` mapping of holidays for the calendar module |
| `quotes.json` | A list of `{"text": "...", "attribution": "..."}` objects for the footer |
| `signature.png` | An image rendered in the bottom-right corner of the report |

### 5. Fonts

The PDF generator looks for JetBrains Mono in standard font directories across Linux, Windows, and macOS. If not found, it falls back to Courier. To install JetBrains Mono:

- **Linux**: `sudo pacman -S ttf-jetbrains-mono` or `sudo apt install fonts-jetbrains-mono`
- **macOS**: `brew install --cask font-jetbrains-mono`
- **Windows**: Download from [JetBrains](https://www.jetbrains.com/lp/mono/) and install to `C:\Windows\Fonts`

### 6. Run

```bash
python main.py              # Generate PDF and print
python main.py --no-print   # Generate PDF only
```

The output is saved to `daily_report.pdf` in the project root.

## Predefined News Feeds

The `news` module ships with a built-in directory of RSS feeds. Use these keys in your config instead of raw URLs:

| Region | Keys |
| :--- | :--- |
| Global | `bbc`, `aljazeera`, `reuters`, `ap`, `npr`, `nyt` |
| Europe | `dw`, `france24`, `euronews`, `guardian` |
| Asia | `cna`, `bangkokpost`, `hindu`, `toi`, `ndtv`, `yourstory`, `scmp`, `technode`, `caixin` |
| Oceania | `abc_au`, `rnz` |
| Latin America | `mercopress`, `batimes` |
| Tech | `techcrunch`, `verge`, `wired`, `arstechnica`, `theregister`, `hackernews`, `noted`, `selfh_st` |

You can also pull from any subreddit:

```yaml
- module: "news"
  title: "REDDIT"
  feeds:
    - subreddit: "selfhosted"
      count: 3
    - subreddit: "linux"
      count: 2
```

**Custom RSS/Atom Feeds**:
You can also add any custom RSS or Atom feed URL in the news module:

```yaml
- module: "news"
  title: "CUSTOM FEED"
  feeds:
    - url: "https://example.com/feed.xml"
      count: 2
```

## Currency Conversion

All financial modules (stocks, crypto, commodities) support automatic currency conversion. Set `master_currency` in `settings` and every price will be converted from its native currency:

```yaml
settings:
  master_currency: "INR"
```

Conversion rates are fetched live from Yahoo Finance.

## Printing

DPPR integrates with CUPS for automated printing. Enable it in `config.yaml`:

```yaml
settings:
  printing:
    enabled: true
    cups_printer: "HP_LaserJet"
```

If `cups_printer` is not set, it defaults to the first available printer. Requires `pycups` (`pip install pycups`) and a working CUPS setup on the host.

## Writing a Custom Module

1. Create a file anywhere inside `modules/`, e.g. `modules/custom/uptime.py`.
2. Implement `get_data(config)`:

```python
import requests

def get_data(config):
    url = config.get('url')
    try:
        resp = requests.get(url, timeout=5)
        return [f"Status: {resp.status_code}"]
    except Exception as e:
        return [f"Error: {e}"]
```

3. Add it to your layout:

```yaml
- module: "uptime"
  title: "UPTIME CHECK"
  config:
    url: "https://example.com"
```

That is it. The engine discovers the module by filename and passes the `config` block to `get_data()`.

## Contributing

If you make a custom module for your own needs, please be kind enough to make a Pull Request and help others out! Your contributions help keep the project diverse and useful for everyone.

## License

MIT
