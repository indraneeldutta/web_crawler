package services

import (
	"bytes"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html"
)

var doctypes = make(map[string]string)

func init() {
	doctypes["HTML 4.01 Strict"] = `"-//W3C//DTD HTML 4.01//EN"`
	doctypes["HTML 4.01 Transitional"] = `"-//W3C//DTD HTML 4.01 Transitional//EN"`
	doctypes["HTML 4.01 Frameset"] = `"-//W3C//DTD HTML 4.01 Frameset//EN"`
	doctypes["XHTML 1.0 Strict"] = `"-//W3C//DTD XHTML 1.0 Strict//EN"`
	doctypes["XHTML 1.0 Transitional"] = `"-//W3C//DTD XHTML 1.0 Transitional//EN"`
	doctypes["XHTML 1.0 Frameset"] = `"-//W3C//DTD XHTML 1.0 Frameset//EN"`
	doctypes["XHTML 1.1"] = `"-//W3C//DTD XHTML 1.1//EN"`
	doctypes["HTML 5"] = "!DOCTYPE html"
	doctypes["HTML5"] = `<html lang="en-US">`
}

func checkDoctype(html string) string {
	var version = "UNKNOWN"

	for doctype, matcher := range doctypes {
		match := strings.Contains(html, matcher)
		if match {
			version = doctype
			break
		}
	}
	return version
}

func DetectHTML(r []byte) (string, error) {
	var htmlVersion string

	reader := bytes.NewReader(r)

	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		return htmlVersion, err
	}

	html, err := doc.Html()
	if err != nil {
		return htmlVersion, err
	}

	htmlVersion = checkDoctype(html)

	return htmlVersion, nil
}

func isTitleElement(n *html.Node) bool {
	return n.Type == html.ElementNode && n.Data == "title"
}

func traverse(n *html.Node) (string, bool) {
	if isTitleElement(n) {
		return n.FirstChild.Data, true
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		result, ok := traverse(c)
		if ok {
			return result, ok
		}
	}

	return "", false
}

func GetHtmlTitle(r []byte) (string, bool) {
	reader := bytes.NewReader(r)
	doc, err := html.Parse(reader)
	if err != nil {
		panic("Fail to parse html")
	}

	return traverse(doc)
}
