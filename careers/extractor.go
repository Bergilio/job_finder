package careers

import (
	"encoding/json"
	"fmt"
	"strings"
)

var jobPostingSchema = json.RawMessage(`{
    "type": "object",
    "properties": {
        "company":          {"type": "string"},
        "source_url":       {"type": "string"},
        "posting_url":      {"type": "string"},
        "title":            {"type": "string"},
        "type":             {"type": "string", "enum": ["part-time", "full-time", "summer-internship", "curricular-internship", "professional-internship", "graduate"]},
        "seniority":        {"type": "string", "enum": ["junior", "mid", "senior", "lead", "principal", "unknown"]},
        "remote":           {"type": "string", "enum": ["remote", "hybrid", "on-site", "unknown"]},
        "locations":        {"type": "array", "items": {"type": "string"}},
        "technologies":     {"type": "array", "items": {"type": "string"}},
        "languages":        {"type": "array", "items": {"type": "string"}},
        "posted_date":      {"type": ["string", "null"]},
        "deadline":         {"type": ["string", "null"]},
        "salary_min":       {"type": ["integer", "null"]},
        "salary_max":       {"type": ["integer", "null"]},
        "salary_currency":  {"type": ["string", "null"]},
        "full_description": {"type": "string"}
    },
    "required": ["company", "source_url", "posting_url", "title", "type", "seniority", "remote", "locations", "technologies", "languages", "posted_date", "deadline", "salary_min", "salary_max", "salary_currency", "full_description"]
}`)

const extractorSysPrompt = `You are a job posting data extractor for the Portuguese job market.
Given the text content of a job posting page, extract all relevant information into the provided JSON schema.

Follow these rules strictly:
- "type" must be one of: "part-time", "full-time", "summer-internship", "curricular-internship", "professional-internship", "graduate"
- "seniority" must be one of: "junior", "mid", "senior", "lead", "principal", "unknown"
- "remote" must be one of: "remote", "hybrid", "on-site", "unknown"
- "technologies" should list specific tools, languages, and frameworks mentioned (e.g. "Go", "PostgreSQL", "Docker")
- "languages" should list spoken languages only (e.g. "Portuguese", "English"), not programming languages
- "full_description" should be the complete job description text, cleaned of HTML artifacts
- "salary_min" and "salary_max" should be annual amounts in the local currency when possible
- "posted_date" and "deadline" should be in ISO 8601 format (YYYY-MM-DD) when possible, otherwise null
- "company" is the name of the hiring company, usually found in the page header
- "locations" is a list of cities where the job is based, e.g. ["Lisbon", "Porto"]
- "remote" should be inferred from phrases like "hybrid", "remote", "on-site", "in person"
- "title" is the job title, usually the largest heading on the page, it can also be the program title

For Portuguese-specific internship types:
- "curricular-internship" is an estágio curricular — part of a university degree, unpaid or minimally paid
- "professional-internship" is an estágio profissional — IEFP programme for recent graduates
- "summer-internship" is a summer placement, typically 1 to 3 months,
- "graduate" is a graduate programme for recent university graduates

If a field cannot be determined from the posting, use null for nullable fields, "unknown" for enum fields, and empty arrays for array fields.
Return ONLY valid JSON matching the provided schema, no explanation, no markdown.`


type RemotePolicy string
const (
	Remote RemotePolicy="remote"
	Hybrid RemotePolicy="hybrid"
	OnSite RemotePolicy="on-site"
	Unknown RemotePolicy="unknown"
)

type TypePolicy string
const(
	PartTime TypePolicy="part-time"
	FullTime TypePolicy="full-time"
	SummerInternship TypePolicy="summer-internship"
	CurricularInternship TypePolicy="curricular-internship"
	ProfessionalInternship TypePolicy="professional-internship"
	Graduate TypePolicy="graduate"
)

type SeniorityLevel string
const (
    Junior      SeniorityLevel = "junior"
    Mid         SeniorityLevel = "mid"
    Senior      SeniorityLevel = "senior"
    Lead        SeniorityLevel = "lead"
    Principal   SeniorityLevel = "principal"
	UnknownSL   SeniorityLevel = "unknown"
)

type JobPosting struct {
	Company string `json:"company"`
	SourceURL string `json:"source_url"`
	PostingURL string `json:"posting_url"`

	Title string `json:"title"`
	Type TypePolicy `json:"type"`
	Seniority SeniorityLevel `json:"seniority"`
	Remote RemotePolicy `json:"remote"`

	Locations []string `json:"locations"`

	Technologies []string `json:"technologies"`
	Languages []string `json:"languages"`

	PostedDate *string `json:"posted_date"`
	Deadline *string `json:"deadline"`

	SalaryMin *int `json:"salary_min"`
	SalaryMax *int `json:"salary_max"`
	SalaryCurrency *string `json:"salary_currency"`

	FullDescription string `json:"full_description"`
}

func (j JobPosting) String() string {
	return fmt.Sprintf(
		"Title:      %s\nCompany:    %s\nType:       %s\nSeniority:  %s\nRemote:     %s\nLocations:  %s\nTech:       %s\nLanguages:  %s\nPosted:     %s\nDeadline:   %s\n",
        j.Title,
        j.Company,
        j.Type,
        j.Seniority,
        j.Remote,
        strings.Join(j.Locations, ", "),
        strings.Join(j.Technologies, ", "),
        strings.Join(j.Languages, ", "),
        ptrStr(j.PostedDate),
        ptrStr(j.Deadline),
	)
}

func ptrStr(str *string) string {
	if str == nil {
		return ""
	}
	return *str
}

func ExtractJobPosting(page string) (JobPosting, error) {
	userPrompt := "Job posting page here:\n" + page

	raw, err := Chat(userPrompt, extractorSysPrompt, jobPostingSchema)
	if err != nil {
		return JobPosting{}, err
	}

	var jobPosting JobPosting
	err = json.Unmarshal([]byte(raw), &jobPosting)
	if err != nil {
		fmt.Printf("%s\n", raw)
		return JobPosting{}, err
	}

	return jobPosting, nil
}