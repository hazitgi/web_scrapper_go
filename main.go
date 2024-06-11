package main

import (
	"fmt"
	"time"
	"github.com/hazitgi/web_scrapper_go/packages"
)

func main() {
	start := time.Now()
	fmt.Println("******** starting ********")
	scrapper.RunScrapper()

	end := time.Since(start)
	fmt.Printf("\nTime taken: %.2f minutes\n", end.Minutes())
}
