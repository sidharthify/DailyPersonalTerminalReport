package system

import (
	"crypto/tls"
	"fmt"
	"mime"
	"net/mail"
	"os"
	"strings"

	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	"github.com/emersion/go-message/charset"
	_ "github.com/emersion/go-message/charset"
)

type EmailInbox struct{}

func (EmailInbox) GetData(cfg map[string]any) []string {
	host := os.Getenv("EMAIL_HOST")
	if host == "" {
		host = "imap.gmail.com"
	}
	user := os.Getenv("EMAIL_USER")
	pass := os.Getenv("EMAIL_PASS")

	if user == "" || pass == "" {
		return []string{"Email Error: Missing EMAIL_USER or EMAIL_PASS in env"}
	}

	count := 3
	if v, ok := cfg["count"].(float64); ok {
		count = int(v)
	}

	if !strings.Contains(host, ":") {
		host += ":993"
	}

	c, err := client.DialTLS(host, &tls.Config{})
	if err != nil {
		return []string{"IMAP Dial Error: " + err.Error()}
	}
	defer c.Logout()

	if err := c.Login(user, pass); err != nil {
		return []string{"IMAP Login Error: " + err.Error()}
	}

	_, err = c.Select("INBOX", false)
	if err != nil {
		return []string{"IMAP Select Error: " + err.Error()}
	}

	criteria := imap.NewSearchCriteria()
	criteria.WithoutFlags = []string{imap.SeenFlag}
	uids, err := c.Search(criteria)
	if err != nil {
		return []string{"IMAP Search Error: " + err.Error()}
	}

	if len(uids) == 0 {
		return []string{"Inbox is empty (no unseen)."}
	}

	start := len(uids) - count
	if start < 0 {
		start = 0
	}
	fetchSet := new(imap.SeqSet)
	for _, uid := range uids[start:] {
		fetchSet.AddNum(uid)
	}

	messages := make(chan *imap.Message, 10)
	done := make(chan error, 1)

	section := &imap.BodySectionName{Peek: true}
	go func() {
		done <- c.Fetch(fetchSet, []imap.FetchItem{section.FetchItem()}, messages)
	}()

	var out []string
	dec := new(mime.WordDecoder)
	dec.CharsetReader = charset.Reader

	for msg := range messages {
		r := msg.GetBody(section)
		if r == nil {
			continue
		}
		m, err := mail.ReadMessage(r)
		if err != nil {
			continue
		}

		subject, _ := dec.DecodeHeader(m.Header.Get("Subject"))
		from, _ := dec.DecodeHeader(m.Header.Get("From"))

		sender := strings.Split(from, "<")[0]
		sender = strings.TrimSpace(strings.ReplaceAll(sender, "\"", ""))

		out = append([]string{fmt.Sprintf("%s: %s", sender, subject)}, out...)
	}

	if err := <-done; err != nil {
		return append(out, "IMAP Fetch Error: "+err.Error())
	}

	if len(out) == 0 {
		return []string{"Inbox is empty."}
	}
	return out
}
