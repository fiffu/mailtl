package extract

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/antchfx/htmlquery"
	"github.com/fiffu/mailtl/lib"
	"github.com/flashmob/go-guerrilla/mail"
	"github.com/oriser/regroup"
	"golang.org/x/net/html"
)

// Common API for extraction

var ErrNoMatches = fmt.Errorf("no data was extracted")

var symbolToCurrency = map[string]string{
	"S$": "SGD",
}

type currencySymbolParser struct {
	Symbol string  `regroup:"symbol"`
	Amount float64 `regroup:"amount"`
}

func (c *currencySymbolParser) Match(text string) (string, float64, error) {
	patt := regroup.MustCompile(`(?P<symbol>[A-Z\W]{1,3}) ?(?P<amount>[\d\.]+)`)
	if err := patt.MatchToTarget(text, c); err != nil {
		return "", 0, err
	}
	return c.Symbol, c.Amount, nil
}

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
	res, err := XPathEach(xml, xpath, func(n *html.Node) string {
		return n.FirstChild.Data
	})
	if err != nil {
		return "", err
	}
	return res[0], nil
}

func XPathEach(xml string, xpath string, callback func(*html.Node) string) ([]string, error) {
	doc, err := htmlquery.Parse(strings.NewReader(xml))
	if err != nil {
		return nil, errNoMatchesf("parser error: %v", err)
	}

	nodes, err := htmlquery.QueryAll(doc, xpath)
	if len(nodes) == 0 {
		return nil, errNoMatchesf("no nodes matched by xpath: '%s'", xpath)
	}
	if err != nil {
		return nil, err
	}

	return lib.Map(nodes, callback), nil
}

func Markup(e *mail.Envelope) string {
	// replace email delimiters
	data := e.Data.String()
	return strings.ReplaceAll(data, "=\n", "")
}

func ParseMoney(text string) (currency string, amount float64, err error) {
	symbol, amount, err := new(currencySymbolParser).Match(text)
	if err != nil {
		return "", 0, err
	}
	if c, ok := symbolToCurrency[symbol]; ok {
		currency = c
	} else {
		currency = symbol
	}
	return
}
