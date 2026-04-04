package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

var client = &http.Client{
	Timeout: 30 * time.Second,
}

func fetchDoc(url string) (*goquery.Document, error) {
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) Chrome/120.0.0.0")

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("Status code error %d %s", res.StatusCode, res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, err
	}

	return doc, nil
}

func getCompaniesInPage(url string, res map[string]string) error {
	doc, err := fetchDoc(url)
	if err != nil {
		return err
	}

	doc.Find("h3.voffset0 a").Each(func (i int, s *goquery.Selection) {

		href, exists := s.Attr("href")

		if exists {
			name := s.Text()
			res[name] = href
		}
	})

	return nil
}

func getNumberOfPages(url string) (int, error) {
	doc, err := fetchDoc(url)
	if err != nil {
		return 0, err
	}

	ret := doc.Find("ul.pagination li:nth-last-child(2)").First()
	return strconv.Atoi(strings.TrimSpace(ret.Text()))
}

func getCompanieWebsite(url string) (string, error) {
	doc, err := fetchDoc(url)
	if err != nil {
		return "", err
	}

	ret := ""
	doc.Find("div.center_mobile.hidden-xs.company_add_details a").Each(func (i int, s *goquery.Selection){
		href, exists := s.Attr("href")

		if exists && s.Text() == "Website" {
			ret = href
		}
	})

	return ret, nil
}


func main() {
	const home string = "https://pt.teamlyzer.com"
	const companySearch string = "https://pt.teamlyzer.com/companies/"
	const pageQuery string = "?page="
	
	n, err := getNumberOfPages(companySearch)
	if err != nil {
		fmt.Println(fmt.Errorf("Unable to extract the amout of pages: %w", err))
		return
	}

	companyUrl := make(map[string]string)
	for i := 1; i <= n; i++ {
		page := companySearch + pageQuery + strconv.Itoa(i)
		err = getCompaniesInPage(page, companyUrl)

		if err != nil {
			fmt.Println(fmt.Errorf("Unable to reach page number %d: %w\n", i, err))
		}

		time.Sleep(200 * time.Millisecond)
	}

	companyWebsite := make(map[string]string)
	for c, url := range companyUrl {
		fullUrl := home + url
		companyWebsite[c], err = getCompanieWebsite(fullUrl)
		
		if err != nil {
			fmt.Println(fmt.Errorf("Unable to extract %s's website: %w", c, err))
		}

		time.Sleep(200 * time.Millisecond)
	}

	for k, v := range companyWebsite {
		fmt.Printf("%s : %s\n", k, v)
	}

	fmt.Printf("Amount of companies: %d\n", len(companyWebsite))
}