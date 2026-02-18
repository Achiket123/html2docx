package converter

import (
	"strings"
	"testing"

	"golang.org/x/net/html"
)

func TestGetAttrMap(t *testing.T) {
	attrs := []html.Attribute{
		{Key: "border", Val: "1"},
		{Key: "Width", Val: "100%"},
		{Key: "ALIGN", Val: "center"},
	}
	m := GetAttrMap(attrs)

	if m["border"] != "1" {
		t.Errorf("expected border=1, got %q", m["border"])
	}
	if m["width"] != "100%" {
		t.Errorf("expected width=100%%, got %q", m["width"])
	}
	if m["align"] != "center" {
		t.Errorf("expected align=center, got %q", m["align"])
	}
}

func TestGetAttrMapEmpty(t *testing.T) {
	m := GetAttrMap(nil)
	if len(m) != 0 {
		t.Errorf("expected empty map, got %d entries", len(m))
	}
}

func TestGetAttrValue(t *testing.T) {
	attrs := []html.Attribute{
		{Key: "href", Val: "https://example.com"},
		{Key: "Class", Val: "link"},
	}

	if v := GetAttrValue(attrs, "href"); v != "https://example.com" {
		t.Errorf("expected href value, got %q", v)
	}
	if v := GetAttrValue(attrs, "class"); v != "link" {
		t.Errorf("expected class value, got %q", v)
	}
	if v := GetAttrValue(attrs, "missing"); v != "" {
		t.Errorf("expected empty for missing attr, got %q", v)
	}
}

func TestExtractText(t *testing.T) {
	tests := []struct {
		name     string
		html     string
		expected string
	}{
		{
			name:     "plain text",
			html:     "<p>Hello World</p>",
			expected: "Hello World",
		},
		{
			name:     "nested elements",
			html:     "<p>Hello <b>bold</b> text</p>",
			expected: "Hello bold text",
		},
		{
			name:     "deeply nested",
			html:     "<div><p><span>Deep</span></p></div>",
			expected: "Deep",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			doc, _ := html.Parse(strings.NewReader(tc.html))
			// Find the body element
			var body *html.Node
			var findBody func(*html.Node)
			findBody = func(n *html.Node) {
				if n.Type == html.ElementNode && n.Data == "body" {
					body = n
					return
				}
				for c := n.FirstChild; c != nil; c = c.NextSibling {
					findBody(c)
				}
			}
			findBody(doc)
			if body == nil || body.FirstChild == nil {
				t.Fatal("could not find body element")
			}
			got := strings.TrimSpace(ExtractText(body.FirstChild))
			if got != tc.expected {
				t.Errorf("expected %q, got %q", tc.expected, got)
			}
		})
	}
}

func TestCollapseWhitespace(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"hello   world", "hello world"},
		{"  hello world", " hello world"},
		{"hello world  ", "hello world "},
		{"  hello   world  ", " hello world "},
		{"\n  hello  \t  world  \n", " hello world "},
		{"singleword", "singleword"},
	}

	for _, tc := range tests {
		got := CollapseWhitespace(tc.input)
		if got != tc.expected {
			t.Errorf("CollapseWhitespace(%q) = %q, want %q", tc.input, got, tc.expected)
		}
	}
}

func TestParseHexToRGB(t *testing.T) {
	tests := []struct {
		hex     string
		r, g, b int
	}{
		{"#FF0000", 255, 0, 0},
		{"#00FF00", 0, 255, 0},
		{"#0000FF", 0, 0, 255},
		{"CCCCCC", 204, 204, 204},
		{"#FFF", 0, 0, 0}, // invalid length returns 0,0,0
		{"", 0, 0, 0},     // empty returns 0,0,0
	}

	for _, tc := range tests {
		r, g, b := ParseHexToRGB(tc.hex)
		if r != tc.r || g != tc.g || b != tc.b {
			t.Errorf("ParseHexToRGB(%q) = (%d,%d,%d), want (%d,%d,%d)", tc.hex, r, g, b, tc.r, tc.g, tc.b)
		}
	}
}
