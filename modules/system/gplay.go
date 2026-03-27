package system

import (
	"fmt"
	"html"
	"io"
	"net/http"
	"os"
	"regexp"
)

type GooglePlay struct{}

func (GooglePlay) GetData(cfg map[string]any) []string {
	pkg := ""
	if v, ok := cfg["package_name"].(string); ok && v != "" {
		pkg = v
	} else {
		pkg = os.Getenv("PACKAGE_NAME")
	}

	if pkg == "" {
		return []string{"Google Play: Package name not configured in layout or PACKAGE_NAME env"}
	}

	u := fmt.Sprintf("https://play.google.com/store/apps/details?id=%s&hl=en&gl=us", pkg)
	req, _ := http.NewRequest("GET", u, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) Chrome/120.0.0.0 Safari/537.36")

	resp, err := httpClient.Do(req)
	if err != nil {
		return []string{fmt.Sprintf("Google Play Error: %v", err)}
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return []string{fmt.Sprintf("Google Play (%s): Not found or rate limited", pkg)}
	}

	bodyBytes, _ := io.ReadAll(resp.Body)
	body := string(bodyBytes)

	installsRe := regexp.MustCompile(`"([0-9,]+\+)",[0-9]+,[0-9]+,"[0-9A-Za-z\+\.]+"`)
	matches := installsRe.FindStringSubmatch(body)
	installs := "Unknown"
	if len(matches) > 1 {
		installs = matches[1]
	} else {
		jsRe := regexp.MustCompile(`"([0-9,]+\+)","[0-9]+ downloads"`)
		matches = jsRe.FindStringSubmatch(body)
		if len(matches) > 1 {
			installs = matches[1]
		} else {
		    // fallback check
		    oldRe := regexp.MustCompile(`>([0-9,]+\+)<[^>]*>Downloads<`)
		    matches = oldRe.FindStringSubmatch(body)
		    if len(matches) > 1 {
		        installs = matches[1]
		    }
		}
	}

	titleRe := regexp.MustCompile(`itemprop="name"[^>]*>(.*?)</span>`)
	tMatch := titleRe.FindStringSubmatch(body)
	title := "App"
	if len(tMatch) > 1 {
		title = html.UnescapeString(tMatch[1])
	}

	return []string{fmt.Sprintf("%s Downloads: %s", title, installs)}
}
