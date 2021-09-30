package services

import (
	"bytes"
	"log"
	"net/http"
	"strings"
	"sync"

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
	var wg sync.WaitGroup
	for _, value := range links {
		wg.Add(1)
		go func(w *sync.WaitGroup, v string) {
			if strings.Contains(v, "http") {
				resp, err := http.Get(v)
				if err != nil {
					log.Println(err)
					inaccLinks = append(inaccLinks, v)
				} else {
					resp.Body.Close()
					if resp.StatusCode != 200 {
						inaccLinks = append(inaccLinks, v)
					}
				}

			} else {
				inaccLinks = append(inaccLinks, v)
			}
			w.Done()
		}(&wg, value)
	}
	wg.Wait()
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
