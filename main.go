package main

import (
	"log"
	"net/http"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func getHttpBody(addr string) io.ReadCloser {
	res, err := http.Get(addr);

	if err != nil {
		log.Fatal(err)
	}

	if res.StatusCode != 200 {
		log.Fatalf("status code error %d %s", res.StatusCode, res.Status)
	}

	return res.Body
}

func getCompaniesInPage(doc *goquery.Document) {
	doc.Find("h3.voffset0 a").Each(func (i int, s *goquery.Selection) {

		href, exists := s.Attr("href")

		if exists {
			name := s.Text()
			fmt.Printf("Company: %s | href: %s\n", name, href)
		}
	})
}


func main() {
	const addr string = "https://pt.teamlyzer.com/companies/"
	const pageQuery string = "?page="
	body := getHttpBody(addr)
	defer body.Close()

	doc, err := goquery.NewDocumentFromReader(body)
	if err != nil {
		log.Fatal(err)
	}

	var n int
	doc.Find("ul.pagination li:nth-last-child(2)").Each(func(i int, s *goquery.Selection) {
		n, err = strconv.Atoi(strings.TrimSpace(s.Text()))
		if err != nil {
			log.Fatal(err)
		}
	})

	for i := 1; i <= n; i++ {
		page := fmt.Sprint(addr, pageQuery, i)

		body = getHttpBody(page)
		defer body.Close()
		
		doc, err = goquery.NewDocumentFromReader(body)
		if err != nil {
			log.Fatal(err)
		}

		getCompaniesInPage(doc)
	}
}