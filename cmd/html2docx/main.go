package main

import (
	"fmt"
	"os"

	"github.com/achiket/html2docx/converter"
)

func main() {
	htmlContents := []string{`<!DOCTYPE html>
<html>
<head>
    <title>1990s Style Raw HTML Site</title>
</head>
<body bgcolor="#FFFFFF" text="#000000" link="#0000FF" vlink="#800080">

    <center>
        <h1>Website Title</h1>
        <table border="0" width="100%">
            <tr>
                <td align="center">
                    <a href="#"><b>Articles</b></a> | 
                    <a href="#"><b>Media</b></a> | 
                    <a href="#"><b>Contact</b></a>
                </td>
            </tr>
        </table>
    </center>
 <hr size="2" width="100%" noshade>
    <table border="0" width="80%" align="center">
        <tr>
            <td>
                <h2>Latest Article</h2>
                <p>This is a paragraph of text. Back then, we used the <b>strong</b> tag for bold and <i>italic</i> tag for emphasis. </p>

                <h3>Understanding Semantic HTML</h3>
                <p>Below is an example of an unordered list, which was one of the few ways to create vertical spacing:</p>
                <ul>
                    <li>Structure (using tables)</li>
                    <li>Meaning (using headers)</li>
                    <li>Accessibility (alt text on images)</li>
                </ul>

                <p>Data was always presented in bordered tables:</p>
                
                <table border="1" cellpadding="5" cellspacing="0" width="100%">
                    <tr bgcolor="#CCCCCC">
                        <th><font face="Arial">Name</font></th>
                        <th><font face="Arial">Position</font></th>
                        <th><font face="Arial">Office</font></th>
                    </tr>
                    <tr>
                        <td>Jane Doe</td>
                        <td>Software Engineer</td>
                        <td>Remote</td>
                    </tr>
                    <tr>
                        <td>John Smith</td>
                        <td>Product Manager</td>
                        <td>New York</td>
                    </tr>
                </table>

                <br>
                <center>
                    <img src="https://via.placeholder.com/150" alt="Old school placeholder" border="2">
                </center>
            </td>
        </tr>
    </table>

    <hr size="2" width="100%" noshade>
	<footer>
    <center>
        <font size="2">
            Copyright &copy; 2026 Raw HTML Company. All rights reserved.<br>
            <i>Best viewed in Netscape Navigator at 800x600 resolution.</i>
        </font>
    </center>
	</footer>

</body>
</html>`}

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
