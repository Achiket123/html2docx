package converter

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestMarkdownConverterBasic(t *testing.T) {
	conv := NewHTMLToMarkdownConverter()
	htmlContents := []string{`<html><body><h1>Title</h1><p>Paragraph.</p></body></html>`}

	md, err := conv.Convert(htmlContents)
	if err != nil {
		t.Fatalf("Convert failed: %v", err)
	}

	if !strings.Contains(md, "# Title") {
		t.Errorf("expected '# Title' in output, got: %s", md)
	}
	if !strings.Contains(md, "Paragraph.") {
		t.Errorf("expected 'Paragraph.' in output, got: %s", md)
	}
}

func TestMarkdownConverterHeadings(t *testing.T) {
	conv := NewHTMLToMarkdownConverter()
	htmlContents := []string{`<html><body>
		<h1>H1</h1><h2>H2</h2><h3>H3</h3>
		<h4>H4</h4><h5>H5</h5><h6>H6</h6>
	</body></html>`}

	md, err := conv.Convert(htmlContents)
	if err != nil {
		t.Fatalf("Convert failed: %v", err)
	}

	for _, prefix := range []string{"# H1", "## H2", "### H3", "#### H4", "##### H5", "###### H6"} {
		if !strings.Contains(md, prefix) {
			t.Errorf("expected %q in output", prefix)
		}
	}
}

func TestMarkdownConverterFormatting(t *testing.T) {
	conv := NewHTMLToMarkdownConverter()
	htmlContents := []string{`<html><body>
		<b>bold</b> <i>italic</i> <u>underline</u>
		<strong>strong</strong> <em>emphasis</em>
	</body></html>`}

	md, err := conv.Convert(htmlContents)
	if err != nil {
		t.Fatalf("Convert failed: %v", err)
	}

	if !strings.Contains(md, "**bold**") {
		t.Errorf("expected **bold**, got: %s", md)
	}
	if !strings.Contains(md, "*italic*") {
		t.Errorf("expected *italic*, got: %s", md)
	}
	if !strings.Contains(md, "<u>underline</u>") {
		t.Errorf("expected <u>underline</u>, got: %s", md)
	}
	if !strings.Contains(md, "**strong**") {
		t.Errorf("expected **strong**, got: %s", md)
	}
	if !strings.Contains(md, "*emphasis*") {
		t.Errorf("expected *emphasis*, got: %s", md)
	}
}

func TestMarkdownConverterLists(t *testing.T) {
	conv := NewHTMLToMarkdownConverter()
	htmlContents := []string{`<html><body>
		<ul><li>Bullet A</li><li>Bullet B</li></ul>
		<ol><li>Number 1</li><li>Number 2</li></ol>
	</body></html>`}

	md, err := conv.Convert(htmlContents)
	if err != nil {
		t.Fatalf("Convert failed: %v", err)
	}

	if !strings.Contains(md, "- Bullet A") {
		t.Errorf("expected '- Bullet A' in output, got: %s", md)
	}
	if !strings.Contains(md, "1. Number 1") {
		t.Errorf("expected '1. Number 1' in output, got: %s", md)
	}
}

func TestMarkdownConverterTable(t *testing.T) {
	conv := NewHTMLToMarkdownConverter()
	htmlContents := []string{`<html><body>
		<table>
			<tr><th>Name</th><th>Age</th></tr>
			<tr><td>Alice</td><td>30</td></tr>
		</table>
	</body></html>`}

	md, err := conv.Convert(htmlContents)
	if err != nil {
		t.Fatalf("Convert failed: %v", err)
	}

	if !strings.Contains(md, "| Name") {
		t.Errorf("expected table header in output, got: %s", md)
	}
	if !strings.Contains(md, "---") {
		t.Errorf("expected table separator in output, got: %s", md)
	}
	if !strings.Contains(md, "Alice") {
		t.Errorf("expected 'Alice' in output, got: %s", md)
	}
}

