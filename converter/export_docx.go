package converter

import (
	"fmt"
	"strings"

	"baliance.com/gooxml/color"
	"baliance.com/gooxml/document"
	"baliance.com/gooxml/measurement"
	"baliance.com/gooxml/schema/soo/ofc/sharedTypes"
	"baliance.com/gooxml/schema/soo/wml"
	"golang.org/x/net/html"
)

// HTMLToDocxConverter converts HTML content to DOCX format.
type HTMLToDocxConverter struct {
	doc *document.Document
}

// NewHTMLToDocxConverter creates a new DOCX converter with default page settings.
func NewHTMLToDocxConverter() *HTMLToDocxConverter {
	doc := document.New()
	section := doc.BodySection()
	section.SetPageMargins(measurement.Inch, measurement.Inch, measurement.Inch, measurement.Inch, 0, 0, 0)
	return &HTMLToDocxConverter{doc: doc}
}

// parseHexColor safely converts a hex color string.
func parseHexColor(hex string) color.Color {
	hex = strings.TrimPrefix(hex, "#")
	if len(hex) != 6 {
		return color.Black
	}
	return color.FromHex(hex)
}

func uint64Ptr(u uint64) *uint64 { return &u }

// Convert parses and converts multiple HTML strings to DOCX content.
func (c *HTMLToDocxConverter) Convert(htmlContents []string) error {
	for i, content := range htmlContents {
		content = UnescapeUnicodeHTML(content)
		root, err := html.Parse(strings.NewReader(content))
		if err != nil {
			return fmt.Errorf("failed to parse index %d: %w", i, err)
		}
		c.walk(root, nil, nil, wml.ST_JcLeft)
		if i < len(htmlContents)-1 {
			c.addPageBreak()
		}
	}
	return nil
}

func (c *HTMLToDocxConverter) walk(n *html.Node, para *document.Paragraph, container interface{}, align wml.ST_Jc) {
	if n.Type == html.TextNode {
		text := strings.TrimSpace(n.Data)
		if text != "" {
			if para == nil {
				p := c.createParagraph(container)
				p.Properties().SetAlignment(align)
				para = &p
			}
			para.AddRun().AddText(n.Data)
		}
		return
	}

	if n.Type == html.ElementNode {
		attrs := GetAttrMap(n.Attr)
		currentAlign := align
		if val, ok := attrs["align"]; ok {
			switch strings.ToLower(val) {
			case "center":
				currentAlign = wml.ST_JcCenter
			case "right":
				currentAlign = wml.ST_JcRight
			case "left":
				currentAlign = wml.ST_JcLeft
			}
		}

		nodeType := EffectiveNodeType(n)
		switch nodeType {
		case "header":
			hdr := c.doc.AddHeader()
			c.doc.BodySection().SetHeader(hdr, wml.ST_HdrFtrDefault)
			c.processChildren(n, nil, hdr, currentAlign)
			return
		case "footer":
			ftr := c.doc.AddFooter()
			c.doc.BodySection().SetFooter(ftr, wml.ST_HdrFtrDefault)
			c.processChildren(n, nil, ftr, currentAlign)
			return
		case "center":
			c.processChildren(n, para, container, wml.ST_JcCenter)
			return
		case "p", "section", "article", "nav", "main", "body", "html":
			p := c.createParagraph(container)
			p.Properties().SetAlignment(currentAlign)
			c.processChildren(n, &p, container, currentAlign)
			return
		case "div", "span":
			c.processChildren(n, para, container, align)
			return
		case "h1", "h2", "h3", "h4", "h5", "h6":
			p := c.createParagraph(container)
			p.SetStyle("Heading" + nodeType[1:])
			p.Properties().SetAlignment(currentAlign)

			c.processChildren(n, &p, container, currentAlign)

			size := 12.0
			switch nodeType {
			case "h1":
				size = 24
			case "h2":
				size = 18
			case "h3":
				size = 14
			case "h4":
				size = 12
			case "h5":
				size = 10
			case "h6":
				size = 8
			}

			for _, r := range p.Runs() {
				r.Properties().SetBold(true)
				r.Properties().SetSize(measurement.Distance(size))
			}
			return
		case "hr":
			p := c.createParagraph(container)
			c.applyHRStyle(&p, attrs)
			return
		case "table":
			c.processTable(n, currentAlign)
			return
		case "ul", "ol":
			c.processList(n, nodeType == "ol", container, currentAlign)
			return
		case "b", "strong":
			c.formatted(n, para, "bold", container, currentAlign)
			return
		case "i", "em":
			c.formatted(n, para, "italic", container, currentAlign)
			return
		case "u":
			c.formatted(n, para, "underline", container, currentAlign)
			return
		case "font":
			c.processFont(n, para, container, currentAlign)
			return
		}
	}
	c.processChildren(n, para, container, align)
}

