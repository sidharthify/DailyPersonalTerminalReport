package wakeup

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/sidharthify/dptr/internal/config"
)

const stateFile = ".local/share/dptr/last_run"

func stateFilePath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, stateFile)
}

func ShouldShow(cfg config.WakeupConfig) (bool, string) {
	if !cfg.Enabled {
		return false, "wakeup disabled in config"
	}

	now := time.Now()
	hour := now.Hour()

	if hour >= cfg.CutoffHour {
		return false, fmt.Sprintf(
			"current hour %02d:xx >= cutoff %02d:00 — too late in the day",
			hour, cfg.CutoffHour,
		)
	}

	lastRun, err := readLastRun()
	if err != nil {
		return true, "first run (no state file)"
	}

	elapsed := now.Sub(lastRun)
	minGap := time.Duration(cfg.MinGapHours) * time.Hour

	if elapsed < minGap {
		remaining := minGap - elapsed
		return false, fmt.Sprintf(
			"only %.1f h since last run (need >= %d h, %.1f h remaining)",
			elapsed.Hours(), cfg.MinGapHours, remaining.Hours(),
		)
	}

	return true, fmt.Sprintf("gap %.1f h >= %d h threshold, hour %02d < cutoff %02d", elapsed.Hours(), cfg.MinGapHours, hour, cfg.CutoffHour)
}

func MarkShown() error {
	path := stateFilePath()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	ts := strconv.FormatInt(time.Now().Unix(), 10)
	return os.WriteFile(path, []byte(ts), 0o644)
}

func Status(cfg config.WakeupConfig) string {
	show, reason := ShouldShow(cfg)
	lastRun, err := readLastRun()
	var lastStr string
	if err != nil {
		lastStr = "never"
	} else {
		lastStr = lastRun.Format("2006-01-02 15:04:05")
	}
	verdict := "NO"
	if show {
		verdict = "YES"
	}
	return fmt.Sprintf(
		"Last shown : %s\nWould show : %s\nReason     : %s\nMin gap    : %d h\nCutoff     : %02d:00\n",
		lastStr, verdict, reason, cfg.MinGapHours, cfg.CutoffHour,
	)
}

func readLastRun() (time.Time, error) {
	data, err := os.ReadFile(stateFilePath())
	if err != nil {
		return time.Time{}, err
	}
	raw := strings.TrimSpace(string(data))
	ts, err := strconv.ParseInt(raw, 10, 64)
	if err != nil {
		return time.Time{}, err
	}
	return time.Unix(ts, 0), nil
}
