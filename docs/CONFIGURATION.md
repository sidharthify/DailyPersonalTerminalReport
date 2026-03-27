# Configuration Guide

DPTR is controlled entirely by `config.yaml`. There are no `.env` files — if a module requires a secret (like an API token or password), you can still export it in your environment (e.g., `EMAIL_PASSWORD=...`) or add it to your systemd service overrides.

## config.yaml Structure

The configuration is divided into four main sections: `user`, `settings`, `wakeup`, and `layout`.

### 1. User Section
Basic personal details passed down to modules.
- `name`: Your name as it appears in the greeting.
- `greeting`: The greeting text (e.g., "Good Morning").
- `lat` / `lon`: Your geographical coordinates (essential for weather, AQI, and astronomy).

### 2. Settings Section
Global settings.
*(Note: `master_currency` used for finances is no longer needed since the finance modules were removed because of missing libraries in Go).*

### 3. Wake-Up Section
Controls when the report is actually displayed on login.
- `enabled`: Set to `true` to use the guard. If `false`, the report displays on *every* invocation.
- `min_gap_hours`: Minimum hours since the last report was shown (default: `6`).
- `cutoff_hour`: The hour in 24h format after which the report won't show (e.g., `15` means 3:00 PM).
- `terminal`: The terminal emulator to launch the report in (e.g., `kitty`, `alacritty`, `wezterm`, `gnome-terminal`).

### 4. Layout Section
An ordered list of modules to run concurrently. Each entry needs:
- `module`: The internal name of the module (e.g., `weather`, `news`, `markets`).
- `title`: The bold header text for that section in the terminal.
- `config`: (Optional) Module-specific parameters like `count`, `symbols`, or `entities`.
- `feeds`: (Specific to the `news` module) A list of predefined feeds, subreddits, or custom RSS URLs.

---

## Environment Variables for Secure Modules

Some modules require authentication. DPTR reads these securely from your environment variables:

### Home Assistant (`ha` module)
- `HA_URL`: The full API URL (e.g., `http://192.168.1.50:8123/api`).
- `HA_ACCESS_TOKEN`: A Long-Lived Access Token created in your HA profile.

### Email Inbox (`email_inbox` module)
- `EMAIL_HOST`: The IMAP server address (e.g., `imap.gmail.com`).
- `EMAIL_USER`: Your full email address.
- `EMAIL_PASS`: Your password or App Password.

### Google Play (`gplay` module)
- `PACKAGE_NAME`: The Android package name to scrape (e.g., `com.example.app`). Can alternatively be set in `config.yaml`.

---

## JSON Data Files

### holidays.json (for `calendar`)
Used by the `calendar` module. A simple mapping of `YYYY-MM-DD` to names. Keep this file in the directory you run the binary from, or use the `holidays_file` config option to point to an absolute path.
```json
{
  "2026-01-01": "New Year's Day",
  "2026-12-25": "Christmas"
}
```

### quotes.json (for `quotes`)
Used by the `quotes` module. A list of objects with text and attribution.
```json
[
  {
    "text": "The only way to do great work is to love what you do.",
    "attribution": "Steve Jobs"
  }
]
```
