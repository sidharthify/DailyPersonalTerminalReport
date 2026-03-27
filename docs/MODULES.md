# Module Catalog

Every module in DPPR is specialized for a single type of data. Detailed configuration options for every module are below.

## Environment

### weather
- **config**: `lat`, `lon` (optional if set globally).
- **Output**: Current temp, humidity, min/max for the day.

### astronomy
- **config**: `lat`, `lon`.
- **Output**: Local sunrise and sunset times.

### openmeteo (AQI)
- **config**: `lat`, `lon`.
- **Output**: US Air Quality Index and pollutant breakdown (PM2.5/PM10).

### breatheoss (AQI)
- **config**: `city` (string like `srinagar`, `jammu_city`). 
- **Output**: Local AQI for regions in Jammu and Kashmir.

---

## News

### news
- **config**: `feeds` (list of `feed`, `subreddit`, or `url`).
- **Predefined Keys**: `bbc`, `techcrunch`, `hackernews`, `verge`, etc. or your own custom feed.
- **Output**: Headings and short descriptions of the latest stories.

### hacker_news / github_trending
- **config**: `count` (number of items), `language` (for github_trending, e.g., `python`, `javascript`).
- **Output**: Top stories or repositories.

---

## Finances

### markets / crypto / commodities
- **config**: `symbols` (list), `master_currency` (auto-inherited).
- **Output**: Asset prices and 24h percentage change.

### actual
- **config**: `main_account_match` (substring of your account name).
- **Output**: Current bank balance using the Actual Budget server.

### calendar
- **config**: `holidays_file` (optional JSON file path).
- **Output**: Upcoming holidays and the date of the next Sunday.

---

## Random

### word_of_the_day
- **config**: (None).
- **Output**: Merriam-Webster's WOTD and its definition.

### daily_joke / fact_of_the_day
- **config**: (None).
- **Output**: Random curated humor or information.

### quotes
- **config**: `quotes_file` (optional JSON file).
- **Output**: A random quote or lyric rendered in the report footer.

---

## Productivity & System

### email_inbox
- **config**: `count` (number of emails).
- **Env**: Requires `EMAIL_USER`, `EMAIL_PASS`, `EMAIL_HOST`.
- **Output**: Recent subject lines.

### gplay (Google Play)
- **config**: `package_name` (optional if set in `.env`).
- **Output**: Download count and daily stats for your Android app.

### system
- **config**: `disks` (list of `path` and `label`).
- **Output**: Disk usage, RAM utilization, CPU stats and temp.

### ha (Home Assistant)
- **config**: `entities` (list of `entity_id` and `label`).
- **Env**: Requires `HA_URL`, `HA_ACCESS_TOKEN`.
- **Output**: Live states of your smart home sensors.

### essentials (Report Review)
- **config**: `show_review` (boolean, default `true`).
- **Output**: A confirmation message: "Report successfully compiled with latest live data."
