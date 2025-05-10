package search

import (
	"fmt"

	"github.com/angkaberapa/Tubes2_BE_BrBaloni-Lulilolli/internal/app/scraper"
)

func DFS(target string, maxrecipe int) (interface{}, error, int) { // ganti return value nya
	elements, err := scraper.LoadElementsFromFile()
	if err != nil {
		return nil, err, 0
	}

	results, nodeCount := findCombinationRouteDFS(elements[target])
	// masukkan kode disini
	fmt.Println("Results:", printInterface(results))
	// Asumsikan Anda sudah memiliki:
	// var modifiedDfsResult interface{} // Hasil dari findCombinationRouteDFSConcurrentModified
	// var targetElementName string      // Misal "Dust" atau "Aquarium"
	// var allElements map[string]*scraper.Element // Dimuat sekali

	nodes, edges, err := TranslateOutputPathToGraph(results, target)
	if err != nil {
		fmt.Println("Error translating:", err)
	} else {
		fmt.Println("Nodes:", nodes)
		fmt.Println("Edges:", edges)
	}
	return results, nil, nodeCount
}

// findCombinationRoute finds a route from basic elements to the target element
func findCombinationRouteDFS(target *scraper.Element) (interface{}, int) {
	combinationsChecked := 0
	elements, err := scraper.LoadElementsFromFile()
	if err != nil {
		return nil, 0
	}

	var dfs func(element *scraper.Element, tier int) interface{}
	dfs = func(element *scraper.Element, tier int) interface{} {
		combinationsChecked++ // Increment for each element visited

		// If the element is a basic element, add it to the route and return true
		if isBasicElement(element) {
			return element.Name
		}

		// Iterate over the combinations that produce this element
		for _, combo := range element.Combinations {
			leftElement := elements[combo.LeftName]
			rightElement := elements[combo.RightName]

			if leftElement.Tier >= tier || rightElement.Tier >= tier {
				// fmt.Println("kombinasi tidak valid")
				continue
			}

			leftRoute := dfs(leftElement, tier-1)
			rightRoute := dfs(rightElement, tier-1)

			if leftRoute != nil && rightRoute != nil {
				var leftRecipePart, rightRecipePart interface{}
				if _, ok := leftRoute.(string); ok {
					leftRecipePart = leftRoute // Ini adalah nama elemen dasar
				} else {
					leftRecipePart = map[string]interface{}{"name": leftElement.Name, "recipe": leftRoute}
				}

				if _, ok := rightRoute.(string); ok {
					rightRecipePart = rightRoute // Ini adalah nama elemen dasar
				} else {
					rightRecipePart = map[string]interface{}{"name": rightElement.Name, "recipe": rightRoute}
				}

				return []interface{}{leftRecipePart, rightRecipePart}
			}
		}

		// If no route is found, return false
		return nil
	}
	route := dfs(target, target.Tier)
	return route, combinationsChecked
}

// ini nyoba2 concurrent
// func DFSConcurrent(target string, maxrecipe int) (interface{}, error, int) {
// 	elements, err := scraper.LoadElementsFromFile()
// 	if err != nil {
// 		return nil, err, 0
// 	}

// 	results, nodeCount := findCombinationRouteDFSConcurrent(elements[target])
// 	fmt.Println("Results:", printInterface(results))
// 	return results, nil, nodeCount
// }

// func findCombinationRouteDFSConcurrent(target *scraper.Element) (interface{}, int) {
// 	elements, err := scraper.LoadElementsFromFile()
// 	if err != nil {
// 		return nil, 0
// 	}

// 	var mu sync.Mutex
// 	var combinationsChecked int
// 	sem := make(chan struct{}, 1000)

// 	var dfs func(element *scraper.Element, tier int) interface{}
// 	dfs = func(element *scraper.Element, tier int) interface{} {
// 		mu.Lock()
// 		combinationsChecked++
// 		mu.Unlock()

// 		if isBasicElement(element) {
// 			return element.Name
// 		}

// 		for _, combo := range element.Combinations {
// 			// Mengambil elemen dari map 'elements' yang diload di awal fungsi ini
// 			leftElementData, okLeft := elements[combo.LeftName]
// 			if !okLeft {
// 				// Mungkin log atau tangani elemen tidak ditemukan
// 				continue
// 			}
// 			rightElementData, okRight := elements[combo.RightName]
// 			if !okRight {
// 				// Mungkin log atau tangani elemen tidak ditemukan
// 				continue
// 			}

// 			// Penting untuk membuat salinan variabel yang akan digunakan dalam closure goroutine
// 			// untuk menghindari masalah race condition pada variabel loop.
// 			currentLeftElement := leftElementData
// 			currentRightElement := rightElementData

// 			if currentLeftElement.Tier >= tier || currentRightElement.Tier >= tier {
// 				continue
// 			}

// 			var leftResult, rightResult interface{}
// 			var wg sync.WaitGroup // WaitGroup untuk menunggu goroutine kiri dan kanan

// 			// Jalankan DFS kiri
// 			select {
// 			case sem <- struct{}{}: // Coba ambil slot dari semaphore
// 				wg.Add(1)
// 				go func(el *scraper.Element, t int) {
// 					defer func() {
// 						<-sem // Lepaskan slot semaphore
// 						wg.Done()
// 					}()
// 					leftResult = dfs(el, t)
// 				}(currentLeftElement, tier-1) // Lewatkan salinan
// 			default:
// 				// Jika semaphore penuh, jalankan secara sinkron
// 				leftResult = dfs(currentLeftElement, tier-1)
// 			}

// 			// Jalankan DFS kanan
// 			select {
// 			case sem <- struct{}{}: // Coba ambil slot dari semaphore
// 				wg.Add(1)
// 				go func(el *scraper.Element, t int) {
// 					defer func() {
// 						<-sem // Lepaskan slot semaphore
// 						wg.Done()
// 					}()
// 					rightResult = dfs(el, t)
// 				}(currentRightElement, tier-1) // Lewatkan salinan
// 			default:
// 				// Jika semaphore penuh, jalankan secara sinkron
// 				rightResult = dfs(currentRightElement, tier-1)
// 			}

// 			wg.Wait() // Tunggu kedua goroutine (jika ada yang diluncurkan) selesai

// 			if leftResult != nil && rightResult != nil {
// 				return []interface{}{
// 					map[string]interface{}{"name": currentLeftElement.Name, "recipe": leftResult},
// 					map[string]interface{}{"name": currentRightElement.Name, "recipe": rightResult},
// 				}
// 			}
// 		}

// 		return nil
// 	}

// 	route := dfs(target, target.Tier)
// 	return route, combinationsChecked
// }
