package scraper

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
	"os"
	"encoding/csv"

	"github.com/PuerkitoBio/goquery"
)

var client = &http.Client{
	Timeout: 60 * time.Second,
}

func FetchDoc(url string) (*goquery.Document, error) {
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

func GetCompaniesInPage(url string, res map[string]string) error {
	doc, err := FetchDoc(url)
	if err != nil {
		return err
	}

	doc.Find("h3.voffset0 a").Each(func (i int, s *goquery.Selection) {

		href, exists := s.Attr("href")

		if exists {
			name := strings.TrimSpace(s.Text())
			res[name] = href
		}
	})

	return nil
}

func GetNumberOfPages(url string) (int, error) {
	doc, err := FetchDoc(url)
	if err != nil {
		return 0, err
	}

	ret := doc.Find("ul.pagination li:nth-last-child(2)").First()
	return strconv.Atoi(strings.TrimSpace(ret.Text()))
}

func GetCompanyWebsite(url string) (string, error) {
	doc, err := FetchDoc(url)
	if err != nil {
		return "", err
	}

	ret := ""
	doc.Find("div.center_mobile.hidden-xs.company_add_details a").Each(func (i int, s *goquery.Selection){
		href, exists := s.Attr("href")

		if exists && s.Text() == strings.TrimSpace("Website") {
			ret = href
		}
	})

	return ret, nil
}

func GetTeamlyzerCompanies() {
	const home string = "https://pt.teamlyzer.com"
	const companySearch string = "https://pt.teamlyzer.com/companies/"
	const pageQuery string = "?page="

	file, err := os.Create("data/companies.csv")
	if err != nil {
		fmt.Println(fmt.Errorf("Unable to create file: %w", err))
		return
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	writer.Write([]string {"Company", "Website"})
	
	n, err := GetNumberOfPages(companySearch)
	if err != nil {
		fmt.Println(fmt.Errorf("Unable to extract the amout of pages: %w", err))
		return
	}

	companyUrl := make(map[string]string)
	for i := 1; i <= n; i++ {
		page := companySearch + pageQuery + strconv.Itoa(i)
		err = GetCompaniesInPage(page, companyUrl)

		if err != nil {
			fmt.Println(fmt.Errorf("Unable to reach page number %d: %w\n", i, err))
		}

		time.Sleep(2 * time.Second)
	}

	for c, url := range companyUrl {
		fullUrl := home + url
		website, err := GetCompanyWebsite(fullUrl)
		
		if err != nil {
			fmt.Println(fmt.Errorf("Unable to extract %s's website: %w", c, err))
		}

		writer.Write([]string {c, website})
		writer.Flush()

		time.Sleep(2 * time.Second)
	}
}