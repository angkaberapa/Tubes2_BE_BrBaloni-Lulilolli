package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

type Element struct {
	Name         string
	CanCreate    []string
	Combinations []*Combination
}
type Combination struct {
	ResultName string
	LeftName   string
	RightName  string
}

// Helper function to check if a string slice contains a specific value
func contains(slice []string, value string) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}
	return false
}

func main() {
	url := "https://little-alchemy.fandom.com/wiki/Elements_(Little_Alchemy_2)"

	// Request HTTP
	res, err := http.Get(url)
	if err != nil {
		log.Fatalf("Failed to fetch page: %v", err)
	}
	defer res.Body.Close()
	fmt.Printf("🔗 Fetching %s...\n", url)

	if res.StatusCode != 200 {
		log.Fatalf("Status code error: %d %s", res.StatusCode, res.Status)
	}
	fmt.Println("✅ Page fetched successfully!")

	// Parse HTML
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatalf("Failed to read HTML: %v", err)
	}
	fmt.Println("📖 Parsing HTML...")

	var elementsMapByName = map[string]*Element{}
	var totalCombinations int = 0

	// Temukan semua nama elemen dulu
	var elementsName []string
	doc.Find("table.list-table tr").Each(func(i int, row *goquery.Selection) {
		cells := row.Find("td")
		if cells.Length() < 2 {
			return
		}
		resultName := cells.Eq(0).ChildrenFiltered("a").First().Text()
		elementsName = append(elementsName, resultName)
	})

	// Temukan semua baris tabel
	doc.Find("table.list-table tr").Each(func(i int, row *goquery.Selection) {
		cells := row.Find("td")
		if cells.Length() < 2 {
			return
		}
		resultName := cells.Eq(0).ChildrenFiltered("a").First().Text()
		var combinations []*Combination
		cells.Eq(1).Find("li").Each(func(_ int, li *goquery.Selection) {
			var parts []string
			li.ChildrenFiltered("a").Each(func(_ int, a *goquery.Selection) {
				parts = append(parts, a.Text())
			})
			if len(parts) == 2 {
				leftName := parts[0]
				rightName := parts[1]
				// check if leftName and rightName is in elementsName
				if contains(elementsName, leftName) && contains(elementsName, rightName) {
					combinations = append(combinations, &Combination{
						ResultName: resultName,
						LeftName:   leftName,
						RightName:  rightName,
					})
					totalCombinations++
				}
				// } else { // ada elemen yang dari Myths and Monsters
				// 	fmt.Println("Combination not found:", leftName, rightName)
				// }
			}
		})

		if len(combinations) > 0 {
			element := &Element{
				Name:         resultName,
				CanCreate:    nil,
				Combinations: combinations,
			}
			elementsMapByName[resultName] = element
		} else { // memang tidak ada yang bisa membuatnya (contoh: time)
			element := &Element{
				Name:         resultName,
				CanCreate:    nil,
				Combinations: nil,
			}
			elementsMapByName[resultName] = element
		}
	})

	// Isi CanCreate untuk setiap elemen
	for _, element := range elementsMapByName {
		for _, combination := range element.Combinations {
			leftElement := elementsMapByName[combination.LeftName]
			rightElement := elementsMapByName[combination.RightName]

			if !contains(leftElement.CanCreate, element.Name) {
				leftElement.CanCreate = append(leftElement.CanCreate, element.Name)
			}

			if !contains(rightElement.CanCreate, element.Name) {
				rightElement.CanCreate = append(rightElement.CanCreate, element.Name)
			}
		}
	}

	// Print total elements and combinations
	fmt.Printf("\nTotal elements: %d\n", len(elementsMapByName))
	fmt.Printf("Total combinations: %d\n\n", totalCombinations)
	fmt.Println("✅ Parsing completed!")

	// contoh pemakaian
	// time := elementsMapByName["Time"]
}
