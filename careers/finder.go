package careers

import (
	"fmt"
	"strings"
	"encoding/json"
)

const finderSysPrompt = `You are a web scraping assistant helping find employment-related pages on company websites.
Given a base URL and a list of links scraped from that website, identify the most relevant URLs for each category.
A single URL may appear in multiple categories if it covers all of them.
Each category may have multiple URLs if the company has separate pages for different programmes or regions.
If a category is not found, return an empty array.
Return ONLY valid JSON matching the provided schema, no explanation, no markdown.

When identifying URLs, look for these patterns:
- "careers_url": links containing words like "careers", "jobs", "work-with-us", "join-us", "vacancies", "opportunities", "trabalha-connosco", "emprego"
- "internships_url": links containing words like "internship", "intern", "estágio", "estagio", "trainee"
- "graduates_url": links containing words like "graduate", "graduates", "junior", "campus", "university", "recent-graduates"

Prefer full absolute URLs over relative paths. If only a relative path is available, do not attempt to resolve it.`

var finderResultSchema = json.RawMessage(`{
    "type": "object",
    "properties": {
        "careers_url":     {"type": "array", "items": {"type": "string"}},
        "internships_url": {"type": "array", "items": {"type": "string"}},
        "graduates_url":   {"type": "array", "items": {"type": "string"}}
    },
    "required": ["careers_url", "internships_url", "graduates_url"]
}`)

type FinderResult struct {
	CareersURL     []string `json:"careers_url"`
    InternshipsURL []string `json:"internships_url"`
    GraduatesURL   []string `json:"graduates_url"`
}

func FindCareersPages(baseURL string, links []string) (FinderResult, error) {
	userPrompt := fmt.Sprintf("Base URL: %s\n\nLinks found on page:\n%s", baseURL, strings.Join(links, "\n"))

	raw, err := Chat(userPrompt, finderSysPrompt, finderResultSchema)
	if err != nil {
		return FinderResult{}, err
	}

	var finderResult FinderResult
	err = json.Unmarshal([]byte(raw), &finderResult)
	if err != nil {
		fmt.Printf("%s\n", raw)
		return FinderResult{}, err
	}

	return finderResult, nil
}