package search

import (
	"fmt"

	"github.com/angkaberapa/Tubes2_BE_BrBaloni-Lulilolli/internal/app/scraper"
)


func DFS(target string, maxrecipe int) (interface{}, error, int) { // ganti return value nya
	elements, err := scraper.LoadElementsFromFile()
	if err != nil {
		return nil, err
	}

	results, nodeCount := findCombinationRouteDFS(elements[target])
	// masukkan kode disini

	return results, nil, nodeCount
}

// findCombinationRoute finds a route from basic elements to the target element
func findCombinationRouteDFS(target *Element) (interface{}, int) {
    combinationsChecked := 0

    var dfs func(element *Element, tier int) interface{}
    dfs = func(element *Element, tier int) interface{} {
        combinationsChecked++ // Increment for each element visited

        // If the element is a basic element, add it to the route and return true
        if isBasicElement(element) {
            return element.Name
        }

        // Iterate over the combinations that produce this element
        for _, combo := range element.Combinations {
            leftElement := combo.Left
            rightElement := combo.Right

            if leftElement.Tier >= tier || rightElement.Tier >= tier {
                // fmt.Println("kombinasi tidak valid")
                continue
            }

            leftRoute := dfs(combo.Left, tier - 1)
            rightRoute := dfs(combo.Right, tier - 1)

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