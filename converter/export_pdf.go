package converter

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/jung-kurt/gofpdf"
	"golang.org/x/net/html"
)

// HTMLToPDFConverter converts HTML content to PDF format using gofpdf.
type HTMLToPDFConverter struct {
	pdf        *gofpdf.Fpdf
	fontStyle  string  // current style: combination of B, I, U
	fontSize   float64 // current font size in pt
	fontFamily string
	tr         func(string) string // UTF-8 translator
	centered   bool                // true when inside a <center> tag
}

// NewHTMLToPDFConverter creates a new PDF converter with A4 page and default margins.
func NewHTMLToPDFConverter() *HTMLToPDFConverter {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetMargins(15, 15, 15)
	pdf.SetAutoPageBreak(true, 15)
	pdf.AddPage()
	pdf.SetFont("Arial", "", 11)
	tr := pdf.UnicodeTranslatorFromDescriptor("")
	return &HTMLToPDFConverter{
		pdf:        pdf,
		fontStyle:  "",
		fontSize:   11,
		fontFamily: "Arial",
		tr:         tr,
	}
}

// Convert parses and converts multiple HTML strings to PDF content.
func (c *HTMLToPDFConverter) Convert(htmlContents []string) error {
	for i, content := range htmlContents {
		root, err := html.Parse(strings.NewReader(content))
		if err != nil {
			return fmt.Errorf("failed to parse HTML index %d: %w", i, err)
		}
		c.walkPDF(root)
		if i < len(htmlContents)-1 {
			c.pdf.AddPage()
		}
	}
	return nil
}

func (c *HTMLToPDFConverter) lineHeight() float64 {
	return c.fontSize * 0.4
}

func (c *HTMLToPDFConverter) applyFont() {
	c.pdf.SetFont(c.fontFamily, c.fontStyle, c.fontSize)
}

func (c *HTMLToPDFConverter) writeText(text string) {
	c.applyFont()
	if c.centered {
		pageW, _ := c.pdf.GetPageSize()
		lMargin, _, rMargin, _ := c.pdf.GetMargins()
		usableW := pageW - lMargin - rMargin
		c.pdf.SetX(lMargin)
		c.pdf.CellFormat(usableW, c.lineHeight(), c.tr(text), "", 1, "C", false, 0, "")
	} else {
		c.pdf.Write(c.lineHeight(), c.tr(text))
	}
}

func (c *HTMLToPDFConverter) walkPDF(n *html.Node) {
	if n.Type == html.TextNode {
		text := n.Data
		if strings.TrimSpace(text) == "" {
			return
		}
		collapsed := CollapseWhitespace(text)
		c.writeText(collapsed)
		return
	}

	if n.Type != html.ElementNode {
		c.processChildrenPDF(n)
		return
	}

	switch n.Data {
	case "head", "title", "style", "script", "meta", "link":
		return
	case "h1":
		c.pdfHeading(n, 22)
	case "h2":
		c.pdfHeading(n, 18)
	case "h3":
		c.pdfHeading(n, 14)
	case "h4":
		c.pdfHeading(n, 12)
	case "h5":
		c.pdfHeading(n, 10)
	case "h6":
		c.pdfHeading(n, 9)
	case "p":
		c.pdfBlock(n, 2, 3)
	case "div", "section", "article", "nav", "main":
		c.processChildrenPDF(n)
	case "center":
		c.processCenterPDF(n)
	case "br":
		c.pdf.Ln(c.lineHeight())
	case "hr":
		c.pdfHR()
	case "b", "strong":
		c.formattedPDF(n, "B")
	case "i", "em":
		c.formattedPDF(n, "I")
	case "u":
		c.formattedPDF(n, "U")
	case "font":
		c.processFontPDF(n)
	case "table":
		c.processTablePDF(n)
	case "ul":
		c.processListPDF(n, false)
	case "ol":
		c.processListPDF(n, true)
	case "a":
		oldStyle := c.fontStyle
		c.fontStyle = c.addStyle(c.fontStyle, "U")
		c.pdf.SetTextColor(0, 0, 255)
		c.processChildrenPDF(n)
		c.pdf.SetTextColor(0, 0, 0)
		c.fontStyle = oldStyle
		c.applyFont()
	case "header":
		c.processChildrenPDF(n)
	case "footer":
		c.processFooterPDF(n)
	case "img":
		// Skip images
	default:
		c.processChildrenPDF(n)
	}
}

func (c *HTMLToPDFConverter) pdfBlock(n *html.Node, spaceBefore, spaceAfter float64) {
	c.pdf.Ln(spaceBefore)
	lMargin, _, _, _ := c.pdf.GetMargins()
	c.pdf.SetX(lMargin)
	c.processChildrenPDF(n)
	c.pdf.Ln(spaceAfter)
}

