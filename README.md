# DPTR: Daily Personal Terminal Report

A modular, Go-native rewrite of [DPPR](https://github.com/FlashWreck/DPPR). Instead of printing a PDF, DPTR opens a terminal window with a richly formatted daily briefing whenever you log into your Linux PC after "waking up" (configurable time gap + hour cutoff).

## How the Wake-Up Guard Works

DPTR runs on every graphical login via a **systemd user service**. Each time it runs it checks two conditions:

1. **Time gap** - has it been at least `min_gap_hours` (default **6 h**) since the last report was shown?
2. **Hour cutoff** - is the current time before `cutoff_hour` (default **15:00 / 3 PM**)?

If both are satisfied, the report opens in your terminal emulator. Both thresholds are fully configurable in `config.yaml`.

State is stored in `~/.local/share/dptr/last_run` (a Unix timestamp).

---

## Documentation

For more detailed information, please refer to the documents in the `docs/` folder:
- **[Getting Started](docs/GETTING_STARTED.md)**: Detailed setup and installation instructions.
- **[Configuration](docs/CONFIGURATION.md)**: Full config.yaml and environment variable options.
- **[Modules](docs/MODULES.md)**: Complete catalog of all available modules.
- **[Development](docs/DEVELOPMENT.md)**: Guide on how to write new Go modules for DPTR.

---

## Build yourself

### NixOS (flakes)

```bash
# Run directly without installing
nix run github:sidharthify/dptr -- --force --config config.yaml

# Or build it locally
nix build .
./result/bin/dptr --force --config config.yaml

# Or open the dev shell
nix develop
go build ./cmd/dptr
```

### NixOS (legacy nix-shell)

```bash
# One-time build
nix-shell -p go --run "go build ./cmd/dptr"

# Or use the dev shell
nix-shell
go build ./cmd/dptr
```

### With Go installed globally

```bash
go build ./cmd/dptr
```

---

## Modules

| Module key | Description |
|---|---|
| `weather` | Temperature, humidity, wind (Open-Meteo) |
| `astronomy` | Sunrise / sunset (Open-Meteo) |
| `openmeteo` | US AQI, PM2.5, PM10 (Open-Meteo Air Quality) |
| `breatheoss` | AQI for Jammu & Kashmir cities (BreatheOSS) |
| `news` | RSS/Atom feeds, Reddit, and 25+ predefined sources |
| `hacker_news` | Top HN stories with score |
| `github_trending` | GitHub repos trending this week |
| `calendar` | Next holiday + days until next Sunday |
| `system` | Disk, RAM, CPU usage, CPU temp, uptime |
| `ha` | Home Assistant entities |
| `email_inbox` | Unread emails via IMAP |
| `gplay` | Google Play Scraper installs |
| `daily_joke` | Random dad joke (icanhazdadjoke) |
| `fact_of_the_day` | Random fact (uselessfacts.jsph.pl) |
| `shower_thoughts` | Top daily posts from r/showerthoughts |
| `word_of_the_day` | Word + definition (daily rotation) |
| `quotes` | Random quote from `quotes.json` (footer) |
| `essentials` | Static confirmation line |

> **Note**: The Finances block (stocks, crypto, forex, actual budget) has been removed because of missing libraries in Go.

### Predefined RSS Feeds

Use these keys in the `news` module feeds list:

| Region | Keys |
|---|---|
| Global | `bbc`, `aljazeera`, `ap`, `npr`, `nyt` |
| Europe | `dw`, `france24`, `euronews`, `guardian` |
| Asia | `cna`, `hindu`, `toi`, `ndtv`, `scmp`, `yourstory`, `technode` |
| Oceania | `abc_au`, `rnz` |
| LatAm | `mercopress`, `batimes` |
| Tech | `techcrunch`, `verge`, `theverge`, `wired`, `arstechnica`, `theregister`, `noted`, `selfh_st` |

Or use any subreddit or custom URL:

```yaml
- module: "news"
  title: "REDDIT"
  feeds:
    - subreddit: "selfhosted"
      count: 3
    - url: "https://example.com/feed.xml"
      count: 2
```

---

## License

MIT
