# job_finder

A pipeline for scraping job postings from Portuguese tech companies, built in Go with a local LLM for intelligent extraction.

## What was built

### Scraper (`scraper/`)
- **`teamlyzer.go`** — scrapes [Teamlyzer](https://pt.teamlyzer.com) to collect company names and their websites, outputting `data/companies.csv`
- **`links.go`** — generic link extractor, fetches a page and returns all `<a href>` links

### Careers (`careers/`)
- **`ollama.go`** — HTTP client for the local Ollama API, wraps the `/api/chat` endpoint with structured JSON output support
- **`finder.go`** — uses the LLM to identify careers, internships, and graduate programme pages from a list of homepage links
- **`job_links.go`** — uses the LLM to identify individual job posting URLs from a list of careers page links
- **`extractor.go`** — uses the LLM to extract structured job posting data from raw page text, outputting a `JobPosting` struct

### Data (`data/`)
- **`companies.csv`** — ~4700 Portuguese tech companies with their websites, scraped from Teamlyzer

## Prerequisites

- [Go 1.21+](https://go.dev/dl/)
- [Ollama](https://ollama.com) running locally with `gemma4:e2b` pulled

```bash
ollama pull gemma4:e2b
ollama serve
```

## How to run

**Run the full pipeline:**

```bash
go run main.go
```

Results are written to `data/jobs.json`.

**Build a binary:**

```bash
go build -o job_finder .
./job_finder
```

## What is yet to be done

- [ ] Tests for each pipeline step
- [ ] Migrate `data/jobs.json` output to a SQLite database for querying and deduplication
- [ ] Rate limiting and retry logic for failed requests
- [ ] Switch from local `gemma4:e2b` to Anthropic/OpenAI API for production runs at scale

## Performance notes

The local `gemma4:e2b` model runs at approximately 43 seconds per job posting on a GTX 1650 (4GB VRAM). For the full dataset of ~4700 companies this is not viable locally. The Anthropic or OpenAI API is recommended for production runs — estimated cost under $1 with GPT-4o-mini batch processing.