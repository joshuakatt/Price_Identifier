package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Please provide a URL as an argument.")
		return
	}

	url := os.Args[1]

	// Adjust the priceKeywords and priceElements based on the specific webpage
	priceKeywords := []string{"price", "cost"}
	priceElements := []string{"span", "div", "p"}

	price := GetPrice(url, priceKeywords, priceElements)
	if price != "" {
		fmt.Println("Price found:", price)
	} else {
		fmt.Println("Price not found")
	}
}

func GetPrice(url string, priceKeywords []string, priceElements []string) string {
	// Send an HTTP GET request to the webpage
	response, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	// Parse the HTML document
	doc, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	// Find the main price/cost component in the webpage
	price := findMainPrice(doc, priceKeywords, priceElements)

	return price
}

func findMainPrice(doc *goquery.Document, priceKeywords []string, priceElements []string) string {
	// Iterate through each keyword and element to find the main price
	for _, keyword := range priceKeywords {
		for _, element := range priceElements {
			price := findPriceWithKeyword(doc, keyword, element)
			if price != "" {
				return price
			}
		}
	}

	return ""
}

func findPriceWithKeyword(doc *goquery.Document, keyword string, element string) string {
	price := ""

	doc.Find(element).Each(func(i int, s *goquery.Selection) {
		text := strings.ToLower(strings.TrimSpace(s.Text()))
		if strings.Contains(text, keyword) {
			price = extractPriceFromText(text)
			if price != "" {
				return
			}
		}
	})

	return price
}

func extractPriceFromText(text string) string {
	// Regular expression to extract a price pattern
	priceRegex := regexp.MustCompile(`\$\d+(\.\d+)?`)

	// Find all occurrences of the price pattern in the text
	matches := priceRegex.FindAllString(text, -1)

	// Return the first matched price or an empty string if not found
	if len(matches) > 0 {
		return matches[0]
	}

	return ""
}
