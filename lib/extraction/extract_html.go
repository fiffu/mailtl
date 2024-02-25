package extraction

import (
	"strings"

	"github.com/antchfx/htmlquery"
)

func ExtractXPath(xml string, xpath string) (string, error) {
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
