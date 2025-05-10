package search

import (
	"fmt"

	"github.com/angkaberapa/Tubes2_BE_BrBaloni-Lulilolli/internal/app/scraper"
)

func DFSMultipleRecipe(target string, maxrecipe int) ([]interface{}, error, int, int) { // ganti return value nya
	elements, err := scraper.LoadElementsFromFile()
	if err != nil {
		return nil, err, 0, 0
	}

	results, nodeCount, totalRecipe := findMultipleRouteDFS(elements[target], maxrecipe)
	// masukkan kode disini
	fmt.Println("Results:", printInterface(results))
	fmt.Println("Total Recipe:", totalRecipe)
	// Asumsikan Anda sudah memiliki:
	// var modifiedDfsResult interface{} // Hasil dari findCombinationRouteDFSConcurrentModified
	// var targetElementName string      // Misal "Dust" atau "Aquarium"
	// var allElements map[string]*scraper.Element // Dimuat sekali

	nodes, edges, err := TranslateMultiRecipeOutputToGraph(target, results)
	if err != nil {
		fmt.Println("Error translating:", err)
	} else {
		fmt.Println("Nodes:", nodes)
		fmt.Println("Edges:", edges)
	}
	return results, nil, nodeCount, totalRecipe
}

// findCombinationRoute finds a route from basic elements to the target element
func findMultipleRouteDFS(target *scraper.Element, maxrecipe int) ([]interface{}, int, int) {
	combinationsChecked := 0
	elements, err := scraper.LoadElementsFromFile()
	if err != nil {
		return nil, 0, 0
	}
	memopath := make(map[string][]interface{})
	memocountrecipe := make(map[string]int)
	var dfs func(element *scraper.Element, tier int) ([]interface{}, int)
	dfs = func(element *scraper.Element, tier int) ([]interface{}, int) {
		var recipes []interface{}
		var recipesCount int = 0
		if cachedRecipes, found := memopath[element.Name]; found {
			return cachedRecipes, memocountrecipe[element.Name]
		}
		combinationsChecked++ // Increment for each element visited

		// If the element is a basic element, add it to the route and return true
		if isBasicElement(element) {
			result := []interface{}{element.Name}
			memopath[element.Name] = result
			memocountrecipe[element.Name] = 1
			return result, 1
		}

		// Iterate over the combinations that produce this element
		for _, combo := range element.Combinations {
			leftElement := elements[combo.LeftName]
			rightElement := elements[combo.RightName]

			if leftElement.Tier >= tier || rightElement.Tier >= tier {
				// fmt.Println("kombinasi tidak valid")
				continue
			}

			leftRoute, totalRecipeLeft := dfs(leftElement, tier-1)
			rightRoute, totalRecipeRight := dfs(rightElement, tier-1)
			if len(leftRoute) == 0 || len(rightRoute) == 0 {
				// fmt.Println("kombinasi tidak valid")
				continue // Kombinasi ini tidak bisa dibuat
			}

			if leftRoute != nil && rightRoute != nil {
				var leftRecipePart, rightRecipePart interface{}
				if len(leftRoute) == 1 {
					if basicName, ok := leftRoute[0].(string); ok {
						// code block
						leftRecipePart = basicName // Ini adalah nama elemen dasar
					} else {
						leftRecipePart = map[string]interface{}{"name": leftElement.Name, "recipe": leftRoute}
					}
				} else {
					leftRecipePart = map[string]interface{}{"name": leftElement.Name, "recipe": leftRoute}
				}

				if len(rightRoute) == 1 {
					if basicName, ok := rightRoute[0].(string); ok {
						// code block
						rightRecipePart = basicName // Ini adalah nama elemen dasar
					} else {
						rightRecipePart = map[string]interface{}{"name": rightElement.Name, "recipe": rightRoute}
					}
				} else {
					rightRecipePart = map[string]interface{}{"name": rightElement.Name, "recipe": rightRoute}
				}

				// bukan return, tapi tambahkan ke recipes
				recipes = append(recipes, []interface{}{leftRecipePart, rightRecipePart})
				recipesCount += totalRecipeLeft * totalRecipeRight

				// kalau sudah mencapai batas, return
				if recipesCount >= maxrecipe {
					return recipes, recipesCount
				}
			}
		}
		memopath[element.Name] = recipes
		memocountrecipe[element.Name] = recipesCount
		// kalau recipesCount kosong
		if recipesCount == 0 {
			return nil, 0
		}
		return recipes, recipesCount
	}
	multipleRoute, totalRecipe := dfs(target, target.Tier)
	return multipleRoute, combinationsChecked, totalRecipe
}