func (c *HTMLToDocxConverter) processFont(n *html.Node, para *document.Paragraph, container interface{}, align wml.ST_Jc) {
	if para == nil {
		p := c.createParagraph(container)
		p.Properties().SetAlignment(align)
		para = &p
	}
	attrs := GetAttrMap(n.Attr)

	for ch := n.FirstChild; ch != nil; ch = ch.NextSibling {
		if ch.Type == html.TextNode {
			text := strings.TrimSpace(ch.Data)
			if text != "" {
				r := para.AddRun()
				r.AddText(ch.Data)

				if val, ok := attrs["size"]; ok {
					sz := 12.0
					switch val {
					case "1":
						sz = 8
					case "2":
						sz = 10
					case "3":
						sz = 12
					case "4":
						sz = 14
					case "5":
						sz = 18
					case "6":
						sz = 24
					case "7":
						sz = 36
					}
					r.Properties().SetSize(measurement.Distance(sz))
				}
				if val, ok := attrs["color"]; ok {
					r.Properties().SetColor(parseHexColor(val))
				}
				if val, ok := attrs["face"]; ok {
					r.Properties().SetFontFamily(val)
				}
			}
		} else {
			c.walk(ch, para, container, align)
		}
	}
}

func (c *HTMLToDocxConverter) processTable(n *html.Node, align wml.ST_Jc) {
	table := c.doc.AddTable()
	table.Properties().SetWidthPercent(100)
	attrs := GetAttrMap(n.Attr)

	if b, ok := attrs["border"]; ok && b != "0" {
		table.Properties().Borders().SetAll(wml.ST_BorderSingle, color.Auto, 1*measurement.Point)
	}

	if align == wml.ST_JcCenter || attrs["align"] == "center" {
		table.Properties().SetAlignment(wml.ST_JcTableCenter)
	}

	var walkTable func(*html.Node)
	walkTable = func(curr *html.Node) {
		for ch := curr.FirstChild; ch != nil; ch = ch.NextSibling {
			if ch.Type == html.ElementNode {
				chType := EffectiveNodeType(ch)
				if chType == "tr" {
					row := table.AddRow()
					trAttrs := GetAttrMap(ch.Attr)

					var collectCells func(*html.Node)
					collectCells = func(cellContainer *html.Node) {
						for cellNode := cellContainer.FirstChild; cellNode != nil; cellNode = cellNode.NextSibling {
							if cellNode.Type == html.ElementNode {
								cellType := EffectiveNodeType(cellNode)
								if cellType == "td" || cellType == "th" {
									cell := row.AddCell()
									tdAttrs := GetAttrMap(cellNode.Attr)

									bg := ""
									if val, ok := tdAttrs["bgcolor"]; ok {
										bg = val
									} else if val, ok := trAttrs["bgcolor"]; ok {
										bg = val
									}
									if bg != "" {
										cell.Properties().SetShading(wml.ST_ShdSolid, parseHexColor(bg), color.Auto)
									}

									p := cell.AddParagraph()
									cellAlign := wml.ST_JcLeft
									if tdAttrs["align"] == "center" {
										cellAlign = wml.ST_JcCenter
									}
									p.Properties().SetAlignment(cellAlign)

									if cellType == "th" {
										p.AddRun().Properties().SetBold(true)
									}
									c.processChildren(cellNode, &p, nil, cellAlign)
								} else if cellType == "div" || cellType == "span" {
									collectCells(cellNode)
								}
							}
						}
					}
					collectCells(ch)
				} else if ch.FirstChild != nil {
					walkTable(ch)
				}
			}
		}
	}
	walkTable(n)
}

