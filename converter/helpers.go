package converter

import (
	"strconv"
	"strings"

	"golang.org/x/net/html"
)

// GetAttrMap converts HTML attributes to a map for easy lookup.
func GetAttrMap(attrs []html.Attribute) map[string]string {
	m := make(map[string]string)
	for _, a := range attrs {
		m[strings.ToLower(a.Key)] = a.Val
	}
	return m
}

// GetAttrValue returns the value of an HTML attribute by key.
func GetAttrValue(attrs []html.Attribute, key string) string {
	for _, attr := range attrs {
		if strings.ToLower(attr.Key) == key {
			return attr.Val
		}
	}
	return ""
}

// ExtractText recursively extracts all text content from an HTML node.
func ExtractText(n *html.Node) string {
	if n.Type == html.TextNode {
		return n.Data
	}
	var sb strings.Builder
	for ch := n.FirstChild; ch != nil; ch = ch.NextSibling {
		sb.WriteString(ExtractText(ch))
	}
	return sb.String()
}

// CollapseWhitespace replaces runs of whitespace with a single space,
// preserving leading/trailing space if present in the original.
func CollapseWhitespace(s string) string {
	hasLeading := len(s) > 0 && (s[0] == ' ' || s[0] == '\t' || s[0] == '\n' || s[0] == '\r')
	hasTrailing := len(s) > 1 && (s[len(s)-1] == ' ' || s[len(s)-1] == '\t' || s[len(s)-1] == '\n' || s[len(s)-1] == '\r')

	words := strings.Fields(s)
	result := strings.Join(words, " ")
	if hasLeading {
		result = " " + result
	}
	if hasTrailing {
		result = result + " "
	}
	return result
}

// ParseHexToRGB converts a hex color string to RGB int values.
func ParseHexToRGB(hex string) (int, int, int) {
	hex = strings.TrimPrefix(hex, "#")
	if len(hex) != 6 {
		return 0, 0, 0
	}
	r, _ := strconv.ParseInt(hex[0:2], 16, 32)
	g, _ := strconv.ParseInt(hex[2:4], 16, 32)
	b, _ := strconv.ParseInt(hex[4:6], 16, 32)
	return int(r), int(g), int(b)
}