func (c *HTMLToPDFConverter) pdfHeading(n *html.Node, size float64) {
	c.pdf.Ln(6)
	oldSize := c.fontSize
	oldStyle := c.fontStyle
	c.fontStyle = c.addStyle(c.fontStyle, "B")
	c.fontSize = size
	c.applyFont()

	text := strings.TrimSpace(ExtractText(n))
	lMargin, _, _, _ := c.pdf.GetMargins()
	c.pdf.SetX(lMargin)
	pageW, _ := c.pdf.GetPageSize()
	_, _, rMargin, _ := c.pdf.GetMargins()
	usableW := pageW - lMargin - rMargin
	align := ""
	if c.centered {
		align = "C"
	}
	c.pdf.MultiCell(usableW, size*0.5, c.tr(text), "", align, false)
	c.pdf.Ln(2)

	c.fontSize = oldSize
	c.fontStyle = oldStyle
	c.applyFont()
}

func (c *HTMLToPDFConverter) pdfHR() {
	c.pdf.Ln(4)
	pageW, _ := c.pdf.GetPageSize()
	lMargin, _, rMargin, _ := c.pdf.GetMargins()
	y := c.pdf.GetY()
	c.pdf.SetDrawColor(128, 128, 128)
	c.pdf.SetLineWidth(0.3)
	c.pdf.Line(lMargin, y, pageW-rMargin, y)
	c.pdf.SetDrawColor(0, 0, 0)
	c.pdf.Ln(4)
}

func (c *HTMLToPDFConverter) formattedPDF(n *html.Node, style string) {
	oldStyle := c.fontStyle
	c.fontStyle = c.addStyle(c.fontStyle, style)
	c.applyFont()
	c.processChildrenPDF(n)
	c.fontStyle = oldStyle
	c.applyFont()
}

func (c *HTMLToPDFConverter) processFontPDF(n *html.Node) {
	attrs := GetAttrMap(n.Attr)
	oldSize := c.fontSize
	oldFamily := c.fontFamily
	oldR, oldG, oldB := c.pdf.GetTextColor()

	if val, ok := attrs["size"]; ok {
		switch val {
		case "1":
			c.fontSize = 8
		case "2":
			c.fontSize = 10
		case "3":
			c.fontSize = 12
		case "4":
			c.fontSize = 14
		case "5":
			c.fontSize = 18
		case "6":
			c.fontSize = 24
		case "7":
			c.fontSize = 36
		}
	}

	if val, ok := attrs["color"]; ok {
		r, g, b := ParseHexToRGB(val)
		c.pdf.SetTextColor(r, g, b)
	}

	if val, ok := attrs["face"]; ok {
		switch strings.ToLower(val) {
		case "arial", "helvetica":
			c.fontFamily = "Arial"
		case "times", "times new roman":
			c.fontFamily = "Times"
		case "courier", "courier new":
			c.fontFamily = "Courier"
		default:
			c.fontFamily = "Arial"
		}
	}

	c.applyFont()
	c.processChildrenPDF(n)

	c.fontSize = oldSize
	c.fontFamily = oldFamily
	c.pdf.SetTextColor(oldR, oldG, oldB)
	c.applyFont()
}

func (c *HTMLToPDFConverter) processFooterPDF(n *html.Node) {
	var lines []string
	var extractLines func(*html.Node)
	extractLines = func(node *html.Node) {
		for ch := node.FirstChild; ch != nil; ch = ch.NextSibling {
			if ch.Type == html.TextNode {
				text := strings.TrimSpace(ch.Data)
				if text != "" {
					lines = append(lines, text)
				}
			} else if ch.Type == html.ElementNode {
				if ch.Data == "br" {
					continue
				}
				extractLines(ch)
			}
		}
	}
	extractLines(n)

	if len(lines) == 0 {
		return
	}

	tr := c.tr
	c.pdf.SetFooterFunc(func() {
		c.pdf.SetY(-20)
		c.pdf.SetFont("Arial", "", 9)
		pageW, _ := c.pdf.GetPageSize()
		lMargin, _, rMargin, _ := c.pdf.GetMargins()
		usableW := pageW - lMargin - rMargin
		for _, line := range lines {
			c.pdf.SetX(lMargin)
			c.pdf.CellFormat(usableW, 4, tr(line), "", 1, "C", false, 0, "")
		}
	})
}

func (c *HTMLToPDFConverter) processCenterPDF(n *html.Node) {
	oldCentered := c.centered
	c.centered = true
	c.processChildrenPDF(n)
	c.centered = oldCentered
}

