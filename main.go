package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"job_finder/careers"
	"job_finder/scraper"
	"os"
)

type Result struct {
	Company  string               `json:"company"`
	Website  string               `json:"website"`
	Postings []careers.JobPosting `json:"postings"`
}

func main() {
	file, err := os.Open("data/companies.csv")
	if err != nil {
		fmt.Println(fmt.Errorf("failed to open companies.csv: %w", err))
		return
	}
	defer file.Close()

	reader := csv.NewReader(file)
	rows, err := reader.ReadAll()
	if err != nil {
		fmt.Println(fmt.Errorf("failed to read companies.csv: %w", err))
		return
	}

	var results []Result

	for i, row := range rows {
		// skip header
		if i == 0 {
			continue
		}

		company := row[0]
		website := row[1]

		if website == "" {
			fmt.Printf("[%d] Skipping %s — no website\n", i, company)
			continue
		}

		fmt.Printf("[%d] Processing %s (%s)\n", i, company, website)

		// Phase 1 — find careers pages
		links, err := scraper.GetLinks(website)
		if err != nil {
			fmt.Printf("[%d] Failed to get links for %s: %v\n", i, company, err)
			continue
		}

		finderResult, err := careers.FindCareersPages(website, links)
		if err != nil {
			fmt.Printf("[%d] Failed to find careers pages for %s: %v\n", i, company, err)
			continue
		}

		// deduplicate careers URLs across all categories
		seen := map[string]bool{}
		var careersURLs []string
		for _, u := range append(append(finderResult.CareersURL, finderResult.InternshipsURL...), finderResult.GraduatesURL...) {
			if !seen[u] {
				seen[u] = true
				careersURLs = append(careersURLs, u)
			}
		}

		if len(careersURLs) == 0 {
			fmt.Printf("[%d] No careers pages found for %s\n", i, company)
			continue
		}

		var postings []careers.JobPosting

		// Phase 2 — find and extract job postings
		for _, careersURL := range careersURLs {
			careersLinks, err := scraper.GetLinks(careersURL)
			if err != nil {
				fmt.Printf("[%d] Failed to get links from careers page %s: %v\n", i, careersURL, err)
				continue
			}

			postingURLs, err := careers.FindJobPostingLinks(careersURL, careersLinks)
			if err != nil {
				fmt.Printf("[%d] Failed to find job posting links for %s: %v\n", i, company, err)
				continue
			}

			for _, postingURL := range postingURLs {
				doc, err := scraper.FetchDoc(postingURL)
				if err != nil {
					fmt.Printf("[%d] Failed to fetch posting %s: %v\n", i, postingURL, err)
					continue
				}

				text := doc.Find("body").Text()
				posting, err := careers.ExtractJobPosting(text)
				if err != nil {
					fmt.Printf("[%d] Failed to extract posting %s: %v\n", i, postingURL, err)
					continue
				}

				posting.PostingURL = postingURL
				posting.SourceURL = careersURL
				posting.Company = company

				fmt.Println(posting.String())
				postings = append(postings, posting)
			}
		}

		results = append(results, Result{
			Company:  company,
			Website:  website,
			Postings: postings,
		})
	}

	// write results to JSON
	output, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		fmt.Println(fmt.Errorf("failed to marshal results: %w", err))
		return
	}

	if err := os.WriteFile("data/jobs.json", output, 0644); err != nil {
		fmt.Println(fmt.Errorf("failed to write jobs.json: %w", err))
		return
	}

	fmt.Printf("\nDone. %d companies processed, results written to data/jobs.json\n", len(results))
}