package converter

import (
	"fmt"
	"os"
	"strings"

	"golang.org/x/net/html"
)

// HTMLToMarkdownConverter converts HTML content to Markdown format.
type HTMLToMarkdownConverter struct {
	markdown  strings.Builder
	listDepth int
}

// NewHTMLToMarkdownConverter creates a new Markdown converter.
func NewHTMLToMarkdownConverter() *HTMLToMarkdownConverter {
	return &HTMLToMarkdownConverter{}
}

// Convert parses and converts multiple HTML strings to a Markdown string.
func (c *HTMLToMarkdownConverter) Convert(htmlContents []string) (string, error) {
	for i, content := range htmlContents {
		root, err := html.Parse(strings.NewReader(content))
		if err != nil {
			return "", fmt.Errorf("failed to parse HTML: %w", err)
		}
		c.walkMD(root)
		if i < len(htmlContents)-1 {
			c.markdown.WriteString("\n\n---\n\n")
		}
	}
	return c.markdown.String(), nil
}

func (c *HTMLToMarkdownConverter) walkMD(n *html.Node) {
	if n.Type == html.TextNode {
		text := n.Data
		if strings.TrimSpace(text) != "" {
			c.markdown.WriteString(text)
		}
		return
	}

	if n.Type == html.ElementNode {
		switch n.Data {
		case "h1":
			c.markdown.WriteString("\n# ")
			c.processChildrenMD(n)
			c.markdown.WriteString("\n\n")
			return
		case "h2":
			c.markdown.WriteString("\n## ")
			c.processChildrenMD(n)
			c.markdown.WriteString("\n\n")
			return
		case "h3":
			c.markdown.WriteString("\n### ")
			c.processChildrenMD(n)
			c.markdown.WriteString("\n\n")
			return
		case "h4":
			c.markdown.WriteString("\n#### ")
			c.processChildrenMD(n)
			c.markdown.WriteString("\n\n")
			return
		case "h5":
			c.markdown.WriteString("\n##### ")
			c.processChildrenMD(n)
			c.markdown.WriteString("\n\n")
			return
		case "h6":
			c.markdown.WriteString("\n###### ")
			c.processChildrenMD(n)
			c.markdown.WriteString("\n\n")
			return
		case "p":
			c.processChildrenMD(n)
			c.markdown.WriteString("\n\n")
			return
		case "br":
			c.markdown.WriteString("  \n")
			return
		case "hr":
			c.markdown.WriteString("\n---\n\n")
			return
		case "b", "strong":
			c.markdown.WriteString("**")
			c.processChildrenMD(n)
			c.markdown.WriteString("**")
			return
		case "i", "em":
			c.markdown.WriteString("*")
			c.processChildrenMD(n)
			c.markdown.WriteString("*")
			return
		case "u":
			c.markdown.WriteString("<u>")
			c.processChildrenMD(n)
			c.markdown.WriteString("</u>")
			return
		case "ul":
			c.markdown.WriteString("\n")
			c.processListMD(n, false)
			c.markdown.WriteString("\n")
			return
		case "ol":
			c.markdown.WriteString("\n")
			c.processListMD(n, true)
			c.markdown.WriteString("\n")
			return
		case "table":
			c.processTableMD(n)
			return
		case "a":
			href := GetAttrValue(n.Attr, "href")
			c.markdown.WriteString("[")
			c.processChildrenMD(n)
			c.markdown.WriteString("](" + href + ")")
			return
		case "img":
			src := GetAttrValue(n.Attr, "src")
			alt := GetAttrValue(n.Attr, "alt")
			c.markdown.WriteString("![" + alt + "](" + src + ")")
			return
		case "code":
			c.markdown.WriteString("`")
			c.processChildrenMD(n)
			c.markdown.WriteString("`")
			return
		case "pre":
			c.markdown.WriteString("\n```\n")
			c.processChildrenMD(n)
			c.markdown.WriteString("\n```\n\n")
			return
		case "blockquote":
			c.markdown.WriteString("\n> ")
			c.processChildrenMD(n)
			c.markdown.WriteString("\n\n")
			return
		}
	}

	c.processChildrenMD(n)
}

func (c *HTMLToMarkdownConverter) processChildrenMD(n *html.Node) {
	for ch := n.FirstChild; ch != nil; ch = ch.NextSibling {
		c.walkMD(ch)
	}
}

func (c *HTMLToMarkdownConverter) processListMD(n *html.Node, ordered bool) {
	index := 1
	for li := n.FirstChild; li != nil; li = li.NextSibling {
		if li.Type == html.ElementNode && li.Data == "li" {
			indent := strings.Repeat("  ", c.listDepth)
			if ordered {
				c.markdown.WriteString(fmt.Sprintf("%s%d. ", indent, index))
				index++
			} else {
				c.markdown.WriteString(indent + "- ")
			}
			c.listDepth++
			c.processChildrenMD(li)
			c.listDepth--
			c.markdown.WriteString("\n")
		}
	}
}

func (c *HTMLToMarkdownConverter) processTableMD(n *html.Node) {
	c.markdown.WriteString("\n")

	var rows [][]*html.Node
	var isHeader []bool

	var collectRows func(*html.Node)
	collectRows = func(node *html.Node) {
		for ch := node.FirstChild; ch != nil; ch = ch.NextSibling {
			if ch.Type == html.ElementNode {
				if ch.Data == "tr" {
					var cells []*html.Node
					header := false
					for cell := ch.FirstChild; cell != nil; cell = cell.NextSibling {
						if cell.Type == html.ElementNode && (cell.Data == "td" || cell.Data == "th") {
							cells = append(cells, cell)
							if cell.Data == "th" {
								header = true
							}
						}
					}
					if len(cells) > 0 {
						rows = append(rows, cells)
						isHeader = append(isHeader, header)
					}
				} else {
					collectRows(ch)
				}
			}
		}
	}
	collectRows(n)

	for i, row := range rows {
		c.markdown.WriteString("| ")
		for _, cell := range row {
			c.processChildrenMD(cell)
			c.markdown.WriteString(" | ")
		}
		c.markdown.WriteString("\n")

		if i == 0 || (i < len(isHeader) && isHeader[i]) {
			c.markdown.WriteString("|")
			for range row {
				c.markdown.WriteString(" --- |")
			}
			c.markdown.WriteString("\n")
		}
	}
	c.markdown.WriteString("\n")
}

// ConvertHTMLToMarkdown is a convenience function that converts HTML to Markdown and saves it to a file.
func ConvertHTMLToMarkdown(htmlContents []string, outputPath string) error {
	conv := NewHTMLToMarkdownConverter()
	markdown, err := conv.Convert(htmlContents)
	if err != nil {
		return fmt.Errorf("failed to convert HTML to Markdown: %w", err)
	}

	markdown = strings.ReplaceAll(markdown, "\n\n\n", "\n\n")
	markdown = strings.TrimSpace(markdown)

	return os.WriteFile(outputPath, []byte(markdown), 0644)
}
