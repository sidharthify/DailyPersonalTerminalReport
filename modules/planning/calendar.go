package planning

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"time"
)

type Calendar struct{}

func (Calendar) GetData(cfg map[string]any) []string {
	holidaysFile := "holidays.json"
	if v, ok := cfg["holidays_file"].(string); ok && v != "" {
		holidaysFile = v
	}

	today := time.Now().Truncate(24 * time.Hour)
	var out []string

	holidayStr := "No upcoming holidays found"
	if data, err := os.ReadFile(holidaysFile); err == nil {
		var holidays map[string]string
		if json.Unmarshal(data, &holidays) == nil {
			type entry struct {
				date time.Time
				name string
			}
			var upcoming []entry
			for ds, name := range holidays {
				t, err := time.Parse("2006-01-02", ds)
				if err != nil {
					continue
				}
				if !t.Before(today) {
					upcoming = append(upcoming, entry{t, name})
				}
			}
			sort.Slice(upcoming, func(i, j int) bool {
				return upcoming[i].date.Before(upcoming[j].date)
			})
			if len(upcoming) > 0 {
				next := upcoming[0]
				days := int(next.date.Sub(today).Hours() / 24)
				if days == 0 {
					holidayStr = "TODAY: " + next.name
				} else {
					holidayStr = fmt.Sprintf("%s  (%s, %d days away)",
						next.name, next.date.Format("02 Jan"), days)
				}
			}
		}
	} else {
		holidayStr = fmt.Sprintf("holidays.json not found (%s)", holidaysFile)
	}
	out = append(out, "Nearest Holiday : "+holidayStr)

	weekday := int(today.Weekday())
	daysToSunday := (7 - weekday) % 7
	if daysToSunday == 0 {
		daysToSunday = 7
	}
	nextSunday := today.AddDate(0, 0, daysToSunday)
	out = append(out, fmt.Sprintf("Next Sunday     : %s  (%d days)", nextSunday.Format("02 Jan"), daysToSunday))

	return out
}
