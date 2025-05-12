package search

import (
	"fmt"

	"github.com/angkaberapa/Tubes2_BE_BrBaloni-Lulilolli/internal/app/scraper"
)

func stringToElement(elements map[string]*scraper.Element, name string) *scraper.Element {
	return elements[name]
}

func isBasicElementByName(element string) bool {
	return element == "Water" || element == "Earth" || element == "Fire" || element == "Air"
}

func isAllBasicElement(elements []string) bool {
	if len(elements) == 0 {
		return false
	}
	for _, element := range elements {
		if !isBasicElementByName(element) {
			return false
		}
	}
	return true
}

func findCombinationRouteBFS(elements map[string]*scraper.Element, target string) ([][]string, int) {
	targetElement := stringToElement(elements, target)
	if isBasicElementByName(targetElement.Name) {
		return [][]string{[]string{targetElement.Name, "-", "-"}}, 1
	}

	combinationsNumber := len(targetElement.Combinations)
	combinationsChecked := 1
	base := [][]string{}
	checker := [][]string{}
	for i := 0; i < combinationsNumber; i++ {
		combinationsChecked += 2

		if stringToElement(elements, targetElement.Combinations[i].LeftName).Tier < targetElement.Tier && stringToElement(elements, targetElement.Combinations[i].RightName).Tier < targetElement.Tier {
			base = append(base, []string{targetElement.Name, targetElement.Combinations[i].LeftName, targetElement.Combinations[i].RightName})
			checker = append(checker, []string{targetElement.Combinations[i].LeftName, targetElement.Combinations[i].RightName})
		} else {
			fmt.Println("kombinasi tidak valid")
		}
	}

	bfsFound := false
	for i := 0; i < len(checker); i++ {
		if isAllBasicElement(checker[i]) {
			bfsFound = true
			break
		}
	}

	for !bfsFound {
		currentCheck := checker[0]
		checker = checker[1:]

		for i := 0; i < len(currentCheck); i++ {
			addToCheck := []string{}

			if !isBasicElementByName(currentCheck[i]) {
				combinationsChecked += 2

				currentElement := stringToElement(elements, currentCheck[i])
				combinationsNumber = len(currentElement.Combinations)

				for i := 0; i < combinationsNumber; i++ {
					addToCheck = append(addToCheck, currentElement.Combinations[i].LeftName)
					addToCheck = append(addToCheck, currentElement.Combinations[i].RightName)

					leftElement := stringToElement(elements, currentElement.Combinations[i].LeftName)
					rightElement := stringToElement(elements, currentElement.Combinations[i].RightName)

					if leftElement.Tier < currentElement.Tier && rightElement.Tier < currentElement.Tier {
						toAdd := []string{currentElement.Name, currentElement.Combinations[i].LeftName, currentElement.Combinations[i].RightName}
						base = append(base, toAdd)
					}
				}
			}

			checker = append(checker, addToCheck)
		}

		for i := 0; i < len(checker); i++ {
			fmt.Println("tes", checker[i])
			if isAllBasicElement(checker[i]) {
				bfsFound = true
				break
			}
		}

		if bfsFound {
			basicElements := [][]string{[]string{"Fire", "-", "-"}, []string{"Water", "-", "-"}, []string{"Earth", "-", "-"}, []string{"Air", "-", "-"}}
			for i := 0; i < len(basicElements); i++ {
				base = append(base, basicElements[i])
			}
			return base, combinationsChecked
		}
	}

	return base, combinationsChecked
}

func createFixList(list [][]string) ([][]string, []string) {
	fixList := [][]string{[]string{"Water", "-", "-"}, []string{"Earth", "-", "-"},
		[]string{"Fire", "-", "-"}, []string{"Air", "-", "-"}}
	exists := []string{"Water", "Earth", "Fire", "Air"}
	for i := len(list) - 4; i >= 0; i-- {
		toBeAdded1 := false
		toBeAdded2 := false
		alreadyExists := false
		for j := 0; j < len(exists); j++ {
			if list[i][0] == exists[j] {
				alreadyExists = true
			}
			if list[i][1] == exists[j] {
				toBeAdded1 = true
			}
			if list[i][2] == exists[j] {
				toBeAdded2 = true
			}
		}
		if toBeAdded1 && toBeAdded2 {
			fixList = append(fixList, list[i])
			if !alreadyExists {
				exists = append(exists, list[i][0])
			}
		}
	}

	output := [][]string{}
	for i := len(fixList) - 1; i >= 0; i-- {
		output = append(output, fixList[i])
	}
	fmt.Println("output", output)
	return output, exists
}

