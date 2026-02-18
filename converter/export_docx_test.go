package converter

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDocxConverterBasic(t *testing.T) {
	conv := NewHTMLToDocxConverter()
	htmlContents := []string{`<html><body><h1>Test Heading</h1><p>Test paragraph.</p></body></html>`}

	if err := conv.Convert(htmlContents); err != nil {
		t.Fatalf("Convert failed: %v", err)
	}

	tmpFile := filepath.Join(t.TempDir(), "test.docx")
	if err := conv.SaveToFile(tmpFile); err != nil {
		t.Fatalf("SaveToFile failed: %v", err)
	}

	info, err := os.Stat(tmpFile)
	if err != nil {
		t.Fatalf("output file not created: %v", err)
	}
	if info.Size() == 0 {
		t.Error("output file is empty")
	}
}

func TestDocxConverterMultiplePages(t *testing.T) {
	conv := NewHTMLToDocxConverter()
	htmlContents := []string{
		`<html><body><h1>Page 1</h1></body></html>`,
		`<html><body><h1>Page 2</h1></body></html>`,
	}

	if err := conv.Convert(htmlContents); err != nil {
		t.Fatalf("Convert failed: %v", err)
	}

	tmpFile := filepath.Join(t.TempDir(), "multi.docx")
	if err := conv.SaveToFile(tmpFile); err != nil {
		t.Fatalf("SaveToFile failed: %v", err)
	}

	info, _ := os.Stat(tmpFile)
	if info.Size() == 0 {
		t.Error("output file is empty")
	}
}

func TestDocxConverterFormatting(t *testing.T) {
	conv := NewHTMLToDocxConverter()
	htmlContents := []string{`<html><body>
		<p>Normal <b>bold</b> <i>italic</i> <u>underline</u></p>
	</body></html>`}

	if err := conv.Convert(htmlContents); err != nil {
		t.Fatalf("Convert failed: %v", err)
	}

	tmpFile := filepath.Join(t.TempDir(), "format.docx")
	if err := conv.SaveToFile(tmpFile); err != nil {
		t.Fatalf("SaveToFile failed: %v", err)
	}
}

func TestDocxConverterTable(t *testing.T) {
	conv := NewHTMLToDocxConverter()
	htmlContents := []string{`<html><body>
		<table border="1">
			<tr><th>Name</th><th>Age</th></tr>
			<tr><td>Alice</td><td>30</td></tr>
		</table>
	</body></html>`}

	if err := conv.Convert(htmlContents); err != nil {
		t.Fatalf("Convert failed: %v", err)
	}

	tmpFile := filepath.Join(t.TempDir(), "table.docx")
	if err := conv.SaveToFile(tmpFile); err != nil {
		t.Fatalf("SaveToFile failed: %v", err)
	}
}

func TestDocxConverterLists(t *testing.T) {
	conv := NewHTMLToDocxConverter()
	htmlContents := []string{`<html><body>
		<ul><li>Item 1</li><li>Item 2</li></ul>
		<ol><li>First</li><li>Second</li></ol>
	</body></html>`}

	if err := conv.Convert(htmlContents); err != nil {
		t.Fatalf("Convert failed: %v", err)
	}

	tmpFile := filepath.Join(t.TempDir(), "lists.docx")
	if err := conv.SaveToFile(tmpFile); err != nil {
		t.Fatalf("SaveToFile failed: %v", err)
	}
}

func TestDocxConverterHeaderFooter(t *testing.T) {
	conv := NewHTMLToDocxConverter()
	htmlContents := []string{`<html><body>
		<header><p>Header text</p></header>
		<p>Body content</p>
		<footer><p>Footer text</p></footer>
	</body></html>`}

	if err := conv.Convert(htmlContents); err != nil {
		t.Fatalf("Convert failed: %v", err)
	}

	tmpFile := filepath.Join(t.TempDir(), "headerfooter.docx")
	if err := conv.SaveToFile(tmpFile); err != nil {
		t.Fatalf("SaveToFile failed: %v", err)
	}
}

func TestDocxConverterInvalidHTML(t *testing.T) {
	conv := NewHTMLToDocxConverter()
	// Go's html parser is lenient, so even malformed HTML should not error
	htmlContents := []string{`<html><body><p>Unclosed paragraph<div>Nested improperly</p></div></body></html>`}

	if err := conv.Convert(htmlContents); err != nil {
		t.Fatalf("Convert should handle malformed HTML gracefully: %v", err)
	}
}

func TestDocxConverterAllHeadings(t *testing.T) {
	conv := NewHTMLToDocxConverter()
	htmlContents := []string{`<html><body>
		<h1>Heading 1</h1>
		<h2>Heading 2</h2>
		<h3>Heading 3</h3>
		<h4>Heading 4</h4>
		<h5>Heading 5</h5>
		<h6>Heading 6</h6>
	</body></html>`}

	if err := conv.Convert(htmlContents); err != nil {
		t.Fatalf("Convert failed: %v", err)
	}

	tmpFile := filepath.Join(t.TempDir(), "headings.docx")
	if err := conv.SaveToFile(tmpFile); err != nil {
		t.Fatalf("SaveToFile failed: %v", err)
	}
}

func TestDocxConverterFontTag(t *testing.T) {
	conv := NewHTMLToDocxConverter()
	htmlContents := []string{`<html><body>
		<font size="5" color="#FF0000" face="Arial">Styled text</font>
	</body></html>`}

	if err := conv.Convert(htmlContents); err != nil {
		t.Fatalf("Convert failed: %v", err)
	}

	tmpFile := filepath.Join(t.TempDir(), "font.docx")
	if err := conv.SaveToFile(tmpFile); err != nil {
		t.Fatalf("SaveToFile failed: %v", err)
	}
}

func TestDocxConverterHorizontalRule(t *testing.T) {
	conv := NewHTMLToDocxConverter()
	htmlContents := []string{`<html><body>
		<p>Before</p>
		<hr size="2" color="#333333">
		<p>After</p>
	</body></html>`}

	if err := conv.Convert(htmlContents); err != nil {
		t.Fatalf("Convert failed: %v", err)
	}

	tmpFile := filepath.Join(t.TempDir(), "hr.docx")
	if err := conv.SaveToFile(tmpFile); err != nil {
		t.Fatalf("SaveToFile failed: %v", err)
	}
}
