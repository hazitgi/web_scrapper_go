package scrapper

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

var HTMLData Response

func FetchData(url string) []byte {
	fmt.Println("Fetching Data from URL")
	res, err := http.Get(url)
	if err != nil {
		fmt.Println("Error while fetching data", err)
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
	HTMLData = &HTMLResponse{Data: body}
	return body
}
