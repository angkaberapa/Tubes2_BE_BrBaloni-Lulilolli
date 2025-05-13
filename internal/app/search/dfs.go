package search

import (
	"fmt"
	"sort"

	"github.com/angkaberapa/Tubes2_BE_BrBaloni-Lulilolli/internal/app/scraper"
)

func DFS(target string, maxrecipe int) (interface{}, error, int) {
	elements, err := scraper.LoadElementsFromFile()
	if err != nil {
		return nil, err, 0
	}

	results, nodeCount := findCombinationRouteDFS(elements[target])
	nodes, edges, err := TranslateOutputPathToGraph(results, target)
	if err != nil {
		fmt.Println("Error translating:", err)
	} else {
		fmt.Println("Nodes:", nodes)
		fmt.Println("Edges:", edges)
	}
	return results, nil, nodeCount
}

func findCombinationRouteDFS(target *scraper.Element) (interface{}, int) {
	combinationsChecked := 0
	elements, err := scraper.LoadElementsFromFile()
	if err != nil {
		return nil, 0
	}

	var dfs func(element *scraper.Element, tier int) interface{}
	dfs = func(element *scraper.Element, tier int) interface{} {
		combinationsChecked++

		if isBasicElement(element) {
			return element.Name
		}
		// sort combinations based on sum of its LeftName.Tier and RightName.Tier
		sort.Slice(element.Combinations, func(i, j int) bool {
			left := elements[element.Combinations[i].LeftName]
			right := elements[element.Combinations[i].RightName]
			return left.Tier+right.Tier < elements[element.Combinations[j].LeftName].Tier+elements[element.Combinations[j].RightName].Tier
		})
		for _, combo := range element.Combinations {
			leftElement := elements[combo.LeftName]
			rightElement := elements[combo.RightName]

			if leftElement.Tier >= tier || rightElement.Tier >= tier {
				continue
			}

			leftRoute := dfs(leftElement, tier-1)
			rightRoute := dfs(rightElement, tier-1)

			if leftRoute != nil && rightRoute != nil {
				var leftRecipePart, rightRecipePart interface{}
				if _, ok := leftRoute.(string); ok {
					leftRecipePart = leftRoute
				} else {
					leftRecipePart = map[string]interface{}{"name": leftElement.Name, "recipe": leftRoute}
				}

				if _, ok := rightRoute.(string); ok {
					rightRecipePart = rightRoute
				} else {
					rightRecipePart = map[string]interface{}{"name": rightElement.Name, "recipe": rightRoute}
				}

				return []interface{}{leftRecipePart, rightRecipePart}
			}
		}
		// tidak ketemu
		return nil
	}
	route := dfs(target, target.Tier)
	return route, combinationsChecked
}
