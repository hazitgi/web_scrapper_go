package scrapper

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func HTMLParser(HTMLData []byte) *goquery.Document {
	fmt.Println("Parsing HTML")
	/* Load HTML Document */
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(HTMLData))
	if err != nil {
		log.Fatal(err)
	}
	return doc
}

// finding from html
func FindIntoDom(HTML *goquery.Document, tag string) *goquery.Selection {
	pageNumbers := HTML.Find(tag)
	return pageNumbers
}

func GenerateData(ParsedHTML *goquery.Document) {
	HtmlElement := FindIntoDom(ParsedHTML, ".page-numbers")
	pageCount := HtmlElement.Eq(HtmlElement.Length() - 2).Text()
	lastPage, err := strconv.Atoi(pageCount)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Page Count: ", lastPage)
	details := []Detail{}

	for i := 1; i <= 1; i++ {
		fmt.Println("Page Number: ", i)
		url := fmt.Sprintf("https://www.qatarcontact.com/search?page=%d", i)
		HTMLData := FetchData(url)
		ParsedHTML := HTMLParser(HTMLData)

		detailsList := ParsedHTML.Find(".atbd_single_listing_wrapper")

		detailsList.Each(func(i int, s *goquery.Selection) {
			emailDom, exits := s.Find(".atbd_listing_title a").Attr("href")
			if !exits {
				return
			}

			emailData := FetchData(emailDom)
			newParsedHTML := HTMLParser(emailData)

			email := newParsedHTML.Find(".atbd_info").Eq(2).Text()
			name := s.Find(".atbd_listing_info .atbd_content_upper .atbd_listing_data_list ul li").Eq(0).Text()
			companyName := s.Find(".atbd_listing_title a").Eq(0).Text()
			category := s.Find(".atbd_listing_category a").Eq(0).Text()
			mobile := s.Find(".atbd_listing_data_list ul li").Eq(1).Text()
			visitCardUrl := func() string {
				if img, exists := s.Find(".atbd_listing_image a:nth-child(2) img").Attr("src"); exists {
					return img
				}
				return ""
			}()
			detail := Detail{
				Name:         strings.TrimSpace(name),
				Mobile:       strings.TrimSpace(mobile),
				CompanyName:  strings.TrimSpace(companyName),
				Category:     strings.TrimSpace(category),
				VisitCardUrl: strings.TrimSpace(visitCardUrl),
				Email:        strings.TrimSpace(email),
			}
			details = append(details, detail)
		})
	}
	saveToJson(details)
	saveAsCsv(details)
}

func saveToJson(details []Detail) {
	body, err := json.MarshalIndent(details, "    ", "")
	if err != nil {
		log.Fatal(err)
	}
	file, err := os.Create("data.json")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	_, err = file.Write(body)

	if err != nil {
		log.Fatal(err)
	}
}

func saveAsCsv(details []Detail) {
	file, err := os.Create("./data.csv")

	if err != nil {
		fmt.Println("Error creating CSV file : ", err)
		return
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Writing CSV header
	header := []string{"email", "name", "company", "category", "mobile", "visit_card_url"}

	if err := writer.Write(header); err != nil {
		fmt.Println("Error writing CSV header : ", err)
		return
	}

	// writing posts to CSV
	for _, detail := range details {
		row := []string{
			detail.Email,
			detail.Name,
			detail.CompanyName,
			detail.Category,
			detail.Mobile,
			detail.VisitCardUrl,
		}

		if err := writer.Write(row); err != nil {
			fmt.Println("Error writing post to CSV:", err)
			return
		}
	}
	fmt.Println("CSV file created successfully!")

}
