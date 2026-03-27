# Module Catalog

Every module in DPTR is specialized for a single type of data. All modules are fetched concurrently. Detailed configuration options for every module are below.

## Environment

### weather
- **config**: `lat`, `lon` (optional if set globally).
- **Output**: Current temp, humidity, wind speed, min/max for the day via Open-Meteo.

### astronomy
- **config**: `lat`, `lon`.
- **Output**: Local sunrise and sunset times via Open-Meteo.

### openmeteo (AQI)
- **config**: `lat`, `lon`.
- **Output**: US Air Quality Index and pollutant breakdown (PM2.5/PM10).

### breatheoss (AQI)
- **config**: `city` (string like `Srinagar`).
- **Output**: Local AQI for regions in Jammu and Kashmir (via BreatheOSS).

---

## News

### news
- **config**: `feeds` (list of `feed`, `subreddit`, or `url`).
- **Predefined Keys**: `bbc`, `techcrunch`, `hackernews`, `verge`, `dw`, `aljazeera` etc. (See README for full list), or your own custom RSS/Atom URL.
- **Output**: Headings and short descriptions of the latest stories.

### hacker_news / github_trending
- **config**: `count` (number of items), `language` (for github_trending, e.g., `go`, `rust`).
- **Output**: Top HN stories or top trending GitHub repositories.

---

## Planning

### calendar
- **config**: `holidays_file` (optional JSON file path).
- **Output**: Upcoming holidays and the date of the (next) Sunday.

---

## Random

### word_of_the_day
- **config**: (None).
- **Output**: A daily seeded esoteric word and its definition.

### daily_joke / fact_of_the_day / shower_thoughts
- **config**: (None for joke/fact), `count` (for shower_thoughts).
- **Output**: Random curated humor, information, or Reddit r/showerthoughts.

### quotes
- **config**: `quotes_file` (optional JSON file).
- **Output**: A random quote or lyric rendered in the report footer.

---

## Productivity & System

### email_inbox
- **config**: `count` (number of emails).
- **Env**: Requires `EMAIL_USER`, `EMAIL_PASS`, `EMAIL_HOST`.
- **Output**: Recent unseen subject lines and senders fetched via IMAP.

### gplay (Google Play)
- **config**: `package_name` (optional if set in `.env`).
- **Output**: Download estimates for your Android app on the Google Play Store (via native Go scraping).

### system
- **config**: `disks` (list of `path` and `label`).
- **Output**: Disk usage, RAM utilization, CPU stats and temp, and system uptime from `/proc` and `/sys`.

### essentials
- **config**: (None).
- **Output**: A confirmation message: "⚡ Report generated. Have a great day!"
