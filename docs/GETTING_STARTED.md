# Getting Started

Welcome to DPPR (Daily Personal Printed Report). This guide will help you get the system running on your local machine or server.

## Prerequisites

- **Python 3.8+**
- **pip** (Python package installer)
- **CUPS** (If you want to use the automated printing feature)

## Installation

1. **Clone the repository**:
   ```bash
   git clone https://github.com/FlashWreck/DPPR
   cd DPPR
   ```

2. **Create a virtual environment**:
   ```bash
   python -m venv .venv
   source .venv/bin/activate  # On Windows: .venv\Scripts\activate
   ```

3. **Install dependencies**:
   ```bash
   pip install -r requirements.txt
   ```

## OS-Specific Requirements

### Linux (Ubuntu/Debian)
If you plan to use printing, you may need the CUPS development headers:
```bash
sudo apt install libcups2-dev
```
For the JetBrains Mono font:
```bash
sudo apt install fonts-jetbrains-mono
```

### macOS
For the JetBrains Mono font:
```bash
brew install --cask font-jetbrains-mono
```

### Windows
Download the JetBrains Mono font from the [official website](https://www.jetbrains.com/lp/mono/) and install it system-wide.

## Your First Run

1. Copy the template config:
   ```bash
   cp config.template.yaml config.yaml
   ```
2. Run a dry-run to see if the PDF generates correctly:
   ```bash
   python main.py --no-print
   ```
3. Check `daily_report.pdf` in the project root.

## Troubleshooting

- **Font not found**: If the PDF looks like standard Courier, ensures JetBrains Mono is installed in a standard system directory.
- **ImportError (pycups)**: This usually happens if the CUPS development libraries are missing on Linux. If you do not need printing, you can ignore this and use `--no-print`.
- **API Errors**: Check your `.env` file credentials.