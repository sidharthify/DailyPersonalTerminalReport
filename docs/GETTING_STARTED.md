# Getting Started

Welcome to **DPTR (Daily Personal Terminal Report)** — a modular, pure Go system that opens a terminal window with a richly formatted daily briefing when you log into your Linux PC.

This completely replaces the older Python/PDF-based `DPPR`.

## Prerequisites

- **Go 1.21+** (if compiling manually)
- **Nix** (optional, but highly recommended for NixOS users)
- A supported terminal emulator (kitty, alacritty, gnome-terminal, etc.)

## Installation

### Method 1: Nix Flakes (Recommended for Nix users)

You can run DPTR directly via flakes without installing it permanently:
```bash
nix run github:sidharthify/dptr -- --force --config config.yaml
```

To install or build locally:
```bash
git clone https://github.com/sidharthify/dptr
cd dptr
nix build .

# The binary is now in result/bin/dptr
./result/bin/dptr --force --config config.template.yaml
```

### Method 2: Standard Go Build

```bash
git clone https://github.com/sidharthify/dptr
cd dptr
go build ./cmd/dptr

# The binary is created as ./dptr
./dptr --force --config config.template.yaml
```

## Systemd Auto-Start Setup

To have DPTR automatically show up when you wake up and log into your PC:

1. Copy the reference config:
   ```bash
   cp config.template.yaml config.yaml
   ```
2. Edit `config.yaml` to your liking (add your coordinates, feeds, preferred terminal).
3. Run the installer:
   ```bash
   chmod +x install/install.sh
   ./install/install.sh
   ```

The installer builds the binary, places it in `~/.local/bin/dptr`, copies your config to `~/.config/dptr/config.yaml`, and enables a systemd user service (`dptr.service`) that triggers on graphical login.

## Your First Run

Run a dry-run to see your DPTR layout immediately, bypassing the wake-up guard:
```bash
./dptr --force --config config.yaml
```

To check how the wake-up guard is currently evaluating:
```bash
./dptr --test-wakeup --config config.yaml
```

## Troubleshooting

- **Report doesn't show on login?**: The wake-up guard might be preventing it if you already logged in recently or it's past your cutoff hour. Check `systemctl --user status dptr.service` or run `dptr --test-wakeup`.
- **Command missing?**: Ensure `~/.local/bin` is in your `$PATH`.