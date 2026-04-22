package scraper

import (
	"github.com/PuerkitoBio/goquery"
)

func GetLinks(url string) ([]string, error) {
	doc, err := FetchDoc(url)
	if err != nil {
		return []string{}, err
	}

	links := []string{}
	doc.Find("a").Each(func (i int, s *goquery.Selection){
		href, exists := s.Attr("href")

		if exists {
			links = append(links, href)
		}
	})

	return links, nil
}
