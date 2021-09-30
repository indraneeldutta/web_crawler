package services

import (
	"bytes"
	"log"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func GetLinks(r []byte, url string) ([]string, []string, error) {
	var internalLinks, externalLinks []string

	reader := bytes.NewReader(r)
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		log.Println(err)
		return internalLinks, externalLinks, err
	}

	var checkerString string
	checker := strings.Split(strings.Split(url, "//")[1], ".")
	if strings.Contains(checker[0], "www") {
		checkerString = checker[1]
	} else {
		checkerString = checker[0]
	}

	doc.Find("[href]").Each(func(index int, item *goquery.Selection) {
		href, _ := item.Attr("href")

		if strings.Contains(href, checkerString) || !strings.Contains(href, "https://") {
			internalLinks = append(internalLinks, href)
		} else {
			externalLinks = append(externalLinks, href)
		}
	})

	return internalLinks, externalLinks, err
}

func CheckInaccessibleLinks(links []string) []string {
	var inaccLinks []string
	for _, value := range links {
		if strings.Contains(value, "http") {
			resp, err := http.Get(value)
			if err != nil {
				log.Println(err)
				inaccLinks = append(inaccLinks, value)
				continue
			}
			resp.Body.Close()
			if resp.StatusCode != 200 {
				inaccLinks = append(inaccLinks, value)
			}
		} else {
			inaccLinks = append(inaccLinks, value)
		}
	}
	return inaccLinks
}

func CheckLoginFormExists(r []byte) bool {
	reader := bytes.NewReader(r)
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		log.Println(err)
	}

	html, err := doc.Html()
	if err != nil {
		log.Println(err)
		return false
	}

	if strings.ContainsAny(html, "sign in") || strings.ContainsAny(html, "login") {
		return true
	}
	return false
}
