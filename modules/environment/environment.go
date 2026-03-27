package environment

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

var httpClient = &http.Client{Timeout: 10 * time.Second}

func getJSON(url string, target any) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("User-Agent", "DPTR/1.0")
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

func strOrDefault(cfg map[string]any, key, def string) string {
	if v, ok := cfg[key]; ok && v != nil {
		return fmt.Sprintf("%v", v)
	}
	return def
}

type Weather struct{}

func (Weather) GetData(cfg map[string]any) []string {
	lat := strOrDefault(cfg, "lat", "")
	lon := strOrDefault(cfg, "lon", "")
	if lat == "" || lon == "" || lat == "<nil>" || lon == "<nil>" {
		return []string{"Weather: Missing coordinates (lat/lon)"}
	}

	url := fmt.Sprintf(
		"https://api.open-meteo.com/v1/forecast?latitude=%s&longitude=%s"+
			"&current=temperature_2m,relative_humidity_2m,weather_code,wind_speed_10m"+
			"&daily=temperature_2m_max,temperature_2m_min&timezone=auto", lat, lon)

	var data map[string]any
	if err := getJSON(url, &data); err != nil {
		return []string{"Weather Error: " + err.Error()}
	}

	current, _ := data["current"].(map[string]any)
	daily, _ := data["daily"].(map[string]any)

	temp := fmtNum(current["temperature_2m"])
	hum := fmtNum(current["relative_humidity_2m"])
	wind := fmtNum(current["wind_speed_10m"])

	maxSlice, _ := daily["temperature_2m_max"].([]any)
	minSlice, _ := daily["temperature_2m_min"].([]any)
	maxT, minT := "N/A", "N/A"
	if len(maxSlice) > 0 {
		maxT = fmtNum(maxSlice[0])
	}
	if len(minSlice) > 0 {
		minT = fmtNum(minSlice[0])
	}

	return []string{
		fmt.Sprintf("Temperature : %s°C  (Lo %s°  Hi %s°)", temp, minT, maxT),
		fmt.Sprintf("Humidity    : %s%%  |  Wind: %s km/h", hum, wind),
	}
}

type Astronomy struct{}

func (Astronomy) GetData(cfg map[string]any) []string {
	lat := strOrDefault(cfg, "lat", "")
	lon := strOrDefault(cfg, "lon", "")
	if lat == "" || lon == "" || lat == "<nil>" {
		return []string{"Astronomy: Missing coordinates (lat/lon)"}
	}

	url := fmt.Sprintf(
		"https://api.open-meteo.com/v1/forecast?latitude=%s&longitude=%s"+
			"&daily=sunrise,sunset&timezone=auto", lat, lon)

	var data map[string]any
	if err := getJSON(url, &data); err != nil {
		return []string{"Astronomy Error: " + err.Error()}
	}

	daily, _ := data["daily"].(map[string]any)
	sunriseSlice, _ := daily["sunrise"].([]any)
	sunsetSlice, _ := daily["sunset"].([]any)

	sunrise, sunset := "N/A", "N/A"
	if len(sunriseSlice) > 0 {
		s := fmt.Sprintf("%v", sunriseSlice[0])
		if t, err := time.Parse("2006-01-02T15:04", s); err == nil {
			sunrise = t.Format("15:04")
		} else {
			sunrise = s
		}
	}
	if len(sunsetSlice) > 0 {
		s := fmt.Sprintf("%v", sunsetSlice[0])
		if t, err := time.Parse("2006-01-02T15:04", s); err == nil {
			sunset = t.Format("15:04")
		} else {
			sunset = s
		}
	}

	return []string{
		fmt.Sprintf("Sunrise : %s", sunrise),
		fmt.Sprintf("Sunset  : %s", sunset),
	}
}

type OpenMeteoAQI struct{}

func (OpenMeteoAQI) GetData(cfg map[string]any) []string {
	lat := strOrDefault(cfg, "lat", "")
	lon := strOrDefault(cfg, "lon", "")
	if lat == "" || lon == "" || lat == "<nil>" {
		return []string{"AQI: Missing coordinates (lat/lon)"}
	}

	url := fmt.Sprintf(
		"https://air-quality-api.open-meteo.com/v1/air-quality"+
			"?latitude=%s&longitude=%s&current=us_aqi,pm10,pm2_5", lat, lon)

	var data map[string]any
	if err := getJSON(url, &data); err != nil {
		return []string{"AQI Error: " + err.Error()}
	}

	current, _ := data["current"].(map[string]any)
	usAQI := current["us_aqi"]
	pm25 := fmtNum(current["pm2_5"])
	pm10 := fmtNum(current["pm10"])

	status := aqiStatus(usAQI)
	return []string{
		fmt.Sprintf("US AQI : %s (%s)  |  PM2.5: %s  |  PM10: %s", fmtNum(usAQI), status, pm25, pm10),
	}
}

func aqiStatus(v any) string {
	f, ok := toFloat(v)
	if !ok {
		return "Unknown"
	}
	switch {
	case f <= 50:
		return "Good"
	case f <= 100:
		return "Moderate"
	case f <= 150:
		return "Unhealthy (Sensitive)"
	case f <= 200:
		return "Unhealthy"
	case f <= 300:
		return "Very Unhealthy"
	default:
		return "Hazardous"
	}
}

type BreatheOSS struct{}

func (BreatheOSS) GetData(cfg map[string]any) []string {
	city := strOrDefault(cfg, "city", "Srinagar")

	path := strings.ToLower(city)
	if path == "jammu" {
		path = "jammu_city"
	}

	url := fmt.Sprintf("https://api.breatheoss.app/aqi/%s", path)
	var data map[string]any
	if err := getJSON(url, &data); err != nil {
		return []string{fmt.Sprintf("%s — AQI: Currently Unavailable", city)}
	}

	aqi := fmtNum(data["aqi"])

	pm25 := "N/A"
	if conc, ok := data["concentrations_raw_ugm3"].(map[string]any); ok {
		pm25 = fmtNum(conc["pm2_5"])
	} else {
		pm25 = fmtNum(data["pm2_5"]) // fallback
	}

	status := aqiStatus(data["aqi"])

	return []string{
		fmt.Sprintf("%s — AQI: %s (%s)  |  PM2.5: %s", city, aqi, status, pm25),
	}
}

func fmtNum(v any) string {
	if v == nil {
		return "N/A"
	}
	switch n := v.(type) {
	case float64:
		if n == float64(int(n)) {
			return fmt.Sprintf("%d", int(n))
		}
		return fmt.Sprintf("%.2f", n)
	case json.Number:
		return n.String()
	default:
		return fmt.Sprintf("%v", v)
	}
}

func toFloat(v any) (float64, bool) {
	if v == nil {
		return 0, false
	}
	switch n := v.(type) {
	case float64:
		return n, true
	case json.Number:
		f, err := n.Float64()
		return f, err == nil
	case int:
		return float64(n), true
	}
	return 0, false
}