func TestMarkdownConverterLinks(t *testing.T) {
	conv := NewHTMLToMarkdownConverter()
	htmlContents := []string{`<html><body>
		<a href="https://example.com">Example</a>
	</body></html>`}

	md, err := conv.Convert(htmlContents)
	if err != nil {
		t.Fatalf("Convert failed: %v", err)
	}

	if !strings.Contains(md, "[Example](https://example.com)") {
		t.Errorf("expected markdown link, got: %s", md)
	}
}

func TestMarkdownConverterImages(t *testing.T) {
	conv := NewHTMLToMarkdownConverter()
	htmlContents := []string{`<html><body>
		<img src="image.png" alt="My Image">
	</body></html>`}

	md, err := conv.Convert(htmlContents)
	if err != nil {
		t.Fatalf("Convert failed: %v", err)
	}

	if !strings.Contains(md, "![My Image](image.png)") {
		t.Errorf("expected markdown image, got: %s", md)
	}
}

func TestMarkdownConverterCodeBlocks(t *testing.T) {
	conv := NewHTMLToMarkdownConverter()
	htmlContents := []string{`<html><body>
		<code>inline code</code>
		<pre>block code</pre>
	</body></html>`}

	md, err := conv.Convert(htmlContents)
	if err != nil {
		t.Fatalf("Convert failed: %v", err)
	}

	if !strings.Contains(md, "`inline code`") {
		t.Errorf("expected inline code, got: %s", md)
	}
	if !strings.Contains(md, "```") {
		t.Errorf("expected code block, got: %s", md)
	}
}

func TestMarkdownConverterBlockquote(t *testing.T) {
	conv := NewHTMLToMarkdownConverter()
	htmlContents := []string{`<html><body>
		<blockquote>Quoted text</blockquote>
	</body></html>`}

	md, err := conv.Convert(htmlContents)
	if err != nil {
		t.Fatalf("Convert failed: %v", err)
	}

	if !strings.Contains(md, "> Quoted text") {
		t.Errorf("expected blockquote, got: %s", md)
	}
}

func TestMarkdownConverterHR(t *testing.T) {
	conv := NewHTMLToMarkdownConverter()
	htmlContents := []string{`<html><body><p>Before</p><hr><p>After</p></body></html>`}

	md, err := conv.Convert(htmlContents)
	if err != nil {
		t.Fatalf("Convert failed: %v", err)
	}

	if !strings.Contains(md, "---") {
		t.Errorf("expected horizontal rule, got: %s", md)
	}
}

func TestMarkdownConverterMultipleContents(t *testing.T) {
	conv := NewHTMLToMarkdownConverter()
	htmlContents := []string{
		`<html><body><h1>Page 1</h1></body></html>`,
		`<html><body><h1>Page 2</h1></body></html>`,
	}

	md, err := conv.Convert(htmlContents)
	if err != nil {
		t.Fatalf("Convert failed: %v", err)
	}

	if !strings.Contains(md, "# Page 1") || !strings.Contains(md, "# Page 2") {
		t.Errorf("expected both pages, got: %s", md)
	}
	// Should have a separator between pages
	if !strings.Contains(md, "---") {
		t.Errorf("expected separator between pages, got: %s", md)
	}
}

func TestConvertHTMLToMarkdownFile(t *testing.T) {
	tmpFile := filepath.Join(t.TempDir(), "output.md")
	htmlContents := []string{`<html><body><h1>File Test</h1><p>Content</p></body></html>`}

	if err := ConvertHTMLToMarkdown(htmlContents, tmpFile); err != nil {
		t.Fatalf("ConvertHTMLToMarkdown failed: %v", err)
	}

	data, err := os.ReadFile(tmpFile)
	if err != nil {
		t.Fatalf("Could not read output: %v", err)
	}

	content := string(data)
	if !strings.Contains(content, "# File Test") {
		t.Errorf("expected '# File Test' in file, got: %s", content)
	}
	// Should not have triple newlines (cleanup check)
	if strings.Contains(content, "\n\n\n") {
		t.Error("output contains triple newlines, cleanup failed")
	}
}
