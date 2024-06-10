package scrapper

import (
// "fmt"
)

func RunScrapper() {
	url := "https://www.qatarcontact.com/search?page=1"

	HTMLData := FetchData(url)
	ParsedHTML := HTMLParser(HTMLData)
	GenerateData(ParsedHTML)

	// page Number
}
