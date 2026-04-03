package main

import (
	"log"
	"net/http"
	"fmt"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func fetchDoc(url string) *goquery.Document {
	res, err := http.Get(url);

	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		log.Fatalf("status code error %d %s", res.StatusCode, res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	return doc
}

func getCompaniesInPage(doc *goquery.Document, res *map[string]string) {
	doc.Find("h3.voffset0 a").Each(func (i int, s *goquery.Selection) {

		href, exists := s.Attr("href")

		if exists {
			name := s.Text()
			(*res)[name] = href
		}
	})
}

func getNumberOfPages(url string) int {
	var n int
	var err error
	doc := fetchDoc(url)

	doc.Find("ul.pagination li:nth-last-child(2)").Each(func(i int, s *goquery.Selection) {
		n, err = strconv.Atoi(strings.TrimSpace(s.Text()))
		if err != nil {
			log.Fatal(err)
		}
	})

	return n
}

func getComanieWebsite(url string) string {
	doc := fetchDoc(url)

	ret := doc.Find("div.center_mobile.hidden-xs.company_add_details a:last-child").First()
	href, exists := ret.Attr("href")

	if exists {
		return href
	}

	return ""
}


func main() {
	const home string = "https://pt.teamlyzer.com"
	const companySearch string = "https://pt.teamlyzer.com/companies/"
	const pageQuery string = "?page="
	
	n := getNumberOfPages(companySearch)

	companyUrl := make(map[string]string)
	for i := 1; i <= n; i++ {
		page := fmt.Sprint(companySearch, pageQuery, i)
		doc := fetchDoc(page)
		getCompaniesInPage(doc, &companyUrl)
	}

	companyWebsite := make(map[string]string)
	for c, url := range companyUrl {
		fullUrl := fmt.Sprint(home, url)
		companyWebsite[c] = getComanieWebsite(fullUrl)
	}

	fmt.Println(companyWebsite)

}