func (c *HTMLToPDFConverter) processTablePDF(n *html.Node) {
	attrs := GetAttrMap(n.Attr)
	borderVal := attrs["border"]
	isDataTable := borderVal != "" && borderVal != "0"

	if !isDataTable {
		c.processTableChildrenAsFlow(n)
		return
	}

	c.pdf.Ln(4)
	pageW, _ := c.pdf.GetPageSize()
	lMargin, _, rMargin, _ := c.pdf.GetMargins()
	usableW := pageW - lMargin - rMargin

	tableW := usableW
	if val, ok := attrs["width"]; ok {
		if strings.HasSuffix(val, "%") {
			pct := strings.TrimSuffix(val, "%")
			if p, err := strconv.ParseFloat(pct, 64); err == nil {
				tableW = usableW * p / 100.0
			}
		}
	}

	type cellData struct {
		text     string
		isHeader bool
	}
	var rows [][]cellData

	var collectRows func(*html.Node)
	collectRows = func(curr *html.Node) {
		for ch := curr.FirstChild; ch != nil; ch = ch.NextSibling {
			if ch.Type == html.ElementNode && ch.Data == "tr" {
				var row []cellData
				for cell := ch.FirstChild; cell != nil; cell = cell.NextSibling {
					if cell.Type == html.ElementNode && (cell.Data == "td" || cell.Data == "th") {
						text := strings.TrimSpace(ExtractText(cell))
						row = append(row, cellData{text: text, isHeader: cell.Data == "th"})
					}
				}
				if len(row) > 0 {
					rows = append(rows, row)
				}
			} else if ch.FirstChild != nil {
				collectRows(ch)
			}
		}
	}
	collectRows(n)

	if len(rows) == 0 {
		return
	}

	maxCols := 0
	for _, row := range rows {
		if len(row) > maxCols {
			maxCols = len(row)
		}
	}
	if maxCols == 0 {
		return
	}
	colW := tableW / float64(maxCols)
	rowH := 7.0

	tableX := lMargin + (usableW-tableW)/2

	for _, row := range rows {
		c.pdf.SetX(tableX)
		for _, cell := range row {
			if cell.isHeader {
				c.pdf.SetFont(c.fontFamily, "B", c.fontSize)
			} else {
				c.pdf.SetFont(c.fontFamily, c.fontStyle, c.fontSize)
			}
			c.pdf.CellFormat(colW, rowH, c.tr(cell.text), "1", 0, "", false, 0, "")
		}
		c.pdf.Ln(rowH)
	}
	c.applyFont()
	c.pdf.Ln(4)
}

func (c *HTMLToPDFConverter) processTableChildrenAsFlow(n *html.Node) {
	for ch := n.FirstChild; ch != nil; ch = ch.NextSibling {
		if ch.Type == html.ElementNode {
			switch ch.Data {
			case "tr":
				c.processTableChildrenAsFlow(ch)
			case "td", "th":
				c.processChildrenPDF(ch)
			case "thead", "tbody", "tfoot":
				c.processTableChildrenAsFlow(ch)
			default:
				c.walkPDF(ch)
			}
		}
	}
}

func (c *HTMLToPDFConverter) processListPDF(n *html.Node, ordered bool) {
	c.pdf.Ln(2)
	index := 1
	lMargin, _, _, _ := c.pdf.GetMargins()
	indent := 10.0

	for li := n.FirstChild; li != nil; li = li.NextSibling {
		if li.Type == html.ElementNode && li.Data == "li" {
			c.pdf.SetX(lMargin + indent)
			c.applyFont()
			if ordered {
				prefix := fmt.Sprintf("%d. ", index)
				c.pdf.Write(c.lineHeight(), prefix)
				index++
			} else {
				c.pdf.Write(c.lineHeight(), "- ")
			}
			text := strings.TrimSpace(ExtractText(li))
			c.pdf.Write(c.lineHeight(), c.tr(text))
			c.pdf.Ln(c.lineHeight() + 1)
		}
	}
	c.pdf.Ln(2)
}

func (c *HTMLToPDFConverter) processChildrenPDF(n *html.Node) {
	for ch := n.FirstChild; ch != nil; ch = ch.NextSibling {
		c.walkPDF(ch)
	}
}

// SaveToFile saves the PDF to a file.
func (c *HTMLToPDFConverter) SaveToFile(filename string) error {
	return c.pdf.OutputFileAndClose(filename)
}

func (c *HTMLToPDFConverter) addStyle(current, add string) string {
	if strings.Contains(current, add) {
		return current
	}
	return current + add
}
