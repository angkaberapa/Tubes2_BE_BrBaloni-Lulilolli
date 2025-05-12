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

func findCombinationRouteBFS(elements map[string]*scraper.Element, target string, numberOfRecipe int) ([]GraphNode, []GraphEdge, int) {
	targetElement := stringToElement(elements, target)
	if isBasicElementByName(targetElement.Name) {
		return []GraphNode{GraphNode{ID: 1, Label: targetElement.Name}}, nil, 1
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

		addToCheck := []string{}
		for i := 0; i < len(currentCheck); i++ {
			if !isBasicElementByName(currentCheck[i]) {
				combinationsChecked += 2

				currentElement := stringToElement(elements, currentCheck[i])
				combinationsNumber = len(currentElement.Combinations)

				for i := 0; i < combinationsNumber; i++ {

					leftElement := stringToElement(elements, currentElement.Combinations[i].LeftName)
					rightElement := stringToElement(elements, currentElement.Combinations[i].RightName)

					if leftElement.Tier < currentElement.Tier && rightElement.Tier < currentElement.Tier {
						addToCheck = append(addToCheck, currentElement.Combinations[i].LeftName)
						addToCheck = append(addToCheck, currentElement.Combinations[i].RightName)
						toAdd := []string{currentElement.Name, currentElement.Combinations[i].LeftName, currentElement.Combinations[i].RightName}

						toBeInserted := true
						for j := 0; j < len(base); j++ {
							if toAdd[0] == base[j][0] {
								if (toAdd[1] == base[j][1] && toAdd[2] == base[j][2]) {
									toBeInserted = false
								}
							}
						}

						if (toBeInserted) {
							base = append(base, toAdd)
						}
					}
				}
			} else {
				addToCheck = append(addToCheck, currentCheck[i])
			}
		}
		
		if len(addToCheck) > 0 {
			checker = append(checker, addToCheck)
		}

		fmt.Println(currentCheck, "->", addToCheck)
		for i := 0; i < len(checker); i++ {
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
			fixList, existsList, counter := createFixList(base)
			m := createMap(fixList, existsList)
			mapp := limitRecipe(counter, m, target, numberOfRecipe)

			if (mapp != nil) {
				nodes, edges := createBFSNode(mapp, existsList)
				return nodes, edges, combinationsChecked
			} else {
				bfsFound = false
			}
		}
	}

	return nil, nil, 0
}

func createFixList(list [][]string) ([][]string, []string, map[string][]int) {
	fixList := [][]string{[]string{"Water", "-", "-"}, []string{"Earth", "-", "-"},
		[]string{"Fire", "-", "-"}, []string{"Air", "-", "-"}}
	exists := []string{"Water", "Earth", "Fire", "Air"}
	counter := map[string][]int{"Water": []int{1}, "Earth": []int{1}, "Fire": []int{1}, "Air": []int{1}}

	for i := len(list) - 5; i >= 0; i-- {
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
				counter[list[i][0]] = []int{}
			}
		}
	}

	output := [][]string{}
	for i := len(fixList) - 1; i >= 0; i-- {
		output = append(output, fixList[i])
	}
	fmt.Println("output", output)

	for i := len(output) - 5; i >= 0; i-- {
		counterLeft := 0
		for j := 0; j < len(counter[output[i][1]]); j++ {
			counterLeft += counter[output[i][1]][j]
		}
		counterRight := 0
		for j := 0; j < len(counter[output[i][2]]); j++ {
			counterRight += counter[output[i][2]][j]
		}
		counter[output[i][0]] = append(counter[output[i][0]], counterLeft * counterRight)
	}
	fmt.Println("counter", counter)
	
	return output, exists, counter
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

	// fmt.Println("toExpand:", name)
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

	// fmt.Println("nodesToAdd", nodesToAdd)
	// fmt.Println("edgesToAdd", edgesToAdd)
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

func findBFS(base [][]string, name string, numberOfRecipe int) ([]GraphNode, []GraphEdge) {
	fixList, existsList, counter := createFixList(base)
	m := createMap(fixList, existsList)
	mapp := limitRecipe(counter, m, name, numberOfRecipe)
	nodes, edges := createBFSNode(mapp, existsList)

	return nodes, edges
}	

func reduceMap(fixList map[string][][]string, name string, pair []string) map[string][][]string {

	fmt.Println("pair", pair)
	output := map[string][][]string{}
	output[name] = [][]string{pair}
	
	for i := 0; i < len(pair); i++ {
		if (!isBasicElementByName(pair[i])) {
			output[pair[i]] = fixList[pair[i]]
			for j := 0; j < len(output[pair[i]]); j++ {
				outputToAdd := reduceMap(fixList, pair[i], output[pair[i]][j])

				for key, value := range outputToAdd {
					output[key] = value
				}
			}
		}
	}

	fmt.Println("output in ReduceMap", output)
	return output
}

func limitRecipe(counter map[string][]int, fixList map[string][][]string, name string, limitRecipe int) map[string][][]string {
	output := map[string][][]string{}
	count := 0

	toList := fixList[name]
	for i := 0; i < len(toList); i++ {

		outputToAdd := reduceMap(fixList, name, toList[i])

		for key, value := range outputToAdd {
			if key == name {
				output[key] = append(output[key], toList[i])
			} else {
				output[key] = value
			}
		}

		if (count >= limitRecipe) {
			fmt.Println("output in limitRecipe", output)
			return output
		}

	}

	return output
}

func BFS(target string, maxrecipe int) ([]GraphNode, []GraphEdge, int, error) {
	elements, err := scraper.LoadElementsFromFile()
	if err != nil {
		return nil, nil, -1, err
	}

	nodes, edges, nodeCount := findCombinationRouteBFS(elements, target, 17)
	// masukkan kode disini
	if nodes == nil {
		return nil, nil, nodeCount, fmt.Errorf("no valid combination found")
	}

	return nodes, edges, nodeCount, nil
}
