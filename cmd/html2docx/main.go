package main

import (
	"fmt"
	"os"

	"github.com/Achiket123/html2docx/converter"
)

func main() {
	htmlContents := []string{}

	content, err := os.ReadFile("test.html")
	if err != nil {
		fmt.Printf("Error reading HTML file: %v\n", err)
		os.Exit(1)
	}
	htmlContents = append(htmlContents, string(content))

	// Convert to Markdown
	mdFile := "output.md"
	if err := converter.ConvertHTMLToMarkdown(htmlContents, mdFile); err != nil {
		fmt.Printf("Error converting to Markdown: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Successfully created: %s\n", mdFile)

	// Create converter and convert HTML to DOCX
	conv := converter.NewHTMLToDocxConverter()
	if err := conv.Convert(htmlContents); err != nil {
		fmt.Printf("Error converting HTML: %v\n", err)
		os.Exit(1)
	}

	docxFile := "output.docx"
	if err := conv.SaveToFile(docxFile); err != nil {
		fmt.Printf("Error saving DOCX: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Successfully created: %s\n", docxFile)

	// Convert HTML directly to PDF (no LibreOffice needed)
	pdfConv := converter.NewHTMLToPDFConverter()
	if err := pdfConv.Convert(htmlContents); err != nil {
		fmt.Printf("Error converting to PDF: %v\n", err)
		os.Exit(1)
	}
	pdfFile := "output.pdf"
	if err := pdfConv.SaveToFile(pdfFile); err != nil {
		fmt.Printf("Error saving PDF: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Successfully created: %s\n", pdfFile)
}