func (c *HTMLToDocxConverter) applyHRStyle(p *document.Paragraph, attrs map[string]string) {
	if p.X().PPr == nil {
		p.X().PPr = wml.NewCT_PPr()
	}
	if p.X().PPr.PBdr == nil {
		p.X().PPr.PBdr = wml.NewCT_PBdr()
	}

	p.X().PPr.PBdr.Top = nil
	p.X().PPr.PBdr.Left = nil
	p.X().PPr.PBdr.Right = nil
	p.X().PPr.PBdr.Between = nil

	p.X().PPr.PBdr.Bottom = wml.NewCT_Border()
	p.X().PPr.PBdr.Bottom.ValAttr = wml.ST_BorderSingle

	thickness := uint64(4)
	if val, ok := attrs["size"]; ok {
		var t uint64
		if _, err := fmt.Sscanf(val, "%d", &t); err == nil && t > 0 {
			thickness = t * 2
		}
	}
	p.X().PPr.PBdr.Bottom.SzAttr = uint64Ptr(thickness)

	colorStr := "808080"
	if val, ok := attrs["color"]; ok {
		clean := strings.TrimPrefix(val, "#")
		if len(clean) == 6 {
			colorStr = clean
		}
	}
	p.X().PPr.PBdr.Bottom.ColorAttr = &wml.ST_HexColor{ST_HexColorRGB: &colorStr}

	if p.X().PPr.Spacing == nil {
		p.X().PPr.Spacing = wml.NewCT_Spacing()
	}
	zero := uint64(0)
	p.X().PPr.Spacing.BeforeAttr = &sharedTypes.ST_TwipsMeasure{ST_UnsignedDecimalNumber: &zero}
	p.X().PPr.Spacing.AfterAttr = &sharedTypes.ST_TwipsMeasure{ST_UnsignedDecimalNumber: &zero}
}

func (c *HTMLToDocxConverter) processChildren(n *html.Node, para *document.Paragraph, container interface{}, align wml.ST_Jc) {
	for ch := n.FirstChild; ch != nil; ch = ch.NextSibling {
		c.walk(ch, para, container, align)
	}
}

func (c *HTMLToDocxConverter) formatted(n *html.Node, para *document.Paragraph, style string, container interface{}, align wml.ST_Jc) {
	if para == nil {
		p := c.createParagraph(container)
		p.Properties().SetAlignment(align)
		para = &p
	}
	for ch := n.FirstChild; ch != nil; ch = ch.NextSibling {
		if ch.Type == html.TextNode {
			r := para.AddRun()
			r.AddText(ch.Data)
			if style == "bold" {
				r.Properties().SetBold(true)
			}
			if style == "italic" {
				r.Properties().SetItalic(true)
			}
			if style == "underline" {
				r.Properties().SetUnderline(wml.ST_UnderlineSingle, color.Auto)
			}
		} else {
			c.walk(ch, para, container, align)
		}
	}
}

func (c *HTMLToDocxConverter) processList(n *html.Node, ordered bool, container interface{}, align wml.ST_Jc) {
	index := 1
	for li := n.FirstChild; li != nil; li = li.NextSibling {
		if li.Type == html.ElementNode && EffectiveNodeType(li) == "li" {
			p := c.createParagraph(container)
			p.Properties().SetAlignment(align)
			prefix := "• "
			if ordered {
				prefix = fmt.Sprintf("%d. ", index)
				index++
			}
			p.AddRun().AddText(prefix)
			c.processChildren(li, &p, container, align)
		}
	}
}

func (c *HTMLToDocxConverter) createParagraph(container interface{}) document.Paragraph {
	if hdr, ok := container.(document.Header); ok {
		return hdr.AddParagraph()
	}
	if ftr, ok := container.(document.Footer); ok {
		return ftr.AddParagraph()
	}
	return c.doc.AddParagraph()
}

func (c *HTMLToDocxConverter) addPageBreak() {
	p := c.doc.AddParagraph()
	run := p.AddRun()
	run.X().EG_RunInnerContent = append(run.X().EG_RunInnerContent, &wml.EG_RunInnerContent{Br: wml.NewCT_Br()})
	run.X().EG_RunInnerContent[len(run.X().EG_RunInnerContent)-1].Br.TypeAttr = wml.ST_BrTypePage
}

// SaveToFile saves the DOCX document to a file.
func (c *HTMLToDocxConverter) SaveToFile(filename string) error {
	return c.doc.SaveToFile(filename)
}
