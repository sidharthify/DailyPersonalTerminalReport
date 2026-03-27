package random

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"
)

var httpClient = &http.Client{Timeout: 10 * time.Second}

func getJSON(u string, headers map[string]string, target any) error {
	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return err
	}
	req.Header.Set("User-Agent", "DPTR/1.0")
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	resp, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	return json.Unmarshal(body, target)
}

func strVal(cfg map[string]any, key, def string) string {
	if v, ok := cfg[key]; ok && v != nil {
		return fmt.Sprintf("%v", v)
	}
	return def
}

type DailyJoke struct{}

func (DailyJoke) GetData(_ map[string]any) []string {
	var result struct {
		Joke string `json:"joke"`
	}
	if err := getJSON("https://icanhazdadjoke.com/", map[string]string{
		"Accept": "application/json",
	}, &result); err != nil {
		return []string{"Joke Error: " + err.Error()}
	}
	joke := result.Joke
	lines := strings.Split(joke, "\r\n")
	if len(lines) == 1 {
		lines = strings.Split(joke, "\n")
	}
	var out []string
	for _, l := range lines {
		if l = strings.TrimSpace(l); l != "" {
			out = append(out, l)
		}
	}
	return out
}

type FactOfTheDay struct{}

func (FactOfTheDay) GetData(_ map[string]any) []string {
	var result struct {
		Text string `json:"text"`
	}
	if err := getJSON("https://uselessfacts.jsph.pl/api/v2/facts/random?language=en",
		nil, &result); err != nil {
		return []string{"Fact Error: " + err.Error()}
	}
	return []string{result.Text}
}

type ShowerThoughts struct{}

type redditPost struct {
	Data struct {
		Title string `json:"title"`
	} `json:"data"`
}

type redditListing struct {
	Data struct {
		Children []redditPost `json:"children"`
	} `json:"data"`
}

func (ShowerThoughts) GetData(cfg map[string]any) []string {
	count := 3
	if v, ok := cfg["count"].(float64); ok {
		count = int(v)
	}

	var listing redditListing
	u := fmt.Sprintf(
		"https://www.reddit.com/r/showerthoughts/top.json?t=day&limit=%d", count)
	if err := getJSON(u, map[string]string{
		"User-Agent": "dptr/1.0",
	}, &listing); err != nil {
		return []string{"Shower Thoughts Error: " + err.Error()}
	}
	var out []string
	for i, child := range listing.Data.Children {
		if i >= count {
			break
		}
		out = append(out, "💭 "+child.Data.Title)
	}
	return out
}

type WordOfTheDay struct{}

func (WordOfTheDay) GetData(_ map[string]any) []string {
	words := []struct{ Word, Def string }{
		{"Ephemeral", "Lasting for a very short time; transitory."},
		{"Sonder", "The realisation that each passerby has a life as vivid as one's own."},
		{"Hiraeth", "A homesickness for a home you cannot return to, or never had."},
		{"Serendipity", "The occurrence of events by pleasant chance."},
		{"Laconic", "Using very few words; brief and concise."},
		{"Petrichor", "The pleasant smell that accompanies the first rain after a dry spell."},
		{"Melancholy", "A feeling of pensive sadness, typically with no obvious cause."},
		{"Eloquent", "Fluent or persuasive in speaking or writing."},
		{"Cogent", "Clear, logical, and convincing."},
		{"Loquacious", "Tending to talk a great deal; talkative."},
		{"Querulous", "Complaining in a petulant or whining manner."},
		{"Perspicacious", "Having a ready insight into things; shrewd."},
		{"Sesquipedalian", "Given to using long words; characterised by long words."},
		{"Sycophant", "A person who acts obsequiously toward someone to gain advantage."},
		{"Penultimate", "Last but one in a series."},
	}
	dayOfYear := time.Now().YearDay()
	w := words[dayOfYear%len(words)]
	return []string{
		fmt.Sprintf("Word : %s", w.Word),
		fmt.Sprintf("Def  : %s", w.Def),
	}
}

type Quotes struct{}

type quoteEntry struct {
	Text        string `json:"text"`
	Attribution string `json:"attribution"`
}

func (Quotes) GetData(cfg map[string]any) []string {
	quotesFile := strVal(cfg, "quotes_file", "quotes.json")

	data, err := os.ReadFile(quotesFile)
	if err != nil {
		return []string{"-- No quotes file found --"}
	}
	var quotes []quoteEntry
	if err := json.Unmarshal(data, &quotes); err != nil || len(quotes) == 0 {
		return []string{"-- quotes.json parse error --"}
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	q := quotes[r.Intn(len(quotes))]

	if q.Attribution != "" {
		return []string{fmt.Sprintf("\"%s\"  — %s", q.Text, q.Attribution)}
	}
	return []string{fmt.Sprintf("\"%s\"", q.Text)}
}
