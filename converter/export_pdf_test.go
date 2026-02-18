package converter

import (
	"os"
	"path/filepath"
	"testing"
)

func TestPDFConverterBasic(t *testing.T) {
	conv := NewHTMLToPDFConverter()
	htmlContents := []string{`<html><body><h1>Test Heading</h1><p>Test paragraph.</p></body></html>`}

	if err := conv.Convert(htmlContents); err != nil {
		t.Fatalf("Convert failed: %v", err)
	}

	tmpFile := filepath.Join(t.TempDir(), "test.pdf")
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

func TestPDFConverterMultiplePages(t *testing.T) {
	conv := NewHTMLToPDFConverter()
	htmlContents := []string{
		`<html><body><h1>Page 1</h1></body></html>`,
		`<html><body><h1>Page 2</h1></body></html>`,
	}

	if err := conv.Convert(htmlContents); err != nil {
		t.Fatalf("Convert failed: %v", err)
	}

	tmpFile := filepath.Join(t.TempDir(), "multi.pdf")
	if err := conv.SaveToFile(tmpFile); err != nil {
		t.Fatalf("SaveToFile failed: %v", err)
	}

	info, _ := os.Stat(tmpFile)
	if info.Size() == 0 {
		t.Error("output file is empty")
	}
}

func TestPDFConverterFormatting(t *testing.T) {
	conv := NewHTMLToPDFConverter()
	htmlContents := []string{`<html><body>
		<p>Normal <b>bold</b> <i>italic</i> <u>underline</u></p>
	</body></html>`}

	if err := conv.Convert(htmlContents); err != nil {
		t.Fatalf("Convert failed: %v", err)
	}

	tmpFile := filepath.Join(t.TempDir(), "format.pdf")
	if err := conv.SaveToFile(tmpFile); err != nil {
		t.Fatalf("SaveToFile failed: %v", err)
	}
}

func TestPDFConverterDataTable(t *testing.T) {
	conv := NewHTMLToPDFConverter()
	htmlContents := []string{`<html><body>
		<table border="1">
			<tr><th>Name</th><th>Age</th></tr>
			<tr><td>Alice</td><td>30</td></tr>
			<tr><td>Bob</td><td>25</td></tr>
		</table>
	</body></html>`}

	if err := conv.Convert(htmlContents); err != nil {
		t.Fatalf("Convert failed: %v", err)
	}

	tmpFile := filepath.Join(t.TempDir(), "data_table.pdf")
	if err := conv.SaveToFile(tmpFile); err != nil {
		t.Fatalf("SaveToFile failed: %v", err)
	}
}

func TestPDFConverterLayoutTable(t *testing.T) {
	conv := NewHTMLToPDFConverter()
	htmlContents := []string{`<html><body>
		<table border="0" width="100%">
			<tr><td><p>Layout content</p></td></tr>
		</table>
	</body></html>`}

	if err := conv.Convert(htmlContents); err != nil {
		t.Fatalf("Convert failed: %v", err)
	}

	tmpFile := filepath.Join(t.TempDir(), "layout_table.pdf")
	if err := conv.SaveToFile(tmpFile); err != nil {
		t.Fatalf("SaveToFile failed: %v", err)
	}
}

func TestPDFConverterLists(t *testing.T) {
	conv := NewHTMLToPDFConverter()
	htmlContents := []string{`<html><body>
		<ul><li>Bullet 1</li><li>Bullet 2</li></ul>
		<ol><li>Number 1</li><li>Number 2</li></ol>
	</body></html>`}

	if err := conv.Convert(htmlContents); err != nil {
		t.Fatalf("Convert failed: %v", err)
	}

	tmpFile := filepath.Join(t.TempDir(), "lists.pdf")
	if err := conv.SaveToFile(tmpFile); err != nil {
		t.Fatalf("SaveToFile failed: %v", err)
	}
}

func TestPDFConverterCentering(t *testing.T) {
	conv := NewHTMLToPDFConverter()
	htmlContents := []string{`<html><body>
		<center>
			<h1>Centered Heading</h1>
			<p>Centered text</p>
		</center>
	</body></html>`}

	if err := conv.Convert(htmlContents); err != nil {
		t.Fatalf("Convert failed: %v", err)
	}

	tmpFile := filepath.Join(t.TempDir(), "center.pdf")
	if err := conv.SaveToFile(tmpFile); err != nil {
		t.Fatalf("SaveToFile failed: %v", err)
	}
}

func TestPDFConverterFooter(t *testing.T) {
	conv := NewHTMLToPDFConverter()
	htmlContents := []string{`<html><body>
		<p>Content</p>
		<footer><center>Copyright 2026</center></footer>
	</body></html>`}

	if err := conv.Convert(htmlContents); err != nil {
		t.Fatalf("Convert failed: %v", err)
	}

	tmpFile := filepath.Join(t.TempDir(), "footer.pdf")
	if err := conv.SaveToFile(tmpFile); err != nil {
		t.Fatalf("SaveToFile failed: %v", err)
	}
}

func TestPDFConverterLinks(t *testing.T) {
	conv := NewHTMLToPDFConverter()
	htmlContents := []string{`<html><body>
		<a href="https://example.com">Click here</a>
	</body></html>`}

	if err := conv.Convert(htmlContents); err != nil {
		t.Fatalf("Convert failed: %v", err)
	}

	tmpFile := filepath.Join(t.TempDir(), "links.pdf")
	if err := conv.SaveToFile(tmpFile); err != nil {
		t.Fatalf("SaveToFile failed: %v", err)
	}
}

func TestPDFConverterFontTag(t *testing.T) {
	conv := NewHTMLToPDFConverter()
	htmlContents := []string{`<html><body>
		<font size="5" color="#FF0000" face="Arial">Red text</font>
		<font face="Courier">Mono text</font>
		<font face="Times New Roman">Serif text</font>
	</body></html>`}

	if err := conv.Convert(htmlContents); err != nil {
		t.Fatalf("Convert failed: %v", err)
	}

	tmpFile := filepath.Join(t.TempDir(), "font.pdf")
	if err := conv.SaveToFile(tmpFile); err != nil {
		t.Fatalf("SaveToFile failed: %v", err)
	}
}

func TestPDFConverterHR(t *testing.T) {
	conv := NewHTMLToPDFConverter()
	htmlContents := []string{`<html><body>
		<p>Before</p><hr><p>After</p>
	</body></html>`}

	if err := conv.Convert(htmlContents); err != nil {
		t.Fatalf("Convert failed: %v", err)
	}

	tmpFile := filepath.Join(t.TempDir(), "hr.pdf")
	if err := conv.SaveToFile(tmpFile); err != nil {
		t.Fatalf("SaveToFile failed: %v", err)
	}
}

func TestPDFConverterAllHeadings(t *testing.T) {
	conv := NewHTMLToPDFConverter()
	htmlContents := []string{`<html><body>
		<h1>H1</h1><h2>H2</h2><h3>H3</h3>
		<h4>H4</h4><h5>H5</h5><h6>H6</h6>
	</body></html>`}

	if err := conv.Convert(htmlContents); err != nil {
		t.Fatalf("Convert failed: %v", err)
	}

	tmpFile := filepath.Join(t.TempDir(), "headings.pdf")
	if err := conv.SaveToFile(tmpFile); err != nil {
		t.Fatalf("SaveToFile failed: %v", err)
	}
}

func TestPDFConverterUTF8(t *testing.T) {
	conv := NewHTMLToPDFConverter()
	htmlContents := []string{`<html><body>
		<p>Copyright &copy; 2026. All rights reserved.</p>
	</body></html>`}

	if err := conv.Convert(htmlContents); err != nil {
		t.Fatalf("Convert failed: %v", err)
	}

	tmpFile := filepath.Join(t.TempDir(), "utf8.pdf")
	if err := conv.SaveToFile(tmpFile); err != nil {
		t.Fatalf("SaveToFile failed: %v", err)
	}
}

func TestPDFConverterSkipsNonVisible(t *testing.T) {
	conv := NewHTMLToPDFConverter()
	htmlContents := []string{`<html>
		<head><title>Title</title><style>body{}</style><script>var x;</script></head>
		<body><p>Visible</p></body>
	</html>`}

	if err := conv.Convert(htmlContents); err != nil {
		t.Fatalf("Convert failed: %v", err)
	}

	tmpFile := filepath.Join(t.TempDir(), "skip.pdf")
	if err := conv.SaveToFile(tmpFile); err != nil {
		t.Fatalf("SaveToFile failed: %v", err)
	}
}

func TestPDFConverterComplexDocument(t *testing.T) {
	conv := NewHTMLToPDFConverter()
	htmlContents := []string{`<!DOCTYPE html>
<html>
<head><title>Test</title></head>
<body>
	<center><h1>Title</h1></center>
	<hr>
	<table border="0" width="80%"><tr><td>
		<h2>Section</h2>
		<p>Text with <b>bold</b> and <i>italic</i>.</p>
		<ul><li>A</li><li>B</li></ul>
		<table border="1"><tr><th>K</th><th>V</th></tr><tr><td>a</td><td>1</td></tr></table>
	</td></tr></table>
	<hr>
	<footer><center><font size="2">Footer &copy; 2026</font></center></footer>
</body>
</html>`}

	if err := conv.Convert(htmlContents); err != nil {
		t.Fatalf("Convert failed: %v", err)
	}

	tmpFile := filepath.Join(t.TempDir(), "complex.pdf")
	if err := conv.SaveToFile(tmpFile); err != nil {
		t.Fatalf("SaveToFile failed: %v", err)
	}

	info, _ := os.Stat(tmpFile)
	if info.Size() < 500 {
		t.Errorf("complex PDF seems too small: %d bytes", info.Size())
	}
}
