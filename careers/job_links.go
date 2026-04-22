package careers

import (
	"encoding/json"
	"fmt"
	"strings"
)

const jobLinksSysPrompt = `You are a web scraping assistant helping identify individual job posting links on company career pages.
Given a base URL and a list of links scraped from a careers page, identify only the links that point to individual job postings.

A link is an individual job posting if it:
- Points to a specific role, position, or programme (e.g. "/jobs/123-junior-engineer", "/p/abc123-software-developer")
- Contains a job title or identifier in the path

A link is NOT an individual job posting if it:
- Is the careers or jobs listing page itself (e.g. "/careers", "/jobs")
- Points to general navigation (e.g. "/about", "/contact", "/home")
- Points to a category or department filter (e.g. "/jobs?department=engineering")
- Is an external link unrelated to job applications

If no individual job posting links are found, return an empty array.
Return ONLY valid JSON matching the provided schema, no explanation, no markdown.`

var jobLinksSchema = json.RawMessage(`{
    "type": "object",
    "properties": {
        "posting_urls": {"type": "array", "items": {"type": "string"}}
    },
    "required": ["posting_urls"]
}`)

type JobLinksResult struct {
	PostingURLs []string `json:"posting_urls"`
}

func FindJobPostingLinks(baseURL string, links []string) ([]string, error) {
	userPrompt := fmt.Sprintf("Base URL: %s\n\nLinks found on careers page:\n%s",
		baseURL, strings.Join(links, "\n"))

	raw, err := Chat(userPrompt, jobLinksSysPrompt, jobLinksSchema)
	if err != nil {
		return nil, err
	}

	var result JobLinksResult
	if err := json.Unmarshal([]byte(raw), &result); err != nil {
		fmt.Printf("RAW: %s\n", raw)
		return nil, err
	}

	return result.PostingURLs, nil
}