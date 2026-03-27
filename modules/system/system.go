package system

import (
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"
)

type System struct{}

func (System) GetData(cfg map[string]any) []string {
	var out []string

	disksRaw, _ := cfg["disks"].([]any)
	for _, d := range disksRaw {
		dm, ok := d.(map[string]any)
		if !ok {
			continue
		}
		path, _ := dm["path"].(string)
		label, _ := dm["label"].(string)
		if path == "" {
			path = "/"
		}
		if label == "" {
			label = "Disk"
		}
		if line, err := diskUsage(path, label); err == nil {
			out = append(out, line)
		}
	}
	if len(disksRaw) == 0 {
		if line, err := diskUsage("/", "Root"); err == nil {
			out = append(out, line)
		}
	}

	out = append(out, ramUsage())

	cpuLine, tempLine := cpuStats()
	out = append(out, cpuLine)
	if tempLine != "" {
		out = append(out, tempLine)
	}

	out = append(out, uptimeStr())

	return out
}

func diskUsage(path, label string) (string, error) {
	total, used, err := statfs(path)
	if err != nil {
		return "", err
	}
	pct := float64(used) / float64(total) * 100
	return fmt.Sprintf("Disk %-6s : %.1f GB / %.1f GB  (%.1f%% used)",
		label,
		float64(used)/(1<<30),
		float64(total)/(1<<30),
		pct,
	), nil
}

func ramUsage() string {
	data, err := os.ReadFile("/proc/meminfo")
	if err != nil {
		return "RAM: Unavailable"
	}
	var total, avail int64
	for _, line := range strings.Split(string(data), "\n") {
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}
		val, _ := strconv.ParseInt(fields[1], 10, 64)
		switch fields[0] {
		case "MemTotal:":
			total = val
		case "MemAvailable:":
			avail = val
		}
	}
	if total == 0 {
		return "RAM: Unavailable"
	}
	used := total - avail
	return fmt.Sprintf("RAM        : %.1f GB / %.1f GB  (%.1f%% used)",
		float64(used)/1048576, float64(total)/1048576, float64(used)/float64(total)*100)
}

func cpuStats() (cpuLine, tempLine string) {
	usage, err := cpuUsage()
	if err != nil {
		cpuLine = "CPU        : Unavailable"
	} else {
		cpuLine = fmt.Sprintf("CPU        : %.1f%%  (%d cores)", usage, runtime.NumCPU())
	}

	for _, zone := range []string{"thermal_zone0", "thermal_zone1", "thermal_zone2"} {
		path := "/sys/class/thermal/" + zone + "/temp"
		data, err := os.ReadFile(path)
		if err != nil {
			continue
		}
		val, err := strconv.ParseFloat(strings.TrimSpace(string(data)), 64)
		if err != nil {
			continue
		}
		tempLine = fmt.Sprintf("CPU Temp   : %.1f°C", val/1000)
		break
	}
	return
}

func cpuUsage() (float64, error) {
	read := func() (total, idle float64, err error) {
		data, err := os.ReadFile("/proc/stat")
		if err != nil {
			return
		}
		for _, line := range strings.Split(string(data), "\n") {
			if !strings.HasPrefix(line, "cpu ") {
				continue
			}
			fields := strings.Fields(line)[1:]
			var vals []float64
			for _, f := range fields {
				v, _ := strconv.ParseFloat(f, 64)
				vals = append(vals, v)
			}
			for _, v := range vals {
				total += v
			}
			if len(vals) > 3 {
				idle = vals[3]
			}
			return
		}
		err = fmt.Errorf("cpu line not found")
		return
	}

	t1, i1, err := read()
	if err != nil {
		return 0, err
	}
	time.Sleep(150 * time.Millisecond)
	t2, i2, err := read()
	if err != nil {
		return 0, err
	}

	dt := t2 - t1
	di := i2 - i1
	if dt == 0 {
		return 0, nil
	}
	return 100 * (1 - di/dt), nil
}

func uptimeStr() string {
	data, err := os.ReadFile("/proc/uptime")
	if err != nil {
		return ""
	}
	fields := strings.Fields(string(data))
	if len(fields) == 0 {
		return ""
	}
	secs, _ := strconv.ParseFloat(fields[0], 64)
	d := time.Duration(secs) * time.Second
	h := int(d.Hours())
	m := int(d.Minutes()) % 60
	return fmt.Sprintf("Uptime     : %dh %dm", h, m)
}

type Essentials struct{}

func (Essentials) GetData(_ map[string]any) []string {
	return []string{"⚡ Report generated. Have a great day!"}
}
