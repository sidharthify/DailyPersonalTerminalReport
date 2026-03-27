package renderer

import (
	"fmt"
	"strings"
	"time"
	"unicode/utf8"
)

const (
	reset   = "\033[0m"
	bold    = "\033[1m"
	dim     = "\033[2m"
	cyan    = "\033[36m"
	yellow  = "\033[33m"
	green   = "\033[32m"
	magenta = "\033[35m"
	white   = "\033[97m"
	bgBlue  = "\033[44m"
)

type Section struct {
	Title string
	Lines []string
}

func RenderReport(userName, greeting string, sections []Section, quote string) {
	width := 72

	printHeader(userName, greeting, width)

	for i, sec := range sections {
		printSection(sec, width)
		_ = i
	}

	if quote != "" {
		printQuote(quote, width)
	}

	printFooter(width)
}

func printHeader(name, greeting string, width int) {
	bar := strings.Repeat("═", width)

	fmt.Printf("\n%s%s%s\n", cyan, bar, reset)

	title := fmt.Sprintf(" %s%s, %s%s%s ", bold+white, greeting, yellow, name, reset)
	printCentered(title, width, bold+bgBlue+white, reset)

	date := time.Now().Format("Monday, 02 January 2006  •  15:04")
	subtitle := fmt.Sprintf("%s%s%s", dim, date, reset)
	printCentered(subtitle, width, "", "")

	fmt.Printf("%s%s%s\n\n", cyan, bar, reset)
}

func printSection(sec Section, width int) {
	if sec.Title != "" {
		bar := strings.Repeat("─", width)
		fmt.Printf("%s%s%s\n", dim, bar, reset)
		fmt.Printf("%s%s  %s  %s\n", bold+magenta, "▌", strings.ToUpper(sec.Title), reset)
	}

	for _, line := range sec.Lines {
		for _, wrapped := range wrapLine(line, width-2) {
			fmt.Printf("  %s\n", wrapped)
		}
	}
	fmt.Println()
}

func printQuote(quote string, width int) {
	bar := strings.Repeat("─", width)
	fmt.Printf("%s%s%s\n", dim, bar, reset)
	fmt.Printf("%s%s  %s  %s\n", bold+green, "▌", "QUOTE OF THE DAY", reset)
	for _, line := range wrapLine(quote, width-4) {
		fmt.Printf("  %s%s%s\n", dim+white, line, reset)
	}
	fmt.Println()
}

func printFooter(width int) {
	bar := strings.Repeat("═", width)
	fmt.Printf("%s%s%s\n", cyan, bar, reset)
	msg := fmt.Sprintf("%s DPTR — Daily Personal Terminal Report %s", dim, reset)
	printCentered(msg, width, "", "")
	fmt.Printf("%s%s%s\n\n", cyan, bar, reset)
}

func printCentered(text string, width int, prefix, suffix string) {
	visual := stripANSI(text)
	pad := (width - utf8.RuneCountInString(visual)) / 2
	if pad < 0 {
		pad = 0
	}
	fmt.Printf("%s%s%s%s\n", prefix, strings.Repeat(" ", pad), text, suffix)
}

func wrapLine(text string, maxWidth int) []string {
	if maxWidth <= 0 {
		maxWidth = 70
	}
	if utf8.RuneCountInString(text) <= maxWidth {
		return []string{text}
	}
	var lines []string
	words := strings.Fields(text)
	cur := ""
	for _, w := range words {
		if cur == "" {
			cur = w
		} else if utf8.RuneCountInString(cur)+1+utf8.RuneCountInString(w) <= maxWidth {
			cur += " " + w
		} else {
			lines = append(lines, cur)
			cur = w
		}
	}
	if cur != "" {
		lines = append(lines, cur)
	}
	return lines
}

func stripANSI(s string) string {
	var out strings.Builder
	inEsc := false
	for _, r := range s {
		if r == '\033' {
			inEsc = true
		}
		if !inEsc {
			out.WriteRune(r)
		}
		if inEsc && r == 'm' {
			inEsc = false
		}
	}
	return out.String()
}