func createMap(list [][]string, mapList []string) map[string][][]string {
	m := make(map[string][][]string)
	for i := 0; i < len(mapList); i++ {
		if !isBasicElementByName(mapList[i]) {
			m[mapList[i]] = [][]string{}
			for j := 0; j < len(list); j++ {
				if list[j][0] == mapList[i] {
					m[mapList[i]] = append(m[mapList[i]], []string{list[j][1], list[j][2]})
				}
			}
		}
	}

	fmt.Println("Map:", m)
	return m
}

func expandNode(graph map[string][][]string, name string, parentID int, ID *int) ([]GraphNode, []GraphEdge) {
	nodesToAdd := []GraphNode{}
	edgesToAdd := []GraphEdge{}

	(*ID)++
	nodesToAdd = append(nodesToAdd, GraphNode{ID: *ID, Label: name})
	edgesToAdd = append(edgesToAdd, GraphEdge{From: parentID, To: *ID})

	fmt.Println("toExpand:", name)
	if !isBasicElementByName(name) {
		edges := graph[name]
		parentID = *ID
		if len(edges) == 1 {
			for j := 0; j < len(edges[0]); j++ {
				nodesToAddRecursive, edgesToAddRecursive := expandNode(graph, edges[0][j], parentID, ID)
				for k := 0; k < len(nodesToAddRecursive); k++ {
					nodesToAdd = append(nodesToAdd, nodesToAddRecursive[k])
				}
				for k := 0; k < len(edgesToAddRecursive); k++ {
					edgesToAdd = append(edgesToAdd, edgesToAddRecursive[k])
				}
			}
		} else {
			for i := 0; i < len(edges); i++ {
				(*ID)++
				nodesToAdd = append(nodesToAdd, GraphNode{ID: *ID, Label: "+"})
				edgesToAdd = append(edgesToAdd, GraphEdge{From: parentID, To: *ID})
				plusID := *ID
				for j := 0; j < len(edges[i]); j++ {
					nodesToAddRecursive, edgesToAddRecursive := expandNode(graph, edges[i][j], plusID, ID)
					for k := 0; k < len(nodesToAddRecursive); k++ {
						nodesToAdd = append(nodesToAdd, nodesToAddRecursive[k])
					}
					for k := 0; k < len(edgesToAddRecursive); k++ {
						edgesToAdd = append(edgesToAdd, edgesToAddRecursive[k])
					}
				}
			}
		}
	}

	fmt.Println("nodesToAdd", nodesToAdd)
	fmt.Println("edgesToAdd", edgesToAdd)
	return nodesToAdd, edgesToAdd
}

func createBFSNode(graph map[string][][]string, list []string) ([]GraphNode, []GraphEdge) {
	ID := -1
	nodes := []GraphNode{}
	edges := []GraphEdge{}
	nodesToAdd, edgesToAdd := expandNode(graph, list[len(list)-1], ID, &ID)
	for j := 0; j < len(nodesToAdd); j++ {
		nodes = append(nodes, nodesToAdd[j])
	}
	for j := 0; j < len(edgesToAdd); j++ {
		edges = append(edges, edgesToAdd[j])
	}

	return nodes, edges[1:]
}

func findBFS(base [][]string) ([]GraphNode, []GraphEdge) {
	fixList, existsList := createFixList(base)
	m := createMap(fixList, existsList)
	nodes, edges := createBFSNode(m, existsList)

	return nodes, edges
}
func BFS(target string, maxrecipe int) ([]GraphNode, []GraphEdge, int, error) {
	elements, err := scraper.LoadElementsFromFile()
	if err != nil {
		return nil, nil, -1, err
	}

	results, nodeCount := findCombinationRouteBFS(elements, target)
	fmt.Println("result", results)
	// masukkan kode disini
	if results == nil {
		return nil, nil, nodeCount, fmt.Errorf("no valid combination found")
	}

	nodes, edges := findBFS(results)
	return nodes, edges, nodeCount, nil
}
