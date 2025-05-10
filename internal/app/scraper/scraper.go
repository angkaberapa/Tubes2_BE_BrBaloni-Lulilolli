package scraper

import (
	"fmt"
	"log"
	"net/http"

	"slices"

	"github.com/PuerkitoBio/goquery"
)

type Element struct {
	Name         string
	CanCreate    []string
	Combinations []*Combination
	Tier         int
}
type Combination struct {
	ResultName string
	LeftName   string
	RightName  string
}

func ScrapeElements() (map[string]*Element, error) {
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
		// buang Time, Ruins, dan Archeologist
		if resultName == "Time" || resultName == "Ruins" || resultName == "Archeologist" {
			return
		}
		elementsName = append(elementsName, resultName)
	})

	// Temukan semua tabel
	doc.Find("table.list-table").Each(func(i int, table *goquery.Selection) {
		// Temukan semua baris tabel

		if i == 1 { // untuk tabel kedua (special elements yaitu time), skip.
			return
		}

		table.Find("tr").Each(func(j int, row *goquery.Selection) {
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
					if slices.Contains(elementsName, leftName) && slices.Contains(elementsName, rightName) {
						combinations = append(combinations, &Combination{
							ResultName: resultName,
							LeftName:   leftName,
							RightName:  rightName,
						})
						totalCombinations++
					}
					// } else { // ada elemen yang dari Myths and Monsters atau dari (Time, Ruins, Archeologist)
					// 	fmt.Println("Combination not found:", leftName, rightName)
					// }
				}
			})
			// buang Ruins dan Archeologist
			if !slices.Contains(elementsName, resultName) {
				return
			}
			currentTier := i
			if i > 1 { // karena untuk tabel kedua (special elements yaitu time) di skip. (ini biar tier nya tidak bertambah aja sih)
				currentTier = i - 1
			}
			if len(combinations) > 0 {
				element := &Element{
					Name:         resultName,
					CanCreate:    nil,
					Combinations: combinations,
					Tier:         currentTier,
				}
				elementsMapByName[resultName] = element
			} else { // memang tidak ada yang bisa membuatnya (contoh: earth)
				element := &Element{
					Name:         resultName,
					CanCreate:    nil,
					Combinations: nil,
					Tier:         currentTier,
				}
				elementsMapByName[resultName] = element
			}
		})
	})

	// Isi CanCreate untuk setiap elemen
	for _, element := range elementsMapByName {
		for _, combination := range element.Combinations {
			leftElement := elementsMapByName[combination.LeftName]
			rightElement := elementsMapByName[combination.RightName]

			if !slices.Contains(leftElement.CanCreate, element.Name) {
				leftElement.CanCreate = append(leftElement.CanCreate, element.Name)
			}

			if !slices.Contains(rightElement.CanCreate, element.Name) {
				rightElement.CanCreate = append(rightElement.CanCreate, element.Name)
			}
		}
	}

	// Print total elements and combinations
	fmt.Printf("\nTotal elements: %d\n", len(elementsMapByName))
	fmt.Printf("Total combinations: %d\n\n", totalCombinations)
	fmt.Println("✅ Parsing completed!")
	return elementsMapByName, nil
	// contoh pemakaian
	// time := elementsMapByName["Time"]
}
