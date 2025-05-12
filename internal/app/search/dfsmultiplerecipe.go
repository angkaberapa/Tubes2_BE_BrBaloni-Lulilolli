package search

import (
	"fmt"
	"sync"

	"github.com/angkaberapa/Tubes2_BE_BrBaloni-Lulilolli/internal/app/scraper"
)

func DFSMultipleRecipe(target string, maxrecipe int) ([]interface{}, error, int, int) {
	elements, err := scraper.LoadElementsFromFile()
	if err != nil {
		return nil, err, 0, 0
	}

	results, nodeCount, totalRecipe := findMultipleRouteDFS(elements[target], maxrecipe)
	fmt.Println("Results:", printInterface(results))
	fmt.Println("Total Recipe:", totalRecipe)

	nodes, edges, err := TranslateMultiRecipeOutputToGraph(target, results)
	if err != nil {
		fmt.Println("Error translating:", err)
	} else {
		fmt.Println("Nodes:", nodes)
		fmt.Println("Edges:", edges)
	}
	return results, nil, nodeCount, totalRecipe
}
func findMultipleRouteDFS(target *scraper.Element, maxrecipe int) ([]interface{}, int, int) {
	var combinationsChecked int
	elements, err := scraper.LoadElementsFromFile()
	if err != nil {
		return nil, 0, 0
	}

	memopath := make(map[string][]interface{})
	memocountrecipe := make(map[string]int)
	var memoMutex sync.Mutex
	var countMutex sync.Mutex

	var dfs func(element *scraper.Element, tier int) ([]interface{}, int)
	dfs = func(element *scraper.Element, tier int) ([]interface{}, int) {
		memoMutex.Lock()
		if cachedRecipes, found := memopath[element.Name]; found {
			count := memocountrecipe[element.Name]
			memoMutex.Unlock()
			return cachedRecipes, count
		}
		memoMutex.Unlock()

		countMutex.Lock()
		combinationsChecked++
		countMutex.Unlock()

		if isBasicElement(element) {
			result := []interface{}{element.Name}
			memoMutex.Lock()
			memopath[element.Name] = result
			memocountrecipe[element.Name] = 1
			memoMutex.Unlock()
			return result, 1
		}

		var mu sync.Mutex
		var wg sync.WaitGroup
		var recipes []interface{}
		var recipesCount int

		for _, combo := range element.Combinations {
			leftElement := elements[combo.LeftName]
			rightElement := elements[combo.RightName]

			if leftElement.Tier >= tier || rightElement.Tier >= tier {
				continue
			}

			wg.Add(2)
			var leftRoute, rightRoute []interface{}
			var leftCount, rightCount int

			go func() {
				defer wg.Done()
				r, c := dfs(leftElement, tier-1)
				mu.Lock()
				leftRoute = r
				leftCount = c
				mu.Unlock()
			}()

			go func() {
				defer wg.Done()
				r, c := dfs(rightElement, tier-1)
				mu.Lock()
				rightRoute = r
				rightCount = c
				mu.Unlock()
			}()

			wg.Wait()

			if len(leftRoute) == 0 || len(rightRoute) == 0 {
				continue
			}

			var leftRecipePart, rightRecipePart interface{}
			if len(leftRoute) == 1 {
				if basicName, ok := leftRoute[0].(string); ok {
					leftRecipePart = basicName
				} else {
					leftRecipePart = map[string]interface{}{"name": leftElement.Name, "recipe": leftRoute}
				}
			} else {
				leftRecipePart = map[string]interface{}{"name": leftElement.Name, "recipe": leftRoute}
			}

			if len(rightRoute) == 1 {
				if basicName, ok := rightRoute[0].(string); ok {
					rightRecipePart = basicName
				} else {
					rightRecipePart = map[string]interface{}{"name": rightElement.Name, "recipe": rightRoute}
				}
			} else {
				rightRecipePart = map[string]interface{}{"name": rightElement.Name, "recipe": rightRoute}
			}

			recipes = append(recipes, []interface{}{leftRecipePart, rightRecipePart})
			recipesCount += leftCount * rightCount

			if recipesCount >= maxrecipe {
				break
			}
		}

		memoMutex.Lock()
		memopath[element.Name] = recipes
		memocountrecipe[element.Name] = recipesCount
		memoMutex.Unlock()

		if recipesCount == 0 {
			return nil, 0
		}
		return recipes, recipesCount
	}

	multipleRoute, totalRecipe := dfs(target, target.Tier)
	return multipleRoute, combinationsChecked, totalRecipe
}
