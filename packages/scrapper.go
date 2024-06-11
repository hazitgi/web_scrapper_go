package scrapper

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
)

type Detail struct {
	Name         string
	Mobile       string
	CompanyName  string
	Category     string
	VisitCardUrl string
	Email        string
}

type Response interface {
	io.Reader
}

type HTMLResponse struct {
	Data []byte
}

var DetailsChannel = make(chan Detail)

// Read method to implement io.Reader interface
func (r *HTMLResponse) Read(p []byte) (n int, err error) {
	copy(p, r.Data)
	if len(r.Data) > len(p) {
		r.Data = r.Data[len(p):]
		return len(p), nil
	}
	n = len(r.Data)
	r.Data = nil
	return n, io.EOF
}

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
	return HTML.Find(tag)
}

func GenerateData(ParsedHTML *goquery.Document, wg *sync.WaitGroup) {
	defer wg.Done()
	HtmlElement := FindIntoDom(ParsedHTML, ".page-numbers")
	pageCount := HtmlElement.Eq(HtmlElement.Length() - 2).Text()
	lastPage, err := strconv.Atoi(pageCount)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Page Count: ", lastPage)

	for i := 1; i < lastPage; i++ {
		wg.Add(1)
		go func(page int) {
			defer wg.Done()

			fmt.Println("Page Number: ", page)
			url := fmt.Sprintf("https://www.qatarcontact.com/search?page=%d", page)
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
				DetailsChannel <- detail
			})
		}(i)
	}
}

func saveToJson(details []Detail) {
	body, err := json.MarshalIndent(details, "", " ")
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

func FetchData(url string) []byte {
	fmt.Println("Fetching Data from URL")
	res, err := http.Get(url)
	if err != nil {
		fmt.Println("Error while fetching data", err)
		return nil
	}

	defer res.Body.Close()

	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	// Read response body
	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println("Error in read response body", err)
	}
	return body
}

func RunScrapper() {
	url := "https://www.qatarcontact.com/search?page=1"
	HTMLData := FetchData(url)
	ParsedHTML := HTMLParser(HTMLData)
	var wg sync.WaitGroup
	wg.Add(1)
	go GenerateData(ParsedHTML, &wg)

	go func() {
		wg.Wait()
		defer close(DetailsChannel)
	}()

	details := []Detail{}
	for detail := range DetailsChannel {
		details = append(details, detail)
	}

	saveToJson(details)
	saveAsCsv(details)
}
