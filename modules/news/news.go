package news

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"
)

var httpClient = &http.Client{Timeout: 12 * time.Second}
var htmlTagRe = regexp.MustCompile(`<[^>]+>`)

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

func truncate(s string, n int) string {
	runes := []rune(s)
	if len(runes) <= n {
		return s
	}
	return string(runes[:n-3]) + "..."
}

type HackerNews struct{}

type hnItem struct {
	Title string `json:"title"`
	Score int    `json:"score"`
	URL   string `json:"url"`
}

func (HackerNews) GetData(cfg map[string]any) []string {
	count := 5
	if v, ok := cfg["count"].(int); ok {
		count = v
	} else if v, ok := cfg["count"].(float64); ok {
		count = int(v)
	}

	var ids []int
	if err := getJSON("https://hacker-news.firebaseio.com/v0/topstories.json", nil, &ids); err != nil {
		return []string{"Hacker News Error: " + err.Error()}
	}
	if len(ids) > count {
		ids = ids[:count]
	}

	type pair struct {
		idx  int
		line string
	}
	ch := make(chan pair, count)
	for i, id := range ids {
		go func(i, id int) {
			var item hnItem
			u := fmt.Sprintf("https://hacker-news.firebaseio.com/v0/item/%d.json", id)
			if err := getJSON(u, nil, &item); err == nil && item.Title != "" {
				ch <- pair{i, fmt.Sprintf("[%d] %s", item.Score, item.Title)}
			} else {
				ch <- pair{i, ""}
			}
		}(i, id)
	}

	results := make([]string, count)
	for range ids {
		p := <-ch
		results[p.idx] = p.line
	}
	var out []string
	for _, l := range results {
		if l != "" {
			out = append(out, l)
		}
	}
	return out
}

type GitHubTrending struct{}

type ghRepo struct {
	FullName    string `json:"full_name"`
	Description string `json:"description"`
	Stars       int    `json:"stargazers_count"`
}

type ghSearchResp struct {
	Items []ghRepo `json:"items"`
}

func (GitHubTrending) GetData(cfg map[string]any) []string {
	count := 5
	if v, ok := cfg["count"].(float64); ok {
		count = int(v)
	}
	lang := ""
	if v, ok := cfg["language"].(string); ok {
		lang = v
	}

	since := time.Now().AddDate(0, 0, -7).Format("2006-01-02")
	query := "created:>" + since
	if lang != "" {
		query += " language:" + lang
	}

	apiURL := "https://api.github.com/search/repositories?q=" +
		url.QueryEscape(query) + "&sort=stars&order=desc"

	var resp ghSearchResp
	if err := getJSON(apiURL, map[string]string{
		"Accept": "application/vnd.github.v3+json",
	}, &resp); err != nil {
		return []string{"GitHub Trending Error: " + err.Error()}
	}

	var out []string
	for i, item := range resp.Items {
		if i >= count {
			break
		}
		out = append(out, fmt.Sprintf("★%d  %s", item.Stars, item.FullName))
		if item.Description != "" {
			out = append(out, "  "+truncate(item.Description, 80))
		}
	}
	return out
}

var predefinedFeeds = map[string]string{
	"bbc":         "http://feeds.bbci.co.uk/news/world/rss.xml",
	"aljazeera":   "https://www.aljazeera.com/xml/rss/all.xml",
	"ap":          "https://apnews.com/feed",
	"npr":         "https://www.npr.org/rss/rss.php?id=1001",
	"nyt":         "https://rss.nytimes.com/services/xml/rss/nyt/HomePage.xml",
	"dw":          "https://rss.dw.com/rdf/rss-en-all",
	"france24":    "https://www.france24.com/en/rss",
	"euronews":    "https://www.euronews.com/rss?format=google-news&level=theme&name=news",
	"guardian":    "https://www.theguardian.com/world/rss",
	"cna":         "https://www.channelnewsasia.com/rss/news/asia/rss.xml",
	"hindu":       "https://www.thehindu.com/news/national/feeder/default.rss",
	"toi":         "https://timesofindia.indiatimes.com/rssfeeds/-2128936835.cms",
	"ndtv":        "https://feeds.feedburner.com/ndtvnews-top-stories",
	"scmp":        "https://www.scmp.com/rss/2/feed.xml",
	"abc_au":      "https://www.abc.net.au/news/feed/2942460/rss.xml",
	"rnz":         "https://www.rnz.co.nz/rss/news.xml",
	"mercopress":  "https://en.mercopress.com/rss/",
	"batimes":     "https://www.batimes.com.ar/rss",
	"techcrunch":  "https://techcrunch.com/feed/",
	"verge":       "https://www.theverge.com/rss/index.xml",
	"theverge":    "https://www.theverge.com/rss/index.xml",
	"wired":       "https://www.wired.com/feed/rss",
	"arstechnica": "https://feeds.arstechnica.com/arstechnica/index",
	"theregister": "https://www.theregister.com/headlines.rss",
	"hackernews":  "https://news.ycombinator.com/rss",
	"noted":       "https://noted.lol/rss/",
	"selfh_st":    "https://selfh.st/rss/",
	"yourstory":   "https://yourstory.com/feed",
	"technode":    "https://technode.com/feed/",
}

