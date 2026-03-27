# Configuration Guide

DPPR is controlled by two main files: `config.yaml` for layout and settings, and `.env` for private credentials.

## config.yaml

The configuration is divided into three sections: `user`, `settings`, and `layout`.

### User Section
- `name`: Your name as it appears in the greeting.
- `greeting`: The greeting text (e.g., "Good Morning").
- `lat` / `lon`: Your geographical coordinates (essential for weather and astronomy).

### Settings Section
- `font`: The primary font name.
- `master_currency`: The currency code (e.g., "USD", "INR") to which all financial data is converted.
- `printing`: 
  - `enabled`: Set to `true` to send the PDF to a printer.
  - `cups_printer`: The name of your CUPS printer (optional).

### Layout Section
This is an ordered list of modules. Each entry needs:
- `module`: The filename of the module (without `.py`).
- `title`: The header text for that section in the report.
- `config`: (Optional) Module-specific parameters like `count` or `symbols`.

---

## .env (Environment Variables)

Never commit this file to version control. It stores your secrets:

### Home Assistant
- `HA_URL`: The full API URL (e.g., `http://192.168.1.50:8123/api`).
- `HA_ACCESS_TOKEN`: A Long-Lived Access Token created in your HA profile.

### Email
- `EMAIL_USER`: Your full email address.
- `EMAIL_PASS`: Your password or App Password (recommended for Gmail).
- `EMAIL_HOST`: The IMAP server address.

### Actual Budget
- `ACTUAL_SERVER_URL`: Your Actual Budget server URL.
- `ACTUAL_PASSWORD`: Your login password.
- `ACTUAL_SYNC_ID`: The "Sync ID" found in your budget settings.
- `ACTUAL_ENCRYPTION_PASSWORD`: Your budget file's encryption password (optional).

### Coordinates (Fallback)
- `LAT`: Default latitude.
- `LON`: Default longitude.

---

## JSON Data Files

### holidays.json
Used by the `calendar` module. A simple mapping of dates to names:
```json
{
  "2026-01-01": "New Year's Day",
  "2026-12-25": "Christmas"
}
```

### quotes.json
Used by the `quotes` module. A list of objects with text and attribution:
```json
[
  {
    "text": "The only way to do great work is to love what you do.",
    "attribution": "Steve Jobs"
  },
  {
    "text": "Strive not to be a success, but rather to be of value.",
    "attribution": "Albert Einstein"
  }
]
```
