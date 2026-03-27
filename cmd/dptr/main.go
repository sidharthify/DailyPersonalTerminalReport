package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/sidharthify/dptr/internal/config"
	"github.com/sidharthify/dptr/internal/renderer"
	"github.com/sidharthify/dptr/internal/runner"
	"github.com/sidharthify/dptr/internal/wakeup"
)

func main() {
	configPath := flag.String("config", "config.yaml", "Path to config.yaml")
	force := flag.Bool("force", false, "Bypass wake-up guard and show report regardless")
	testWakeup := flag.Bool("test-wakeup", false, "Print wake-up guard status and exit")
	noWakeup := flag.Bool("no-wakeup-check", false, "Skip guard but still mark shown (for scripting)")
	showInTerm := flag.Bool("terminal", false, "Open report in a new terminal window (used internally by the service)")
	flag.Parse()

	if !filepath.IsAbs(*configPath) {
		home, _ := os.UserHomeDir()
		defaultCfg := filepath.Join(home, ".config", "dptr", "config.yaml")

		if _, err := os.Stat(*configPath); err == nil {
			// keep *configPath
		} else if _, err := os.Stat(defaultCfg); err == nil {
			*configPath = defaultCfg
		} else {
			exe, _ := os.Executable()
			exeCfg := filepath.Join(filepath.Dir(exe), *configPath)
			if _, err := os.Stat(exeCfg); err == nil {
				*configPath = exeCfg
			} else {
				// Default to ~/.config/... (so the error message is helpful)
				*configPath = defaultCfg
			}
		}
	}
	projectDir := filepath.Dir(*configPath)

	cfg, err := config.Load(*configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "dptr: %v\n", err)
		os.Exit(1)
	}

	if *testWakeup {
		fmt.Print(wakeup.Status(cfg.Wakeup))
		return
	}

	if !*force && !*noWakeup {
		show, reason := wakeup.ShouldShow(cfg.Wakeup)
		if !show {
			_ = reason
			return
		}
	}

	if *showInTerm {
		openInTerminal(cfg.Wakeup.Terminal, os.Args)
		return
	}

	if !*force {
		if err := wakeup.MarkShown(); err != nil {
			fmt.Fprintf(os.Stderr, "dptr: warning — could not update state file: %v\n", err)
		}
	}

	sections, quote := runner.RunModules(cfg, projectDir)

	renderer.RenderReport(
		cfg.User.Name,
		cfg.User.Greeting,
		sections,
		quote,
	)

	fmt.Print("\nPress ENTER to close...")
	fmt.Scanln()
}

func openInTerminal(term string, originalArgs []string) {
	exe, err := os.Executable()
	if err != nil {
		exe = originalArgs[0]
	}

	args := []string{exe}
	for _, a := range originalArgs[1:] {
		if a != "--terminal" && a != "-terminal" {
			args = append(args, a)
		}
	}
	args = append(args, "--no-wakeup-check")

	var cmd *exec.Cmd
	switch term {
	case "kitty":
		cmd = exec.Command("kitty", append([]string{"--"}, args...)...)
	case "alacritty":
		cmd = exec.Command("alacritty", append([]string{"-e"}, args...)...)
	case "gnome-terminal":
		cmd = exec.Command("gnome-terminal", append([]string{"--"}, args...)...)
	case "xterm":
		cmd = exec.Command("xterm", append([]string{"-e"}, args...)...)
	case "konsole":
		cmd = exec.Command("konsole", append([]string{"-e"}, args...)...)
	case "wezterm":
		cmd = exec.Command("wezterm", append([]string{"start", "--"}, args...)...)
	default:
		cmd = exec.Command(term, append([]string{"--"}, args...)...)
	}

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		fmt.Fprintf(os.Stderr, "dptr: could not open terminal %q: %v — rendering inline\n", term, err)
	}
}
