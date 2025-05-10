package search

import (
	"sync"

	"github.com/angkaberapa/Tubes2_BE_BrBaloni-Lulilolli/internal/app/scraper"
)

func DFS(target string, maxrecipe int) (interface{}, error, int) { // ganti return value nya
	elements, err := scraper.LoadElementsFromFile()
	if err != nil {
		return nil, err, 0
	}

	results, nodeCount := findCombinationRouteDFS(elements[target])
	// masukkan kode disini

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
				return []interface{}{leftRoute, rightRoute}
			}

		}

		// If no route is found, return false
		return nil
	}

	route := dfs(target, target.Tier)

	return route, combinationsChecked
}

func DFSConcurrent(target string, maxrecipe int) (interface{}, error, int) {
	elements, err := scraper.LoadElementsFromFile()
	if err != nil {
		return nil, err, 0
	}

	results, nodeCount := findCombinationRouteDFSConcurrent(elements[target])
	return results, nil, nodeCount
}

func findCombinationRouteDFSConcurrent(target *scraper.Element) (interface{}, int) {
	elements, err := scraper.LoadElementsFromFile()
	if err != nil {
		return nil, 0
	}

	var mu sync.Mutex
	var combinationsChecked int
	sem := make(chan struct{}, 100) // maksimal 2 goroutine aktif

	var dfs func(element *scraper.Element, tier int) interface{}
	dfs = func(element *scraper.Element, tier int) interface{} {
		mu.Lock()
		combinationsChecked++
		mu.Unlock()

		if isBasicElement(element) {
			return element.Name
		}

		for _, combo := range element.Combinations {
			leftElement := elements[combo.LeftName]
			rightElement := elements[combo.RightName]

			if leftElement.Tier >= tier || rightElement.Tier >= tier {
				continue
			}

			var leftResult, rightResult interface{}

			// jalankan DFS kiri dalam goroutine (jika tersedia)
			done := make(chan struct{})

			select {
			case sem <- struct{}{}:
				go func() {
					leftResult = dfs(leftElement, tier-1)
					<-sem
					close(done)
				}()
				<-done
			default:
				leftResult = dfs(leftElement, tier-1)
			}

			// kanan selalu langsung
			rightResult = dfs(rightElement, tier-1)

			if leftResult != nil && rightResult != nil {
				return []interface{}{leftResult, rightResult}
			}
		}

		return nil
	}

	route := dfs(target, target.Tier)
	return route, combinationsChecked
}

// package search

// import (
// 	"fmt"
// 	"sync"

// 	"github.com/angkaberapa/Tubes2_BE_BrBaloni-Lulilolli/internal/app/scraper"
// )

// type Recipe struct {
// 	Steps []string
// }

// func DFS(target string, max int) ([]Recipe, error) {
// 	elements, err := scraper.LoadElementsFromFile()
// 	if err != nil {
// 		return nil, err
// 	}

// 	start, ok := elements[target]
// 	if !ok {
// 		return nil, fmt.Errorf("element %s not found", target)
// 	}

// 	var results []Recipe
// 	resultsCh := make(chan Recipe)
// 	doneCh := make(chan struct{})
// 	var wg sync.WaitGroup

// 	// DFS rekursif pakai goroutine
// 	var dfs func(current *scraper.Element, path []string)
// 	dfs = func(current *scraper.Element, path []string) {
// 		defer wg.Done()

// 		for _, comb := range current.Combinations {
// 			step := fmt.Sprintf("%s + %s = %s", comb.LeftName, comb.RightName, comb.ResultName)
// 			newPath := append([]string{step}, path...)

// 			left := elements[comb.LeftName]
// 			right := elements[comb.RightName]

// 			if len(left.Combinations) == 0 && len(right.Combinations) == 0 {
// 				// Leaf node: end of recipe path
// 				select {
// 				case resultsCh <- Recipe{Steps: newPath}:
// 				case <-doneCh:
// 					return
// 				}
// 			} else {
// 				// Explore left and right in parallel
// 				// wg.Add(2)
// 				// go dfs(left, newPath)
// 				// go dfs(right, newPath)
// 				dfs(left, newPath)
// 				dfs(right, newPath)
// 			}
// 		}
// 	}

// 	wg.Add(1)
// 	go dfs(start, []string{})

// 	// Collector goroutine
// 	go func() {
// 		wg.Wait()
// 		close(resultsCh)
// 	}()

// 	for recipe := range resultsCh {
// 		results = append(results, recipe)
// 		if len(results) >= max {
// 			close(doneCh) // Sinyal semua goroutine untuk berhenti
// 			break
// 		}
// 	}

// 	return results, nil
// }
