package main

import (
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type AmazonItem struct {
	Title string
	Price string
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Please provide a URL as an argument.")
		return
	}

	url := os.Args[1]

	item, err := GetAmazonItem(url)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	if item.Price != "" {
		fmt.Println("Price found:", item.Price)
	} else {
		fmt.Println("Price not found")
	}
}

func GetAmazonItem(url string) (*AmazonItem, error) {
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	doc, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		return nil, err
	}

	title := findTitle(doc)
	price := findPrice(doc)

	return &AmazonItem{
		Title: title,
		Price: price,
	}, nil
}

func findTitle(doc *goquery.Document) string {
	return doc.Find("#productTitle").Text()
}

func findPrice(doc *goquery.Document) string {
	priceRegex := regexp.MustCompile(`\$[\d.,]+`)
	var price string

	doc.Find(".a-price, .a-offscreen").Each(func(i int, s *goquery.Selection) {
		text := strings.TrimSpace(s.Text())
		if priceRegex.MatchString(text) {
			price = strings.TrimPrefix(text, "$")
			if price != "" {
				return
			}
		}
	})

	return price
}