type News struct{}

type rssRoot struct {
	XMLName xml.Name   `xml:"rss"`
	Channel rssChannel `xml:"channel"`
}

type atomRoot struct {
	XMLName xml.Name    `xml:"feed"`
	Entries []atomEntry `xml:"entry"`
}

type rssChannel struct {
	Items []rssItem `xml:"item"`
}

type rssItem struct {
	Title       string `xml:"title"`
	Description string `xml:"description"`
}

type atomEntry struct {
	Title   string `xml:"title"`
	Summary string `xml:"summary"`
}

func (News) GetData(cfg map[string]any) []string {
	feedsRaw, _ := cfg["feeds"].([]any)
	var out []string

	for _, f := range feedsRaw {
		entry, ok := f.(map[string]any)
		if !ok {
			continue
		}

		feedKey, _ := entry["feed"].(string)
		subreddit, _ := entry["subreddit"].(string)
		customURL, _ := entry["url"].(string)
		count := 1
		if v, ok := entry["count"].(float64); ok {
			count = int(v)
		}

		var feedURL, displayTitle string
		switch {
		case subreddit != "":
			feedURL = "https://www.reddit.com/r/" + subreddit + "/top.rss?t=day"
			displayTitle = "r/" + subreddit
		case feedKey != "" && predefinedFeeds[strings.ToLower(feedKey)] != "":
			feedURL = predefinedFeeds[strings.ToLower(feedKey)]
			displayTitle = strings.ToUpper(feedKey)
		case customURL != "":
			feedURL = customURL
			displayTitle, _ = entry["title"].(string)
			if displayTitle == "" {
				displayTitle = "Feed"
			}
		default:
			out = append(out, fmt.Sprintf("Feed error: Unknown key '%s'", feedKey))
			continue
		}

		items, err := fetchFeedItems(feedURL, count)
		if err != nil {
			out = append(out, fmt.Sprintf("Feed error (%s): %s", displayTitle, err))
			continue
		}
		out = append(out, items...)
	}
	return out
}

func fetchFeedItems(feedURL string, count int) ([]string, error) {
	req, _ := http.NewRequest("GET", feedURL, nil)
	req.Header.Set("User-Agent", "DPTR/1.0")
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var rss rssRoot
	if err := xml.Unmarshal(body, &rss); err == nil && len(rss.Channel.Items) > 0 {
		return extractRSS(rss.Channel.Items, count), nil
	}

	var atom atomRoot
	if err := xml.Unmarshal(body, &atom); err == nil && len(atom.Entries) > 0 {
		return extractAtom(atom.Entries, count), nil
	}

	return nil, fmt.Errorf("could not parse feed")
}

func extractRSS(items []rssItem, count int) []string {
	var out []string
	for i, item := range items {
		if i >= count {
			break
		}
		t := cleanText(item.Title)
		if t != "" {
			out = append(out, ">> "+truncate(t, 110))
		}
		d := cleanText(item.Description)
		if d != "" {
			out = append(out, "   "+truncate(d, 120))
		}
	}
	return out
}

func extractAtom(entries []atomEntry, count int) []string {
	var out []string
	for i, e := range entries {
		if i >= count {
			break
		}
		t := cleanText(e.Title)
		if t != "" {
			out = append(out, ">> "+truncate(t, 110))
		}
		d := cleanText(e.Summary)
		if d != "" {
			out = append(out, "   "+truncate(d, 120))
		}
	}
	return out
}

func cleanText(s string) string {
	s = html.UnescapeString(s)
	s = htmlTagRe.ReplaceAllString(s, "")
	return strings.Join(strings.Fields(s), " ")
}
