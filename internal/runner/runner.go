package runner

import (
	"github.com/sidharthify/dptr/internal/config"
	"github.com/sidharthify/dptr/internal/renderer"
	"github.com/sidharthify/dptr/modules/environment"
	"github.com/sidharthify/dptr/modules/news"
	"github.com/sidharthify/dptr/modules/planning"
	"github.com/sidharthify/dptr/modules/random"
	"github.com/sidharthify/dptr/modules/system"
	"sync"
)

type Module interface {
	GetData(cfg map[string]any) []string
}

func moduleFactory(name string) Module {
	switch name {
	case "weather":
		return &environment.Weather{}
	case "astronomy":
		return &environment.Astronomy{}
	case "openmeteo":
		return &environment.OpenMeteoAQI{}
	case "breatheoss":
		return &environment.BreatheOSS{}
	case "news":
		return &news.News{}
	case "hacker_news":
		return &news.HackerNews{}
	case "github_trending":
		return &news.GitHubTrending{}
	case "calendar":
		return &planning.Calendar{}
	case "system":
		return &system.System{}
	case "ha":
		return &system.HomeAssistant{}
	case "email_inbox":
		return &system.EmailInbox{}
	case "gplay":
		return &system.GooglePlay{}
	case "daily_joke":
		return &random.DailyJoke{}
	case "word_of_the_day":
		return &random.WordOfTheDay{}
	case "fact_of_the_day":
		return &random.FactOfTheDay{}
	case "shower_thoughts":
		return &random.ShowerThoughts{}
	case "quotes":
		return &random.Quotes{}
	case "essentials":
		return &system.Essentials{}
	default:
		return nil
	}
}

func RunModules(cfg *config.Config, _ string) ([]renderer.Section, string) {
	type result struct {
		index int
		sec   renderer.Section
		quote string
	}

	results := make([]result, len(cfg.Layout))
	var wg sync.WaitGroup

	for i, entry := range cfg.Layout {
		wg.Add(1)
		go func(idx int, e config.LayoutEntry) {
			defer wg.Done()

			mod := moduleFactory(e.Module)
			if mod == nil {
				results[idx] = result{
					index: idx,
					sec: renderer.Section{
						Title: e.Title,
						Lines: []string{"⚠ Unknown module: " + e.Module},
					},
				}
				return
			}

			modCfg := buildModuleConfig(cfg, e)
			lines := mod.GetData(modCfg)

			if e.Module == "quotes" {
				quote := ""
				if len(lines) > 0 {
					quote = lines[0]
				}
				results[idx] = result{index: idx, quote: quote}
				return
			}

			results[idx] = result{
				index: idx,
				sec: renderer.Section{
					Title: e.Title,
					Lines: lines,
				},
			}
		}(i, entry)
	}

	wg.Wait()

	var sections []renderer.Section
	quote := ""
	for _, r := range results {
		if r.quote != "" {
			quote = r.quote
		} else if r.sec.Title != "" || len(r.sec.Lines) > 0 {
			sections = append(sections, r.sec)
		}
	}
	return sections, quote
}

func buildModuleConfig(cfg *config.Config, e config.LayoutEntry) map[string]any {
	mc := make(map[string]any)
	for k, v := range e.Config {
		mc[k] = v
	}
	if cfg.Settings.MasterCurrency != "" {
		if _, ok := mc["master_currency"]; !ok {
			mc["master_currency"] = cfg.Settings.MasterCurrency
		}
	}
	if lat, ok := cfg.User.Extra["lat"]; ok {
		if _, exists := mc["lat"]; !exists {
			mc["lat"] = lat
		}
	}
	if lon, ok := cfg.User.Extra["lon"]; ok {
		if _, exists := mc["lon"]; !exists {
			mc["lon"] = lon
		}
	}
	if e.Module == "news" && len(e.Feeds) > 0 {
		mc["feeds"] = e.Feeds
	}
	return mc
}
