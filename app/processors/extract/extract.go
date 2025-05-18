package extract

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/antchfx/htmlquery"
	"github.com/fiffu/mailtl/lib"
	"github.com/flashmob/go-guerrilla/mail"
)

// Common API for extraction

var ErrNoMatches = fmt.Errorf("no data was extracted")

func errNoMatchesf(reason string, args ...any) error {
	reason = fmt.Sprintf(reason, args...)
	return fmt.Errorf("%w (%s)", ErrNoMatches, reason)
}

func Fingerprint(envelope *mail.Envelope) string {
	rcptTo := strings.Join(lib.StringsOf(lib.IndirectsOf(envelope.RcptTo)), "")
	s := strings.Join(
		[]string{
			envelope.MailFrom.String(),
			rcptTo,
			envelope.Subject,
		},
		"",
	)
	h := sha256.New()
	h.Write([]byte(s))
	digest := h.Sum(nil)
	return hex.EncodeToString(digest)
}

func XPath(xml string, xpath string) (string, error) {
	doc, err := htmlquery.Parse(strings.NewReader(xml))
	if err != nil {
		return "", errNoMatchesf("parser error: %v", err)
	}

	nodes, err := htmlquery.QueryAll(doc, xpath)
	if len(nodes) == 0 {
		return "", errNoMatchesf("no nodes matched by xpath: '%s'", xpath)
	}
	if err != nil {
		return "", err
	}

	return nodes[0].FirstChild.Data, nil
}

func Markup(e *mail.Envelope) string {
	// replace email delimiters
	data := e.Data.String()
	return strings.ReplaceAll(data, "=\n", "")
}
